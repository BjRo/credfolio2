import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";

// Track Document and Page renders
const mockOnDocumentLoadSuccess = vi.fn();

// Store page callbacks so tests can trigger them
let capturedPageProps: Map<
  number,
  {
    onGetTextSuccess?: (textContent: unknown) => void;
    customTextRenderer?: (params: { str: string; itemIndex: number }) => string;
    onRenderTextLayerSuccess?: () => void;
  }
>;

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
      onGetTextSuccess,
      customTextRenderer,
      onRenderTextLayerSuccess,
    }: {
      pageNumber: number;
      scale?: number;
      width?: number;
      renderTextLayer?: boolean;
      renderAnnotationLayer?: boolean;
      loading?: React.ReactNode;
      onGetTextSuccess?: (textContent: unknown) => void;
      customTextRenderer?: (params: { str: string; itemIndex: number }) => string;
      onRenderTextLayerSuccess?: () => void;
    }) => {
      // Store callbacks for test access
      const React = require("react");
      React.useEffect(() => {
        if (capturedPageProps) {
          capturedPageProps.set(pageNumber, {
            onGetTextSuccess,
            customTextRenderer,
            onRenderTextLayerSuccess,
          });
        }
      });

      return (
        <div
          data-testid={`pdf-page-${pageNumber}`}
          data-scale={scale}
          data-text-layer={renderTextLayer}
          data-has-custom-renderer={!!customTextRenderer}
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
    capturedPageProps = new Map();
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

  it("renders the document container", () => {
    render(<PDFViewer fileUrl="https://example.com/test.pdf" />);
    const document = screen.getByTestId("pdf-document");
    expect(document).toBeInTheDocument();
  });
});

describe("PDFViewer highlight integration", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    capturedPageProps = new Map();
  });

  it("does not set customTextRenderer when no highlightText", () => {
    render(<PDFViewer fileUrl="https://example.com/test.pdf" />);

    const page1 = screen.getByTestId("pdf-page-1");
    expect(page1).toHaveAttribute("data-has-custom-renderer", "false");
  });

  it("sets customTextRenderer when highlightText is provided", () => {
    render(<PDFViewer fileUrl="https://example.com/test.pdf" highlightText="hello" />);

    const page1 = screen.getByTestId("pdf-page-1");
    expect(page1).toHaveAttribute("data-has-custom-renderer", "true");
  });

  it("calls onHighlightResult callback when match is found", () => {
    const onHighlightResult = vi.fn();
    render(
      <PDFViewer
        fileUrl="https://example.com/test.pdf"
        highlightText="hello"
        onHighlightResult={onHighlightResult}
      />
    );

    // Simulate text extraction completing for page 1
    const page1Props = capturedPageProps.get(1);
    expect(page1Props?.onGetTextSuccess).toBeInstanceOf(Function);

    page1Props?.onGetTextSuccess?.({
      items: [
        {
          str: "hello world",
          dir: "ltr",
          width: 100,
          height: 12,
          transform: [1, 0, 0, 1, 0, 0],
          fontName: "Arial",
          hasEOL: false,
        },
      ],
      styles: {},
    });

    expect(onHighlightResult).toHaveBeenCalledWith(true);
  });

  it("passes onGetTextSuccess to each page", () => {
    render(<PDFViewer fileUrl="https://example.com/test.pdf" highlightText="test" />);

    expect(capturedPageProps.get(1)?.onGetTextSuccess).toBeInstanceOf(Function);
    expect(capturedPageProps.get(2)?.onGetTextSuccess).toBeInstanceOf(Function);
    expect(capturedPageProps.get(3)?.onGetTextSuccess).toBeInstanceOf(Function);
  });

  it("passes onRenderTextLayerSuccess to each page", () => {
    render(<PDFViewer fileUrl="https://example.com/test.pdf" highlightText="test" />);

    expect(capturedPageProps.get(1)?.onRenderTextLayerSuccess).toBeInstanceOf(Function);
    expect(capturedPageProps.get(2)?.onRenderTextLayerSuccess).toBeInstanceOf(Function);
    expect(capturedPageProps.get(3)?.onRenderTextLayerSuccess).toBeInstanceOf(Function);
  });
});
