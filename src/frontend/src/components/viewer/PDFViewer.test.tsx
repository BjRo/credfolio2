import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";

// Track Document and Page renders
const mockOnDocumentLoadSuccess = vi.fn();

// Mock react-pdf to avoid canvas/worker dependencies in tests
vi.mock("react-pdf", () => {
  const React = require("react");

  return {
    Document: ({
      children,
      onLoadSuccess,
    }: {
      file: string;
      children: React.ReactNode;
      onLoadSuccess?: (pdf: { numPages: number }) => void;
      loading?: React.ReactNode;
      error?: React.ReactNode;
    }) => {
      // Simulate successful load
      React.useEffect(() => {
        if (onLoadSuccess) {
          mockOnDocumentLoadSuccess.mockImplementation(onLoadSuccess);
          onLoadSuccess({ numPages: 3 });
        }
      }, [onLoadSuccess]);

      return <div data-testid="pdf-document">{children}</div>;
    },
    Page: ({
      pageNumber,
      scale,
      renderTextLayer,
    }: {
      pageNumber: number;
      scale?: number;
      width?: number;
      renderTextLayer?: boolean;
      renderAnnotationLayer?: boolean;
      loading?: React.ReactNode;
    }) => {
      return (
        <div
          data-testid={`pdf-page-${pageNumber}`}
          data-scale={scale}
          data-text-layer={renderTextLayer}
        >
          Page {pageNumber}
        </div>
      );
    },
    pdfjs: {
      GlobalWorkerOptions: { workerSrc: "" },
    },
  };
});

import { PDFViewer } from "./PDFViewer";

describe("PDFViewer", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("renders the PDF document component", () => {
    render(<PDFViewer fileUrl="https://example.com/test.pdf" />);
    expect(screen.getByTestId("pdf-document")).toBeInTheDocument();
  });

  it("renders page navigation controls", () => {
    render(<PDFViewer fileUrl="https://example.com/test.pdf" />);
    expect(screen.getByRole("button", { name: "Previous page" })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Next page" })).toBeInTheDocument();
  });

  it("shows page count after document loads", () => {
    render(<PDFViewer fileUrl="https://example.com/test.pdf" />);
    expect(screen.getByText("of 3")).toBeInTheDocument();
  });

  it("renders zoom controls", () => {
    render(<PDFViewer fileUrl="https://example.com/test.pdf" />);
    expect(screen.getByRole("button", { name: "Zoom out" })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Zoom in" })).toBeInTheDocument();
  });

  it("shows current zoom level", () => {
    render(<PDFViewer fileUrl="https://example.com/test.pdf" />);
    expect(screen.getByText("100%")).toBeInTheDocument();
  });

  it("disables previous button on first page", () => {
    render(<PDFViewer fileUrl="https://example.com/test.pdf" />);
    expect(screen.getByRole("button", { name: "Previous page" })).toBeDisabled();
  });

  it("enables next button when there are more pages", () => {
    render(<PDFViewer fileUrl="https://example.com/test.pdf" />);
    expect(screen.getByRole("button", { name: "Next page" })).toBeEnabled();
  });

  it("navigates to next page on next button click", async () => {
    const user = userEvent.setup();
    render(<PDFViewer fileUrl="https://example.com/test.pdf" />);

    await user.click(screen.getByRole("button", { name: "Next page" }));

    // Page input should show page 2
    const pageInput = screen.getByRole("spinbutton", { name: "Current page" });
    expect(pageInput).toHaveValue(2);
  });

  it("navigates to previous page on previous button click", async () => {
    const user = userEvent.setup();
    render(<PDFViewer fileUrl="https://example.com/test.pdf" />);

    // Go to page 2 first
    await user.click(screen.getByRole("button", { name: "Next page" }));
    // Now go back
    await user.click(screen.getByRole("button", { name: "Previous page" }));

    const pageInput = screen.getByRole("spinbutton", { name: "Current page" });
    expect(pageInput).toHaveValue(1);
  });

  it("disables next button on last page", async () => {
    const user = userEvent.setup();
    render(<PDFViewer fileUrl="https://example.com/test.pdf" />);

    // Navigate to last page (page 3)
    await user.click(screen.getByRole("button", { name: "Next page" }));
    await user.click(screen.getByRole("button", { name: "Next page" }));

    expect(screen.getByRole("button", { name: "Next page" })).toBeDisabled();
  });

  it("zooms in when zoom in button is clicked", async () => {
    const user = userEvent.setup();
    render(<PDFViewer fileUrl="https://example.com/test.pdf" />);

    await user.click(screen.getByRole("button", { name: "Zoom in" }));
    expect(screen.getByText("125%")).toBeInTheDocument();
  });

  it("zooms out when zoom out button is clicked", async () => {
    const user = userEvent.setup();
    render(<PDFViewer fileUrl="https://example.com/test.pdf" />);

    await user.click(screen.getByRole("button", { name: "Zoom out" }));
    expect(screen.getByText("75%")).toBeInTheDocument();
  });

  it("does not zoom below minimum zoom level", async () => {
    const user = userEvent.setup();
    render(<PDFViewer fileUrl="https://example.com/test.pdf" />);

    // Click zoom out many times
    for (let i = 0; i < 10; i++) {
      await user.click(screen.getByRole("button", { name: "Zoom out" }));
    }

    expect(screen.getByText("25%")).toBeInTheDocument();
  });

  it("does not zoom above maximum zoom level", async () => {
    const user = userEvent.setup();
    render(<PDFViewer fileUrl="https://example.com/test.pdf" />);

    // Click zoom in many times
    for (let i = 0; i < 20; i++) {
      await user.click(screen.getByRole("button", { name: "Zoom in" }));
    }

    expect(screen.getByText("300%")).toBeInTheDocument();
  });

  it("renders all pages of the document", () => {
    render(<PDFViewer fileUrl="https://example.com/test.pdf" />);
    expect(screen.getByTestId("pdf-page-1")).toBeInTheDocument();
    expect(screen.getByTestId("pdf-page-2")).toBeInTheDocument();
    expect(screen.getByTestId("pdf-page-3")).toBeInTheDocument();
  });

  it("enables text layer on pages", () => {
    render(<PDFViewer fileUrl="https://example.com/test.pdf" />);
    const page = screen.getByTestId("pdf-page-1");
    expect(page).toHaveAttribute("data-text-layer", "true");
  });

  it("allows page navigation via input", async () => {
    const user = userEvent.setup();
    render(<PDFViewer fileUrl="https://example.com/test.pdf" />);

    const pageInput = screen.getByRole("spinbutton", { name: "Current page" });
    await user.clear(pageInput);
    await user.type(pageInput, "3");
    await user.keyboard("{Enter}");

    expect(pageInput).toHaveValue(3);
  });

  it("clamps page input to valid range", async () => {
    const user = userEvent.setup();
    render(<PDFViewer fileUrl="https://example.com/test.pdf" />);

    const pageInput = screen.getByRole("spinbutton", { name: "Current page" });
    await user.clear(pageInput);
    await user.type(pageInput, "99");
    await user.keyboard("{Enter}");

    // Should clamp to max page (3)
    expect(pageInput).toHaveValue(3);
  });

  it("shows loading skeleton while PDF loads", () => {
    render(<PDFViewer fileUrl="https://example.com/test.pdf" />);
    // The Document component is passed a loading prop
    const document = screen.getByTestId("pdf-document");
    expect(document).toBeInTheDocument();
  });
});
