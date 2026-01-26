"use client";

import { AlertTriangle } from "lucide-react";
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
import { DeleteExperienceDocument } from "@/graphql/generated/graphql";

interface DeleteExperienceDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  experienceId: string;
  experienceTitle: string;
  companyName: string;
  onSuccess?: () => void;
}

export function DeleteExperienceDialog({
  open,
  onOpenChange,
  experienceId,
  experienceTitle,
  companyName,
  onSuccess,
}: DeleteExperienceDialogProps) {
  const [error, setError] = useState<string | null>(null);
  const [result, deleteExperience] = useMutation(DeleteExperienceDocument);

  const handleDelete = async () => {
    setError(null);

    try {
      const response = await deleteExperience({ id: experienceId });

      if (response.error) {
        setError(response.error.message);
        return;
      }

      if (!response.data?.deleteExperience.success) {
        setError("Failed to delete experience");
        return;
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
          <DialogTitle className="flex items-center gap-2">
            <AlertTriangle className="h-5 w-5 text-destructive" />
            Delete Experience
          </DialogTitle>
          <DialogDescription>
            Are you sure you want to delete this experience? This action cannot be undone.
          </DialogDescription>
        </DialogHeader>

        <div className="py-4">
          <p className="text-sm text-muted-foreground">You are about to delete:</p>
          <p className="mt-1 font-medium">{experienceTitle}</p>
          <p className="text-sm text-muted-foreground">{companyName}</p>
        </div>

        {error && (
          <div className="rounded-md bg-destructive/15 p-3 text-sm text-destructive">{error}</div>
        )}

        <DialogFooter>
          <Button type="button" variant="outline" onClick={handleCancel} disabled={result.fetching}>
            Cancel
          </Button>
          <Button
            type="button"
            variant="destructive"
            onClick={handleDelete}
            disabled={result.fetching}
          >
            {result.fetching ? "Deleting..." : "Delete"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
