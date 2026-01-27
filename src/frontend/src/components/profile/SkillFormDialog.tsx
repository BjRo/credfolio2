"use client";

import { useState } from "react";
import { useMutation } from "urql";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  CreateSkillDocument,
  SkillCategory,
  UpdateSkillDocument,
} from "@/graphql/generated/graphql";

const CATEGORY_OPTIONS: { value: SkillCategory; label: string }[] = [
  { value: SkillCategory.Technical, label: "Technical" },
  { value: SkillCategory.Soft, label: "Soft" },
  { value: SkillCategory.Domain, label: "Domain" },
];

interface SkillFormDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  userId: string;
  skill?: {
    id: string;
    name: string;
    category: SkillCategory;
  };
  defaultCategory?: SkillCategory;
  onSuccess?: () => void;
}

export function SkillFormDialog({
  open,
  onOpenChange,
  userId,
  skill,
  defaultCategory,
  onSuccess,
}: SkillFormDialogProps) {
  const mode = skill?.id ? "edit" : "create";
  const [name, setName] = useState(skill?.name ?? "");
  const [category, setCategory] = useState<SkillCategory>(
    skill?.category ?? defaultCategory ?? SkillCategory.Technical
  );
  const [error, setError] = useState<string | null>(null);

  const [createResult, createSkill] = useMutation(CreateSkillDocument);
  const [updateResult, updateSkill] = useMutation(UpdateSkillDocument);

  const isSubmitting = createResult.fetching || updateResult.fetching;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    const trimmedName = name.trim();
    if (!trimmedName) {
      setError("Skill name is required");
      return;
    }

    try {
      if (mode === "create") {
        const result = await createSkill({
          userId,
          input: { name: trimmedName, category },
        });

        if (result.error) {
          setError(result.error.message);
          return;
        }

        const response = result.data?.createSkill;
        if (response?.__typename === "SkillValidationError") {
          setError(response.message);
          return;
        }
      } else if (skill) {
        const result = await updateSkill({
          id: skill.id,
          input: { name: trimmedName, category },
        });

        if (result.error) {
          setError(result.error.message);
          return;
        }

        const response = result.data?.updateSkill;
        if (response?.__typename === "SkillValidationError") {
          setError(response.message);
          return;
        }
      }

      onSuccess?.();
      onOpenChange(false);
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred");
    }
  };

  const handleCancel = () => {
    setError(null);
    onOpenChange(false);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>{mode === "create" ? "Add Skill" : "Edit Skill"}</DialogTitle>
          <DialogDescription>
            {mode === "create"
              ? "Add a new skill to your profile."
              : "Update the details of this skill."}
          </DialogDescription>
        </DialogHeader>

        {error && (
          <div className="rounded-md bg-destructive/15 p-3 text-sm text-destructive">{error}</div>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="skill-name">Skill Name</Label>
            <Input
              id="skill-name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="e.g., React, Project Management, Data Analysis"
              required
              autoFocus
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="skill-category">Category</Label>
            <select
              id="skill-category"
              value={category}
              onChange={(e) => setCategory(e.target.value as SkillCategory)}
              className="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-base shadow-sm transition-colors placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50 md:text-sm"
            >
              {CATEGORY_OPTIONS.map((opt) => (
                <option key={opt.value} value={opt.value}>
                  {opt.label}
                </option>
              ))}
            </select>
          </div>

          <DialogFooter>
            <Button type="button" variant="outline" onClick={handleCancel} disabled={isSubmitting}>
              Cancel
            </Button>
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting ? "Saving..." : mode === "create" ? "Add Skill" : "Save Changes"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
