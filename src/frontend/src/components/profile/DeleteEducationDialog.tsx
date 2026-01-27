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
import { DeleteEducationDocument } from "@/graphql/generated/graphql";

interface DeleteEducationDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  educationId: string;
  degree: string;
  institutionName: string;
  onSuccess?: () => void;
}

export function DeleteEducationDialog({
  open,
  onOpenChange,
  educationId,
  degree,
  institutionName,
  onSuccess,
}: DeleteEducationDialogProps) {
  const [error, setError] = useState<string | null>(null);
  const [result, deleteEducation] = useMutation(DeleteEducationDocument);

  const handleDelete = async () => {
    setError(null);

    try {
      const response = await deleteEducation({ id: educationId });

      if (response.error) {
        setError(response.error.message);
        return;
      }

      if (!response.data?.deleteEducation.success) {
        setError("Failed to delete education entry");
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
            Delete Education
          </DialogTitle>
          <DialogDescription>
            Are you sure you want to delete this education entry? This action cannot be undone.
          </DialogDescription>
        </DialogHeader>

        <div className="py-4">
          <p className="text-sm text-muted-foreground">You are about to delete:</p>
          <p className="mt-1 font-medium">{degree}</p>
          <p className="text-sm text-muted-foreground">{institutionName}</p>
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
