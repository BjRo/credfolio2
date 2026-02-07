"use client";

import { useCallback, useMemo, useRef, useState } from "react";
import { Document, Page, pdfjs } from "react-pdf";
import "react-pdf/dist/Page/TextLayer.css";
import "react-pdf/dist/Page/AnnotationLayer.css";
import { ChevronLeft, ChevronRight, Minus, Plus } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useTextHighlight } from "@/hooks/useTextHighlight";

pdfjs.GlobalWorkerOptions.workerSrc = "/pdf.worker.min.mjs";

const MIN_ZOOM = 0.25;
const MAX_ZOOM = 3.0;
const ZOOM_STEP = 0.25;

const MOBILE_BREAKPOINT = 768;
const MOBILE_DEFAULT_SCALE = 0.5;

interface PDFViewerProps {
  fileUrl: string;
  highlightText?: string;
  onHighlightResult?: (found: boolean) => void;
  toolbarLeft?: React.ReactNode;
}

function LoadingSkeleton() {
  return (
    <div className="flex justify-center py-8">
      <div className="w-[612px] h-[792px] animate-pulse bg-muted rounded" />
    </div>
  );
}

function ErrorDisplay({ error }: { error?: Error }) {
  return (
    <div className="flex justify-center py-8">
      <div className="text-center p-8">
        <p className="text-destructive font-medium">Failed to load PDF</p>
        <p className="text-sm text-muted-foreground mt-1">
          The document could not be loaded. Please try again.
        </p>
        {error && <p className="text-xs text-muted-foreground mt-2 font-mono">{error.message}</p>}
      </div>
    </div>
  );
}

function getInitialScale() {
  if (typeof window === "undefined") return 1.0;
  return window.innerWidth < MOBILE_BREAKPOINT ? MOBILE_DEFAULT_SCALE : 1.0;
}

export function PDFViewer({
  fileUrl,
  highlightText,
  onHighlightResult,
  toolbarLeft,
}: PDFViewerProps) {
  const [numPages, setNumPages] = useState<number>(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [scale, setScale] = useState(getInitialScale);
  const [pageInputValue, setPageInputValue] = useState("1");
  const [loadError, setLoadError] = useState<Error | null>(null);
  const pageRefs = useRef<Map<number, HTMLDivElement>>(new Map());

  const { getOnGetTextSuccess, getCustomTextRenderer, getOnRenderTextLayerSuccess } =
    useTextHighlight({
      highlightText,
      numPages,
      onHighlightResult,
      pageRefs,
    });

  const pageCallbacks = useMemo(() => {
    const callbacks = new Map<
      number,
      {
        onGetTextSuccess: ReturnType<typeof getOnGetTextSuccess>;
        onRenderTextLayerSuccess: ReturnType<typeof getOnRenderTextLayerSuccess>;
        customTextRenderer: ReturnType<typeof getCustomTextRenderer>;
      }
    >();
    for (let i = 1; i <= numPages; i++) {
      callbacks.set(i, {
        onGetTextSuccess: getOnGetTextSuccess(i),
        onRenderTextLayerSuccess: getOnRenderTextLayerSuccess(i),
        customTextRenderer: getCustomTextRenderer(i),
      });
    }
    return callbacks;
  }, [numPages, getOnGetTextSuccess, getOnRenderTextLayerSuccess, getCustomTextRenderer]);

  const onDocumentLoadSuccess = useCallback(({ numPages }: { numPages: number }) => {
    pageRefs.current.clear();
    setNumPages(numPages);
  }, []);

  const onDocumentLoadError = useCallback((error: Error) => {
    setLoadError(error);
  }, []);

  const goToPage = useCallback(
    (page: number) => {
      const clamped = Math.max(1, Math.min(page, numPages));
      setCurrentPage(clamped);
      setPageInputValue(String(clamped));

      const pageEl = pageRefs.current.get(clamped);
      if (pageEl) {
        pageEl.scrollIntoView({ behavior: "smooth", block: "start" });
      }
    },
    [numPages]
  );

  const handlePageInputChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    setPageInputValue(e.target.value);
  }, []);

  const handlePageInputSubmit = useCallback(
    (e: React.KeyboardEvent<HTMLInputElement> | React.FocusEvent<HTMLInputElement>) => {
      if ("key" in e && e.key !== "Enter") return;

      const parsed = parseInt(pageInputValue, 10);
      if (!Number.isNaN(parsed)) {
        goToPage(parsed);
      } else {
        setPageInputValue(String(currentPage));
      }
    },
    [pageInputValue, goToPage, currentPage]
  );

  const zoomIn = useCallback(() => {
    setScale((prev) => Math.min(prev + ZOOM_STEP, MAX_ZOOM));
  }, []);

  const zoomOut = useCallback(() => {
    setScale((prev) => Math.max(prev - ZOOM_STEP, MIN_ZOOM));
  }, []);

  const zoomPercent = Math.round(scale * 100);

  return (
    <div className="flex flex-col h-full">
      {/* Toolbar */}
      <div className="flex flex-wrap items-center border-b bg-background px-4 py-2 sticky top-0 z-10 gap-y-1">
        {/* Left: back button + title (passed from parent) */}
        {toolbarLeft && (
          <div className="flex items-center gap-3 flex-1 min-w-0 md:flex-none md:mr-4">
            {toolbarLeft}
          </div>
        )}

        {/* Right: page nav + zoom — pushed right on desktop */}
        <div className="flex items-center justify-between flex-1 gap-2">
          {/* Page navigation */}
          <div className="flex items-center gap-1">
            <Button
              variant="ghost"
              size="icon-sm"
              onClick={() => goToPage(currentPage - 1)}
              disabled={currentPage <= 1}
              aria-label="Previous page"
            >
              <ChevronLeft className="size-4" />
            </Button>

            <div className="flex items-center gap-1 text-sm">
              <input
                type="number"
                min={1}
                max={numPages || 1}
                value={pageInputValue}
                onChange={handlePageInputChange}
                onKeyDown={handlePageInputSubmit}
                onBlur={handlePageInputSubmit}
                aria-label="Current page"
                className="w-12 h-7 text-center text-sm border rounded bg-background tabular-nums"
              />
              <span className="text-muted-foreground">of {numPages || "..."}</span>
            </div>

            <Button
              variant="ghost"
              size="icon-sm"
              onClick={() => goToPage(currentPage + 1)}
              disabled={currentPage >= numPages}
              aria-label="Next page"
            >
              <ChevronRight className="size-4" />
            </Button>
          </div>

          {/* Zoom controls */}
          <div className="flex items-center gap-1">
            <Button
              variant="ghost"
              size="icon-sm"
              onClick={zoomOut}
              disabled={scale <= MIN_ZOOM}
              aria-label="Zoom out"
            >
              <Minus className="size-4" />
            </Button>

            <span className="text-sm text-muted-foreground w-12 text-center tabular-nums">
              {zoomPercent}%
            </span>

            <Button
              variant="ghost"
              size="icon-sm"
              onClick={zoomIn}
              disabled={scale >= MAX_ZOOM}
              aria-label="Zoom in"
            >
              <Plus className="size-4" />
            </Button>
          </div>
        </div>
      </div>

      {/* PDF Content */}
      <div className="flex-1 overflow-auto bg-neutral-700">
        <Document
          file={fileUrl}
          onLoadSuccess={onDocumentLoadSuccess}
          onLoadError={onDocumentLoadError}
          loading={<LoadingSkeleton />}
          error={<ErrorDisplay error={loadError ?? undefined} />}
        >
          <div className="flex flex-col items-center gap-4 py-4">
            {Array.from({ length: numPages }, (_, i) => i + 1).map((pageNumber) => (
              <div
                key={pageNumber}
                ref={(el) => {
                  if (el) {
                    pageRefs.current.set(pageNumber, el);
                  }
                }}
                className="shadow-md bg-white" // Intentionally bg-white: PDF pages are white documents even in dark mode
              >
                <Page
                  pageNumber={pageNumber}
                  scale={scale}
                  renderTextLayer={true}
                  renderAnnotationLayer={true}
                  loading={<div className="w-[612px] h-[792px] animate-pulse bg-muted rounded" />}
                  onGetTextSuccess={pageCallbacks.get(pageNumber)?.onGetTextSuccess}
                  customTextRenderer={pageCallbacks.get(pageNumber)?.customTextRenderer}
                  onRenderTextLayerSuccess={pageCallbacks.get(pageNumber)?.onRenderTextLayerSuccess}
                />
              </div>
            ))}
          </div>
        </Document>
      </div>
    </div>
  );
}
