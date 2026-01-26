import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it } from "vitest";
import { ProfileHeader } from "./ProfileHeader";
import type { ProfileData } from "./types";

const mockProfileData: ProfileData = {
  name: "John Doe",
  email: "john@example.com",
  phone: "+1 555-123-4567",
  location: "San Francisco, CA",
  summary: "Experienced software engineer with 10 years of experience.",
  experience: [],
  education: [],
  skills: [],
  extractedAt: "2024-01-01T00:00:00Z",
  confidence: 0.95,
};

const longSummary =
  "Experienced software engineer with over 10 years of experience in building scalable web applications. " +
  "Proficient in TypeScript, React, Node.js, and cloud technologies. " +
  "Led multiple successful product launches and mentored junior developers. " +
  "Strong background in agile methodologies and test-driven development.";

describe("ProfileHeader", () => {
  it("renders the name", () => {
    render(<ProfileHeader data={mockProfileData} />);
    expect(screen.getByRole("heading", { level: 1 })).toHaveTextContent("John Doe");
  });

  it("renders the email when provided", () => {
    render(<ProfileHeader data={mockProfileData} />);
    expect(screen.getByText("john@example.com")).toBeInTheDocument();
  });

  it("renders the phone when provided", () => {
    render(<ProfileHeader data={mockProfileData} />);
    expect(screen.getByText("+1 555-123-4567")).toBeInTheDocument();
  });

  it("renders the location when provided", () => {
    render(<ProfileHeader data={mockProfileData} />);
    expect(screen.getByText("San Francisco, CA")).toBeInTheDocument();
  });

  it("renders the summary when provided", () => {
    render(<ProfileHeader data={mockProfileData} />);
    expect(
      screen.getByText("Experienced software engineer with 10 years of experience.")
    ).toBeInTheDocument();
  });

  it("renders the confidence score", () => {
    render(<ProfileHeader data={mockProfileData} />);
    expect(screen.getByText("Confidence: 95%")).toBeInTheDocument();
  });

  it("does not render email when not provided", () => {
    const dataWithoutEmail = { ...mockProfileData, email: null };
    render(<ProfileHeader data={dataWithoutEmail} />);
    expect(screen.queryByText("john@example.com")).not.toBeInTheDocument();
  });

  it("does not render summary when not provided", () => {
    const dataWithoutSummary = { ...mockProfileData, summary: null };
    render(<ProfileHeader data={dataWithoutSummary} />);
    expect(screen.queryByText("Summary")).not.toBeInTheDocument();
  });

  describe("Avatar", () => {
    it("renders avatar with initials from name", () => {
      render(<ProfileHeader data={mockProfileData} />);
      expect(screen.getByText("JD")).toBeInTheDocument();
    });

    it("renders avatar with single initial for single name", () => {
      const singleNameData = { ...mockProfileData, name: "Madonna" };
      render(<ProfileHeader data={singleNameData} />);
      expect(screen.getByText("M")).toBeInTheDocument();
    });

    it("limits initials to two characters", () => {
      const longNameData = { ...mockProfileData, name: "John Paul Jones Smith" };
      render(<ProfileHeader data={longNameData} />);
      expect(screen.getByText("JP")).toBeInTheDocument();
    });

    it("has accessible label for avatar", () => {
      render(<ProfileHeader data={mockProfileData} />);
      expect(screen.getByLabelText("Avatar for John Doe")).toBeInTheDocument();
    });
  });

  describe("Expandable Summary", () => {
    it("shows expand button for long summaries", () => {
      const dataWithLongSummary = { ...mockProfileData, summary: longSummary };
      render(<ProfileHeader data={dataWithLongSummary} />);
      expect(screen.getByRole("button", { name: /show more/i })).toBeInTheDocument();
    });

    it("does not show expand button for short summaries", () => {
      render(<ProfileHeader data={mockProfileData} />);
      expect(screen.queryByRole("button", { name: /show more/i })).not.toBeInTheDocument();
    });

    it("toggles summary expansion when clicking button", async () => {
      const user = userEvent.setup();
      const dataWithLongSummary = { ...mockProfileData, summary: longSummary };
      render(<ProfileHeader data={dataWithLongSummary} />);

      const button = screen.getByRole("button", { name: /show more/i });
      expect(button).toHaveAttribute("aria-expanded", "false");

      await user.click(button);
      expect(screen.getByRole("button", { name: /show less/i })).toBeInTheDocument();
      expect(screen.getByRole("button", { name: /show less/i })).toHaveAttribute(
        "aria-expanded",
        "true"
      );

      await user.click(screen.getByRole("button", { name: /show less/i }));
      expect(screen.getByRole("button", { name: /show more/i })).toBeInTheDocument();
    });
  });

  describe("Contact Links", () => {
    it("renders email as mailto link", () => {
      render(<ProfileHeader data={mockProfileData} />);
      const emailLink = screen.getByRole("link", { name: "john@example.com" });
      expect(emailLink).toHaveAttribute("href", "mailto:john@example.com");
    });

    it("renders phone as tel link", () => {
      render(<ProfileHeader data={mockProfileData} />);
      const phoneLink = screen.getByRole("link", { name: "+1 555-123-4567" });
      expect(phoneLink).toHaveAttribute("href", "tel:+1 555-123-4567");
    });
  });
});
