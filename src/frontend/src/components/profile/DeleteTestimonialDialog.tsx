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
import { DeleteTestimonialDocument } from "@/graphql/generated/graphql";

interface DeleteTestimonialDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  testimonialId: string;
  quote: string;
  authorName: string;
  onSuccess?: () => void;
}

export function DeleteTestimonialDialog({
  open,
  onOpenChange,
  testimonialId,
  quote,
  authorName,
  onSuccess,
}: DeleteTestimonialDialogProps) {
  const [error, setError] = useState<string | null>(null);
  const [result, deleteTestimonial] = useMutation(DeleteTestimonialDocument);

  const handleDelete = async () => {
    setError(null);

    try {
      const response = await deleteTestimonial({ id: testimonialId });

      if (response.error) {
        setError(response.error.message);
        return;
      }

      if (!response.data?.deleteTestimonial.success) {
        setError("Failed to delete testimonial");
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

  // Truncate quote for display if too long
  const displayQuote = quote.length > 100 ? `${quote.substring(0, 100)}...` : quote;

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <AlertTriangle className="h-5 w-5 text-destructive" />
            Delete Testimonial
          </DialogTitle>
          <DialogDescription>
            Are you sure you want to delete this testimonial? This action cannot be undone.
          </DialogDescription>
        </DialogHeader>

        <div className="py-4">
          <p className="text-sm text-muted-foreground">You are about to delete:</p>
          <p className="mt-1 font-medium italic">&ldquo;{displayQuote}&rdquo;</p>
          <p className="text-sm text-muted-foreground mt-1">— {authorName}</p>
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
