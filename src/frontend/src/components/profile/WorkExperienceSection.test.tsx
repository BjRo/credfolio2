import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import type { WorkExperience } from "./types";
import { WorkExperienceSection } from "./WorkExperienceSection";

const mockExperience: WorkExperience[] = [
  {
    company: "TechCorp",
    title: "Senior Engineer",
    location: "San Francisco, CA",
    startDate: "Jan 2020",
    endDate: "Dec 2023",
    isCurrent: false,
    description: "Led development of core platform features.",
  },
  {
    company: "StartupCo",
    title: "Software Engineer",
    location: "New York, NY",
    startDate: "Jun 2018",
    endDate: null,
    isCurrent: true,
    description: null,
  },
];

describe("WorkExperienceSection", () => {
  it("renders the section heading", () => {
    render(<WorkExperienceSection experience={mockExperience} />);
    expect(screen.getByRole("heading", { level: 2, name: "Work Experience" })).toBeInTheDocument();
  });

  it("renders job titles", () => {
    render(<WorkExperienceSection experience={mockExperience} />);
    expect(screen.getByText("Senior Engineer")).toBeInTheDocument();
    expect(screen.getByText("Software Engineer")).toBeInTheDocument();
  });

  it("renders company names", () => {
    render(<WorkExperienceSection experience={mockExperience} />);
    expect(screen.getByText("TechCorp")).toBeInTheDocument();
    expect(screen.getByText("StartupCo")).toBeInTheDocument();
  });

  it("renders locations when provided", () => {
    render(<WorkExperienceSection experience={mockExperience} />);
    expect(screen.getByText("San Francisco, CA")).toBeInTheDocument();
    expect(screen.getByText("New York, NY")).toBeInTheDocument();
  });

  it("renders date ranges", () => {
    render(<WorkExperienceSection experience={mockExperience} />);
    expect(screen.getByText("Jan 2020 - Dec 2023")).toBeInTheDocument();
  });

  it("shows Present for current positions", () => {
    render(<WorkExperienceSection experience={mockExperience} />);
    expect(screen.getByText("Jun 2018 - Present")).toBeInTheDocument();
  });

  it("renders description when provided", () => {
    render(<WorkExperienceSection experience={mockExperience} />);
    expect(screen.getByText("Led development of core platform features.")).toBeInTheDocument();
  });

  it("returns null when experience array is empty", () => {
    const { container } = render(<WorkExperienceSection experience={[]} />);
    expect(container.firstChild).toBeNull();
  });
});
