import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import type {
  DetectionProgressProps,
  DetectionResultsProps,
  DocumentDetectionResult,
  ExtractionProgressProps,
  ExtractionResults,
  ExtractionReviewProps,
  ProcessDocumentIds,
} from "./types";
import { UploadFlow } from "./UploadFlow";

const mockPush = vi.fn();
vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: mockPush }),
}));

// Mock the child components to isolate UploadFlow logic
vi.mock("./DocumentUpload", () => ({
  DocumentUpload: ({
    onUploadComplete,
  }: {
    userId: string;
    onUploadComplete: (fileId: string, fileName: string) => void;
  }) => (
    <div data-testid="document-upload">
      <button
        type="button"
        data-testid="simulate-upload"
        onClick={() => onUploadComplete("file-456", "test-document.pdf")}
      >
        Simulate Upload
      </button>
    </div>
  ),
}));

const mockDetection: DocumentDetectionResult = {
  hasCareerInfo: true,
  hasTestimonial: true,
  testimonialAuthor: "Jane Doe",
  confidence: 0.92,
  summary: "A hybrid document",
  documentTypeHint: "HYBRID",
  fileId: "file-456",
};

vi.mock("./DetectionProgress", () => ({
  DetectionProgress: ({ onDetectionComplete, onError }: DetectionProgressProps) => (
    <div data-testid="detection-progress">
      <button
        type="button"
        data-testid="simulate-detection-complete"
        onClick={() => onDetectionComplete(mockDetection)}
      >
        Complete Detection
      </button>
      <button
        type="button"
        data-testid="simulate-detection-error"
        onClick={() => onError("Detection failed")}
      >
        Fail Detection
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

const mockExtractionResults: ExtractionResults = {
  resume: {
    id: "resume-1",
    status: "COMPLETED",
    extractedData: {
      name: "John Doe",
      email: "john@example.com",
      phone: null,
      location: "Berlin",
      summary: "A developer",
      extractedAt: "2025-01-01T00:00:00Z",
      confidence: 0.95,
    },
    errorMessage: null,
  },
  referenceLetter: null,
};

const mockProcessIds: ProcessDocumentIds = {
  resumeId: "resume-1",
  referenceLetterID: null,
};

vi.mock("./ExtractionProgress", () => ({
  ExtractionProgress: ({ onComplete, onError }: ExtractionProgressProps) => (
    <div data-testid="extraction-progress">
      <button
        type="button"
        data-testid="simulate-extraction-complete"
        onClick={() => onComplete(mockExtractionResults, mockProcessIds)}
      >
        Complete Extraction
      </button>
      <button
        type="button"
        data-testid="simulate-extraction-error"
        onClick={() => onError("Extraction failed")}
      >
        Fail Extraction
      </button>
    </div>
  ),
}));

vi.mock("./ExtractionReview", () => ({
  ExtractionReview: ({ onImportComplete, onBack }: ExtractionReviewProps) => (
    <div data-testid="extraction-review">
      <button
        type="button"
        data-testid="simulate-import"
        onClick={() => onImportComplete("profile-1")}
      >
        Import
      </button>
      <button type="button" data-testid="back-from-review" onClick={onBack}>
        Back
      </button>
    </div>
  ),
}));

// Helper: navigate from upload → detect → review-detection
async function navigateToDetectionResults() {
  fireEvent.click(screen.getByTestId("simulate-upload"));
  await waitFor(() => {
    expect(screen.getByTestId("detection-progress")).toBeInTheDocument();
  });
  fireEvent.click(screen.getByTestId("simulate-detection-complete"));
  await waitFor(() => {
    expect(screen.getByTestId("detection-results")).toBeInTheDocument();
  });
}

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

  it("transitions to detect step when upload completes", async () => {
    render(<UploadFlow />);

    fireEvent.click(screen.getByTestId("simulate-upload"));

    await waitFor(() => {
      expect(screen.getByTestId("detection-progress")).toBeInTheDocument();
      const analyzeStep = screen.getByText("Analyze").closest("li");
      expect(analyzeStep).toHaveAttribute("aria-current", "step");
    });
  });

  it("transitions to review-detection step when detection completes", async () => {
    render(<UploadFlow />);

    fireEvent.click(screen.getByTestId("simulate-upload"));
    await waitFor(() => {
      expect(screen.getByTestId("detection-progress")).toBeInTheDocument();
    });

    fireEvent.click(screen.getByTestId("simulate-detection-complete"));

    await waitFor(() => {
      const reviewStep = screen.getByText("Review Detection").closest("li");
      expect(reviewStep).toHaveAttribute("aria-current", "step");
    });
  });

  it("shows DetectionResults component after detection completes", async () => {
    render(<UploadFlow />);

    await navigateToDetectionResults();

    expect(screen.getByTestId("detection-filename")).toHaveTextContent("test-document.pdf");
  });

  it("hides DocumentUpload after upload completes", async () => {
    render(<UploadFlow />);

    fireEvent.click(screen.getByTestId("simulate-upload"));

    await waitFor(() => {
      expect(screen.queryByTestId("document-upload")).not.toBeInTheDocument();
    });
  });

  it("transitions to extract step when user proceeds", async () => {
    render(<UploadFlow />);

    await navigateToDetectionResults();

    fireEvent.click(screen.getByTestId("proceed-button"));

    await waitFor(() => {
      const extractStep = screen.getByText("Extract").closest("li");
      expect(extractStep).toHaveAttribute("aria-current", "step");
    });
  });

  it("returns to upload step when user cancels", async () => {
    render(<UploadFlow />);

    await navigateToDetectionResults();

    fireEvent.click(screen.getByTestId("cancel-button"));

    await waitFor(() => {
      expect(screen.getByTestId("document-upload")).toBeInTheDocument();
      const uploadStep = screen.getByText("Upload").closest("li");
      expect(uploadStep).toHaveAttribute("aria-current", "step");
    });
  });

  it("shows ExtractionProgress in extract step", async () => {
    render(<UploadFlow />);

    await navigateToDetectionResults();

    fireEvent.click(screen.getByTestId("proceed-button"));

    await waitFor(() => {
      expect(screen.getByTestId("extraction-progress")).toBeInTheDocument();
      const extractStep = screen.getByText("Extract").closest("li");
      expect(extractStep).toHaveAttribute("aria-current", "step");
    });
  });

  it("transitions to review-results when extraction completes", async () => {
    render(<UploadFlow />);

    await navigateToDetectionResults();

    fireEvent.click(screen.getByTestId("proceed-button"));
    await waitFor(() => {
      expect(screen.getByTestId("extraction-progress")).toBeInTheDocument();
    });

    fireEvent.click(screen.getByTestId("simulate-extraction-complete"));

    await waitFor(() => {
      expect(screen.getByTestId("extraction-review")).toBeInTheDocument();
      const reviewStep = screen.getByText("Review Results").closest("li");
      expect(reviewStep).toHaveAttribute("aria-current", "step");
    });
  });

  it("redirects to profile when import completes", async () => {
    render(<UploadFlow />);

    await navigateToDetectionResults();

    fireEvent.click(screen.getByTestId("proceed-button"));
    await waitFor(() => {
      expect(screen.getByTestId("extraction-progress")).toBeInTheDocument();
    });
    fireEvent.click(screen.getByTestId("simulate-extraction-complete"));
    await waitFor(() => {
      expect(screen.getByTestId("extraction-review")).toBeInTheDocument();
    });

    // Click import
    fireEvent.click(screen.getByTestId("simulate-import"));

    await waitFor(() => {
      const importStep = screen.getByText("Import").closest("li");
      expect(importStep).toHaveAttribute("aria-current", "step");
    });

    expect(mockPush).toHaveBeenCalledWith("/profile/profile-1");
  });

  it("returns to review-detection when back is clicked from review-results", async () => {
    render(<UploadFlow />);

    await navigateToDetectionResults();

    fireEvent.click(screen.getByTestId("proceed-button"));
    await waitFor(() => {
      expect(screen.getByTestId("extraction-progress")).toBeInTheDocument();
    });
    fireEvent.click(screen.getByTestId("simulate-extraction-complete"));
    await waitFor(() => {
      expect(screen.getByTestId("extraction-review")).toBeInTheDocument();
    });

    // Click back
    fireEvent.click(screen.getByTestId("back-from-review"));

    await waitFor(() => {
      expect(screen.getByTestId("detection-results")).toBeInTheDocument();
      const reviewDetectionStep = screen.getByText("Review Detection").closest("li");
      expect(reviewDetectionStep).toHaveAttribute("aria-current", "step");
    });
  });
});
