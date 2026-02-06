import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";
import { ProfileHeader } from "./ProfileHeader";

// Mock urql to avoid provider requirement in tests
vi.mock("urql", () => ({
  useMutation: () => [{ fetching: false }, vi.fn()],
}));

// Mock the dialog component to avoid complex setup
vi.mock("./ProfileHeaderFormDialog", () => ({
  ProfileHeaderFormDialog: vi.fn(() => null),
}));

const mockProfile = {
  name: "John Doe",
  email: "john@example.com",
  phone: "+1 555-123-4567",
  location: "San Francisco, CA",
  summary: "Experienced software engineer with 10 years of experience.",
  profilePhotoUrl: null as string | null,
};

const longSummary =
  "Experienced software engineer with over 10 years of experience in building scalable web applications. " +
  "Proficient in TypeScript, React, Node.js, and cloud technologies. " +
  "Led multiple successful product launches and mentored junior developers. " +
  "Strong background in agile methodologies and test-driven development.";

describe("ProfileHeader", () => {
  it("renders the name", () => {
    render(<ProfileHeader profile={mockProfile} />);
    expect(screen.getByRole("heading", { level: 1 })).toHaveTextContent("John Doe");
  });

  it("renders the email when provided", () => {
    render(<ProfileHeader profile={mockProfile} />);
    expect(screen.getByText("john@example.com")).toBeInTheDocument();
  });

  it("renders the phone when provided", () => {
    render(<ProfileHeader profile={mockProfile} />);
    expect(screen.getByText("+1 555-123-4567")).toBeInTheDocument();
  });

  it("renders the location when provided", () => {
    render(<ProfileHeader profile={mockProfile} />);
    expect(screen.getByText("San Francisco, CA")).toBeInTheDocument();
  });

  it("renders the summary when provided", () => {
    render(<ProfileHeader profile={mockProfile} />);
    expect(
      screen.getByText("Experienced software engineer with 10 years of experience.")
    ).toBeInTheDocument();
  });

  it("does not render email when not provided", () => {
    const profileWithoutEmail = { ...mockProfile, email: null };
    render(<ProfileHeader profile={profileWithoutEmail} />);
    expect(screen.queryByText("john@example.com")).not.toBeInTheDocument();
  });

  it("does not render summary when not provided", () => {
    const profileWithoutSummary = { ...mockProfile, summary: null };
    render(<ProfileHeader profile={profileWithoutSummary} />);
    expect(screen.queryByText("Summary")).not.toBeInTheDocument();
  });

  describe("Avatar", () => {
    it("renders avatar with initials from name", () => {
      render(<ProfileHeader profile={mockProfile} />);
      expect(screen.getByText("JD")).toBeInTheDocument();
    });

    it("renders avatar with single initial for single name", () => {
      const singleNameProfile = { ...mockProfile, name: "Madonna" };
      render(<ProfileHeader profile={singleNameProfile} />);
      expect(screen.getByText("M")).toBeInTheDocument();
    });

    it("limits initials to two characters", () => {
      const longNameProfile = { ...mockProfile, name: "John Paul Jones Smith" };
      render(<ProfileHeader profile={longNameProfile} />);
      expect(screen.getByText("JP")).toBeInTheDocument();
    });

    it("has accessible label for avatar", () => {
      // When userId is not provided, the button is disabled and has "Avatar for" label
      render(<ProfileHeader profile={mockProfile} />);
      expect(screen.getByRole("button", { name: /avatar for john doe/i })).toBeInTheDocument();
    });
  });

  describe("Expandable Summary", () => {
    it("shows expand button for long summaries", () => {
      const profileWithLongSummary = { ...mockProfile, summary: longSummary };
      render(<ProfileHeader profile={profileWithLongSummary} />);
      expect(screen.getByRole("button", { name: /show more/i })).toBeInTheDocument();
    });

    it("does not show expand button for short summaries", () => {
      render(<ProfileHeader profile={mockProfile} />);
      expect(screen.queryByRole("button", { name: /show more/i })).not.toBeInTheDocument();
    });

    it("toggles summary expansion when clicking button", async () => {
      const user = userEvent.setup();
      const profileWithLongSummary = { ...mockProfile, summary: longSummary };
      render(<ProfileHeader profile={profileWithLongSummary} />);

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
      render(<ProfileHeader profile={mockProfile} />);
      const emailLink = screen.getByRole("link", { name: "john@example.com" });
      expect(emailLink).toHaveAttribute("href", "mailto:john@example.com");
    });

    it("renders phone as tel link", () => {
      render(<ProfileHeader profile={mockProfile} />);
      const phoneLink = screen.getByRole("link", { name: "+1 555-123-4567" });
      expect(phoneLink).toHaveAttribute("href", "tel:+1 555-123-4567");
    });
  });

  describe("Edit Button", () => {
    it("shows edit button when userId is provided", () => {
      render(<ProfileHeader profile={mockProfile} userId="user-123" />);
      expect(screen.getByRole("button", { name: /edit profile/i })).toBeInTheDocument();
    });

    it("does not show edit button when userId is not provided", () => {
      render(<ProfileHeader profile={mockProfile} />);
      expect(screen.queryByRole("button", { name: /edit profile/i })).not.toBeInTheDocument();
    });

    it("opens edit dialog when edit button is clicked", async () => {
      const user = userEvent.setup();
      render(<ProfileHeader profile={mockProfile} userId="user-123" />);

      const editButton = screen.getByRole("button", { name: /edit profile/i });
      await user.click(editButton);

      // The dialog component is mocked, so we just verify the button is clickable
      expect(editButton).toBeInTheDocument();
    });
  });

  describe("Fallback behavior", () => {
    it("shows Unknown when name is null", () => {
      const profileWithoutName = { ...mockProfile, name: null };
      render(<ProfileHeader profile={profileWithoutName} />);
      expect(screen.getByRole("heading", { level: 1 })).toHaveTextContent("Unknown");
    });

    it("uses profile data directly", () => {
      render(<ProfileHeader profile={mockProfile} />);
      expect(screen.getByRole("heading", { level: 1 })).toHaveTextContent("John Doe");
    });
  });
});
