"use client";

import { useState } from "react";
import { useMutation } from "urql";

import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { CreateExperienceDocument, UpdateExperienceDocument } from "@/graphql/generated/graphql";
import { WorkExperienceForm, type WorkExperienceFormData } from "./WorkExperienceForm";

interface WorkExperienceFormDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  userId: string;
  experience?: {
    id: string;
    company: string;
    title: string;
    location?: string | null;
    startDate?: string | null;
    endDate?: string | null;
    isCurrent: boolean;
    description?: string | null;
    highlights: string[];
  };
  onSuccess?: () => void;
}

export function WorkExperienceFormDialog({
  open,
  onOpenChange,
  userId,
  experience,
  onSuccess,
}: WorkExperienceFormDialogProps) {
  const [error, setError] = useState<string | null>(null);
  const mode = experience?.id ? "edit" : "create";

  const [createResult, createExperience] = useMutation(CreateExperienceDocument);
  const [updateResult, updateExperience] = useMutation(UpdateExperienceDocument);

  const isSubmitting = createResult.fetching || updateResult.fetching;

  const handleSubmit = async (data: WorkExperienceFormData) => {
    setError(null);

    try {
      if (mode === "create") {
        const result = await createExperience({
          userId,
          input: {
            company: data.company,
            title: data.title,
            location: data.location || null,
            startDate: data.startDate || null,
            endDate: data.endDate || null,
            isCurrent: data.isCurrent,
            description: data.description || null,
            highlights: data.highlights.length > 0 ? data.highlights : null,
          },
        });

        if (result.error) {
          setError(result.error.message);
          return;
        }

        const response = result.data?.createExperience;
        if (response?.__typename === "ExperienceValidationError") {
          setError(response.message);
          return;
        }
      } else if (experience) {
        const result = await updateExperience({
          id: experience.id,
          input: {
            company: data.company,
            title: data.title,
            location: data.location || null,
            startDate: data.startDate || null,
            endDate: data.endDate || null,
            isCurrent: data.isCurrent,
            description: data.description || null,
            highlights: data.highlights.length > 0 ? data.highlights : null,
          },
        });

        if (result.error) {
          setError(result.error.message);
          return;
        }

        const response = result.data?.updateExperience;
        if (response?.__typename === "ExperienceValidationError") {
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
      <DialogContent className="sm:max-w-[600px] max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>
            {mode === "create" ? "Add Work Experience" : "Edit Work Experience"}
          </DialogTitle>
          <DialogDescription>
            {mode === "create"
              ? "Add a new position to your work history."
              : "Update the details of this position."}
          </DialogDescription>
        </DialogHeader>

        {error && (
          <div className="rounded-md bg-destructive/15 p-3 text-sm text-destructive">{error}</div>
        )}

        <WorkExperienceForm
          initialData={
            experience
              ? {
                  company: experience.company,
                  title: experience.title,
                  location: experience.location ?? "",
                  startDate: experience.startDate ?? "",
                  endDate: experience.endDate ?? "",
                  isCurrent: experience.isCurrent,
                  description: experience.description ?? "",
                  highlights: experience.highlights,
                }
              : undefined
          }
          onSubmit={handleSubmit}
          onCancel={handleCancel}
          isSubmitting={isSubmitting}
          mode={mode}
        />
      </DialogContent>
    </Dialog>
  );
}
