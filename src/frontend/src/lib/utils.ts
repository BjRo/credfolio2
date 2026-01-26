import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

/**
 * Formats a date string for display.
 * Handles ISO format (2020-01-15, 2020-01-15T00:00:00Z) and human-readable formats.
 * Returns null if the input looks like invalid/garbage data.
 */
export function formatDate(dateStr: string | null | undefined): string | null {
  if (!dateStr) return null;

  // Check for obviously invalid data (contains JSON characters or looks like garbage)
  if (
    dateStr.includes("{") ||
    dateStr.includes("}") ||
    dateStr.includes("]") ||
    dateStr.includes('"')
  ) {
    return null;
  }

  // Check for zero/null dates (like -0001-11-30 which Go uses for zero time)
  if (dateStr.startsWith("-0001") || dateStr.startsWith("0001-01-01")) {
    return null;
  }

  // If it's "Present" or similar, return as-is
  if (dateStr.toLowerCase() === "present" || dateStr.toLowerCase() === "current") {
    return "Present";
  }

  // Try to parse ISO format (YYYY-MM-DD or YYYY-MM-DDTHH:mm:ss)
  const isoMatch = dateStr.match(/^-?(\d{4})-(\d{2})-(\d{2})/);
  if (isoMatch) {
    const year = isoMatch[1];
    const month = parseInt(isoMatch[2], 10);
    const monthNames = [
      "Jan",
      "Feb",
      "Mar",
      "Apr",
      "May",
      "Jun",
      "Jul",
      "Aug",
      "Sep",
      "Oct",
      "Nov",
      "Dec",
    ];

    // Validate month
    if (month >= 1 && month <= 12) {
      return `${monthNames[month - 1]} ${year}`;
    }
    // If month is invalid, just return the year
    return year;
  }

  // If it's already human-readable, return as-is
  return dateStr;
}

/**
 * Parse a date string to a Date object for calculations.
 * Handles ISO format (2020-01-15) and human-readable format (Jan 2020).
 * Returns null if the date cannot be parsed.
 */
export function parseToDate(dateStr: string | null | undefined): Date | null {
  if (!dateStr) return null;

  // Skip invalid data
  if (
    dateStr.includes("{") ||
    dateStr.includes("}") ||
    dateStr.startsWith("-0001") ||
    dateStr.startsWith("0001-01-01")
  ) {
    return null;
  }

  // Handle "Present" or "current"
  if (dateStr.toLowerCase() === "present" || dateStr.toLowerCase() === "current") {
    return new Date();
  }

  // Try ISO format (YYYY-MM-DD)
  const isoMatch = dateStr.match(/^(\d{4})-(\d{2})-(\d{2})/);
  if (isoMatch) {
    const year = parseInt(isoMatch[1], 10);
    const month = parseInt(isoMatch[2], 10) - 1; // JS months are 0-indexed
    const day = parseInt(isoMatch[3], 10);
    return new Date(year, month, day);
  }

  // Try human-readable format (e.g., "Jan 2020", "January 2020", "2020")
  const monthNames: Record<string, number> = {
    jan: 0,
    january: 0,
    feb: 1,
    february: 1,
    mar: 2,
    march: 2,
    apr: 3,
    april: 3,
    may: 4,
    jun: 5,
    june: 5,
    jul: 6,
    july: 6,
    aug: 7,
    august: 7,
    sep: 8,
    sept: 8,
    september: 8,
    oct: 9,
    october: 9,
    nov: 10,
    november: 10,
    dec: 11,
    december: 11,
  };

  // Match "Month Year" or "Month, Year"
  const monthYearMatch = dateStr.match(/^([a-zA-Z]+),?\s*(\d{4})$/);
  if (monthYearMatch) {
    const monthKey = monthYearMatch[1].toLowerCase();
    const month = monthNames[monthKey];
    const year = parseInt(monthYearMatch[2], 10);
    if (month !== undefined && !Number.isNaN(year)) {
      return new Date(year, month, 1);
    }
  }

  // Match just year "2020"
  const yearMatch = dateStr.match(/^(\d{4})$/);
  if (yearMatch) {
    const year = parseInt(yearMatch[1], 10);
    return new Date(year, 0, 1);
  }

  return null;
}

/**
 * Calculate the duration in months between two dates.
 * Returns null if either date cannot be parsed.
 */
export function calculateDurationMonths(
  startDate: string | null | undefined,
  endDate: string | null | undefined,
  isCurrent: boolean
): number | null {
  const start = parseToDate(startDate);
  const end = isCurrent ? new Date() : parseToDate(endDate);

  if (!start || !end) return null;

  const months =
    (end.getFullYear() - start.getFullYear()) * 12 + (end.getMonth() - start.getMonth());

  return Math.max(0, months);
}

/**
 * Format a duration in months as a human-readable string.
 * Examples: "2 yrs 3 mos", "8 mos", "1 yr 1 mo"
 */
export function formatDuration(months: number): string {
  if (months < 1) return "< 1 mo";

  const years = Math.floor(months / 12);
  const remainingMonths = months % 12;

  const parts: string[] = [];

  if (years > 0) {
    parts.push(`${years} ${years === 1 ? "yr" : "yrs"}`);
  }

  if (remainingMonths > 0) {
    parts.push(`${remainingMonths} ${remainingMonths === 1 ? "mo" : "mos"}`);
  }

  return parts.join(" ");
}

interface DateRange {
  startDate?: string | null | undefined;
  endDate?: string | null | undefined;
  isCurrent: boolean;
}

/**
 * Calculate total tenure for multiple experiences, merging overlapping periods.
 * Returns null if no valid date ranges are found.
 */
export function calculateTotalTenure(experiences: DateRange[]): number | null {
  // Convert to parsed date ranges
  const ranges: { start: Date; end: Date }[] = [];

  for (const exp of experiences) {
    const start = parseToDate(exp.startDate);
    const end = exp.isCurrent ? new Date() : parseToDate(exp.endDate);

    if (start && end && start <= end) {
      ranges.push({ start, end });
    }
  }

  if (ranges.length === 0) return null;

  // Sort by start date
  ranges.sort((a, b) => a.start.getTime() - b.start.getTime());

  // Merge overlapping ranges
  const merged: { start: Date; end: Date }[] = [];
  let current = ranges[0];

  for (let i = 1; i < ranges.length; i++) {
    const next = ranges[i];
    if (next.start <= current.end) {
      // Overlapping - extend current range if needed
      current.end = new Date(Math.max(current.end.getTime(), next.end.getTime()));
    } else {
      // Non-overlapping - push current and start new
      merged.push(current);
      current = next;
    }
  }
  merged.push(current);

  // Calculate total months
  let totalMonths = 0;
  for (const range of merged) {
    const months =
      (range.end.getFullYear() - range.start.getFullYear()) * 12 +
      (range.end.getMonth() - range.start.getMonth());
    totalMonths += Math.max(0, months);
  }

  return totalMonths;
}
