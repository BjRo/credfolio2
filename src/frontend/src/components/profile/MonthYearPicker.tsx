"use client";

import { cn } from "@/lib/utils";

const MONTHS = [
  { value: "01", label: "January" },
  { value: "02", label: "February" },
  { value: "03", label: "March" },
  { value: "04", label: "April" },
  { value: "05", label: "May" },
  { value: "06", label: "June" },
  { value: "07", label: "July" },
  { value: "08", label: "August" },
  { value: "09", label: "September" },
  { value: "10", label: "October" },
  { value: "11", label: "November" },
  { value: "12", label: "December" },
];

// Generate years from 1970 to current year + 10
const currentYear = new Date().getFullYear();
const YEARS = Array.from({ length: currentYear - 1970 + 11 }, (_, i) =>
  String(currentYear + 10 - i)
);

export interface MonthYearValue {
  month: string;
  year: string;
}

interface MonthYearPickerProps {
  value?: MonthYearValue;
  onChange?: (value: MonthYearValue) => void;
  disabled?: boolean;
  className?: string;
  placeholder?: {
    month?: string;
    year?: string;
  };
}

/**
 * Parses a date string into MonthYearValue.
 * Handles formats:
 * - ISO full: "2020-01-15", "2020-01-15T00:00:00Z"
 * - ISO year-month: "2020-01"
 * - Human readable: "Jan 2020", "January 2020"
 * - Slash format: "01/2020", "1/2020"
 * - Year only: "2020"
 */
export function parseDateString(dateStr: string | null | undefined): MonthYearValue {
  if (!dateStr) return { month: "", year: "" };

  const trimmed = dateStr.trim();

  // Handle "Present" or "Current"
  if (trimmed.toLowerCase() === "present" || trimmed.toLowerCase() === "current") {
    return { month: "", year: "" };
  }

  // Try to parse ISO format with day (YYYY-MM-DD or YYYY-MM-DDTHH:mm:ss)
  const isoFullMatch = trimmed.match(/^(\d{4})-(\d{2})-\d{2}/);
  if (isoFullMatch) {
    return { month: isoFullMatch[2], year: isoFullMatch[1] };
  }

  // Try to parse ISO year-month format (YYYY-MM)
  const isoYearMonthMatch = trimmed.match(/^(\d{4})-(\d{2})$/);
  if (isoYearMonthMatch) {
    return { month: isoYearMonthMatch[2], year: isoYearMonthMatch[1] };
  }

  // Try to parse slash format (MM/YYYY or M/YYYY)
  const slashMatch = trimmed.match(/^(\d{1,2})\/(\d{4})$/);
  if (slashMatch) {
    const month = slashMatch[1].padStart(2, "0");
    return { month, year: slashMatch[2] };
  }

  // Try to parse "Jan 2020" or "January 2020" format (month name + year)
  const monthYearMatch = trimmed.match(/^([A-Za-z]+)\s+(\d{4})$/);
  if (monthYearMatch) {
    const monthStr = monthYearMatch[1].toLowerCase();
    const year = monthYearMatch[2];
    const monthIndex = MONTHS.findIndex((m) =>
      m.label.toLowerCase().startsWith(monthStr.slice(0, 3))
    );
    if (monthIndex !== -1) {
      return { month: MONTHS[monthIndex].value, year };
    }
  }

  // Try to parse year-only format
  const yearMatch = trimmed.match(/^(\d{4})$/);
  if (yearMatch) {
    return { month: "", year: yearMatch[1] };
  }

  // Debug: log unrecognized format
  if (typeof window !== "undefined") {
    console.warn(`[parseDateString] Unrecognized date format: "${dateStr}"`);
  }

  return { month: "", year: "" };
}

/**
 * Formats a MonthYearValue into a date string like "Jan 2020".
 */
export function formatMonthYear(value: MonthYearValue): string {
  if (!value.year) return "";

  if (!value.month) return value.year;

  const month = MONTHS.find((m) => m.value === value.month);
  if (!month) return value.year;

  return `${month.label.slice(0, 3)} ${value.year}`;
}

export function MonthYearPicker({
  value = { month: "", year: "" },
  onChange,
  disabled = false,
  className,
  placeholder = { month: "Month", year: "Year" },
}: MonthYearPickerProps) {
  const handleMonthChange = (month: string) => {
    onChange?.({ ...value, month });
  };

  const handleYearChange = (year: string) => {
    onChange?.({ ...value, year });
  };

  const selectClassName = cn(
    "h-9 rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50",
    disabled && "opacity-50 cursor-not-allowed"
  );

  return (
    <div className={cn("flex gap-2", className)}>
      <select
        value={value.month}
        onChange={(e) => handleMonthChange(e.target.value)}
        disabled={disabled}
        className={cn(selectClassName, "flex-1")}
        aria-label="Month"
      >
        <option value="">{placeholder.month}</option>
        {MONTHS.map((month) => (
          <option key={month.value} value={month.value}>
            {month.label}
          </option>
        ))}
      </select>
      <select
        value={value.year}
        onChange={(e) => handleYearChange(e.target.value)}
        disabled={disabled}
        className={cn(selectClassName, "w-24")}
        aria-label="Year"
      >
        <option value="">{placeholder.year}</option>
        {YEARS.map((year) => (
          <option key={year} value={year}>
            {year}
          </option>
        ))}
      </select>
    </div>
  );
}
