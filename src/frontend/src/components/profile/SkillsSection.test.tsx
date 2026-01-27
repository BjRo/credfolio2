import { render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import { ExperienceSource, SkillCategory } from "@/graphql/generated/graphql";
import { SkillsSection } from "./SkillsSection";

// Mock dialog components to avoid urql provider requirement in tests
vi.mock("./SkillFormDialog", () => ({
  SkillFormDialog: () => null,
}));
vi.mock("./DeleteSkillDialog", () => ({
  DeleteSkillDialog: () => null,
}));

const mockExtractedSkills = ["JavaScript", "TypeScript", "React", "Node.js", "GraphQL"];

const mockProfileSkills = [
  {
    __typename: "ProfileSkill" as const,
    id: "1",
    name: "Go",
    normalizedName: "go",
    category: SkillCategory.Technical,
    displayOrder: 0,
    source: ExperienceSource.Manual,
    createdAt: "2024-01-01T00:00:00Z",
    updatedAt: "2024-01-01T00:00:00Z",
  },
  {
    __typename: "ProfileSkill" as const,
    id: "2",
    name: "Leadership",
    normalizedName: "leadership",
    category: SkillCategory.Soft,
    displayOrder: 1,
    source: ExperienceSource.Manual,
    createdAt: "2024-01-01T00:00:00Z",
    updatedAt: "2024-01-01T00:00:00Z",
  },
];

describe("SkillsSection", () => {
  it("renders the section heading with extracted skills", () => {
    render(<SkillsSection extractedSkills={mockExtractedSkills} />);
    expect(screen.getByRole("heading", { level: 2, name: "Skills" })).toBeInTheDocument();
  });

  it("renders extracted skills as badges", () => {
    render(<SkillsSection extractedSkills={mockExtractedSkills} />);
    expect(screen.getByText("JavaScript")).toBeInTheDocument();
    expect(screen.getByText("TypeScript")).toBeInTheDocument();
    expect(screen.getByText("React")).toBeInTheDocument();
    expect(screen.getByText("Node.js")).toBeInTheDocument();
    expect(screen.getByText("GraphQL")).toBeInTheDocument();
  });

  it("returns null when no skills at all and not editable", () => {
    const { container } = render(<SkillsSection />);
    expect(container.firstChild).toBeNull();
  });

  it("renders profile skills grouped by category", () => {
    render(<SkillsSection profileSkills={mockProfileSkills} />);
    expect(screen.getByText("Go")).toBeInTheDocument();
    expect(screen.getByText("Leadership")).toBeInTheDocument();
    expect(screen.getByText("Technical")).toBeInTheDocument();
    expect(screen.getByText("Soft Skills")).toBeInTheDocument();
  });

  it("shows extracted skills when no profile skills exist", () => {
    render(<SkillsSection extractedSkills={mockExtractedSkills} />);
    expect(screen.getByText("JavaScript")).toBeInTheDocument();
  });

  it("shows profile skills instead of extracted when both exist", () => {
    render(
      <SkillsSection profileSkills={mockProfileSkills} extractedSkills={mockExtractedSkills} />
    );
    // Profile skills are shown
    expect(screen.getByText("Go")).toBeInTheDocument();
    // Extracted skills are NOT shown when profile skills exist
    expect(screen.queryByText("JavaScript")).not.toBeInTheDocument();
  });

  describe("Accessibility", () => {
    it("renders extracted skills in an accessible list", () => {
      render(<SkillsSection extractedSkills={mockExtractedSkills} />);
      expect(screen.getByRole("list", { name: "Skills list" })).toBeInTheDocument();
    });

    it("renders each extracted skill as a list item", () => {
      render(<SkillsSection extractedSkills={mockExtractedSkills} />);
      const listItems = screen.getAllByRole("listitem");
      expect(listItems).toHaveLength(5);
    });

    it("renders profile skills in category-labeled lists", () => {
      render(<SkillsSection profileSkills={mockProfileSkills} />);
      expect(screen.getByRole("list", { name: "Technical skills" })).toBeInTheDocument();
      expect(screen.getByRole("list", { name: "Soft Skills skills" })).toBeInTheDocument();
    });
  });

  describe("Editable mode", () => {
    it("shows add button when userId is provided", () => {
      render(<SkillsSection userId="user-1" />);
      expect(screen.getByRole("button", { name: "Add skill" })).toBeInTheDocument();
    });

    it("shows empty state when editable with no skills", () => {
      render(<SkillsSection userId="user-1" />);
      expect(screen.getByText(/No skills yet/)).toBeInTheDocument();
    });

    it("does not show add button when not editable", () => {
      render(<SkillsSection extractedSkills={mockExtractedSkills} />);
      expect(screen.queryByRole("button", { name: "Add skill" })).not.toBeInTheDocument();
    });
  });
});
