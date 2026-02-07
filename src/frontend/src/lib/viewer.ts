const MAX_HIGHLIGHT_LENGTH = 500;

/**
 * Build a URL for the source document viewer.
 * Used by TestimonialsSection and ValidationPopover to create deep links.
 */
export function buildViewerUrl(letterId: string, highlightText?: string): string {
  const params = new URLSearchParams();
  params.set("letterId", letterId);
  if (highlightText) {
    const truncated =
      highlightText.length > MAX_HIGHLIGHT_LENGTH
        ? highlightText.slice(0, MAX_HIGHLIGHT_LENGTH)
        : highlightText;
    params.set("highlight", truncated);
  }
  return `/viewer?${params.toString()}`;
}
