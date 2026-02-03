"use client";

import { CheckCircle2 } from "lucide-react";
import { cn } from "@/lib/utils";

interface CredibilityDotsProps {
  /** Number of validations (reference letters backing this item) */
  count: number;
  /** Additional CSS classes */
  className?: string;
}

/**
 * Displays a verified checkmark indicating validation from reference letters.
 * Shows a single CheckCircle2 icon when there's at least one validation.
 */
export function CredibilityDots({ count, className }: CredibilityDotsProps) {
  if (count === 0) {
    return null;
  }

  const label = `Validated by ${count} reference letter${count === 1 ? "" : "s"}`;

  return (
    <span
      role="img"
      className={cn("inline-flex items-center ml-1.5", className)}
      aria-label={label}
      title={label}
    >
      <CheckCircle2 className="h-3.5 w-3.5 text-green-500" aria-hidden="true" />
    </span>
  );
}
