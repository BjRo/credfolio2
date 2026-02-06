import { render, screen } from "@testing-library/react";
import { useQuery } from "urql";
import { beforeEach, describe, expect, it, type Mock, vi } from "vitest";
import Home from "./page";

// Mock urql's useQuery
vi.mock("urql", () => ({
  useQuery: vi.fn(),
}));

// Mock next/navigation
const mockPush = vi.fn();
vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: mockPush }),
}));

const mockUseQuery = useQuery as Mock;

describe("Home Page", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("shows loading state while fetching profile", () => {
    mockUseQuery.mockReturnValue([{ fetching: true, data: undefined, error: undefined }]);
    render(<Home />);
    expect(screen.getByRole("status")).toBeInTheDocument();
  });

  it("redirects to upload page when user has no profile", () => {
    mockUseQuery.mockReturnValue([
      { fetching: false, data: { profileByUserId: null }, error: undefined },
    ]);
    render(<Home />);
    expect(mockPush).toHaveBeenCalledWith("/upload");
  });

  it("redirects to upload page when profile query returns no data", () => {
    mockUseQuery.mockReturnValue([{ fetching: false, data: undefined, error: undefined }]);
    render(<Home />);
    expect(mockPush).toHaveBeenCalledWith("/upload");
  });

  it("redirects to profile page when profile exists", () => {
    mockUseQuery.mockReturnValue([
      {
        fetching: false,
        data: { profileByUserId: { id: "profile-123" } },
        error: undefined,
      },
    ]);
    render(<Home />);
    expect(mockPush).toHaveBeenCalledWith("/profile/profile-123");
  });

  it("shows redirecting message when profile exists", () => {
    mockUseQuery.mockReturnValue([
      {
        fetching: false,
        data: { profileByUserId: { id: "profile-456" } },
        error: undefined,
      },
    ]);
    render(<Home />);
    expect(screen.getByText(/Redirecting/i)).toBeInTheDocument();
  });

  it("redirects to upload page when query returns an error", () => {
    mockUseQuery.mockReturnValue([
      { fetching: false, data: undefined, error: new Error("Network error") },
    ]);
    render(<Home />);
    expect(mockPush).toHaveBeenCalledWith("/upload");
  });
});
