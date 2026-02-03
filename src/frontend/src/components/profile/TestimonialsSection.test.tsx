import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";
import { TestimonialRelationship } from "@/graphql/generated/graphql";
import { TestimonialsSection } from "./TestimonialsSection";

// Mock dialog components to avoid urql provider requirement in tests
vi.mock("./AuthorEditModal", () => ({
  AuthorEditModal: () => null,
}));
vi.mock("./DeleteTestimonialDialog", () => ({
  DeleteTestimonialDialog: () => null,
}));

// Mock next/image
vi.mock("next/image", () => ({
  default: ({ src, alt, ...props }: { src: string; alt: string }) => (
    // biome-ignore lint/performance/noImgElement: this is a mock for testing
    <img src={src} alt={alt} {...props} />
  ),
}));

const mockTestimonials = [
  {
    __typename: "Testimonial" as const,
    id: "1",
    quote: "Great team player with excellent leadership skills.",
    author: {
      __typename: "Author" as const,
      id: "author-1",
      name: "John Manager",
      title: "Engineering Manager",
      company: "Acme Corp",
      linkedInUrl: null,
      imageUrl: null,
    },
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
    author: {
      __typename: "Author" as const,
      id: "author-2",
      name: "Sarah Peer",
      title: "Senior Engineer",
      company: "Acme Corp",
      linkedInUrl: null,
      imageUrl: null,
    },
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
    author: {
      __typename: "Author" as const,
      id: "author-1",
      name: "John Manager",
      title: "Engineering Manager",
      company: "Acme Corp",
      linkedInUrl: null,
      imageUrl: null,
    },
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
    author: {
      __typename: "Author" as const,
      id: "author-1",
      name: "John Manager",
      title: "Engineering Manager",
      company: "Acme Corp",
      linkedInUrl: null,
      imageUrl: null,
    },
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
    author: {
      __typename: "Author" as const,
      id: "author-2",
      name: "Sarah Peer",
      title: "Senior Engineer",
      company: "Acme Corp",
      linkedInUrl: null,
      imageUrl: null,
    },
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

  describe("Kebab menu", () => {
    it("displays kebab menu for each quote that has a source", () => {
      render(<TestimonialsSection testimonials={mockTestimonialsWithSourceBadge} />);
      // mockTestimonialsWithSourceBadge has 2 testimonials from 2 different authors,
      // but only one (John Manager) has a file, so only one menu
      const menuButtons = screen.getAllByRole("button", { name: /more actions/i });
      expect(menuButtons).toHaveLength(1);
    });

    it("kebab menu is hidden by default and visible on hover", () => {
      render(<TestimonialsSection testimonials={mockTestimonialsWithSourceBadge} />);
      const menuButton = screen.getByRole("button", { name: /more actions/i });
      // Check that the menu button's container has opacity-0 class (hidden by default)
      const menuContainer = menuButton.parentElement;
      expect(menuContainer).toHaveClass("opacity-0");
      // Check it has hover visibility classes (now on quote level)
      expect(menuContainer).toHaveClass("group-hover/quote:opacity-100");
    });

    it("kebab menu contains 'View source document' option", async () => {
      const user = userEvent.setup();
      render(<TestimonialsSection testimonials={mockTestimonialsWithSourceBadge} />);

      const menuButton = screen.getByRole("button", { name: /more actions/i });
      await user.click(menuButton);

      expect(screen.getByRole("menuitem", { name: /view source document/i })).toBeInTheDocument();
    });

    it("'View source document' links to the PDF file URL", async () => {
      const user = userEvent.setup();
      render(<TestimonialsSection testimonials={mockTestimonialsWithSourceBadge} />);

      const menuButton = screen.getByRole("button", { name: /more actions/i });
      await user.click(menuButton);

      const viewSourceItem = screen.getByRole("menuitem", { name: /view source document/i });
      // The menuitem should be a link
      expect(viewSourceItem.closest("a")).toHaveAttribute(
        "href",
        "https://example.com/reference-letter.pdf"
      );
    });

    it("'View source document' opens in a new tab", async () => {
      const user = userEvent.setup();
      render(<TestimonialsSection testimonials={mockTestimonialsWithSourceBadge} />);

      const menuButton = screen.getByRole("button", { name: /more actions/i });
      await user.click(menuButton);

      const viewSourceItem = screen.getByRole("menuitem", { name: /view source document/i });
      const link = viewSourceItem.closest("a");
      expect(link).toHaveAttribute("target", "_blank");
      expect(link).toHaveAttribute("rel", "noopener noreferrer");
    });

    it("does not display kebab menu when no testimonials have source files", () => {
      render(<TestimonialsSection testimonials={mockTestimonials} />);
      expect(screen.queryByRole("button", { name: /more actions/i })).not.toBeInTheDocument();
    });
  });

  describe("Grouping by author", () => {
    const mockTestimonialsFromSameAuthor = [
      {
        __typename: "Testimonial" as const,
        id: "1",
        quote: "First quote from John about leadership.",
        author: {
          __typename: "Author" as const,
          id: "author-1",
          name: "John Manager",
          title: "Engineering Manager",
          company: "Acme Corp",
          linkedInUrl: "https://linkedin.com/in/johnmanager",
          imageUrl: null,
        },
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
            url: "https://example.com/letter1.pdf",
          },
        },
      },
      {
        __typename: "Testimonial" as const,
        id: "2",
        quote: "Second quote from John about teamwork.",
        author: {
          __typename: "Author" as const,
          id: "author-1",
          name: "John Manager",
          title: "Engineering Manager",
          company: "Acme Corp",
          linkedInUrl: "https://linkedin.com/in/johnmanager",
          imageUrl: null,
        },
        authorName: "John Manager",
        authorTitle: "Engineering Manager",
        authorCompany: "Acme Corp",
        relationship: TestimonialRelationship.Manager,
        createdAt: "2024-01-02T00:00:00Z",
        validatedSkills: [],
        referenceLetter: {
          __typename: "ReferenceLetter" as const,
          id: "ref-2",
          file: {
            __typename: "File" as const,
            id: "file-2",
            url: "https://example.com/letter2.pdf",
          },
        },
      },
      {
        __typename: "Testimonial" as const,
        id: "3",
        quote: "Quote from a different author.",
        author: {
          __typename: "Author" as const,
          id: "author-2",
          name: "Sarah Peer",
          title: "Senior Engineer",
          company: "Tech Inc",
          linkedInUrl: null,
          imageUrl: null,
        },
        authorName: "Sarah Peer",
        authorTitle: "Senior Engineer",
        authorCompany: "Tech Inc",
        relationship: TestimonialRelationship.Peer,
        createdAt: "2024-01-03T00:00:00Z",
        validatedSkills: [],
        referenceLetter: null,
      },
    ];

    it("groups testimonials from the same author together", () => {
      render(<TestimonialsSection testimonials={mockTestimonialsFromSameAuthor} />);
      // Should show author name only once per group
      const johnNames = screen.getAllByText("John Manager");
      expect(johnNames).toHaveLength(1);
      // But both quotes from John should be visible
      expect(screen.getByText("First quote from John about leadership.")).toBeInTheDocument();
      expect(screen.getByText("Second quote from John about teamwork.")).toBeInTheDocument();
    });

    it("shows author info once per group", () => {
      render(<TestimonialsSection testimonials={mockTestimonialsFromSameAuthor} />);
      // John's title/company should appear once
      const johnAttribution = screen.getAllByText("Engineering Manager at Acme Corp");
      expect(johnAttribution).toHaveLength(1);
      // Sarah's title/company should appear once
      expect(screen.getByText("Senior Engineer at Tech Inc")).toBeInTheDocument();
    });

    it("shows LinkedIn link when author has linkedInUrl", () => {
      render(<TestimonialsSection testimonials={mockTestimonialsFromSameAuthor} />);
      const linkedInLink = screen.getByRole("link", { name: /linkedin/i });
      expect(linkedInLink).toHaveAttribute("href", "https://linkedin.com/in/johnmanager");
      expect(linkedInLink).toHaveAttribute("target", "_blank");
    });

    it("displays one kebab menu per quote with source documents", () => {
      render(<TestimonialsSection testimonials={mockTestimonialsFromSameAuthor} />);
      // John's 2 quotes both have reference letters with files, so 2 kebab menus
      // Sarah's quote has no reference letter, so no kebab menu
      const menuButtons = screen.getAllByRole("button", { name: /more actions/i });
      expect(menuButtons).toHaveLength(2);
    });

    it("kebab menu shows source document option for each quote", async () => {
      const user = userEvent.setup();
      render(<TestimonialsSection testimonials={mockTestimonialsFromSameAuthor} />);

      const menuButtons = screen.getAllByRole("button", { name: /more actions/i });
      await user.click(menuButtons[0]);

      // Should show "View source document" option
      const viewSourceItem = screen.getByRole("menuitem", { name: /view source document/i });
      expect(viewSourceItem).toBeInTheDocument();
    });

    it("shows relationship badge in author group header", () => {
      render(<TestimonialsSection testimonials={mockTestimonialsFromSameAuthor} />);
      // Manager badge for John
      expect(screen.getByText("Manager")).toBeInTheDocument();
      // Peer badge for Sarah
      expect(screen.getByText("Peer")).toBeInTheDocument();
    });
  });

  describe("Expand/collapse functionality", () => {
    const mockManyQuotesFromSameAuthor = [
      {
        __typename: "Testimonial" as const,
        id: "1",
        quote: "First visible quote.",
        author: {
          __typename: "Author" as const,
          id: "author-1",
          name: "Prolific Author",
          title: "CEO",
          company: "Big Corp",
          linkedInUrl: null,
          imageUrl: null,
        },
        authorName: "Prolific Author",
        authorTitle: "CEO",
        authorCompany: "Big Corp",
        relationship: TestimonialRelationship.Manager,
        createdAt: "2024-01-01T00:00:00Z",
        validatedSkills: [],
        referenceLetter: null,
      },
      {
        __typename: "Testimonial" as const,
        id: "2",
        quote: "Second visible quote.",
        author: {
          __typename: "Author" as const,
          id: "author-1",
          name: "Prolific Author",
          title: "CEO",
          company: "Big Corp",
          linkedInUrl: null,
          imageUrl: null,
        },
        authorName: "Prolific Author",
        authorTitle: "CEO",
        authorCompany: "Big Corp",
        relationship: TestimonialRelationship.Manager,
        createdAt: "2024-01-02T00:00:00Z",
        validatedSkills: [],
        referenceLetter: null,
      },
      {
        __typename: "Testimonial" as const,
        id: "3",
        quote: "Third collapsed quote.",
        author: {
          __typename: "Author" as const,
          id: "author-1",
          name: "Prolific Author",
          title: "CEO",
          company: "Big Corp",
          linkedInUrl: null,
          imageUrl: null,
        },
        authorName: "Prolific Author",
        authorTitle: "CEO",
        authorCompany: "Big Corp",
        relationship: TestimonialRelationship.Manager,
        createdAt: "2024-01-03T00:00:00Z",
        validatedSkills: [],
        referenceLetter: null,
      },
      {
        __typename: "Testimonial" as const,
        id: "4",
        quote: "Fourth collapsed quote.",
        author: {
          __typename: "Author" as const,
          id: "author-1",
          name: "Prolific Author",
          title: "CEO",
          company: "Big Corp",
          linkedInUrl: null,
          imageUrl: null,
        },
        authorName: "Prolific Author",
        authorTitle: "CEO",
        authorCompany: "Big Corp",
        relationship: TestimonialRelationship.Manager,
        createdAt: "2024-01-04T00:00:00Z",
        validatedSkills: [],
        referenceLetter: null,
      },
    ];

    it("shows first two quotes expanded by default when author has 3+ quotes", () => {
      render(<TestimonialsSection testimonials={mockManyQuotesFromSameAuthor} />);
      // First two quotes should be visible
      expect(screen.getByText("First visible quote.")).toBeInTheDocument();
      expect(screen.getByText("Second visible quote.")).toBeInTheDocument();
      // Third and fourth quotes should be hidden initially
      expect(screen.queryByText("Third collapsed quote.")).not.toBeInTheDocument();
      expect(screen.queryByText("Fourth collapsed quote.")).not.toBeInTheDocument();
    });

    it("shows 'Show X more' button when quotes are collapsed", () => {
      render(<TestimonialsSection testimonials={mockManyQuotesFromSameAuthor} />);
      expect(screen.getByRole("button", { name: /show 2 more/i })).toBeInTheDocument();
    });

    it("expands all quotes when 'Show more' is clicked", async () => {
      const user = userEvent.setup();
      render(<TestimonialsSection testimonials={mockManyQuotesFromSameAuthor} />);

      await user.click(screen.getByRole("button", { name: /show 2 more/i }));

      // All quotes should now be visible
      expect(screen.getByText("First visible quote.")).toBeInTheDocument();
      expect(screen.getByText("Second visible quote.")).toBeInTheDocument();
      expect(screen.getByText("Third collapsed quote.")).toBeInTheDocument();
      expect(screen.getByText("Fourth collapsed quote.")).toBeInTheDocument();
    });

    it("shows 'Show less' button after expanding", async () => {
      const user = userEvent.setup();
      render(<TestimonialsSection testimonials={mockManyQuotesFromSameAuthor} />);

      await user.click(screen.getByRole("button", { name: /show 2 more/i }));

      expect(screen.getByRole("button", { name: /show less/i })).toBeInTheDocument();
    });

    it("collapses quotes when 'Show less' is clicked", async () => {
      const user = userEvent.setup();
      render(<TestimonialsSection testimonials={mockManyQuotesFromSameAuthor} />);

      // Expand
      await user.click(screen.getByRole("button", { name: /show 2 more/i }));
      // Collapse
      await user.click(screen.getByRole("button", { name: /show less/i }));

      // Back to initial state
      expect(screen.getByText("First visible quote.")).toBeInTheDocument();
      expect(screen.getByText("Second visible quote.")).toBeInTheDocument();
      expect(screen.queryByText("Third collapsed quote.")).not.toBeInTheDocument();
      expect(screen.queryByText("Fourth collapsed quote.")).not.toBeInTheDocument();
    });

    it("shows all quotes expanded when author has only 1-2 quotes", () => {
      const twoQuotes = mockManyQuotesFromSameAuthor.slice(0, 2);
      render(<TestimonialsSection testimonials={twoQuotes} />);

      expect(screen.getByText("First visible quote.")).toBeInTheDocument();
      expect(screen.getByText("Second visible quote.")).toBeInTheDocument();
      // No "Show more" button when all quotes are visible
      expect(screen.queryByRole("button", { name: /show.*more/i })).not.toBeInTheDocument();
    });
  });

  describe("Validated skills in grouped testimonials", () => {
    const mockGroupedWithSkills = [
      {
        __typename: "Testimonial" as const,
        id: "1",
        quote: "Quote about leadership skills.",
        author: {
          __typename: "Author" as const,
          id: "author-1",
          name: "John Manager",
          title: "Engineering Manager",
          company: "Acme Corp",
          linkedInUrl: null,
          imageUrl: null,
        },
        authorName: "John Manager",
        authorTitle: "Engineering Manager",
        authorCompany: "Acme Corp",
        relationship: TestimonialRelationship.Manager,
        createdAt: "2024-01-01T00:00:00Z",
        validatedSkills: [
          { __typename: "ProfileSkill" as const, id: "skill-1", name: "Leadership" },
        ],
        referenceLetter: null,
      },
      {
        __typename: "Testimonial" as const,
        id: "2",
        quote: "Quote about teamwork skills.",
        author: {
          __typename: "Author" as const,
          id: "author-1",
          name: "John Manager",
          title: "Engineering Manager",
          company: "Acme Corp",
          linkedInUrl: null,
          imageUrl: null,
        },
        authorName: "John Manager",
        authorTitle: "Engineering Manager",
        authorCompany: "Acme Corp",
        relationship: TestimonialRelationship.Manager,
        createdAt: "2024-01-02T00:00:00Z",
        validatedSkills: [{ __typename: "ProfileSkill" as const, id: "skill-2", name: "Teamwork" }],
        referenceLetter: null,
      },
    ];

    it("displays validated skills for each quote within a group", () => {
      render(<TestimonialsSection testimonials={mockGroupedWithSkills} />);
      expect(screen.getByRole("button", { name: "Leadership" })).toBeInTheDocument();
      expect(screen.getByRole("button", { name: "Teamwork" })).toBeInTheDocument();
    });

    it("calls onSkillClick when a skill in a grouped quote is clicked", async () => {
      const user = userEvent.setup();
      const onSkillClick = vi.fn();
      render(
        <TestimonialsSection testimonials={mockGroupedWithSkills} onSkillClick={onSkillClick} />
      );

      await user.click(screen.getByRole("button", { name: "Leadership" }));
      expect(onSkillClick).toHaveBeenCalledWith("skill-1");

      await user.click(screen.getByRole("button", { name: "Teamwork" }));
      expect(onSkillClick).toHaveBeenCalledWith("skill-2");
    });
  });

  describe("Fallback for testimonials without author entity", () => {
    const mockTestimonialsWithoutAuthor = [
      {
        __typename: "Testimonial" as const,
        id: "1",
        quote: "Legacy testimonial without author entity.",
        author: null,
        authorName: "Legacy Author",
        authorTitle: "Old Title",
        authorCompany: "Old Company",
        relationship: TestimonialRelationship.Manager,
        createdAt: "2024-01-01T00:00:00Z",
        validatedSkills: [],
        referenceLetter: null,
      },
    ];

    it("falls back to legacy author fields when author entity is null", () => {
      render(<TestimonialsSection testimonials={mockTestimonialsWithoutAuthor} />);
      expect(screen.getByText("Legacy Author")).toBeInTheDocument();
      expect(screen.getByText("Old Title at Old Company")).toBeInTheDocument();
    });
  });

  describe("Quote styling", () => {
    const mockTestimonialsFromSameAuthor = [
      {
        __typename: "Testimonial" as const,
        id: "1",
        quote: "First quote from John.",
        author: {
          __typename: "Author" as const,
          id: "author-1",
          name: "John Manager",
          title: "Engineering Manager",
          company: "Acme Corp",
          linkedInUrl: null,
          imageUrl: null,
        },
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
        quote: "Second quote from John.",
        author: {
          __typename: "Author" as const,
          id: "author-1",
          name: "John Manager",
          title: "Engineering Manager",
          company: "Acme Corp",
          linkedInUrl: null,
          imageUrl: null,
        },
        authorName: "John Manager",
        authorTitle: "Engineering Manager",
        authorCompany: "Acme Corp",
        relationship: TestimonialRelationship.Manager,
        createdAt: "2024-01-02T00:00:00Z",
        validatedSkills: [],
        referenceLetter: null,
      },
    ];

    it("renders continuous left border on quotes wrapper, not individual quotes", () => {
      render(<TestimonialsSection testimonials={mockTestimonialsFromSameAuthor} />);

      // Find the quotes wrapper with the continuous border
      const quotesWrapper = document.querySelector("[data-testid='quotes-wrapper']");
      expect(quotesWrapper).toBeInTheDocument();
      expect(quotesWrapper).toHaveClass("border-l-2");

      // Individual quote items should NOT have the border class
      const quoteItems = document.querySelectorAll("[data-testid='quote-item']");
      expect(quoteItems.length).toBe(2);
      quoteItems.forEach((item) => {
        expect(item).not.toHaveClass("border-l-2");
      });
    });

    it("renders quote text with minimal spacing from opening quote mark", () => {
      render(<TestimonialsSection testimonials={mockTestimonialsFromSameAuthor} />);

      // The opening quote mark should be positioned close to the text
      const quoteMarks = document.querySelectorAll("[data-testid='opening-quote-mark']");
      expect(quoteMarks.length).toBe(2);
      // Quote mark should not have excessive left padding
      quoteMarks.forEach((mark) => {
        expect(mark).not.toHaveClass("pl-4");
      });
    });
  });

  describe("Unknown author handling", () => {
    const mockUnknownAuthorTestimonial = [
      {
        __typename: "Testimonial" as const,
        id: "1",
        quote: "Quote from an unknown author.",
        author: {
          __typename: "Author" as const,
          id: "author-1",
          name: "unknown",
          title: null,
          company: null,
          linkedInUrl: null,
          imageUrl: null,
        },
        authorName: "unknown",
        authorTitle: null,
        authorCompany: null,
        relationship: TestimonialRelationship.Manager,
        createdAt: "2024-01-01T00:00:00Z",
        validatedSkills: [],
        referenceLetter: null,
      },
    ];

    const mockEmptyNameAuthorTestimonial = [
      {
        __typename: "Testimonial" as const,
        id: "1",
        quote: "Quote from author with empty name.",
        author: {
          __typename: "Author" as const,
          id: "author-1",
          name: "",
          title: null,
          company: null,
          linkedInUrl: null,
          imageUrl: null,
        },
        authorName: "",
        authorTitle: null,
        authorCompany: null,
        relationship: TestimonialRelationship.Manager,
        createdAt: "2024-01-01T00:00:00Z",
        validatedSkills: [],
        referenceLetter: null,
      },
    ];

    it("shows 'Unknown Author' when author name is 'unknown'", () => {
      render(<TestimonialsSection testimonials={mockUnknownAuthorTestimonial} />);
      expect(screen.getByText("Unknown Author")).toBeInTheDocument();
    });

    it("shows 'Unknown Author' when author name is empty", () => {
      render(<TestimonialsSection testimonials={mockEmptyNameAuthorTestimonial} />);
      expect(screen.getByText("Unknown Author")).toBeInTheDocument();
    });

    it("shows 'Author not detected' banner for unknown authors", () => {
      render(<TestimonialsSection testimonials={mockUnknownAuthorTestimonial} />);
      expect(screen.getByText("Author not detected")).toBeInTheDocument();
    });

    it("applies dashed border styling to unknown author cards", () => {
      render(<TestimonialsSection testimonials={mockUnknownAuthorTestimonial} />);
      const card = screen.getByText("Unknown Author").closest("[class*='rounded-lg']");
      expect(card).toHaveClass("border-dashed");
    });

    it("shows 'Add details' button when editable and author is unknown", () => {
      render(
        <TestimonialsSection
          testimonials={mockUnknownAuthorTestimonial}
          onAuthorUpdated={() => {}}
        />
      );
      expect(screen.getByRole("button", { name: "Add details" })).toBeInTheDocument();
    });

    it("does not show 'Add details' button when not editable", () => {
      render(<TestimonialsSection testimonials={mockUnknownAuthorTestimonial} />);
      expect(screen.queryByRole("button", { name: "Add details" })).not.toBeInTheDocument();
    });

    it("shows question mark avatar for unknown authors", () => {
      render(<TestimonialsSection testimonials={mockUnknownAuthorTestimonial} />);
      expect(screen.getByText("?")).toBeInTheDocument();
    });
  });

  describe("Edit author functionality", () => {
    const mockTestimonialWithAuthor = [
      {
        __typename: "Testimonial" as const,
        id: "1",
        quote: "Great leadership skills.",
        author: {
          __typename: "Author" as const,
          id: "author-1",
          name: "John Manager",
          title: "Engineering Manager",
          company: "Acme Corp",
          linkedInUrl: null,
          imageUrl: null,
        },
        authorName: "John Manager",
        authorTitle: "Engineering Manager",
        authorCompany: "Acme Corp",
        relationship: TestimonialRelationship.Manager,
        createdAt: "2024-01-01T00:00:00Z",
        validatedSkills: [],
        referenceLetter: null,
      },
    ];

    it("shows edit author menu button when editable", () => {
      render(
        <TestimonialsSection testimonials={mockTestimonialWithAuthor} onAuthorUpdated={() => {}} />
      );
      // The edit author button has aria-label "Edit author"
      expect(screen.getByRole("button", { name: "Edit author" })).toBeInTheDocument();
    });

    it("hides edit author menu button when not editable", () => {
      render(<TestimonialsSection testimonials={mockTestimonialWithAuthor} />);
      expect(screen.queryByRole("button", { name: "Edit author" })).not.toBeInTheDocument();
    });

    it("edit author menu is hidden by default and visible on hover", () => {
      render(
        <TestimonialsSection testimonials={mockTestimonialWithAuthor} onAuthorUpdated={() => {}} />
      );
      const menuButton = screen.getByRole("button", { name: "Edit author" });
      const menuContainer = menuButton.parentElement;
      expect(menuContainer).toHaveClass("opacity-0");
      expect(menuContainer).toHaveClass("group-hover/card:opacity-100");
    });

    it("edit author menu contains 'Edit author' option", async () => {
      const user = userEvent.setup();
      render(
        <TestimonialsSection testimonials={mockTestimonialWithAuthor} onAuthorUpdated={() => {}} />
      );

      const menuButton = screen.getByRole("button", { name: "Edit author" });
      await user.click(menuButton);

      expect(screen.getByRole("menuitem", { name: /edit author/i })).toBeInTheDocument();
    });

    it("does not show edit author menu for legacy testimonials without author entity", () => {
      const legacyTestimonial = [
        {
          __typename: "Testimonial" as const,
          id: "1",
          quote: "Legacy testimonial.",
          author: null,
          authorName: "Legacy Author",
          authorTitle: "Old Title",
          authorCompany: "Old Company",
          relationship: TestimonialRelationship.Manager,
          createdAt: "2024-01-01T00:00:00Z",
          validatedSkills: [],
          referenceLetter: null,
        },
      ];
      render(<TestimonialsSection testimonials={legacyTestimonial} onAuthorUpdated={() => {}} />);
      expect(screen.queryByRole("button", { name: "Edit author" })).not.toBeInTheDocument();
    });
  });

  describe("Author image display", () => {
    const mockTestimonialWithImage = [
      {
        __typename: "Testimonial" as const,
        id: "1",
        quote: "Quote from author with image.",
        author: {
          __typename: "Author" as const,
          id: "author-1",
          name: "Jane Doe",
          title: "CTO",
          company: "Tech Corp",
          linkedInUrl: null,
          imageUrl: "https://example.com/jane.jpg",
        },
        authorName: "Jane Doe",
        authorTitle: "CTO",
        authorCompany: "Tech Corp",
        relationship: TestimonialRelationship.Manager,
        createdAt: "2024-01-01T00:00:00Z",
        validatedSkills: [],
        referenceLetter: null,
      },
    ];

    const mockTestimonialWithoutImage = [
      {
        __typename: "Testimonial" as const,
        id: "1",
        quote: "Quote from author without image.",
        author: {
          __typename: "Author" as const,
          id: "author-1",
          name: "John Smith",
          title: "Manager",
          company: "Acme Corp",
          linkedInUrl: null,
          imageUrl: null,
        },
        authorName: "John Smith",
        authorTitle: "Manager",
        authorCompany: "Acme Corp",
        relationship: TestimonialRelationship.Manager,
        createdAt: "2024-01-01T00:00:00Z",
        validatedSkills: [],
        referenceLetter: null,
      },
    ];

    it("displays author image when imageUrl is provided", () => {
      render(<TestimonialsSection testimonials={mockTestimonialWithImage} />);
      const image = screen.getByAltText(/photo of jane doe/i);
      expect(image).toBeInTheDocument();
      expect(image).toHaveAttribute("src", expect.stringContaining("jane.jpg"));
    });

    it("displays initials avatar when author has no image", () => {
      render(<TestimonialsSection testimonials={mockTestimonialWithoutImage} />);
      // Should show initials "JS" for "John Smith"
      expect(screen.getByText("JS")).toBeInTheDocument();
    });

    it("does not show initials when author has image", () => {
      render(<TestimonialsSection testimonials={mockTestimonialWithImage} />);
      // Should not show initials "JD" when image is present
      expect(screen.queryByText("JD")).not.toBeInTheDocument();
    });
  });

  describe("Delete testimonial functionality", () => {
    const mockTestimonialForDelete = [
      {
        __typename: "Testimonial" as const,
        id: "1",
        quote: "Quote to delete.",
        author: {
          __typename: "Author" as const,
          id: "author-1",
          name: "John Manager",
          title: "Engineering Manager",
          company: "Acme Corp",
          linkedInUrl: null,
          imageUrl: null,
        },
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
            url: "https://example.com/letter.pdf",
          },
        },
      },
    ];

    it("shows delete option in quote menu when editable", async () => {
      const user = userEvent.setup();
      render(
        <TestimonialsSection
          testimonials={mockTestimonialForDelete}
          onTestimonialDeleted={() => {}}
        />
      );

      const menuButton = screen.getByRole("button", { name: /more actions/i });
      await user.click(menuButton);

      expect(screen.getByRole("menuitem", { name: /delete/i })).toBeInTheDocument();
    });

    it("does not show delete option when not editable", async () => {
      const user = userEvent.setup();
      render(<TestimonialsSection testimonials={mockTestimonialForDelete} />);

      const menuButton = screen.getByRole("button", { name: /more actions/i });
      await user.click(menuButton);

      // Should only have "View source document", not delete
      expect(screen.queryByRole("menuitem", { name: /delete/i })).not.toBeInTheDocument();
    });
  });
});
