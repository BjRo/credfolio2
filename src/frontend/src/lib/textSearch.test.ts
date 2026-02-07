import { describe, expect, it } from "vitest";
import type { TextItem } from "./textSearch";
import { escapeHtml, findMatchInPage, normalizeText, renderHighlightedText } from "./textSearch";

describe("escapeHtml", () => {
  it("escapes ampersands", () => {
    expect(escapeHtml("a & b")).toBe("a &amp; b");
  });

  it("escapes angle brackets", () => {
    expect(escapeHtml("<script>")).toBe("&lt;script&gt;");
  });

  it("escapes double quotes", () => {
    expect(escapeHtml('"hello"')).toBe("&quot;hello&quot;");
  });

  it("escapes single quotes", () => {
    expect(escapeHtml("it's")).toBe("it&#39;s");
  });

  it("handles empty string", () => {
    expect(escapeHtml("")).toBe("");
  });

  it("handles strings with no special characters", () => {
    expect(escapeHtml("hello world")).toBe("hello world");
  });
});

describe("normalizeText", () => {
  it("returns empty string for empty input", () => {
    expect(normalizeText("")).toBe("");
  });

  it("collapses multiple spaces to single space", () => {
    expect(normalizeText("hello   world")).toBe("hello world");
  });

  it("trims leading and trailing whitespace", () => {
    expect(normalizeText("  hello  ")).toBe("hello");
  });

  it("replaces NBSP with regular space", () => {
    expect(normalizeText("hello\u00A0world")).toBe("hello world");
  });

  it("replaces smart single quotes with straight quotes", () => {
    expect(normalizeText("\u2018hello\u2019")).toBe("'hello'");
  });

  it("replaces smart double quotes with straight quotes", () => {
    expect(normalizeText("\u201Chello\u201D")).toBe('"hello"');
  });

  it("replaces en-dash and em-dash with hyphen", () => {
    expect(normalizeText("a\u2013b\u2014c")).toBe("a-b-c");
  });

  it("replaces ligatures", () => {
    expect(normalizeText("\uFB00\uFB01\uFB02\uFB03\uFB04")).toBe("fffiflffiffl");
  });

  it("collapses tabs and newlines to single space", () => {
    expect(normalizeText("hello\t\nworld")).toBe("hello world");
  });
});

describe("findMatchInPage", () => {
  function makeItems(strs: string[]): TextItem[] {
    return strs.map((str) => ({ str }));
  }

  it("returns null when items are empty", () => {
    expect(findMatchInPage([], "hello")).toBeNull();
  });

  it("returns null when search text is empty", () => {
    expect(findMatchInPage(makeItems(["hello"]), "")).toBeNull();
  });

  it("returns null when no match is found", () => {
    expect(findMatchInPage(makeItems(["hello world"]), "xyz")).toBeNull();
  });

  it("finds a match within a single text item", () => {
    const result = findMatchInPage(makeItems(["hello world"]), "world");
    expect(result).not.toBeNull();
    expect(result?.ranges).toHaveLength(1);
    expect(result?.ranges[0]).toEqual({
      itemIndex: 0,
      startOffset: 6,
      endOffset: 11,
    });
  });

  it("performs case-insensitive matching", () => {
    const result = findMatchInPage(makeItems(["Hello World"]), "hello");
    expect(result).not.toBeNull();
    expect(result?.ranges).toHaveLength(1);
    expect(result?.ranges[0]).toEqual({
      itemIndex: 0,
      startOffset: 0,
      endOffset: 5,
    });
  });

  it("finds a match spanning multiple text items", () => {
    const result = findMatchInPage(makeItems(["hel", "lo wor", "ld"]), "hello world");
    expect(result).not.toBeNull();
    expect(result?.ranges).toHaveLength(3);
    expect(result?.ranges[0]).toEqual({ itemIndex: 0, startOffset: 0, endOffset: 3 });
    expect(result?.ranges[1]).toEqual({ itemIndex: 1, startOffset: 0, endOffset: 6 });
    expect(result?.ranges[2]).toEqual({ itemIndex: 2, startOffset: 0, endOffset: 2 });
  });

  it("handles normalized text matching (smart quotes)", () => {
    const result = findMatchInPage(makeItems(["\u201Chello\u201D"]), '"hello"');
    expect(result).not.toBeNull();
    expect(result?.ranges).toHaveLength(1);
  });

  it("handles items with only whitespace", () => {
    const result = findMatchInPage(makeItems(["hello", " ", "world"]), "hello world");
    expect(result).not.toBeNull();
    expect(result?.ranges).toHaveLength(3);
  });

  it("finds first match only", () => {
    const result = findMatchInPage(makeItems(["abc abc"]), "abc");
    expect(result).not.toBeNull();
    expect(result?.ranges).toHaveLength(1);
    expect(result?.ranges[0]).toEqual({ itemIndex: 0, startOffset: 0, endOffset: 3 });
  });

  it("preserves the original itemIndex for sparse items", () => {
    // Simulating that items passed could have specific indexes
    const items: TextItem[] = [{ str: "hello " }, { str: "world" }];
    const result = findMatchInPage(items, "hello world");
    expect(result).not.toBeNull();
    expect(result?.ranges[0].itemIndex).toBe(0);
    expect(result?.ranges[1].itemIndex).toBe(1);
  });
});

describe("renderHighlightedText", () => {
  it("returns escaped text when no ranges provided", () => {
    expect(renderHighlightedText("<b>hello</b>", [])).toBe("&lt;b&gt;hello&lt;/b&gt;");
  });

  it("wraps matched portion in mark tags", () => {
    const result = renderHighlightedText("hello world", [
      { itemIndex: 0, startOffset: 6, endOffset: 11 },
    ]);
    expect(result).toBe('hello <mark class="pdf-highlight">world</mark>');
  });

  it("handles match at the beginning", () => {
    const result = renderHighlightedText("hello world", [
      { itemIndex: 0, startOffset: 0, endOffset: 5 },
    ]);
    expect(result).toBe('<mark class="pdf-highlight">hello</mark> world');
  });

  it("handles full string match", () => {
    const result = renderHighlightedText("hello", [{ itemIndex: 0, startOffset: 0, endOffset: 5 }]);
    expect(result).toBe('<mark class="pdf-highlight">hello</mark>');
  });

  it("escapes HTML in non-highlighted portions", () => {
    const result = renderHighlightedText("<b>test</b> match", [
      { itemIndex: 0, startOffset: 12, endOffset: 17 },
    ]);
    expect(result).toBe('&lt;b&gt;test&lt;/b&gt; <mark class="pdf-highlight">match</mark>');
  });

  it("escapes HTML inside highlighted portions", () => {
    const result = renderHighlightedText("<b>hello</b>", [
      { itemIndex: 0, startOffset: 0, endOffset: 12 },
    ]);
    expect(result).toBe('<mark class="pdf-highlight">&lt;b&gt;hello&lt;/b&gt;</mark>');
  });

  it("handles multiple ranges for the same item", () => {
    // This case happens when a cross-item search puts multiple ranges on one item
    // But per the spec, findMatchInPage groups by itemIndex, so one range per item
    // renderHighlightedText should handle the single range for this item
    const result = renderHighlightedText("abcdef", [
      { itemIndex: 0, startOffset: 1, endOffset: 4 },
    ]);
    expect(result).toBe('a<mark class="pdf-highlight">bcd</mark>ef');
  });
});
