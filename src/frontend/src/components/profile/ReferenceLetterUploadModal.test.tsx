import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { ReferenceLetterUploadModal } from "./ReferenceLetterUploadModal";

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

describe("ReferenceLetterUploadModal", () => {
  const defaultProps = {
    open: true,
    onOpenChange: vi.fn(),
    userId: "test-user-id",
  };

  beforeEach(() => {
    vi.clearAllMocks();
    MockXHR.reset();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it("renders modal when open is true", () => {
    render(<ReferenceLetterUploadModal {...defaultProps} />);
    expect(screen.getByRole("dialog")).toBeInTheDocument();
    expect(screen.getByText("Add Reference Letter")).toBeInTheDocument();
  });

  it("does not render when open is false", () => {
    render(<ReferenceLetterUploadModal {...defaultProps} open={false} />);
    expect(screen.queryByRole("dialog")).not.toBeInTheDocument();
  });

  it("shows drag and drop zone with file type hints", () => {
    render(<ReferenceLetterUploadModal {...defaultProps} />);
    expect(screen.getByText(/drag and drop/i)).toBeInTheDocument();
    expect(screen.getByText(/pdf, docx, txt/i)).toBeInTheDocument();
  });

  it("shows file input for selecting files", () => {
    render(<ReferenceLetterUploadModal {...defaultProps} />);
    const input = screen.getByTestId("file-input");
    expect(input).toBeInTheDocument();
    expect(input).toHaveAttribute("type", "file");
  });

  it("calls onOpenChange when cancel button is clicked", () => {
    const onOpenChange = vi.fn();
    render(<ReferenceLetterUploadModal {...defaultProps} onOpenChange={onOpenChange} />);

    const cancelButton = screen.getByRole("button", { name: /cancel/i });
    fireEvent.click(cancelButton);

    expect(onOpenChange).toHaveBeenCalledWith(false);
  });

  it("shows error for invalid file type", async () => {
    render(<ReferenceLetterUploadModal {...defaultProps} />);

    const input = screen.getByTestId("file-input");
    const invalidFile = new File(["test"], "test.exe", { type: "application/x-executable" });

    fireEvent.change(input, { target: { files: [invalidFile] } });

    await waitFor(() => {
      expect(screen.getByText(/invalid file type/i)).toBeInTheDocument();
    });
  });

  it("shows error for file exceeding size limit", async () => {
    render(<ReferenceLetterUploadModal {...defaultProps} />);

    const input = screen.getByTestId("file-input");
    // Create a file larger than 10MB
    const largeFile = new File(["x"], "large.pdf", {
      type: "application/pdf",
    });

    Object.defineProperty(largeFile, "size", { value: 11 * 1024 * 1024 });

    fireEvent.change(input, { target: { files: [largeFile] } });

    await waitFor(() => {
      expect(screen.getByText(/file too large/i)).toBeInTheDocument();
    });
  });

  it("shows uploading state with progress when valid file selected", async () => {
    render(<ReferenceLetterUploadModal {...defaultProps} />);

    const input = screen.getByTestId("file-input");
    const validFile = new File(["test content"], "reference.pdf", {
      type: "application/pdf",
    });

    fireEvent.change(input, { target: { files: [validFile] } });

    // The uploading state should be shown immediately
    await waitFor(() => {
      expect(screen.getByText(/uploading/i)).toBeInTheDocument();
    });
  });

  it("initiates upload and shows progress for valid file", async () => {
    render(<ReferenceLetterUploadModal {...defaultProps} />);

    const input = screen.getByTestId("file-input");
    const validFile = new File(["test content"], "reference.pdf", {
      type: "application/pdf",
    });

    fireEvent.change(input, { target: { files: [validFile] } });

    // Should show uploading state
    await waitFor(() => {
      expect(screen.getByText(/uploading/i)).toBeInTheDocument();
    });

    // XHR should be initialized
    const xhr = MockXHR.getLastInstance();
    expect(xhr).toBeDefined();
    expect(xhr.open).toHaveBeenCalledWith("POST", expect.any(String));
    expect(xhr.send).toHaveBeenCalled();
  });

  it("shows error and retry button when network error occurs", async () => {
    render(<ReferenceLetterUploadModal {...defaultProps} />);

    const input = screen.getByTestId("file-input");
    const validFile = new File(["test content"], "reference.pdf", {
      type: "application/pdf",
    });

    fireEvent.change(input, { target: { files: [validFile] } });

    // Wait for XHR to be set up
    await waitFor(() => {
      const xhr = MockXHR.getLastInstance();
      expect(xhr).toBeDefined();
      expect(xhr.addEventListener).toHaveBeenCalled();
    });

    // Simulate network error
    const xhr = MockXHR.getLastInstance();
    const errorCall = xhr.addEventListener.mock.calls.find(
      (call: [string, unknown]) => call[0] === "error"
    );
    expect(errorCall).toBeDefined();
    const errorHandler = errorCall[1] as () => void;
    errorHandler();

    // Should show error state with retry button
    await waitFor(() => {
      expect(screen.getByText(/upload failed/i)).toBeInTheDocument();
      expect(screen.getByText("Try again")).toBeInTheDocument();
    });
  });

  it("has retry button when validation error occurs", async () => {
    render(<ReferenceLetterUploadModal {...defaultProps} />);

    const input = screen.getByTestId("file-input");
    const invalidFile = new File(["test"], "test.exe", { type: "application/x-executable" });

    fireEvent.change(input, { target: { files: [invalidFile] } });

    // Wait for error message to appear
    await waitFor(() => {
      expect(screen.getByText(/invalid file type/i)).toBeInTheDocument();
    });

    // Use findByText since the button text isn't inside a semantic button role for screen readers
    expect(screen.getByText("Try again")).toBeInTheDocument();
  });

  it("resets error state when try again is clicked", async () => {
    render(<ReferenceLetterUploadModal {...defaultProps} />);

    const input = screen.getByTestId("file-input");
    const invalidFile = new File(["test"], "test.exe", { type: "application/x-executable" });

    fireEvent.change(input, { target: { files: [invalidFile] } });

    await waitFor(() => {
      expect(screen.getByText(/invalid file type/i)).toBeInTheDocument();
    });

    // Find and click the try again button by text
    const tryAgainButton = screen.getByText("Try again");
    fireEvent.click(tryAgainButton);

    // Error should be cleared
    await waitFor(() => {
      expect(screen.queryByText(/invalid file type/i)).not.toBeInTheDocument();
    });
  });

  it("accepts file types: pdf, docx, txt", () => {
    render(<ReferenceLetterUploadModal {...defaultProps} />);
    const input = screen.getByTestId("file-input");
    expect(input).toHaveAttribute(
      "accept",
      "application/pdf,application/vnd.openxmlformats-officedocument.wordprocessingml.document,text/plain"
    );
  });

  it("disables cancel button while processing", async () => {
    render(<ReferenceLetterUploadModal {...defaultProps} />);

    const input = screen.getByTestId("file-input");
    const validFile = new File(["test content"], "reference.pdf", {
      type: "application/pdf",
    });

    fireEvent.change(input, { target: { files: [validFile] } });

    // Wait for XHR to be set up
    await waitFor(() => {
      const xhr = MockXHR.getLastInstance();
      expect(xhr).toBeDefined();
    });

    // Simulate successful upload
    const xhr = MockXHR.getLastInstance();
    xhr.status = 200;
    xhr.responseText = JSON.stringify({
      data: {
        uploadFile: {
          __typename: "UploadFileResult",
          file: { id: "file-id", filename: "reference.pdf" },
          referenceLetter: { id: "ref-letter-id", status: "PENDING" },
        },
      },
    });

    // Trigger load event
    const loadHandler = xhr.addEventListener.mock.calls.find(
      (call: [string, unknown]) => call[0] === "load"
    )?.[1] as () => void;
    expect(loadHandler).toBeDefined();
    loadHandler();

    // Wait for processing state
    await waitFor(() => {
      expect(screen.getByText(/processing reference letter/i)).toBeInTheDocument();
    });

    // Cancel button should be disabled during processing
    const cancelButton = screen.getByRole("button", { name: /cancel/i });
    expect(cancelButton).toBeDisabled();
  });
});
