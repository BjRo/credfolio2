import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { DocumentUpload } from "./DocumentUpload";

// Mock XMLHttpRequest for upload tests
class MockXHR {
  static instances: MockXHR[] = [];
  status = 200;
  responseText = "";
  upload = {
    addEventListener: vi.fn(),
  };
  addEventListener = vi.fn();
  open = vi.fn();
  send = vi.fn();

  constructor() {
    MockXHR.instances.push(this);
  }

  static reset() {
    MockXHR.instances = [];
  }

  static getLastInstance() {
    return MockXHR.instances[MockXHR.instances.length - 1];
  }
}

vi.stubGlobal("XMLHttpRequest", MockXHR);

describe("DocumentUpload", () => {
  const defaultProps = {
    userId: "test-user-id",
    onDetectionComplete: vi.fn(),
  };

  beforeEach(() => {
    vi.clearAllMocks();
    MockXHR.reset();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it("renders drop zone with upload instructions", () => {
    render(<DocumentUpload {...defaultProps} />);
    expect(screen.getByText(/drag and drop/i)).toBeInTheDocument();
    expect(screen.getByText(/pdf, docx, txt/i)).toBeInTheDocument();
  });

  it("renders file input with correct accept types", () => {
    render(<DocumentUpload {...defaultProps} />);
    const input = screen.getByTestId("file-input");
    expect(input).toBeInTheDocument();
    expect(input).toHaveAttribute("type", "file");
    expect(input).toHaveAttribute(
      "accept",
      "application/pdf,application/vnd.openxmlformats-officedocument.wordprocessingml.document,text/plain"
    );
  });

  it("shows error for invalid file type", async () => {
    render(<DocumentUpload {...defaultProps} />);

    const input = screen.getByTestId("file-input");
    const invalidFile = new File(["test"], "test.exe", { type: "application/x-executable" });

    fireEvent.change(input, { target: { files: [invalidFile] } });

    await waitFor(() => {
      expect(screen.getByText(/invalid file type/i)).toBeInTheDocument();
    });
  });

  it("shows error for file exceeding size limit", async () => {
    render(<DocumentUpload {...defaultProps} />);

    const input = screen.getByTestId("file-input");
    const largeFile = new File(["x"], "large.pdf", { type: "application/pdf" });
    Object.defineProperty(largeFile, "size", { value: 11 * 1024 * 1024 });

    fireEvent.change(input, { target: { files: [largeFile] } });

    await waitFor(() => {
      expect(screen.getByText(/file too large/i)).toBeInTheDocument();
    });
  });

  it("shows uploading state when valid file selected", async () => {
    render(<DocumentUpload {...defaultProps} />);

    const input = screen.getByTestId("file-input");
    const validFile = new File(["test content"], "resume.pdf", { type: "application/pdf" });

    fireEvent.change(input, { target: { files: [validFile] } });

    await waitFor(() => {
      expect(screen.getByText(/uploading/i)).toBeInTheDocument();
    });
  });

  it("initiates XHR upload with detectDocumentContent mutation", async () => {
    render(<DocumentUpload {...defaultProps} />);

    const input = screen.getByTestId("file-input");
    const validFile = new File(["test content"], "resume.pdf", { type: "application/pdf" });

    fireEvent.change(input, { target: { files: [validFile] } });

    await waitFor(() => {
      const xhr = MockXHR.getLastInstance();
      expect(xhr).toBeDefined();
      expect(xhr.open).toHaveBeenCalledWith("POST", expect.any(String));
      expect(xhr.send).toHaveBeenCalled();
    });

    // Verify the mutation is detectDocumentContent
    const xhr = MockXHR.getLastInstance();
    const sentFormData = xhr.send.mock.calls[0][0] as FormData;
    const operations = sentFormData.get("operations") as string;
    expect(operations).toContain("detectDocumentContent");
  });

  it("calls onDetectionComplete with detection result on success", async () => {
    const onDetectionComplete = vi.fn();
    render(<DocumentUpload {...defaultProps} onDetectionComplete={onDetectionComplete} />);

    const input = screen.getByTestId("file-input");
    const validFile = new File(["test content"], "resume.pdf", { type: "application/pdf" });

    fireEvent.change(input, { target: { files: [validFile] } });

    await waitFor(() => {
      const xhr = MockXHR.getLastInstance();
      expect(xhr).toBeDefined();
      expect(xhr.addEventListener).toHaveBeenCalled();
    });

    const xhr = MockXHR.getLastInstance();
    xhr.status = 200;
    xhr.responseText = JSON.stringify({
      data: {
        detectDocumentContent: {
          __typename: "DetectDocumentContentResult",
          detection: {
            hasCareerInfo: true,
            hasTestimonial: false,
            testimonialAuthor: null,
            confidence: 0.95,
            summary: "A professional resume",
            documentTypeHint: "RESUME",
            fileId: "file-123",
          },
        },
      },
    });

    const loadHandler = xhr.addEventListener.mock.calls.find(
      (call: [string, unknown]) => call[0] === "load"
    )?.[1] as () => void;
    expect(loadHandler).toBeDefined();
    loadHandler();

    await waitFor(() => {
      expect(onDetectionComplete).toHaveBeenCalledWith(
        {
          hasCareerInfo: true,
          hasTestimonial: false,
          testimonialAuthor: null,
          confidence: 0.95,
          summary: "A professional resume",
          documentTypeHint: "RESUME",
          fileId: "file-123",
        },
        "resume.pdf"
      );
    });
  });

  it("shows error for FileValidationError response", async () => {
    render(<DocumentUpload {...defaultProps} />);

    const input = screen.getByTestId("file-input");
    const validFile = new File(["test content"], "resume.pdf", { type: "application/pdf" });

    fireEvent.change(input, { target: { files: [validFile] } });

    await waitFor(() => {
      const xhr = MockXHR.getLastInstance();
      expect(xhr).toBeDefined();
      expect(xhr.addEventListener).toHaveBeenCalled();
    });

    const xhr = MockXHR.getLastInstance();
    xhr.status = 200;
    xhr.responseText = JSON.stringify({
      data: {
        detectDocumentContent: {
          __typename: "FileValidationError",
          message: "File content is corrupted",
          field: "file",
        },
      },
    });

    const loadHandler = xhr.addEventListener.mock.calls.find(
      (call: [string, unknown]) => call[0] === "load"
    )?.[1] as () => void;
    loadHandler();

    await waitFor(() => {
      expect(screen.getByText("File content is corrupted")).toBeInTheDocument();
    });
  });

  it("shows error for network failure", async () => {
    render(<DocumentUpload {...defaultProps} />);

    const input = screen.getByTestId("file-input");
    const validFile = new File(["test content"], "resume.pdf", { type: "application/pdf" });

    fireEvent.change(input, { target: { files: [validFile] } });

    await waitFor(() => {
      const xhr = MockXHR.getLastInstance();
      expect(xhr).toBeDefined();
      expect(xhr.addEventListener).toHaveBeenCalled();
    });

    const xhr = MockXHR.getLastInstance();
    const errorHandler = xhr.addEventListener.mock.calls.find(
      (call: [string, unknown]) => call[0] === "error"
    )?.[1] as () => void;
    expect(errorHandler).toBeDefined();
    errorHandler();

    await waitFor(() => {
      expect(screen.getByText(/upload failed/i)).toBeInTheDocument();
    });
  });

  it("resets to idle when try again is clicked after error", async () => {
    render(<DocumentUpload {...defaultProps} />);

    const input = screen.getByTestId("file-input");
    const invalidFile = new File(["test"], "test.exe", { type: "application/x-executable" });

    fireEvent.change(input, { target: { files: [invalidFile] } });

    await waitFor(() => {
      expect(screen.getByText(/invalid file type/i)).toBeInTheDocument();
    });

    const tryAgainButton = screen.getByText("Try again");
    fireEvent.click(tryAgainButton);

    await waitFor(() => {
      expect(screen.queryByText(/invalid file type/i)).not.toBeInTheDocument();
      expect(screen.getByText(/drag and drop/i)).toBeInTheDocument();
    });
  });

  it("disables file input during upload", async () => {
    render(<DocumentUpload {...defaultProps} />);

    const input = screen.getByTestId("file-input");
    const validFile = new File(["test content"], "resume.pdf", { type: "application/pdf" });

    fireEvent.change(input, { target: { files: [validFile] } });

    await waitFor(() => {
      expect(screen.getByTestId("file-input")).toBeDisabled();
    });
  });

  it("calls onError callback when error occurs", async () => {
    const onError = vi.fn();
    render(<DocumentUpload {...defaultProps} onError={onError} />);

    const input = screen.getByTestId("file-input");
    const invalidFile = new File(["test"], "test.exe", { type: "application/x-executable" });

    fireEvent.change(input, { target: { files: [invalidFile] } });

    await waitFor(() => {
      expect(onError).toHaveBeenCalledWith(expect.stringContaining("Invalid file type"));
    });
  });
});
