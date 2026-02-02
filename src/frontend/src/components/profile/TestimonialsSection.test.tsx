import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";
import { TestimonialRelationship } from "@/graphql/generated/graphql";
import { TestimonialsSection } from "./TestimonialsSection";

const mockTestimonials = [
  {
    __typename: "Testimonial" as const,
    id: "1",
    quote: "Great team player with excellent leadership skills.",
    authorName: "John Manager",
    authorTitle: "Engineering Manager",
    authorCompany: "Acme Corp",
    relationship: TestimonialRelationship.Manager,
    createdAt: "2024-01-01T00:00:00Z",
    validatedSkills: [],
    referenceLetter: null,
  },
  {
    __typename: "Testimonial" as const,
    id: "2",
    quote: "A brilliant collaborator who consistently delivers high-quality work.",
    authorName: "Sarah Peer",
    authorTitle: "Senior Engineer",
    authorCompany: "Acme Corp",
    relationship: TestimonialRelationship.Peer,
    createdAt: "2024-01-02T00:00:00Z",
    validatedSkills: [],
    referenceLetter: null,
  },
];

const mockTestimonialsWithSkills = [
  {
    __typename: "Testimonial" as const,
    id: "1",
    quote: "Great team player with excellent leadership skills.",
    authorName: "John Manager",
    authorTitle: "Engineering Manager",
    authorCompany: "Acme Corp",
    relationship: TestimonialRelationship.Manager,
    createdAt: "2024-01-01T00:00:00Z",
    validatedSkills: [
      { __typename: "ProfileSkill" as const, id: "skill-1", name: "Leadership" },
      { __typename: "ProfileSkill" as const, id: "skill-2", name: "Team Management" },
    ],
    referenceLetter: null,
  },
];

const mockTestimonialsWithSourceBadge = [
  {
    __typename: "Testimonial" as const,
    id: "1",
    quote: "Great team player with excellent leadership skills.",
    authorName: "John Manager",
    authorTitle: "Engineering Manager",
    authorCompany: "Acme Corp",
    relationship: TestimonialRelationship.Manager,
    createdAt: "2024-01-01T00:00:00Z",
    validatedSkills: [],
    referenceLetter: {
      __typename: "ReferenceLetter" as const,
      id: "ref-1",
      file: {
        __typename: "File" as const,
        id: "file-1",
        url: "https://example.com/reference-letter.pdf",
      },
    },
  },
  {
    __typename: "Testimonial" as const,
    id: "2",
    quote: "A brilliant collaborator who consistently delivers high-quality work.",
    authorName: "Sarah Peer",
    authorTitle: "Senior Engineer",
    authorCompany: "Acme Corp",
    relationship: TestimonialRelationship.Peer,
    createdAt: "2024-01-02T00:00:00Z",
    validatedSkills: [],
    referenceLetter: {
      __typename: "ReferenceLetter" as const,
      id: "ref-2",
      file: null, // No file attached
    },
  },
];

describe("TestimonialsSection", () => {
  it("renders the section heading", () => {
    render(<TestimonialsSection testimonials={mockTestimonials} />);
    expect(screen.getByRole("heading", { level: 2, name: "What Others Say" })).toBeInTheDocument();
  });

  it("returns null when no testimonials and no onAddReference", () => {
    const { container } = render(<TestimonialsSection testimonials={[]} />);
    expect(container.firstChild).toBeNull();
  });

  it("renders testimonial quotes", () => {
    render(<TestimonialsSection testimonials={mockTestimonials} />);
    expect(
      screen.getByText("Great team player with excellent leadership skills.")
    ).toBeInTheDocument();
    expect(
      screen.getByText("A brilliant collaborator who consistently delivers high-quality work.")
    ).toBeInTheDocument();
  });

  it("renders author names", () => {
    render(<TestimonialsSection testimonials={mockTestimonials} />);
    expect(screen.getByText("John Manager")).toBeInTheDocument();
    expect(screen.getByText("Sarah Peer")).toBeInTheDocument();
  });

  it("renders author title and company", () => {
    render(<TestimonialsSection testimonials={mockTestimonials} />);
    expect(screen.getByText("Engineering Manager at Acme Corp")).toBeInTheDocument();
    expect(screen.getByText("Senior Engineer at Acme Corp")).toBeInTheDocument();
  });

  it("renders relationship badges", () => {
    render(<TestimonialsSection testimonials={mockTestimonials} />);
    expect(screen.getByText("Manager")).toBeInTheDocument();
    expect(screen.getByText("Peer")).toBeInTheDocument();
  });

  describe("Empty state", () => {
    it("shows empty state when onAddReference is provided", () => {
      render(<TestimonialsSection testimonials={[]} onAddReference={() => {}} />);
      expect(screen.getByText("No testimonials yet.")).toBeInTheDocument();
      expect(screen.getByText(/Add a reference letter/)).toBeInTheDocument();
    });

    it("shows add button in empty state", () => {
      render(<TestimonialsSection testimonials={[]} onAddReference={() => {}} />);
      expect(screen.getByRole("button", { name: "Add Reference Letter" })).toBeInTheDocument();
    });

    it("calls onAddReference when add button is clicked", async () => {
      const user = userEvent.setup();
      const onAddReference = vi.fn();
      render(<TestimonialsSection testimonials={[]} onAddReference={onAddReference} />);

      await user.click(screen.getByRole("button", { name: "Add Reference Letter" }));
      expect(onAddReference).toHaveBeenCalled();
    });
  });

  describe("Loading state", () => {
    it("shows loading skeleton when isLoading is true", () => {
      render(<TestimonialsSection testimonials={[]} isLoading={true} onAddReference={() => {}} />);
      // Loading state shows animated skeleton
      const skeletons = document.querySelectorAll(".animate-pulse");
      expect(skeletons.length).toBeGreaterThan(0);
    });
  });

  describe("Header add button", () => {
    it("shows add button in header when onAddReference is provided", () => {
      render(<TestimonialsSection testimonials={mockTestimonials} onAddReference={() => {}} />);
      expect(screen.getByRole("button", { name: "Add reference letter" })).toBeInTheDocument();
    });

    it("calls onAddReference when header add button is clicked", async () => {
      const user = userEvent.setup();
      const onAddReference = vi.fn();
      render(
        <TestimonialsSection testimonials={mockTestimonials} onAddReference={onAddReference} />
      );

      await user.click(screen.getByRole("button", { name: "Add reference letter" }));
      expect(onAddReference).toHaveBeenCalled();
    });

    it("hides add button in header when onAddReference is not provided", () => {
      render(<TestimonialsSection testimonials={mockTestimonials} />);
      expect(
        screen.queryByRole("button", { name: "Add reference letter" })
      ).not.toBeInTheDocument();
    });
  });

  describe("Validated skills", () => {
    it("displays validated skills when present", () => {
      render(<TestimonialsSection testimonials={mockTestimonialsWithSkills} />);
      expect(screen.getByText("Validates:")).toBeInTheDocument();
      expect(screen.getByRole("button", { name: "Leadership" })).toBeInTheDocument();
      expect(screen.getByRole("button", { name: "Team Management" })).toBeInTheDocument();
    });

    it("does not display validated skills section when empty", () => {
      render(<TestimonialsSection testimonials={mockTestimonials} />);
      expect(screen.queryByText("Validates:")).not.toBeInTheDocument();
    });

    it("calls onSkillClick when a skill is clicked", async () => {
      const user = userEvent.setup();
      const onSkillClick = vi.fn();
      render(
        <TestimonialsSection
          testimonials={mockTestimonialsWithSkills}
          onSkillClick={onSkillClick}
        />
      );

      await user.click(screen.getByRole("button", { name: "Leadership" }));
      expect(onSkillClick).toHaveBeenCalledWith("skill-1");
    });

    it("renders skill separators between multiple skills", () => {
      render(<TestimonialsSection testimonials={mockTestimonialsWithSkills} />);
      // Check that there's a separator (·) between skills
      const separator = screen.getByText("·");
      expect(separator).toBeInTheDocument();
    });
  });

  describe("Source badge", () => {
    it("displays source badge when testimonial has a reference letter with file", () => {
      render(<TestimonialsSection testimonials={mockTestimonialsWithSourceBadge} />);
      // Should show exactly one source badge (the one with file)
      const sourceLinks = screen.getAllByRole("link", { name: /view source/i });
      expect(sourceLinks).toHaveLength(1);
    });

    it("source badge links to the PDF file URL", () => {
      render(<TestimonialsSection testimonials={mockTestimonialsWithSourceBadge} />);
      const sourceLink = screen.getByRole("link", { name: /view source/i });
      expect(sourceLink).toHaveAttribute("href", "https://example.com/reference-letter.pdf");
    });

    it("source badge opens in a new tab", () => {
      render(<TestimonialsSection testimonials={mockTestimonialsWithSourceBadge} />);
      const sourceLink = screen.getByRole("link", { name: /view source/i });
      expect(sourceLink).toHaveAttribute("target", "_blank");
      expect(sourceLink).toHaveAttribute("rel", "noopener noreferrer");
    });

    it("does not display source badge when reference letter has no file", () => {
      // The second testimonial in mockTestimonialsWithSourceBadge has no file
      render(<TestimonialsSection testimonials={mockTestimonialsWithSourceBadge} />);
      // Only one source link should exist (for the first testimonial)
      const sourceLinks = screen.getAllByRole("link", { name: /view source/i });
      expect(sourceLinks).toHaveLength(1);
    });

    it("does not display source badge when there is no reference letter", () => {
      render(<TestimonialsSection testimonials={mockTestimonials} />);
      expect(screen.queryByRole("link", { name: /view source/i })).not.toBeInTheDocument();
    });
  });
});
