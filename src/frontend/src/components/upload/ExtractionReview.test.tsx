import { fireEvent, render, screen, waitFor, within } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { ExtractionReview } from "./ExtractionReview";
import type { ExtractionResults, ProcessDocumentIds } from "./types";

vi.mock("./FeedbackForm", () => ({
  FeedbackForm: ({ userId, fileId }: { userId: string; fileId: string }) => (
    <div data-testid="feedback-form">
      Feedback for {userId} / {fileId}
    </div>
  ),
}));

const RESUME_DATA: ExtractionResults["resume"] = {
  id: "resume-1",
  status: "COMPLETED",
  extractedData: {
    name: "John Doe",
    email: "john@example.com",
    phone: "+49 123 456",
    location: "Berlin, Germany",
    summary: "Experienced software engineer with 10 years of experience.",
    experiences: [
      {
        company: "Acme Corp",
        title: "Senior Engineer",
        location: "Berlin",
        startDate: "2020-01",
        endDate: null,
        isCurrent: true,
        description: "Led the backend team",
      },
      {
        company: "StartupCo",
        title: "Software Engineer",
        location: "Munich",
        startDate: "2018-06",
        endDate: "2019-12",
        isCurrent: false,
        description: null,
      },
    ],
    educations: [
      {
        institution: "MIT",
        degree: "Bachelor of Science",
        field: "Computer Science",
        startDate: "2014-09",
        endDate: "2018-05",
        gpa: "3.9",
        achievements: null,
      },
    ],
    skills: ["Go", "TypeScript", "PostgreSQL"],
    extractedAt: "2025-01-01T00:00:00Z",
    confidence: 0.95,
  },
  errorMessage: null,
};

const LETTER_DATA: ExtractionResults["referenceLetter"] = {
  id: "letter-1",
  status: "COMPLETED",
  extractedData: {
    author: {
      name: "Jane Smith",
      title: "CTO",
      company: "TechCo",
      relationship: "Former Manager",
    },
    testimonials: [
      {
        quote: "John is an exceptional engineer.",
        skillsMentioned: ["Go", "TypeScript"],
      },
      { quote: "He excels at problem-solving.", skillsMentioned: null },
    ],
    skillMentions: [
      { skill: "Go", quote: "Expert-level Go programming", context: "technical skills" },
    ],
    experienceMentions: [
      { company: "TechCo", role: "Senior Engineer", quote: "Led the backend team" },
    ],
    discoveredSkills: [
      {
        skill: "System Design",
        quote: "Designed the entire microservices architecture",
        context: "architecture",
      },
    ],
    metadata: { extractedAt: "2025-01-01", modelVersion: "v1", processingTimeMs: 5000 },
  },
};

const PROCESS_IDS: ProcessDocumentIds = {
  resumeId: "resume-1",
  referenceLetterID: "letter-1",
};

describe("ExtractionReview", () => {
  const defaultProps = {
    userId: "user-1",
    fileId: "file-1",
    results: { resume: RESUME_DATA, referenceLetter: LETTER_DATA } as ExtractionResults,
    processDocumentIds: PROCESS_IDS,
    onImportComplete: vi.fn(),
    onBack: vi.fn(),
  };

  let fetchSpy: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    fetchSpy = vi.fn().mockResolvedValue({
      json: () =>
        Promise.resolve({
          data: {
            importDocumentResults: {
              __typename: "ImportDocumentResultsResult",
              profile: { id: "profile-1" },
              importedCount: { experiences: 2, educations: 1, skills: 3, testimonials: 2 },
            },
          },
        }),
    });
    vi.stubGlobal("fetch", fetchSpy);
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe("Career info section", () => {
    it("shows extracted name", () => {
      render(<ExtractionReview {...defaultProps} />);
      expect(screen.getByText("John Doe")).toBeInTheDocument();
    });

    it("shows extracted email", () => {
      render(<ExtractionReview {...defaultProps} />);
      expect(screen.getByText("john@example.com")).toBeInTheDocument();
    });

    it("shows extracted location", () => {
      render(<ExtractionReview {...defaultProps} />);
      expect(screen.getByText("Berlin, Germany")).toBeInTheDocument();
    });

    it("shows profile summary", () => {
      render(<ExtractionReview {...defaultProps} />);
      expect(
        screen.getByText("Experienced software engineer with 10 years of experience.")
      ).toBeInTheDocument();
    });

    it("does not show career info section when resume is null", () => {
      render(
        <ExtractionReview
          {...defaultProps}
          results={{ resume: null, referenceLetter: LETTER_DATA }}
        />
      );
      expect(screen.queryByText("John Doe")).not.toBeInTheDocument();
    });
  });

  describe("Experiences section", () => {
    it("renders each experience with a checkbox", () => {
      render(<ExtractionReview {...defaultProps} />);
      expect(screen.getByText("Senior Engineer")).toBeInTheDocument();
      expect(screen.getByText("Software Engineer")).toBeInTheDocument();

      const group = screen.getByRole("group", { name: "Work experiences" });
      const checkboxes = within(group).getAllByRole("checkbox");
      // 2 card checkboxes + 2 inner Checkbox elements (aria-hidden)
      expect(checkboxes.length).toBeGreaterThanOrEqual(2);
    });

    it("shows experiences pre-selected", () => {
      render(<ExtractionReview {...defaultProps} />);
      const group = screen.getByRole("group", { name: "Work experiences" });
      const cards = within(group).getAllByRole("checkbox", { checked: true });
      expect(cards.length).toBeGreaterThanOrEqual(2);
    });

    it("toggles experience selection on click", () => {
      render(<ExtractionReview {...defaultProps} />);
      const group = screen.getByRole("group", { name: "Work experiences" });
      const firstCard = within(group).getAllByRole("checkbox")[0];

      fireEvent.click(firstCard);
      expect(firstCard).toHaveAttribute("aria-checked", "false");

      fireEvent.click(firstCard);
      expect(firstCard).toHaveAttribute("aria-checked", "true");
    });
  });

  describe("Education section", () => {
    it("renders education entries with checkboxes", () => {
      render(<ExtractionReview {...defaultProps} />);
      expect(screen.getByText("MIT")).toBeInTheDocument();
      expect(screen.getByText(/Bachelor of Science/)).toBeInTheDocument();
    });
  });

  describe("Skills section", () => {
    it("renders each skill as a selectable chip", () => {
      render(<ExtractionReview {...defaultProps} />);
      const group = screen.getByRole("group", { name: "Skills" });
      expect(within(group).getByText("Go")).toBeInTheDocument();
      expect(within(group).getByText("TypeScript")).toBeInTheDocument();
      expect(within(group).getByText("PostgreSQL")).toBeInTheDocument();
    });

    it("toggles skill selection on click", () => {
      render(<ExtractionReview {...defaultProps} />);
      const group = screen.getByRole("group", { name: "Skills" });
      const goChip = within(group)
        .getAllByRole("checkbox")
        .find((el) => el.textContent?.includes("Go") && !el.textContent?.includes("TypeScript"));
      expect(goChip).toBeDefined();

      if (goChip) {
        fireEvent.click(goChip);
        expect(goChip).toHaveAttribute("aria-checked", "false");
      }
    });
  });

  describe("Testimonial section", () => {
    it("shows author name and relationship", () => {
      render(<ExtractionReview {...defaultProps} />);
      expect(screen.getByText("Jane Smith")).toBeInTheDocument();
      expect(screen.getAllByText(/Former Manager/).length).toBeGreaterThanOrEqual(1);
    });

    it("shows testimonial quotes with checkboxes", () => {
      render(<ExtractionReview {...defaultProps} />);
      expect(screen.getByText(/John is an exceptional engineer/)).toBeInTheDocument();
      expect(screen.getByText(/He excels at problem-solving/)).toBeInTheDocument();
    });

    it("does not show testimonial section when referenceLetter is null", () => {
      render(
        <ExtractionReview
          {...defaultProps}
          results={{ resume: RESUME_DATA, referenceLetter: null }}
        />
      );
      expect(screen.queryByText("Jane Smith")).not.toBeInTheDocument();
    });
  });

  describe("Discovered skills section", () => {
    it("shows discovered skills", () => {
      render(<ExtractionReview {...defaultProps} />);
      expect(screen.getByText("System Design")).toBeInTheDocument();
    });

    it("discovered skills are NOT pre-selected", () => {
      render(<ExtractionReview {...defaultProps} />);
      const group = screen.getByRole("group", { name: "Discovered skills" });
      const card = within(group).getByRole("checkbox");
      expect(card).toHaveAttribute("aria-checked", "false");
    });
  });

  describe("Selection counter", () => {
    it("shows total selected item count", () => {
      render(<ExtractionReview {...defaultProps} />);
      // 2 experiences + 1 education + 3 skills + 2 testimonials + 0 discovered = 8
      expect(screen.getByText("8 item(s) selected")).toBeInTheDocument();
    });

    it("updates count when items are deselected", () => {
      render(<ExtractionReview {...defaultProps} />);
      expect(screen.getByText("8 item(s) selected")).toBeInTheDocument();

      // Deselect first experience
      const expGroup = screen.getByRole("group", { name: "Work experiences" });
      const firstExpCard = within(expGroup).getAllByRole("checkbox")[0];
      fireEvent.click(firstExpCard);

      expect(screen.getByText("7 item(s) selected")).toBeInTheDocument();
    });
  });

  describe("Import action", () => {
    it("shows Import Selected button", () => {
      render(<ExtractionReview {...defaultProps} />);
      expect(screen.getByRole("button", { name: /Import Selected/i })).toBeInTheDocument();
    });

    it("sends selection indices in mutation", async () => {
      render(<ExtractionReview {...defaultProps} />);

      fireEvent.click(screen.getByRole("button", { name: /Import Selected/i }));

      await waitFor(() => {
        expect(fetchSpy).toHaveBeenCalledTimes(1);
      });

      const callBody = JSON.parse(fetchSpy.mock.calls[0][1].body);
      expect(callBody.variables.userId).toBe("user-1");
      expect(callBody.variables.input.resumeId).toBe("resume-1");
      expect(callBody.variables.input.referenceLetterID).toBe("letter-1");
      // Selection indices should be arrays
      expect(callBody.variables.input.selectedExperienceIndices).toEqual(
        expect.arrayContaining([0, 1])
      );
      expect(callBody.variables.input.selectedEducationIndices).toEqual(
        expect.arrayContaining([0])
      );
      expect(callBody.variables.input.selectedSkills).toEqual(
        expect.arrayContaining(["Go", "TypeScript", "PostgreSQL"])
      );
      expect(callBody.variables.input.selectedTestimonialIndices).toEqual(
        expect.arrayContaining([0, 1])
      );
    });

    it("shows loading state during import", async () => {
      fetchSpy.mockReturnValueOnce(new Promise(() => {})); // never resolves
      render(<ExtractionReview {...defaultProps} />);

      fireEvent.click(screen.getByRole("button", { name: /Import Selected/i }));

      await waitFor(() => {
        expect(screen.getByRole("button", { name: /Importing/i })).toBeDisabled();
      });
    });

    it("calls onImportComplete with profile ID on success", async () => {
      render(<ExtractionReview {...defaultProps} />);

      fireEvent.click(screen.getByRole("button", { name: /Import Selected/i }));

      await waitFor(() => {
        expect(defaultProps.onImportComplete).toHaveBeenCalledWith("profile-1");
      });
    });

    it("shows error message on import failure", async () => {
      fetchSpy.mockResolvedValueOnce({
        json: () =>
          Promise.resolve({
            data: {
              importDocumentResults: {
                __typename: "ImportDocumentResultsError",
                message: "Resume not in completed state",
                field: "resumeId",
              },
            },
          }),
      });

      render(<ExtractionReview {...defaultProps} />);
      fireEvent.click(screen.getByRole("button", { name: /Import Selected/i }));

      await waitFor(() => {
        expect(screen.getByText("Resume not in completed state")).toBeInTheDocument();
      });
    });

    it("disables import when all extractions failed", () => {
      const failedResults: ExtractionResults = {
        resume: { id: "resume-1", status: "FAILED", extractedData: null, errorMessage: "Failed" },
        referenceLetter: null,
      };

      render(<ExtractionReview {...defaultProps} results={failedResults} />);

      expect(screen.getByRole("button", { name: /Import Selected/i })).toBeDisabled();
    });

    it("disables import when nothing is selected", () => {
      // Resume-only with empty arrays
      const emptyResume: ExtractionResults = {
        resume: {
          id: "resume-1",
          status: "COMPLETED",
          extractedData: {
            name: "John",
            email: null,
            phone: null,
            location: null,
            summary: null,
            experiences: [],
            educations: [],
            skills: [],
            extractedAt: "2025-01-01",
            confidence: 0.9,
          },
          errorMessage: null,
        },
        referenceLetter: null,
      };

      render(<ExtractionReview {...defaultProps} results={emptyResume} />);
      expect(screen.getByRole("button", { name: /Import Selected/i })).toBeDisabled();
    });
  });

  describe("Feedback", () => {
    it("shows 'Something doesn't look right?' link", () => {
      render(<ExtractionReview {...defaultProps} />);
      expect(
        screen.getByRole("button", { name: /Something doesn't look right/i })
      ).toBeInTheDocument();
    });

    it("shows FeedbackForm when link is clicked", () => {
      render(<ExtractionReview {...defaultProps} />);

      fireEvent.click(screen.getByRole("button", { name: /Something doesn't look right/i }));

      expect(screen.getByTestId("feedback-form")).toBeInTheDocument();
    });

    it("does not block import while feedback form is open", () => {
      render(<ExtractionReview {...defaultProps} />);

      fireEvent.click(screen.getByRole("button", { name: /Something doesn't look right/i }));

      expect(screen.getByRole("button", { name: /Import Selected/i })).toBeEnabled();
    });
  });

  describe("Navigation", () => {
    it("shows Back button", () => {
      render(<ExtractionReview {...defaultProps} />);
      expect(screen.getByRole("button", { name: /Back/i })).toBeInTheDocument();
    });

    it("calls onBack when back button is clicked", () => {
      render(<ExtractionReview {...defaultProps} />);

      fireEvent.click(screen.getByRole("button", { name: /Back/i }));

      expect(defaultProps.onBack).toHaveBeenCalledTimes(1);
    });
  });

  describe("Partial results", () => {
    it("shows warning when one extraction failed", () => {
      const partialResults: ExtractionResults = {
        resume: RESUME_DATA,
        referenceLetter: {
          id: "letter-1",
          status: "FAILED",
          extractedData: null,
        },
      };

      render(<ExtractionReview {...defaultProps} results={partialResults} />);

      expect(screen.getByText(/extraction could not be completed/i)).toBeInTheDocument();
    });

    it("still allows import when career info succeeded but letter failed", () => {
      const partialResults: ExtractionResults = {
        resume: RESUME_DATA,
        referenceLetter: {
          id: "letter-1",
          status: "FAILED",
          extractedData: null,
        },
      };

      render(<ExtractionReview {...defaultProps} results={partialResults} />);

      expect(screen.getByRole("button", { name: /Import Selected/i })).toBeEnabled();
    });
  });
});
