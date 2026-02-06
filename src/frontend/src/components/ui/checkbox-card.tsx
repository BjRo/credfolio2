"use client";

import type * as React from "react";
import { Checkbox } from "@/components/ui/checkbox";

interface CheckboxCardProps {
  checked: boolean;
  onToggle: () => void;
  disabled: boolean;
  selectedClassName: string;
  borderStyle?: string;
  unselectedClassName?: string;
  children: React.ReactNode;
}

function CheckboxCard({
  checked,
  onToggle,
  disabled,
  selectedClassName,
  borderStyle = "border",
  unselectedClassName = "bg-card border-border hover:bg-muted/50",
  children,
}: CheckboxCardProps) {
  return (
    // biome-ignore lint/a11y/useSemanticElements: Custom styled checkbox card with inner Checkbox component
    <div
      role="checkbox"
      aria-checked={checked}
      tabIndex={disabled ? -1 : 0}
      onClick={() => !disabled && onToggle()}
      onKeyDown={(e) => {
        if (!disabled && (e.key === " " || e.key === "Enter")) {
          e.preventDefault();
          onToggle();
        }
      }}
      className={`flex items-start gap-3 p-4 rounded-lg ${borderStyle} transition-colors w-full text-left ${
        disabled
          ? "bg-muted cursor-not-allowed border-border"
          : checked
            ? `${selectedClassName} cursor-pointer`
            : `${unselectedClassName} cursor-pointer`
      }`}
    >
      <Checkbox
        checked={checked}
        onCheckedChange={() => onToggle()}
        onClick={(e) => e.stopPropagation()}
        disabled={disabled}
        className="mt-1"
        tabIndex={-1}
        aria-hidden="true"
      />
      <div className="flex-1">{children}</div>
    </div>
  );
}

export { CheckboxCard };
export type { CheckboxCardProps };
