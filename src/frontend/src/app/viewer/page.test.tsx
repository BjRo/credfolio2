import { act, render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { useQuery } from "urql";
import { beforeEach, describe, expect, it, type Mock, vi } from "vitest";
import ViewerPage from "./page";

// Mock urql
vi.mock("urql", () => ({
  useQuery: vi.fn(),
}));

// Mock next/navigation
const mockBack = vi.fn();
const mockPush = vi.fn();
let mockSearchParams = new URLSearchParams();
vi.mock("next/navigation", () => ({
  useRouter: () => ({ back: mockBack, push: mockPush }),
  useSearchParams: () => mockSearchParams,
}));

// Capture the onHighlightResult callback from PDFViewer renders
let capturedOnHighlightResult: ((found: boolean) => void) | undefined;

// Mock next/dynamic to return a synchronous mock PDFViewer
// This avoids react-pdf browser dependencies and async React.lazy resolution
vi.mock("next/dynamic", () => ({
  __esModule: true,
  default: () => {
    return function MockPDFViewer(props: {
      fileUrl: string;
      highlightText?: string;
      onHighlightResult?: (found: boolean) => void;
    }) {
      capturedOnHighlightResult = props.onHighlightResult;
      return (
        <div
          data-testid="pdf-viewer"
          data-file-url={props.fileUrl}
          data-highlight={props.highlightText || ""}
        />
      );
    };
  },
}));

const mockUseQuery = useQuery as Mock;

const mockLetterData = {
  referenceLetter: {
    id: "letter-1",
    title: "Reference Letter",
    authorName: "Jane Smith",
    authorTitle: "Engineering Manager",
    organization: "Acme Corp",
    file: {
      id: "file-1",
      url: "https://storage.example.com/presigned-url",
      filename: "reference.pdf",
      contentType: "application/pdf",
    },
  },
};

describe("ViewerPage", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    capturedOnHighlightResult = undefined;
    mockSearchParams = new URLSearchParams();
  });

  describe("missing or invalid letterId", () => {
    it("shows error when letterId is missing", () => {
      mockSearchParams = new URLSearchParams();
      mockUseQuery.mockReturnValue([{ fetching: false, data: undefined, error: undefined }]);
      render(<ViewerPage />);
      expect(screen.getByText("Document not found")).toBeInTheDocument();
    });

    it("shows error when letterId is not a valid UUID", () => {
      mockSearchParams = new URLSearchParams({ letterId: "not-a-uuid" });
      mockUseQuery.mockReturnValue([{ fetching: false, data: undefined, error: undefined }]);
      render(<ViewerPage />);
      expect(screen.getByText("Document not found")).toBeInTheDocument();
    });
  });

  describe("loading state", () => {
    it("shows loading skeleton while fetching", () => {
      mockSearchParams = new URLSearchParams({
        letterId: "550e8400-e29b-41d4-a716-446655440000",
      });
      mockUseQuery.mockReturnValue([{ fetching: true, data: undefined, error: undefined }]);
      render(<ViewerPage />);
      expect(screen.getByRole("status")).toBeInTheDocument();
    });
  });

  describe("error states", () => {
    it("shows error when GraphQL query fails", () => {
      mockSearchParams = new URLSearchParams({
        letterId: "550e8400-e29b-41d4-a716-446655440000",
      });
      mockUseQuery.mockReturnValue([
        { fetching: false, data: undefined, error: new Error("Network error") },
      ]);
      render(<ViewerPage />);
      expect(screen.getByText("Failed to load document")).toBeInTheDocument();
    });

    it("shows error when reference letter not found", () => {
      mockSearchParams = new URLSearchParams({
        letterId: "550e8400-e29b-41d4-a716-446655440000",
      });
      mockUseQuery.mockReturnValue([
        { fetching: false, data: { referenceLetter: null }, error: undefined },
      ]);
      render(<ViewerPage />);
      expect(screen.getByText("Document not found")).toBeInTheDocument();
    });

    it("shows error when file is missing", () => {
      mockSearchParams = new URLSearchParams({
        letterId: "550e8400-e29b-41d4-a716-446655440000",
      });
      mockUseQuery.mockReturnValue([
        {
          fetching: false,
          data: { referenceLetter: { ...mockLetterData.referenceLetter, file: null } },
          error: undefined,
        },
      ]);
      render(<ViewerPage />);
      expect(screen.getByText("Document file unavailable")).toBeInTheDocument();
    });
  });

  describe("success states", () => {
    it("renders PDFViewer with file URL when letter is loaded", () => {
      mockSearchParams = new URLSearchParams({
        letterId: "550e8400-e29b-41d4-a716-446655440000",
      });
      mockUseQuery.mockReturnValue([{ fetching: false, data: mockLetterData, error: undefined }]);
      render(<ViewerPage />);
      const viewer = screen.getByTestId("pdf-viewer");
      expect(viewer).toBeInTheDocument();
      expect(viewer).toHaveAttribute("data-file-url", "https://storage.example.com/presigned-url");
    });

    it("passes highlight text to PDFViewer when provided", () => {
      mockSearchParams = new URLSearchParams({
        letterId: "550e8400-e29b-41d4-a716-446655440000",
        highlight: "excellent engineer",
      });
      mockUseQuery.mockReturnValue([{ fetching: false, data: mockLetterData, error: undefined }]);
      render(<ViewerPage />);
      const viewer = screen.getByTestId("pdf-viewer");
      expect(viewer).toHaveAttribute("data-highlight", "excellent engineer");
    });

    it("does not pass highlight text when not provided", () => {
      mockSearchParams = new URLSearchParams({
        letterId: "550e8400-e29b-41d4-a716-446655440000",
      });
      mockUseQuery.mockReturnValue([{ fetching: false, data: mockLetterData, error: undefined }]);
      render(<ViewerPage />);
      const viewer = screen.getByTestId("pdf-viewer");
      expect(viewer).toHaveAttribute("data-highlight", "");
    });
  });

  describe("toolbar", () => {
    it("shows document title from letter title", () => {
      mockSearchParams = new URLSearchParams({
        letterId: "550e8400-e29b-41d4-a716-446655440000",
      });
      mockUseQuery.mockReturnValue([{ fetching: false, data: mockLetterData, error: undefined }]);
      render(<ViewerPage />);
      expect(screen.getByText("Reference Letter")).toBeInTheDocument();
    });

    it("shows author name fallback when title is missing", () => {
      mockSearchParams = new URLSearchParams({
        letterId: "550e8400-e29b-41d4-a716-446655440000",
      });
      const letterNoTitle = {
        referenceLetter: {
          ...mockLetterData.referenceLetter,
          title: null,
        },
      };
      mockUseQuery.mockReturnValue([{ fetching: false, data: letterNoTitle, error: undefined }]);
      render(<ViewerPage />);
      expect(screen.getByText("Reference from Jane Smith")).toBeInTheDocument();
    });

    it("shows subtitle with author title and organization", () => {
      mockSearchParams = new URLSearchParams({
        letterId: "550e8400-e29b-41d4-a716-446655440000",
      });
      mockUseQuery.mockReturnValue([{ fetching: false, data: mockLetterData, error: undefined }]);
      render(<ViewerPage />);
      expect(screen.getByText("Engineering Manager, Acme Corp")).toBeInTheDocument();
    });

    it("navigates back when back button is clicked and history exists", async () => {
      const user = userEvent.setup();
      Object.defineProperty(window.history, "length", { value: 2, writable: true });
      mockSearchParams = new URLSearchParams({
        letterId: "550e8400-e29b-41d4-a716-446655440000",
      });
      mockUseQuery.mockReturnValue([{ fetching: false, data: mockLetterData, error: undefined }]);
      render(<ViewerPage />);
      await user.click(screen.getByLabelText("Go back"));
      expect(mockBack).toHaveBeenCalled();
      Object.defineProperty(window.history, "length", { value: 1, writable: true });
    });

    it("navigates to home when back button is clicked with no history", async () => {
      const user = userEvent.setup();
      Object.defineProperty(window.history, "length", { value: 1, writable: true });
      mockSearchParams = new URLSearchParams({
        letterId: "550e8400-e29b-41d4-a716-446655440000",
      });
      mockUseQuery.mockReturnValue([{ fetching: false, data: mockLetterData, error: undefined }]);
      render(<ViewerPage />);
      await user.click(screen.getByLabelText("Go back"));
      expect(mockPush).toHaveBeenCalledWith("/");
    });
  });

  describe("info banner", () => {
    it("shows banner when highlight text not found", () => {
      mockSearchParams = new URLSearchParams({
        letterId: "550e8400-e29b-41d4-a716-446655440000",
        highlight: "some text",
      });
      mockUseQuery.mockReturnValue([{ fetching: false, data: mockLetterData, error: undefined }]);
      render(<ViewerPage />);

      act(() => {
        capturedOnHighlightResult?.(false);
      });

      expect(
        screen.getByText("Could not locate exact quote — showing full document")
      ).toBeInTheDocument();
    });

    it("does not show banner when highlight text is found", () => {
      mockSearchParams = new URLSearchParams({
        letterId: "550e8400-e29b-41d4-a716-446655440000",
        highlight: "some text",
      });
      mockUseQuery.mockReturnValue([{ fetching: false, data: mockLetterData, error: undefined }]);
      render(<ViewerPage />);

      act(() => {
        capturedOnHighlightResult?.(true);
      });

      expect(
        screen.queryByText("Could not locate exact quote — showing full document")
      ).not.toBeInTheDocument();
    });

    it("does not show banner when no highlight param is provided", () => {
      mockSearchParams = new URLSearchParams({
        letterId: "550e8400-e29b-41d4-a716-446655440000",
      });
      mockUseQuery.mockReturnValue([{ fetching: false, data: mockLetterData, error: undefined }]);
      render(<ViewerPage />);

      expect(
        screen.queryByText("Could not locate exact quote — showing full document")
      ).not.toBeInTheDocument();
    });

    it("dismisses banner when X is clicked", async () => {
      const user = userEvent.setup();
      mockSearchParams = new URLSearchParams({
        letterId: "550e8400-e29b-41d4-a716-446655440000",
        highlight: "some text",
      });
      mockUseQuery.mockReturnValue([{ fetching: false, data: mockLetterData, error: undefined }]);
      render(<ViewerPage />);

      act(() => {
        capturedOnHighlightResult?.(false);
      });

      expect(
        screen.getByText("Could not locate exact quote — showing full document")
      ).toBeInTheDocument();

      await user.click(screen.getByLabelText("Dismiss banner"));

      expect(
        screen.queryByText("Could not locate exact quote — showing full document")
      ).not.toBeInTheDocument();
    });
  });
});
