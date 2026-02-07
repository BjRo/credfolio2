import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import type { TextContent } from "react-pdf";
import type { PageMatchResult, TextItem } from "@/lib/textSearch";
import { findMatchInPage, renderHighlightedText } from "@/lib/textSearch";

interface UseTextHighlightOptions {
  highlightText?: string;
  numPages: number;
  onHighlightResult?: (found: boolean) => void;
  pageRefs: React.RefObject<Map<number, HTMLDivElement>>;
}

interface CustomTextRendererParams {
  str: string;
  itemIndex: number;
}

interface UseTextHighlightReturn {
  getOnGetTextSuccess: (pageNumber: number) => (textContent: TextContent) => void;
  getCustomTextRenderer: (
    pageNumber: number
  ) => ((params: CustomTextRendererParams) => string) | undefined;
  getOnRenderTextLayerSuccess: (pageNumber: number) => () => void;
}

function isTextItem(item: unknown): item is { str: string } {
  return typeof item === "object" && item !== null && "str" in item;
}

export function useTextHighlight({
  highlightText,
  numPages,
  onHighlightResult,
  pageRefs,
}: UseTextHighlightOptions): UseTextHighlightReturn {
  // Match data stored in ref for synchronous access by customTextRenderer
  // Map<pageNumber, Map<itemIndex, PageMatchResult>>
  const matchDataRef = useRef(new Map<number, Map<number, PageMatchResult>>());
  const [renderKey, setRenderKey] = useState(0);
  const firstMatchPageRef = useRef<number | null>(null);
  const hasScrolledRef = useRef(false);
  const hasReportedRef = useRef(false);
  const searchedPagesRef = useRef(new Set<number>());

  // Reset when highlightText changes — highlightText is intentionally in deps
  // to trigger reset on search text change (refs are the side effects we want)
  // biome-ignore lint/correctness/useExhaustiveDependencies: highlightText drives the reset
  useEffect(() => {
    matchDataRef.current.clear();
    firstMatchPageRef.current = null;
    hasScrolledRef.current = false;
    hasReportedRef.current = false;
    searchedPagesRef.current.clear();
    setRenderKey((k) => k + 1);
  }, [highlightText]);

  const getOnGetTextSuccess = useCallback(
    (pageNumber: number) => (textContent: TextContent) => {
      if (!highlightText) return;

      searchedPagesRef.current.add(pageNumber);

      // Build TextItem array, preserving original indices (including TextMarkedContent gaps)
      const items: TextItem[] = [];
      const indexMap: number[] = []; // indexMap[i] = original index in textContent.items

      for (let i = 0; i < textContent.items.length; i++) {
        const item = textContent.items[i];
        if (isTextItem(item)) {
          indexMap.push(i);
          items.push({ str: item.str });
        }
      }

      const result = findMatchInPage(items, highlightText);

      if (result) {
        // Map itemIndex from filtered array back to original array index
        const pageMatchMap = new Map<number, PageMatchResult>();
        for (const range of result.ranges) {
          const originalIndex = indexMap[range.itemIndex];
          const existing = pageMatchMap.get(originalIndex);
          if (existing) {
            existing.ranges.push({ ...range, itemIndex: originalIndex });
          } else {
            pageMatchMap.set(originalIndex, {
              ranges: [{ ...range, itemIndex: originalIndex }],
            });
          }
        }
        matchDataRef.current.set(pageNumber, pageMatchMap);

        if (firstMatchPageRef.current === null) {
          firstMatchPageRef.current = pageNumber;
        }

        // Report found
        if (!hasReportedRef.current) {
          hasReportedRef.current = true;
          onHighlightResult?.(true);
        }
      }

      // Check if all pages searched with no match
      if (!hasReportedRef.current && searchedPagesRef.current.size >= numPages) {
        hasReportedRef.current = true;
        onHighlightResult?.(false);
      }

      // Trigger re-render of text layer
      setRenderKey((k) => k + 1);
    },
    [highlightText, numPages, onHighlightResult]
  );

  const getCustomTextRenderer = useMemo(() => {
    if (!highlightText) return () => undefined;

    // Including renderKey in deps (and referencing it here) forces useMemo to return
    // a new function identity after setRenderKey increments, which triggers react-pdf
    // to re-invoke customTextRenderer with the updated match data from refs.
    void renderKey;

    return (pageNumber: number) => {
      return ({ str, itemIndex }: CustomTextRendererParams): string => {
        // Only look at match data for this specific page to avoid cross-page false positives
        const pageMatchMap = matchDataRef.current.get(pageNumber);
        if (pageMatchMap) {
          const matchResult = pageMatchMap.get(itemIndex);
          if (matchResult) {
            return renderHighlightedText(str, matchResult.ranges);
          }
        }
        return str;
      };
    };
  }, [highlightText, renderKey]);

  const getOnRenderTextLayerSuccess = useCallback(
    (pageNumber: number) => () => {
      if (!highlightText || hasScrolledRef.current) return;
      if (firstMatchPageRef.current !== pageNumber) return;

      const pageEl = pageRefs.current?.get(pageNumber);
      if (!pageEl) return;

      const highlightEl = pageEl.querySelector(".pdf-highlight");
      if (highlightEl) {
        hasScrolledRef.current = true;
        highlightEl.scrollIntoView({ behavior: "smooth", block: "center" });
      }
    },
    [highlightText, pageRefs]
  );

  return {
    getOnGetTextSuccess,
    getCustomTextRenderer,
    getOnRenderTextLayerSuccess,
  };
}
