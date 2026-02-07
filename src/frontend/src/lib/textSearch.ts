/**
 * Pure utility functions for PDF text search and highlighting.
 * No React dependencies — all logic is testable in isolation.
 */

export interface TextItem {
  str: string;
}

export interface HighlightRange {
  itemIndex: number;
  startOffset: number;
  endOffset: number;
}

export interface PageMatchResult {
  ranges: HighlightRange[];
}

// Maps normalized character index → original character index
interface NormalizeResult {
  normalized: string;
  mapping: number[]; // mapping[normalizedIdx] = originalCharIdx
}

const LIGATURE_MAP: Record<string, string> = {
  "\uFB00": "ff",
  "\uFB01": "fi",
  "\uFB02": "fl",
  "\uFB03": "ffi",
  "\uFB04": "ffl",
};

/** Escape HTML special characters for safe innerHTML rendering. */
export function escapeHtml(text: string): string {
  return text
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#39;");
}

/** Normalize text: collapse whitespace, trim, replace ligatures, smart quotes, dashes, NBSP. */
export function normalizeText(text: string): string {
  return normalizeWithMapping(text, true).normalized;
}

function normalizeWithMapping(text: string, trim = false): NormalizeResult {
  const normalized: string[] = [];
  const mapping: number[] = [];

  let lastWasSpace = trim; // when trimming, suppress leading space

  for (let i = 0; i < text.length; i++) {
    let ch = text[i];

    // Replace NBSP
    if (ch === "\u00A0") ch = " ";
    // Replace smart quotes
    else if (ch === "\u2018" || ch === "\u2019") ch = "'";
    else if (ch === "\u201C" || ch === "\u201D") ch = '"';
    // Replace dashes
    else if (ch === "\u2013" || ch === "\u2014") ch = "-";
    // Replace ligatures
    else if (LIGATURE_MAP[ch]) {
      const expansion = LIGATURE_MAP[ch];
      for (let j = 0; j < expansion.length; j++) {
        normalized.push(expansion[j]);
        mapping.push(i);
      }
      lastWasSpace = false;
      continue;
    }

    // Collapse whitespace
    if (/\s/.test(ch)) {
      if (!lastWasSpace) {
        normalized.push(" ");
        mapping.push(i);
        lastWasSpace = true;
      }
      continue;
    }

    normalized.push(ch);
    mapping.push(i);
    lastWasSpace = false;
  }

  // Trim trailing space
  if (trim && normalized.length > 0 && normalized[normalized.length - 1] === " ") {
    normalized.pop();
    mapping.pop();
  }

  return { normalized: normalized.join(""), mapping };
}

interface GlobalMappingEntry {
  itemIndex: number;
  originalCharIdx: number;
}

/**
 * Find the first occurrence of searchText across concatenated text items.
 * Returns per-item HighlightRanges mapped back to original character positions.
 */
export function findMatchInPage(items: TextItem[], searchText: string): PageMatchResult | null {
  if (items.length === 0 || !searchText) return null;

  const normalizedSearch = normalizeText(searchText);
  if (!normalizedSearch) return null;

  // Build concatenated normalized page text with global mapping
  const pageTextParts: string[] = [];
  const globalMapping: GlobalMappingEntry[] = [];

  for (let itemIdx = 0; itemIdx < items.length; itemIdx++) {
    const { normalized, mapping } = normalizeWithMapping(items[itemIdx].str);

    for (let j = 0; j < normalized.length; j++) {
      pageTextParts.push(normalized[j]);
      globalMapping.push({
        itemIndex: itemIdx,
        originalCharIdx: mapping[j],
      });
    }
  }

  const pageText = pageTextParts.join("");
  const matchStart = pageText.toLowerCase().indexOf(normalizedSearch.toLowerCase());
  if (matchStart === -1) return null;

  const matchEnd = matchStart + normalizedSearch.length;

  // Group matched global indices by itemIndex
  const rangeMap = new Map<number, { minOriginal: number; maxOriginal: number }>();

  for (let gi = matchStart; gi < matchEnd; gi++) {
    const entry = globalMapping[gi];
    const existing = rangeMap.get(entry.itemIndex);
    if (existing) {
      existing.minOriginal = Math.min(existing.minOriginal, entry.originalCharIdx);
      existing.maxOriginal = Math.max(existing.maxOriginal, entry.originalCharIdx);
    } else {
      rangeMap.set(entry.itemIndex, {
        minOriginal: entry.originalCharIdx,
        maxOriginal: entry.originalCharIdx,
      });
    }
  }

  const ranges: HighlightRange[] = [];
  // Sort by itemIndex to maintain order
  const sortedEntries = [...rangeMap.entries()].sort(([a], [b]) => a - b);

  for (const [itemIndex, { minOriginal, maxOriginal }] of sortedEntries) {
    ranges.push({
      itemIndex,
      startOffset: minOriginal,
      endOffset: maxOriginal + 1,
    });
  }

  return { ranges };
}

/**
 * Produce highlighted HTML for a single text item.
 * Wraps matched portions in <mark class="pdf-highlight">.
 * All text is HTML-escaped for XSS safety.
 */
export function renderHighlightedText(str: string, ranges: HighlightRange[]): string {
  if (ranges.length === 0) {
    return escapeHtml(str);
  }

  // Use the first range for this item (there should only be one per item from findMatchInPage)
  const range = ranges[0];
  const before = str.slice(0, range.startOffset);
  const match = str.slice(range.startOffset, range.endOffset);
  const after = str.slice(range.endOffset);

  return (
    escapeHtml(before) +
    `<mark class="pdf-highlight">${escapeHtml(match)}</mark>` +
    escapeHtml(after)
  );
}
