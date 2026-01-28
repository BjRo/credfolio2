import { render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import { ExperienceSource } from "@/graphql/generated/graphql";
import { EducationSection } from "./EducationSection";
import type { ProfileEducation } from "./types";

// Mock dialog components to avoid urql provider requirement in tests
vi.mock("./EducationFormDialog", () => ({
  EducationFormDialog: () => null,
}));
vi.mock("./DeleteEducationDialog", () => ({
  DeleteEducationDialog: () => null,
}));

const mockProfileEducations: ProfileEducation[] = [
  {
    id: "edu-1",
    institution: "Stanford University",
    degree: "Master of Science",
    field: "Computer Science",
    startDate: "Sep 2016",
    endDate: "Jun 2018",
    isCurrent: false,
    gpa: "3.9",
    description: "Dean's List, Research Assistant",
    displayOrder: 0,
    source: ExperienceSource.Manual,
    createdAt: "2026-01-01T00:00:00Z",
    updatedAt: "2026-01-01T00:00:00Z",
  },
  {
    id: "edu-2",
    institution: "UC Berkeley",
    degree: "Bachelor of Science",
    field: "Electrical Engineering",
    startDate: "Sep 2012",
    endDate: "May 2016",
    isCurrent: false,
    gpa: null,
    description: null,
    displayOrder: 1,
    source: ExperienceSource.Manual,
    createdAt: "2026-01-01T00:00:00Z",
    updatedAt: "2026-01-01T00:00:00Z",
  },
];

describe("EducationSection", () => {
  it("renders the section heading", () => {
    render(<EducationSection profileEducations={mockProfileEducations} />);
    expect(screen.getByRole("heading", { level: 2, name: "Education" })).toBeInTheDocument();
  });

  it("renders institution names", () => {
    render(<EducationSection profileEducations={mockProfileEducations} />);
    expect(screen.getByText("Stanford University")).toBeInTheDocument();
    expect(screen.getByText("UC Berkeley")).toBeInTheDocument();
  });

  it("renders degree and field", () => {
    render(<EducationSection profileEducations={mockProfileEducations} />);
    expect(screen.getByText("Master of Science in Computer Science")).toBeInTheDocument();
    expect(screen.getByText("Bachelor of Science in Electrical Engineering")).toBeInTheDocument();
  });

  it("renders GPA when provided", () => {
    render(<EducationSection profileEducations={mockProfileEducations} />);
    expect(screen.getByText("GPA:")).toBeInTheDocument();
    expect(screen.getByText("3.9")).toBeInTheDocument();
  });

  it("renders date ranges", () => {
    render(<EducationSection profileEducations={mockProfileEducations} />);
    expect(screen.getByText("Sep 2016 - Jun 2018")).toBeInTheDocument();
    expect(screen.getByText("Sep 2012 - May 2016")).toBeInTheDocument();
  });

  it("renders description when provided", () => {
    render(<EducationSection profileEducations={mockProfileEducations} />);
    expect(screen.getByText("Dean's List, Research Assistant")).toBeInTheDocument();
  });

  it("returns null when profileEducations is empty and not editable", () => {
    const { container } = render(<EducationSection profileEducations={[]} />);
    expect(container.firstChild).toBeNull();
  });

  it("shows add button when userId is provided", () => {
    render(<EducationSection profileEducations={[]} userId="user-123" />);
    expect(screen.getByRole("button", { name: "Add education" })).toBeInTheDocument();
  });

  it("shows empty state with add prompt when editable and no entries", () => {
    render(<EducationSection profileEducations={[]} userId="user-123" />);
    expect(screen.getByText(/No education entries yet/)).toBeInTheDocument();
  });

  it("does not show add button when userId is not provided", () => {
    render(<EducationSection profileEducations={mockProfileEducations} />);
    expect(screen.queryByRole("button", { name: "Add education" })).not.toBeInTheDocument();
  });

  it("shows action menu for profile educations when editable", () => {
    render(<EducationSection profileEducations={mockProfileEducations} userId="user-123" />);
    expect(screen.getByText("Stanford University")).toBeInTheDocument();
    // Action menu button should be present (at least the desktop one)
    const actionButtons = screen.getAllByRole("button", { name: "More actions" });
    expect(actionButtons.length).toBeGreaterThan(0);
  });

  describe("Hover-only kebab menu", () => {
    it("hides action menu by default until hover/focus", () => {
      render(<EducationSection profileEducations={mockProfileEducations} userId="user-123" />);

      // Action buttons should exist in the DOM for accessibility
      const actionButtons = screen.getAllByRole("button", { name: "More actions" });
      expect(actionButtons.length).toBeGreaterThan(0);

      // Check the action menu container has hidden-until-hover classes
      const actionButton = actionButtons[0];
      const menuContainer = actionButton.parentElement;
      expect(menuContainer).toHaveClass("opacity-0");
      expect(menuContainer).toHaveClass("group-hover/card:opacity-100");
      expect(menuContainer).toHaveClass("group-focus-within/card:opacity-100");
    });

    it("makes menu visible on focus for keyboard accessibility", () => {
      render(<EducationSection profileEducations={mockProfileEducations} userId="user-123" />);

      const actionButtons = screen.getAllByRole("button", { name: "More actions" });
      const actionButton = actionButtons[0];
      const menuContainer = actionButton.parentElement;

      // Container should have focus-within classes for keyboard accessibility
      expect(menuContainer).toHaveClass("focus-within:opacity-100");
    });
  });
});
