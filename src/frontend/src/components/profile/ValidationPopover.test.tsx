import { render, screen } from "@testing-library/react";
import type { ReactNode } from "react";
import { useQuery } from "urql";
import { beforeEach, describe, expect, it, type Mock, vi } from "vitest";
import { AuthorRelationship } from "@/graphql/generated/graphql";
import { ValidationPopover } from "./ValidationPopover";

// Mock urql's useQuery
vi.mock("urql", () => ({
  useQuery: vi.fn(),
}));

// Mock HoverCard to render content inline (avoids Portal/jsdom issues)
vi.mock("@/components/ui/hover-card", () => ({
  HoverCard: ({ children }: { children: ReactNode }) => <div>{children}</div>,
  HoverCardTrigger: ({ children }: { children: ReactNode }) => <div>{children}</div>,
  HoverCardContent: ({ children }: { children: ReactNode }) => <div>{children}</div>,
}));

const mockUseQuery = useQuery as Mock;

const mockValidationsWithFile = [
  {
    __typename: "SkillValidation" as const,
    id: "v1",
    quoteSnippet: "mentoring engineers and contributing hands-on",
    createdAt: "2024-01-01T00:00:00Z",
    referenceLetter: {
      __typename: "ReferenceLetter" as const,
      id: "letter-1",
      file: { __typename: "File" as const, id: "file-1", url: "https://example.com/file.pdf" },
      extractedData: {
        __typename: "ExtractedLetterData" as const,
        author: {
          __typename: "ExtractedAuthor" as const,
          name: "Amit Matani",
          title: "CEO",
          company: "Wellfound",
          relationship: AuthorRelationship.Manager,
        },
      },
    },
  },
];

const mockValidationsWithoutFile = [
  {
    __typename: "SkillValidation" as const,
    id: "v2",
    quoteSnippet: "excellent problem solving skills",
    createdAt: "2024-01-01T00:00:00Z",
    referenceLetter: {
      __typename: "ReferenceLetter" as const,
      id: "letter-2",
      file: null,
      extractedData: {
        __typename: "ExtractedLetterData" as const,
        author: {
          __typename: "ExtractedAuthor" as const,
          name: "Jane Doe",
          title: "CTO",
          company: "TechCo",
          relationship: AuthorRelationship.Peer,
        },
      },
    },
  },
];

const mockValidationsWithoutQuote = [
  {
    __typename: "SkillValidation" as const,
    id: "v3",
    quoteSnippet: null,
    createdAt: "2024-01-01T00:00:00Z",
    referenceLetter: {
      __typename: "ReferenceLetter" as const,
      id: "letter-3",
      file: { __typename: "File" as const, id: "file-3", url: "https://example.com/file3.pdf" },
      extractedData: {
        __typename: "ExtractedLetterData" as const,
        author: {
          __typename: "ExtractedAuthor" as const,
          name: "Bob Smith",
          title: "VP",
          company: "BigCo",
          relationship: AuthorRelationship.Manager,
        },
      },
    },
  },
];

function setupMockQuery(validations: typeof mockValidationsWithFile) {
  mockUseQuery.mockImplementation(({ pause }: { pause: boolean }) => {
    if (pause) {
      return [{ data: undefined, fetching: false, error: undefined }];
    }
    return [
      {
        data: { skillValidations: validations },
        fetching: false,
        error: undefined,
      },
    ];
  });
}

describe("ValidationPopover", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("renders children directly when validationCount is 0", () => {
    render(
      <ValidationPopover itemId="skill-1" type="skill" itemName="Mentoring" validationCount={0}>
        <span>Mentoring</span>
      </ValidationPopover>
    );

    expect(screen.getByText("Mentoring")).toBeInTheDocument();
  });

  describe("with validations loaded", () => {
    it("renders per-validation 'View in source' link when reference letter has a file", () => {
      setupMockQuery(mockValidationsWithFile);

      // Render just the content since HoverCard trigger requires user interaction
      render(
        <ValidationPopover itemId="skill-1" type="skill" itemName="Mentoring" validationCount={1}>
          <span>Mentoring</span>
        </ValidationPopover>
      );

      const link = screen.getByRole("link", { name: /view in source/i });
      expect(link).toBeInTheDocument();
      expect(link).toHaveAttribute("target", "_blank");
      expect(link).toHaveAttribute("rel", "noopener noreferrer");
      expect(link).toHaveAttribute(
        "href",
        expect.stringContaining("/viewer?letterId=letter-1&highlight=")
      );
    });

    it("does not render 'View in source' link when reference letter has no file", () => {
      setupMockQuery(mockValidationsWithoutFile);

      render(
        <ValidationPopover itemId="skill-1" type="skill" itemName="Mentoring" validationCount={1}>
          <span>Mentoring</span>
        </ValidationPopover>
      );

      expect(screen.queryByRole("link", { name: /view in source/i })).not.toBeInTheDocument();
    });

    it("links to viewer without highlight when validation has no quote snippet", () => {
      setupMockQuery(mockValidationsWithoutQuote);

      render(
        <ValidationPopover itemId="skill-1" type="skill" itemName="Mentoring" validationCount={1}>
          <span>Mentoring</span>
        </ValidationPopover>
      );

      const link = screen.getByRole("link", { name: /view in source/i });
      expect(link).toHaveAttribute("href", "/viewer?letterId=letter-3");
    });

    it("does not render 'View full testimonials' link", () => {
      setupMockQuery(mockValidationsWithFile);

      render(
        <ValidationPopover itemId="skill-1" type="skill" itemName="Mentoring" validationCount={1}>
          <span>Mentoring</span>
        </ValidationPopover>
      );

      expect(screen.queryByText(/view full testimonials/i)).not.toBeInTheDocument();
    });

    it("renders multiple validations each with their own source link", () => {
      const mixedValidations = [...mockValidationsWithFile, ...mockValidationsWithoutFile];
      setupMockQuery(mixedValidations);

      render(
        <ValidationPopover itemId="skill-1" type="skill" itemName="Mentoring" validationCount={2}>
          <span>Mentoring</span>
        </ValidationPopover>
      );

      // Only one link should appear (the one with a file)
      const links = screen.getAllByRole("link", { name: /view in source/i });
      expect(links).toHaveLength(1);
    });
  });
});
