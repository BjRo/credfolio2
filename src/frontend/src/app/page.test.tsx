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

// Mock the ResumeUpload component to avoid rendering its internals
vi.mock("@/components", () => ({
  ResumeUpload: ({ userId }: { userId: string }) => (
    <div data-testid="resume-upload" data-user-id={userId}>
      Resume Upload Mock
    </div>
  ),
}));

const mockUseQuery = useQuery as Mock;

describe("Home Page", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("shows loading state while fetching resumes", () => {
    mockUseQuery.mockReturnValue([{ fetching: true, data: undefined, error: undefined }]);
    render(<Home />);
    expect(screen.getByRole("status")).toBeInTheDocument();
  });

  it("shows upload form when user has no resumes", () => {
    mockUseQuery.mockReturnValue([{ fetching: false, data: { resumes: [] }, error: undefined }]);
    render(<Home />);
    expect(screen.getByTestId("resume-upload")).toBeInTheDocument();
    expect(screen.getByRole("heading", { name: /Upload Your Resume/i })).toBeInTheDocument();
  });

  it("shows upload form when user has no completed resumes", () => {
    mockUseQuery.mockReturnValue([
      {
        fetching: false,
        data: {
          resumes: [
            { id: "resume-1", status: "PENDING" },
            { id: "resume-2", status: "FAILED" },
          ],
        },
        error: undefined,
      },
    ]);
    render(<Home />);
    expect(screen.getByTestId("resume-upload")).toBeInTheDocument();
  });

  it("redirects to profile page when a completed resume exists", () => {
    mockUseQuery.mockReturnValue([
      {
        fetching: false,
        data: {
          resumes: [
            { id: "resume-1", status: "PENDING" },
            { id: "resume-completed", status: "COMPLETED" },
          ],
        },
        error: undefined,
      },
    ]);
    render(<Home />);
    expect(mockPush).toHaveBeenCalledWith("/profile/resume-completed");
  });

  it("redirects to the first completed resume when multiple exist", () => {
    mockUseQuery.mockReturnValue([
      {
        fetching: false,
        data: {
          resumes: [
            { id: "resume-first", status: "COMPLETED" },
            { id: "resume-second", status: "COMPLETED" },
          ],
        },
        error: undefined,
      },
    ]);
    render(<Home />);
    expect(mockPush).toHaveBeenCalledWith("/profile/resume-first");
  });

  it("shows upload form when query returns an error", () => {
    mockUseQuery.mockReturnValue([
      { fetching: false, data: undefined, error: new Error("Network error") },
    ]);
    render(<Home />);
    expect(screen.getByTestId("resume-upload")).toBeInTheDocument();
  });

  it("passes the demo user ID to ResumeUpload", () => {
    mockUseQuery.mockReturnValue([{ fetching: false, data: { resumes: [] }, error: undefined }]);
    render(<Home />);
    expect(screen.getByTestId("resume-upload")).toHaveAttribute(
      "data-user-id",
      "00000000-0000-0000-0000-000000000001"
    );
  });
});
