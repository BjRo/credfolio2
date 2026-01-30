"use client";

import { cn } from "@/lib/utils";

interface CredibilityDotsProps {
  /** Number of validations (reference letters backing this item) */
  count: number;
  /** Additional CSS classes */
  className?: string;
}

/**
 * Displays credibility dots indicating how many reference letters validate an item.
 * - 0 dots: No validations (only from resume)
 * - 1 dot: 1 reference letter validates this
 * - 2 dots: 2 reference letters validate this
 * - 3 dots: 3+ reference letters validate this
 */
export function CredibilityDots({ count, className }: CredibilityDotsProps) {
  if (count === 0) {
    return null;
  }

  // Cap at 3 dots for display
  const displayCount = Math.min(count, 3);

  const label = `Validated by ${count} reference letter${count === 1 ? "" : "s"}`;

  return (
    <span
      role="img"
      className={cn("inline-flex items-center gap-0.5 ml-1.5", className)}
      aria-label={label}
      title={label}
    >
      {Array.from({ length: displayCount }).map((_, i) => (
        <span
          // biome-ignore lint/suspicious/noArrayIndexKey: Static dots that never reorder
          key={i}
          className="inline-block w-1.5 h-1.5 rounded-full bg-green-500"
          aria-hidden="true"
        />
      ))}
    </span>
  );
}
