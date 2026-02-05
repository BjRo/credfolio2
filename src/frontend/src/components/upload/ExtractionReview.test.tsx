import { fireEvent, render, screen, waitFor } from "@testing-library/react";
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
    extractedAt: "2025-01-01T00:00:00Z",
    confidence: 0.95,
  },
  errorMessage: null,
};

const LETTER_DATA: ExtractionResults["referenceLetter"] = {
  id: "letter-1",
  status: "COMPLETED",
  extractedData: {
    author: { name: "Jane Smith", title: "CTO", company: "TechCo", relationship: "Former Manager" },
    testimonials: [
      { quote: "John is an exceptional engineer.", skillsMentioned: ["Go", "TypeScript"] },
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
              importedCount: { experiences: 3, educations: 2, skills: 5 },
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

  describe("Testimonial section", () => {
    it("shows author name and relationship", () => {
      render(<ExtractionReview {...defaultProps} />);
      expect(screen.getByText("Jane Smith")).toBeInTheDocument();
      expect(screen.getByText(/Former Manager/)).toBeInTheDocument();
    });

    it("shows testimonial quotes", () => {
      render(<ExtractionReview {...defaultProps} />);
      expect(screen.getByText(/John is an exceptional engineer/)).toBeInTheDocument();
      expect(screen.getByText(/He excels at problem-solving/)).toBeInTheDocument();
    });

    it("shows skill mentions", () => {
      render(<ExtractionReview {...defaultProps} />);
      expect(screen.getByText(/Expert-level Go programming/)).toBeInTheDocument();
    });

    it("shows experience mentions", () => {
      render(<ExtractionReview {...defaultProps} />);
      expect(screen.getByText(/Led the backend team/)).toBeInTheDocument();
    });

    it("shows discovered skills", () => {
      render(<ExtractionReview {...defaultProps} />);
      expect(screen.getByText("System Design")).toBeInTheDocument();
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

  describe("Import action", () => {
    it("shows import button", () => {
      render(<ExtractionReview {...defaultProps} />);
      expect(screen.getByRole("button", { name: /Import to profile/i })).toBeInTheDocument();
    });

    it("calls importDocumentResults mutation when import button is clicked", async () => {
      render(<ExtractionReview {...defaultProps} />);

      fireEvent.click(screen.getByRole("button", { name: /Import to profile/i }));

      await waitFor(() => {
        expect(fetchSpy).toHaveBeenCalledTimes(1);
      });

      const callBody = JSON.parse(fetchSpy.mock.calls[0][1].body);
      expect(callBody.variables.userId).toBe("user-1");
      expect(callBody.variables.input.resumeId).toBe("resume-1");
      expect(callBody.variables.input.referenceLetterID).toBe("letter-1");
    });

    it("shows loading state during import", async () => {
      fetchSpy.mockReturnValueOnce(new Promise(() => {})); // never resolves
      render(<ExtractionReview {...defaultProps} />);

      fireEvent.click(screen.getByRole("button", { name: /Import to profile/i }));

      await waitFor(() => {
        expect(screen.getByRole("button", { name: /Importing/i })).toBeDisabled();
      });
    });

    it("calls onImportComplete with profile ID on success", async () => {
      render(<ExtractionReview {...defaultProps} />);

      fireEvent.click(screen.getByRole("button", { name: /Import to profile/i }));

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
      fireEvent.click(screen.getByRole("button", { name: /Import to profile/i }));

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

      expect(screen.getByRole("button", { name: /Import to profile/i })).toBeDisabled();
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

      expect(screen.getByRole("button", { name: /Import to profile/i })).toBeEnabled();
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

      expect(screen.getByRole("button", { name: /Import to profile/i })).toBeEnabled();
    });
  });
});
