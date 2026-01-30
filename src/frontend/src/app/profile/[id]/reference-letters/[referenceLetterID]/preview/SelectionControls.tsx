"use client";

import { Button } from "@/components/ui/button";

interface SelectionControlsProps {
  selectedCount: number;
  totalCount: number;
  onSelectAll: () => void;
  onDeselectAll: () => void;
  disabled?: boolean;
}

export function SelectionControls({
  selectedCount,
  totalCount,
  onSelectAll,
  onDeselectAll,
  disabled = false,
}: SelectionControlsProps) {
  const allSelected = selectedCount === totalCount;
  const noneSelected = selectedCount === 0;

  return (
    <div className="flex items-center gap-3">
      <span className="text-sm text-muted-foreground">
        {selectedCount} of {totalCount} selected
      </span>
      {!allSelected && (
        <Button
          variant="ghost"
          size="sm"
          onClick={onSelectAll}
          disabled={disabled}
          className="text-xs"
        >
          Select All
        </Button>
      )}
      {!noneSelected && (
        <Button
          variant="ghost"
          size="sm"
          onClick={onDeselectAll}
          disabled={disabled}
          className="text-xs"
        >
          Deselect All
        </Button>
      )}
    </div>
  );
}
