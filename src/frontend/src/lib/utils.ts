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
