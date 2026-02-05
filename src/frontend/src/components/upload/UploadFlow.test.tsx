import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import type { DocumentDetectionResult } from "./types";
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

  it("shows placeholder for unimplemented steps", async () => {
    render(<UploadFlow />);

    const simulateButton = screen.getByTestId("simulate-detection");
    fireEvent.click(simulateButton);

    await waitFor(() => {
      expect(screen.getByText(/test-document\.pdf/)).toBeInTheDocument();
      expect(screen.getByText(/coming soon/i)).toBeInTheDocument();
    });
  });

  it("hides DocumentUpload after detection completes", async () => {
    render(<UploadFlow />);

    const simulateButton = screen.getByTestId("simulate-detection");
    fireEvent.click(simulateButton);

    await waitFor(() => {
      expect(screen.queryByTestId("document-upload")).not.toBeInTheDocument();
    });
  });
});
