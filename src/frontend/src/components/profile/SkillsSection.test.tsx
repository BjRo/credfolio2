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
  it("renders the section heading", () => {
    render(<SkillsSection profileSkills={mockProfileSkills} />);
    expect(screen.getByRole("heading", { level: 2, name: "Skills" })).toBeInTheDocument();
  });

  it("returns null when no skills and not editable", () => {
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

  describe("Accessibility", () => {
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
      render(<SkillsSection profileSkills={mockProfileSkills} />);
      expect(screen.queryByRole("button", { name: "Add skill" })).not.toBeInTheDocument();
    });
  });
});
