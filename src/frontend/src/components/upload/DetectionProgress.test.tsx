import { act, render, screen } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { DetectionProgress } from "./DetectionProgress";

describe("DetectionProgress", () => {
  let mockFetch: ReturnType<typeof vi.fn>;
  const defaultProps = {
    fileId: "file-123",
    onDetectionComplete: vi.fn(),
    onError: vi.fn(),
  };

  beforeEach(() => {
    vi.useFakeTimers();
    vi.clearAllMocks();
    mockFetch = vi.fn();
    vi.stubGlobal("fetch", mockFetch);
  });

  afterEach(() => {
    vi.useRealTimers();
    vi.restoreAllMocks();
  });

  function mockFetchSequence(...responses: unknown[]) {
    for (const response of responses) {
      mockFetch.mockResolvedValueOnce({
        json: () => Promise.resolve(response),
      });
    }
  }

  it("renders analyzing message", () => {
    mockFetchSequence({ data: { documentDetectionStatus: { status: "PENDING" } } });

    render(<DetectionProgress {...defaultProps} />);
    expect(screen.getByText("Analyzing document...")).toBeInTheDocument();
  });

  it("renders loading spinner", () => {
    mockFetchSequence({ data: { documentDetectionStatus: { status: "PENDING" } } });

    render(<DetectionProgress {...defaultProps} />);
    expect(screen.getByRole("img", { name: /loading spinner/i })).toBeInTheDocument();
  });

  it("calls onDetectionComplete when detection succeeds", async () => {
    const onDetectionComplete = vi.fn();
    const detection = {
      hasCareerInfo: true,
      hasTestimonial: false,
      testimonialAuthor: null,
      confidence: 0.95,
      summary: "A resume document.",
      documentTypeHint: "RESUME",
      fileId: "file-123",
    };

    mockFetchSequence({
      data: {
        documentDetectionStatus: {
          fileId: "file-123",
          status: "COMPLETED",
          detection,
          error: null,
        },
      },
    });

    render(<DetectionProgress {...defaultProps} onDetectionComplete={onDetectionComplete} />);

    await act(async () => {
      await vi.advanceTimersByTimeAsync(2000);
    });

    expect(onDetectionComplete).toHaveBeenCalledWith(detection);
  });

  it("calls onError when detection fails", async () => {
    const onError = vi.fn();

    mockFetchSequence({
      data: {
        documentDetectionStatus: {
          fileId: "file-123",
          status: "FAILED",
          detection: null,
          error: "text extraction failed",
        },
      },
    });

    render(<DetectionProgress {...defaultProps} onError={onError} />);

    await act(async () => {
      await vi.advanceTimersByTimeAsync(2000);
    });

    expect(onError).toHaveBeenCalledWith("text extraction failed");
  });

  it("shows error message when detection fails", async () => {
    mockFetchSequence({
      data: {
        documentDetectionStatus: {
          fileId: "file-123",
          status: "FAILED",
          detection: null,
          error: "text extraction failed",
        },
      },
    });

    render(<DetectionProgress {...defaultProps} />);

    await act(async () => {
      await vi.advanceTimersByTimeAsync(2000);
    });

    expect(screen.getByText("text extraction failed")).toBeInTheDocument();
  });

  it("continues polling on network errors", async () => {
    const detection = {
      hasCareerInfo: true,
      hasTestimonial: false,
      testimonialAuthor: null,
      confidence: 0.95,
      summary: "A resume.",
      documentTypeHint: "RESUME",
      fileId: "file-123",
    };

    const onDetectionComplete = vi.fn();

    mockFetch.mockRejectedValueOnce(new Error("Network error")).mockResolvedValueOnce({
      json: () =>
        Promise.resolve({
          data: {
            documentDetectionStatus: {
              fileId: "file-123",
              status: "COMPLETED",
              detection,
              error: null,
            },
          },
        }),
    });

    render(<DetectionProgress {...defaultProps} onDetectionComplete={onDetectionComplete} />);

    // First poll — network error, swallowed
    await act(async () => {
      await vi.advanceTimersByTimeAsync(2000);
    });

    // Second poll — success
    await act(async () => {
      await vi.advanceTimersByTimeAsync(2000);
    });

    expect(onDetectionComplete).toHaveBeenCalledWith(detection);
  });
});
