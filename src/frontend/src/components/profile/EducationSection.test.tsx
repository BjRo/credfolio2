import { render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import { EducationSection } from "./EducationSection";
import type { Education } from "./types";

// Mock dialog components to avoid urql provider requirement in tests
vi.mock("./EducationFormDialog", () => ({
  EducationFormDialog: () => null,
}));
vi.mock("./DeleteEducationDialog", () => ({
  DeleteEducationDialog: () => null,
}));

const mockEducation: Education[] = [
  {
    institution: "Stanford University",
    degree: "Master of Science",
    field: "Computer Science",
    startDate: "Sep 2016",
    endDate: "Jun 2018",
    gpa: "3.9",
    achievements: "Dean's List, Research Assistant",
  },
  {
    institution: "UC Berkeley",
    degree: "Bachelor of Science",
    field: "Electrical Engineering",
    startDate: "Sep 2012",
    endDate: "May 2016",
    gpa: null,
    achievements: null,
  },
];

describe("EducationSection", () => {
  it("renders the section heading", () => {
    render(<EducationSection education={mockEducation} />);
    expect(screen.getByRole("heading", { level: 2, name: "Education" })).toBeInTheDocument();
  });

  it("renders institution names", () => {
    render(<EducationSection education={mockEducation} />);
    expect(screen.getByText("Stanford University")).toBeInTheDocument();
    expect(screen.getByText("UC Berkeley")).toBeInTheDocument();
  });

  it("renders degree and field", () => {
    render(<EducationSection education={mockEducation} />);
    expect(screen.getByText("Master of Science in Computer Science")).toBeInTheDocument();
    expect(screen.getByText("Bachelor of Science in Electrical Engineering")).toBeInTheDocument();
  });

  it("renders GPA when provided", () => {
    render(<EducationSection education={mockEducation} />);
    expect(screen.getByText("GPA:")).toBeInTheDocument();
    expect(screen.getByText("3.9")).toBeInTheDocument();
  });

  it("renders date ranges", () => {
    render(<EducationSection education={mockEducation} />);
    expect(screen.getByText("Sep 2016 - Jun 2018")).toBeInTheDocument();
    expect(screen.getByText("Sep 2012 - May 2016")).toBeInTheDocument();
  });

  it("renders achievements when provided", () => {
    render(<EducationSection education={mockEducation} />);
    expect(screen.getByText("Dean's List, Research Assistant")).toBeInTheDocument();
  });

  it("returns null when education array is empty and not editable", () => {
    const { container } = render(<EducationSection education={[]} />);
    expect(container.firstChild).toBeNull();
  });

  it("shows add button when userId is provided", () => {
    render(<EducationSection education={[]} userId="user-123" />);
    expect(screen.getByRole("button", { name: "Add education" })).toBeInTheDocument();
  });

  it("shows empty state with add prompt when editable and no entries", () => {
    render(<EducationSection education={[]} userId="user-123" />);
    expect(screen.getByText(/No education entries yet/)).toBeInTheDocument();
  });

  it("does not show add button when userId is not provided", () => {
    render(<EducationSection education={mockEducation} />);
    expect(screen.queryByRole("button", { name: "Add education" })).not.toBeInTheDocument();
  });

  it("shows action menu for profile educations", () => {
    render(
      <EducationSection
        profileEducations={[
          {
            id: "edu-1",
            institution: "MIT",
            degree: "PhD",
            field: "Physics",
            isCurrent: false,
            displayOrder: 0,
            source: "MANUAL" as const,
            createdAt: "2026-01-01",
            updatedAt: "2026-01-01",
          },
        ]}
        userId="user-123"
      />
    );
    expect(screen.getByText("MIT")).toBeInTheDocument();
    // Action menu button should be present (at least the desktop one)
    const actionButtons = screen.getAllByRole("button", { name: "More actions" });
    expect(actionButtons.length).toBeGreaterThan(0);
  });
});
