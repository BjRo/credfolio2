import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";
import { ExperienceSource } from "@/graphql/generated/graphql";
import type { ProfileExperience } from "./types";
import { WorkExperienceSection } from "./WorkExperienceSection";

// Mock dialog components to avoid urql provider requirement in tests
vi.mock("./WorkExperienceFormDialog", () => ({
  WorkExperienceFormDialog: () => null,
}));
vi.mock("./DeleteExperienceDialog", () => ({
  DeleteExperienceDialog: () => null,
}));

const mockProfileExperiences: ProfileExperience[] = [
  {
    id: "exp-1",
    company: "TechCorp",
    title: "Senior Engineer",
    location: "San Francisco, CA",
    startDate: "Jan 2020",
    endDate: "Dec 2023",
    isCurrent: false,
    description: "Led development of core platform features.",
    highlights: [],
    displayOrder: 0,
    source: ExperienceSource.Manual,
    createdAt: "2026-01-01T00:00:00Z",
    updatedAt: "2026-01-01T00:00:00Z",
  },
  {
    id: "exp-2",
    company: "StartupCo",
    title: "Software Engineer",
    location: "New York, NY",
    startDate: "Jun 2018",
    endDate: null,
    isCurrent: true,
    description: null,
    highlights: [],
    displayOrder: 1,
    source: ExperienceSource.Manual,
    createdAt: "2026-01-01T00:00:00Z",
    updatedAt: "2026-01-01T00:00:00Z",
  },
];

const longDescription =
  "Led development of core platform features including authentication, authorization, and billing systems. " +
  "Managed a team of 5 engineers and coordinated with product managers to deliver features on time. " +
  "Implemented CI/CD pipelines and improved deployment frequency by 300%.";

describe("WorkExperienceSection", () => {
  it("renders the section heading", () => {
    render(<WorkExperienceSection profileExperiences={mockProfileExperiences} />);
    expect(screen.getByRole("heading", { level: 2, name: "Work Experience" })).toBeInTheDocument();
  });

  it("renders job titles", () => {
    render(<WorkExperienceSection profileExperiences={mockProfileExperiences} />);
    expect(screen.getByText("Senior Engineer")).toBeInTheDocument();
    expect(screen.getByText("Software Engineer")).toBeInTheDocument();
  });

  it("renders company names", () => {
    render(<WorkExperienceSection profileExperiences={mockProfileExperiences} />);
    expect(screen.getByText("TechCorp")).toBeInTheDocument();
    expect(screen.getByText("StartupCo")).toBeInTheDocument();
  });

  it("renders locations when provided", () => {
    render(<WorkExperienceSection profileExperiences={mockProfileExperiences} />);
    expect(screen.getByText("San Francisco, CA")).toBeInTheDocument();
    expect(screen.getByText("New York, NY")).toBeInTheDocument();
  });

  it("renders date ranges", () => {
    render(<WorkExperienceSection profileExperiences={mockProfileExperiences} />);
    // Date appears twice (mobile + desktop views)
    expect(screen.getAllByText("Jan 2020 - Dec 2023").length).toBeGreaterThan(0);
  });

  it("shows Present for current positions", () => {
    render(<WorkExperienceSection profileExperiences={mockProfileExperiences} />);
    // Date appears twice (mobile + desktop views)
    expect(screen.getAllByText("Jun 2018 - Present").length).toBeGreaterThan(0);
  });

  it("renders description when provided", () => {
    render(<WorkExperienceSection profileExperiences={mockProfileExperiences} />);
    expect(screen.getByText("Led development of core platform features.")).toBeInTheDocument();
  });

  it("returns null when profileExperiences is empty and not editable", () => {
    const { container } = render(<WorkExperienceSection profileExperiences={[]} />);
    expect(container.firstChild).toBeNull();
  });

  describe("Current Job Highlighting", () => {
    it("shows Current badge for current positions", () => {
      render(<WorkExperienceSection profileExperiences={mockProfileExperiences} />);
      expect(screen.getByText("Current")).toBeInTheDocument();
    });

    it("does not show Current badge for past positions", () => {
      const pastExperiences: ProfileExperience[] = [
        {
          id: "exp-past",
          company: "OldCo",
          title: "Developer",
          location: null,
          startDate: "2015",
          endDate: "2018",
          isCurrent: false,
          description: null,
          highlights: [],
          displayOrder: 0,
          source: ExperienceSource.Manual,
          createdAt: "2026-01-01T00:00:00Z",
          updatedAt: "2026-01-01T00:00:00Z",
        },
      ];
      render(<WorkExperienceSection profileExperiences={pastExperiences} />);
      expect(screen.queryByText("Current")).not.toBeInTheDocument();
    });
  });

  describe("Expandable Description", () => {
    it("shows expand button for long descriptions", () => {
      const experiencesWithLongDesc: ProfileExperience[] = [
        {
          id: "exp-long",
          company: "TestCo",
          title: "Engineer",
          location: null,
          startDate: "2020",
          endDate: null,
          isCurrent: true,
          description: longDescription,
          highlights: [],
          displayOrder: 0,
          source: ExperienceSource.Manual,
          createdAt: "2026-01-01T00:00:00Z",
          updatedAt: "2026-01-01T00:00:00Z",
        },
      ];
      render(<WorkExperienceSection profileExperiences={experiencesWithLongDesc} />);
      expect(screen.getByRole("button", { name: /show more/i })).toBeInTheDocument();
    });

    it("does not show expand button for short descriptions", () => {
      render(<WorkExperienceSection profileExperiences={mockProfileExperiences} />);
      expect(screen.queryByRole("button", { name: /show more/i })).not.toBeInTheDocument();
    });

    it("toggles description expansion when clicking button", async () => {
      const user = userEvent.setup();
      const experiencesWithLongDesc: ProfileExperience[] = [
        {
          id: "exp-long",
          company: "TestCo",
          title: "Engineer",
          location: null,
          startDate: "2020",
          endDate: null,
          isCurrent: true,
          description: longDescription,
          highlights: [],
          displayOrder: 0,
          source: ExperienceSource.Manual,
          createdAt: "2026-01-01T00:00:00Z",
          updatedAt: "2026-01-01T00:00:00Z",
        },
      ];
      render(<WorkExperienceSection profileExperiences={experiencesWithLongDesc} />);

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

  describe("Company grouping", () => {
    it("groups multiple roles at the same company", () => {
      const multiRoleExperiences: ProfileExperience[] = [
        {
          id: "exp-sr",
          company: "TechCorp",
          title: "Senior Engineer",
          location: "San Francisco, CA",
          startDate: "Jan 2022",
          endDate: null,
          isCurrent: true,
          description: null,
          highlights: [],
          displayOrder: 0,
          source: ExperienceSource.Manual,
          createdAt: "2026-01-01T00:00:00Z",
          updatedAt: "2026-01-01T00:00:00Z",
        },
        {
          id: "exp-jr",
          company: "TechCorp",
          title: "Junior Engineer",
          location: "San Francisco, CA",
          startDate: "Jan 2020",
          endDate: "Dec 2021",
          isCurrent: false,
          description: null,
          highlights: [],
          displayOrder: 1,
          source: ExperienceSource.Manual,
          createdAt: "2026-01-01T00:00:00Z",
          updatedAt: "2026-01-01T00:00:00Z",
        },
      ];
      render(<WorkExperienceSection profileExperiences={multiRoleExperiences} />);

      // Company name should appear once as a group header
      expect(screen.getAllByText("TechCorp")).toHaveLength(1);
      // Both role titles should be rendered
      expect(screen.getByText("Senior Engineer")).toBeInTheDocument();
      expect(screen.getByText("Junior Engineer")).toBeInTheDocument();
    });
  });

  describe("Hover-only kebab menu", () => {
    it("hides action menu by default until hover/focus", () => {
      render(
        <WorkExperienceSection profileExperiences={mockProfileExperiences} userId="user-123" />
      );

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
      render(
        <WorkExperienceSection profileExperiences={mockProfileExperiences} userId="user-123" />
      );

      const actionButtons = screen.getAllByRole("button", { name: "More actions" });
      const actionButton = actionButtons[0];
      const menuContainer = actionButton.parentElement;

      // Container should have focus-within classes for keyboard accessibility
      expect(menuContainer).toHaveClass("focus-within:opacity-100");
    });
  });
});
