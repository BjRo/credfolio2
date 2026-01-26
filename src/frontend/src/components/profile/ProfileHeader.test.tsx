import { render, screen } from "@testing-library/react";
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
});
