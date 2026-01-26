import { Check } from "lucide-react";
import type * as React from "react";

import { cn } from "@/lib/utils";

interface CheckboxProps extends Omit<React.ComponentProps<"button">, "onChange"> {
  checked?: boolean;
  onCheckedChange?: (checked: boolean) => void;
}

function Checkbox({ className, checked = false, onCheckedChange, ...props }: CheckboxProps) {
  return (
    // biome-ignore lint/a11y/useSemanticElements: Custom styled checkbox using button with proper ARIA role
    <button
      type="button"
      role="checkbox"
      aria-checked={checked}
      data-state={checked ? "checked" : "unchecked"}
      data-slot="checkbox"
      className={cn(
        "peer h-4 w-4 shrink-0 rounded-sm border border-primary shadow focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50 data-[state=checked]:bg-primary data-[state=checked]:text-primary-foreground",
        className
      )}
      onClick={() => onCheckedChange?.(!checked)}
      {...props}
    >
      {checked && (
        <span className="flex items-center justify-center text-current">
          <Check className="h-3 w-3" />
        </span>
      )}
    </button>
  );
}

export { Checkbox };
