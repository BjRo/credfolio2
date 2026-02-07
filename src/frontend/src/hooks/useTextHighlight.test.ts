import { act, renderHook } from "@testing-library/react";
import type { TextContent } from "react-pdf";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { useTextHighlight } from "./useTextHighlight";

function makeTextContent(strings: string[]): TextContent {
  return {
    items: strings.map((str) => ({
      str,
      dir: "ltr" as const,
      width: 100,
      height: 12,
      transform: [1, 0, 0, 1, 0, 0] as [number, number, number, number, number, number],
      fontName: "Arial",
      hasEOL: false,
    })),
    styles: {},
  };
}

describe("useTextHighlight", () => {
  const defaultProps = {
    highlightText: undefined as string | undefined,
    numPages: 3,
    onHighlightResult: vi.fn(),
    pageRefs: { current: new Map<number, HTMLDivElement>() },
  };

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("returns undefined customTextRenderer when no highlightText", () => {
    const { result } = renderHook(() => useTextHighlight(defaultProps));
    expect(result.current.customTextRenderer).toBeUndefined();
  });

  it("returns a customTextRenderer function when highlightText is provided", () => {
    const { result } = renderHook(() =>
      useTextHighlight({ ...defaultProps, highlightText: "hello" })
    );
    expect(result.current.customTextRenderer).toBeInstanceOf(Function);
  });

  it("returns getOnGetTextSuccess factory", () => {
    const { result } = renderHook(() =>
      useTextHighlight({ ...defaultProps, highlightText: "hello" })
    );
    expect(result.current.getOnGetTextSuccess).toBeInstanceOf(Function);
    expect(result.current.getOnGetTextSuccess(1)).toBeInstanceOf(Function);
  });

  it("returns getOnRenderTextLayerSuccess factory", () => {
    const { result } = renderHook(() =>
      useTextHighlight({ ...defaultProps, highlightText: "hello" })
    );
    expect(result.current.getOnRenderTextLayerSuccess).toBeInstanceOf(Function);
    expect(result.current.getOnRenderTextLayerSuccess(1)).toBeInstanceOf(Function);
  });

  it("customTextRenderer returns plain text before onGetTextSuccess fires", () => {
    const { result } = renderHook(() =>
      useTextHighlight({ ...defaultProps, highlightText: "hello" })
    );

    // biome-ignore lint/style/noNonNullAssertion: asserted above
    const renderer = result.current.customTextRenderer!;
    const output = renderer({ str: "hello world", itemIndex: 0 });
    expect(output).toBe("hello world");
  });

  it("customTextRenderer returns highlighted HTML after onGetTextSuccess", () => {
    const { result } = renderHook(() =>
      useTextHighlight({ ...defaultProps, highlightText: "world" })
    );

    const textContent = makeTextContent(["hello world"]);

    // Trigger onGetTextSuccess for page 1
    act(() => {
      result.current.getOnGetTextSuccess(1)(textContent);
    });

    // biome-ignore lint/style/noNonNullAssertion: asserted above
    const renderer = result.current.customTextRenderer!;
    const output = renderer({ str: "hello world", itemIndex: 0 });
    expect(output).toContain('<mark class="pdf-highlight">world</mark>');
  });

  it("calls onHighlightResult(true) when match is found", () => {
    const onHighlightResult = vi.fn();
    const { result } = renderHook(() =>
      useTextHighlight({
        ...defaultProps,
        highlightText: "world",
        onHighlightResult,
      })
    );

    const textContent = makeTextContent(["hello world"]);

    act(() => {
      result.current.getOnGetTextSuccess(1)(textContent);
    });

    expect(onHighlightResult).toHaveBeenCalledWith(true);
  });

  it("calls onHighlightResult(false) after all pages searched with no match", () => {
    const onHighlightResult = vi.fn();
    const { result } = renderHook(() =>
      useTextHighlight({
        ...defaultProps,
        highlightText: "xyz",
        numPages: 2,
        onHighlightResult,
      })
    );

    const textContent1 = makeTextContent(["hello"]);
    const textContent2 = makeTextContent(["world"]);

    act(() => {
      result.current.getOnGetTextSuccess(1)(textContent1);
    });
    act(() => {
      result.current.getOnGetTextSuccess(2)(textContent2);
    });

    expect(onHighlightResult).toHaveBeenCalledWith(false);
  });

  it("resets state when highlightText changes", () => {
    const props = { ...defaultProps, highlightText: "hello" };
    const { result, rerender } = renderHook((p) => useTextHighlight(p), {
      initialProps: props,
    });

    const textContent = makeTextContent(["hello world"]);

    act(() => {
      result.current.getOnGetTextSuccess(1)(textContent);
    });

    // Verify match exists
    // biome-ignore lint/style/noNonNullAssertion: asserted above
    const renderer1 = result.current.customTextRenderer!;
    expect(renderer1({ str: "hello world", itemIndex: 0 })).toContain("pdf-highlight");

    // Change search text — should reset
    rerender({ ...props, highlightText: "xyz" });

    // biome-ignore lint/style/noNonNullAssertion: asserted above
    const renderer2 = result.current.customTextRenderer!;
    // Before new onGetTextSuccess fires, no match data → plain text
    expect(renderer2({ str: "hello world", itemIndex: 0 })).toBe("hello world");
  });

  it("returns undefined customTextRenderer when highlightText is empty string", () => {
    const { result } = renderHook(() => useTextHighlight({ ...defaultProps, highlightText: "" }));
    expect(result.current.customTextRenderer).toBeUndefined();
  });

  it("does not call onHighlightResult(true) multiple times for multiple matching pages", () => {
    const onHighlightResult = vi.fn();
    const { result } = renderHook(() =>
      useTextHighlight({
        ...defaultProps,
        highlightText: "hello",
        numPages: 2,
        onHighlightResult,
      })
    );

    const textContent1 = makeTextContent(["hello"]);
    const textContent2 = makeTextContent(["hello again"]);

    act(() => {
      result.current.getOnGetTextSuccess(1)(textContent1);
    });
    act(() => {
      result.current.getOnGetTextSuccess(2)(textContent2);
    });

    // Should be called with true exactly once (first match)
    const trueCalls = onHighlightResult.mock.calls.filter(([val]: [boolean]) => val === true);
    expect(trueCalls).toHaveLength(1);
  });

  it("skips TextMarkedContent items when computing match indices", () => {
    const { result } = renderHook(() =>
      useTextHighlight({ ...defaultProps, highlightText: "world" })
    );

    // Simulate textContent with a TextMarkedContent item (no str property)
    const textContent: TextContent = {
      items: [
        {
          str: "hello ",
          dir: "ltr" as const,
          width: 100,
          height: 12,
          transform: [1, 0, 0, 1, 0, 0] as [number, number, number, number, number, number],
          fontName: "Arial",
          hasEOL: false,
        },
        // TextMarkedContent - has 'type' but no 'str'
        { id: 1, type: "beginMarkedContentProps" as const, tag: "Span" },
        {
          str: "world",
          dir: "ltr" as const,
          width: 100,
          height: 12,
          transform: [1, 0, 0, 1, 0, 0] as [number, number, number, number, number, number],
          fontName: "Arial",
          hasEOL: false,
        },
      ],
      styles: {},
    };

    act(() => {
      result.current.getOnGetTextSuccess(1)(textContent);
    });

    // itemIndex 2 is the "world" item (index in the full items array, not filtered)
    // biome-ignore lint/style/noNonNullAssertion: asserted above
    const renderer = result.current.customTextRenderer!;
    const output = renderer({ str: "world", itemIndex: 2 });
    expect(output).toContain('<mark class="pdf-highlight">world</mark>');

    // itemIndex 0 should not be highlighted (it's "hello ")
    const output0 = renderer({ str: "hello ", itemIndex: 0 });
    expect(output0).not.toContain("pdf-highlight");
  });
});
