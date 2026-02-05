import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import type { DetectionResultsProps, DocumentDetectionResult } from "./types";
import { UploadFlow } from "./UploadFlow";

// Mock the child components to isolate UploadFlow logic
vi.mock("./DocumentUpload", () => ({
  DocumentUpload: ({
    onDetectionComplete,
  }: {
    userId: string;
    onDetectionComplete: (detection: DocumentDetectionResult, fileName: string) => void;
  }) => (
    <div data-testid="document-upload">
      <button
        type="button"
        data-testid="simulate-detection"
        onClick={() =>
          onDetectionComplete(
            {
              hasCareerInfo: true,
              hasTestimonial: true,
              testimonialAuthor: "Jane Doe",
              confidence: 0.92,
              summary: "A hybrid document",
              documentTypeHint: "HYBRID",
              fileId: "file-456",
            },
            "test-document.pdf"
          )
        }
      >
        Simulate Detection
      </button>
    </div>
  ),
}));

vi.mock("./DetectionResults", () => ({
  DetectionResults: ({ detection, fileName, onProceed, onCancel }: DetectionResultsProps) => (
    <div data-testid="detection-results">
      <span data-testid="detection-filename">{fileName}</span>
      <span data-testid="detection-summary">{detection.summary}</span>
      <button type="button" data-testid="proceed-button" onClick={() => onProceed(true, true)}>
        Proceed
      </button>
      <button type="button" data-testid="cancel-button" onClick={onCancel}>
        Cancel
      </button>
    </div>
  ),
}));

describe("UploadFlow", () => {
  it("renders StepIndicator with upload as current step", () => {
    render(<UploadFlow />);
    const uploadStep = screen.getByText("Upload").closest("li");
    expect(uploadStep).toHaveAttribute("aria-current", "step");
  });

  it("renders DocumentUpload component in upload step", () => {
    render(<UploadFlow />);
    expect(screen.getByTestId("document-upload")).toBeInTheDocument();
  });

  it("transitions to review-detection step when detection completes", async () => {
    render(<UploadFlow />);

    const simulateButton = screen.getByTestId("simulate-detection");
    fireEvent.click(simulateButton);

    await waitFor(() => {
      const reviewStep = screen.getByText("Review Detection").closest("li");
      expect(reviewStep).toHaveAttribute("aria-current", "step");
    });
  });

  it("shows DetectionResults component after detection completes", async () => {
    render(<UploadFlow />);

    fireEvent.click(screen.getByTestId("simulate-detection"));

    await waitFor(() => {
      expect(screen.getByTestId("detection-results")).toBeInTheDocument();
      expect(screen.getByTestId("detection-filename")).toHaveTextContent("test-document.pdf");
    });
  });

  it("hides DocumentUpload after detection completes", async () => {
    render(<UploadFlow />);

    fireEvent.click(screen.getByTestId("simulate-detection"));

    await waitFor(() => {
      expect(screen.queryByTestId("document-upload")).not.toBeInTheDocument();
    });
  });

  it("transitions to extract step when user proceeds", async () => {
    render(<UploadFlow />);

    fireEvent.click(screen.getByTestId("simulate-detection"));

    await waitFor(() => {
      expect(screen.getByTestId("detection-results")).toBeInTheDocument();
    });

    fireEvent.click(screen.getByTestId("proceed-button"));

    await waitFor(() => {
      const extractStep = screen.getByText("Extract").closest("li");
      expect(extractStep).toHaveAttribute("aria-current", "step");
    });
  });

  it("returns to upload step when user cancels", async () => {
    render(<UploadFlow />);

    fireEvent.click(screen.getByTestId("simulate-detection"));

    await waitFor(() => {
      expect(screen.getByTestId("detection-results")).toBeInTheDocument();
    });

    fireEvent.click(screen.getByTestId("cancel-button"));

    await waitFor(() => {
      expect(screen.getByTestId("document-upload")).toBeInTheDocument();
      const uploadStep = screen.getByText("Upload").closest("li");
      expect(uploadStep).toHaveAttribute("aria-current", "step");
    });
  });
});
