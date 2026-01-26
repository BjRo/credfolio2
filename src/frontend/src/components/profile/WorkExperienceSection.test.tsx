import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it } from "vitest";
import type { WorkExperience } from "./types";
import { WorkExperienceSection } from "./WorkExperienceSection";

const mockExperiences: WorkExperience[] = [
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

const longDescription =
  "Led development of core platform features including authentication, authorization, and billing systems. " +
  "Managed a team of 5 engineers and coordinated with product managers to deliver features on time. " +
  "Implemented CI/CD pipelines and improved deployment frequency by 300%.";

describe("WorkExperienceSection", () => {
  it("renders the section heading", () => {
    render(<WorkExperienceSection experiences={mockExperiences} />);
    expect(screen.getByRole("heading", { level: 2, name: "Work Experience" })).toBeInTheDocument();
  });

  it("renders job titles", () => {
    render(<WorkExperienceSection experiences={mockExperiences} />);
    expect(screen.getByText("Senior Engineer")).toBeInTheDocument();
    expect(screen.getByText("Software Engineer")).toBeInTheDocument();
  });

  it("renders company names", () => {
    render(<WorkExperienceSection experiences={mockExperiences} />);
    expect(screen.getByText("TechCorp")).toBeInTheDocument();
    expect(screen.getByText("StartupCo")).toBeInTheDocument();
  });

  it("renders locations when provided", () => {
    render(<WorkExperienceSection experiences={mockExperiences} />);
    expect(screen.getByText("San Francisco, CA")).toBeInTheDocument();
    expect(screen.getByText("New York, NY")).toBeInTheDocument();
  });

  it("renders date ranges", () => {
    render(<WorkExperienceSection experiences={mockExperiences} />);
    expect(screen.getByText("Jan 2020 - Dec 2023")).toBeInTheDocument();
  });

  it("shows Present for current positions", () => {
    render(<WorkExperienceSection experiences={mockExperiences} />);
    expect(screen.getByText("Jun 2018 - Present")).toBeInTheDocument();
  });

  it("renders description when provided", () => {
    render(<WorkExperienceSection experiences={mockExperiences} />);
    expect(screen.getByText("Led development of core platform features.")).toBeInTheDocument();
  });

  it("returns null when experiences array is empty and not editable", () => {
    const { container } = render(<WorkExperienceSection experiences={[]} />);
    expect(container.firstChild).toBeNull();
  });

  describe("Current Job Highlighting", () => {
    it("shows Current badge for current positions", () => {
      render(<WorkExperienceSection experiences={mockExperiences} />);
      expect(screen.getByText("Current")).toBeInTheDocument();
    });

    it("does not show Current badge for past positions", () => {
      const pastExperiences: WorkExperience[] = [
        {
          company: "OldCo",
          title: "Developer",
          location: null,
          startDate: "2015",
          endDate: "2018",
          isCurrent: false,
          description: null,
        },
      ];
      render(<WorkExperienceSection experiences={pastExperiences} />);
      expect(screen.queryByText("Current")).not.toBeInTheDocument();
    });
  });

  describe("Expandable Description", () => {
    it("shows expand button for long descriptions", () => {
      const experiencesWithLongDesc: WorkExperience[] = [
        {
          company: "TestCo",
          title: "Engineer",
          location: null,
          startDate: "2020",
          endDate: null,
          isCurrent: true,
          description: longDescription,
        },
      ];
      render(<WorkExperienceSection experiences={experiencesWithLongDesc} />);
      expect(screen.getByRole("button", { name: /show more/i })).toBeInTheDocument();
    });

    it("does not show expand button for short descriptions", () => {
      render(<WorkExperienceSection experiences={mockExperiences} />);
      expect(screen.queryByRole("button", { name: /show more/i })).not.toBeInTheDocument();
    });

    it("toggles description expansion when clicking button", async () => {
      const user = userEvent.setup();
      const experiencesWithLongDesc: WorkExperience[] = [
        {
          company: "TestCo",
          title: "Engineer",
          location: null,
          startDate: "2020",
          endDate: null,
          isCurrent: true,
          description: longDescription,
        },
      ];
      render(<WorkExperienceSection experiences={experiencesWithLongDesc} />);

      const button = screen.getByRole("button", { name: /show more/i });
      expect(button).toHaveAttribute("aria-expanded", "false");

      await user.click(button);
      expect(screen.getByRole("button", { name: /show less/i })).toBeInTheDocument();
      expect(screen.getByRole("button", { name: /show less/i })).toHaveAttribute(
        "aria-expanded",
        "true"
      );
    });
  });
});
