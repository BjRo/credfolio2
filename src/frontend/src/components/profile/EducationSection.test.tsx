import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { EducationSection } from "./EducationSection";
import type { Education } from "./types";

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
    expect(screen.getByText("GPA: 3.9")).toBeInTheDocument();
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

  it("returns null when education array is empty", () => {
    const { container } = render(<EducationSection education={[]} />);
    expect(container.firstChild).toBeNull();
  });
});
