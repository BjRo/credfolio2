"use client";

import { Plus, X } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { cn } from "@/lib/utils";

interface HighlightsEditorProps {
  highlights: string[];
  onChange: (highlights: string[]) => void;
  maxHighlights?: number;
  disabled?: boolean;
  className?: string;
  placeholder?: string;
}

export function HighlightsEditor({
  highlights,
  onChange,
  maxHighlights = 10,
  disabled = false,
  className,
  placeholder = "Add an achievement or responsibility...",
}: HighlightsEditorProps) {
  const handleAdd = () => {
    if (highlights.length < maxHighlights) {
      onChange([...highlights, ""]);
    }
  };

  const handleRemove = (index: number) => {
    onChange(highlights.filter((_, i) => i !== index));
  };

  const handleChange = (index: number, value: string) => {
    const updated = [...highlights];
    updated[index] = value;
    onChange(updated);
  };

  const handleKeyDown = (e: React.KeyboardEvent, index: number) => {
    // Add new highlight on Enter if current one is not empty
    if (e.key === "Enter" && highlights[index].trim() !== "") {
      e.preventDefault();
      if (highlights.length < maxHighlights) {
        onChange([...highlights, ""]);
        // Focus will be handled by useEffect in the new item
      }
    }
    // Remove empty highlight on Backspace
    if (e.key === "Backspace" && highlights[index] === "" && highlights.length > 1) {
      e.preventDefault();
      handleRemove(index);
    }
  };

  return (
    <div className={cn("space-y-2", className)}>
      {highlights.map((highlight, index) => (
        // biome-ignore lint/suspicious/noArrayIndexKey: Items are edited in-place by index, no stable ID available
        <div key={index} className="flex items-center gap-2">
          <span className="text-muted-foreground text-sm">•</span>
          <Input
            value={highlight}
            onChange={(e) => handleChange(index, e.target.value)}
            onKeyDown={(e) => handleKeyDown(e, index)}
            disabled={disabled}
            placeholder={placeholder}
            className="flex-1"
          />
          <Button
            type="button"
            variant="ghost"
            size="icon-sm"
            onClick={() => handleRemove(index)}
            disabled={disabled || highlights.length === 1}
            className="shrink-0"
            aria-label="Remove highlight"
          >
            <X className="h-4 w-4" />
          </Button>
        </div>
      ))}
      {highlights.length < maxHighlights && (
        <Button
          type="button"
          variant="outline"
          size="sm"
          onClick={handleAdd}
          disabled={disabled}
          className="w-full"
        >
          <Plus className="h-4 w-4 mr-1" />
          Add highlight
        </Button>
      )}
    </div>
  );
}
