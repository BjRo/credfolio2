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
import { CreateEducationDocument, UpdateEducationDocument } from "@/graphql/generated/graphql";
import { EducationForm, type EducationFormData } from "./EducationForm";

interface EducationFormDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  userId: string;
  education?: {
    id: string;
    institution: string;
    degree: string;
    field?: string | null;
    startDate?: string | null;
    endDate?: string | null;
    isCurrent: boolean;
    description?: string | null;
    gpa?: string | null;
  };
  mode?: "create" | "edit";
  onSuccess?: () => void;
}

export function EducationFormDialog({
  open,
  onOpenChange,
  userId,
  education,
  mode: modeOverride,
  onSuccess,
}: EducationFormDialogProps) {
  const [error, setError] = useState<string | null>(null);
  const mode = modeOverride ?? (education?.id ? "edit" : "create");

  const [createResult, createEducation] = useMutation(CreateEducationDocument);
  const [updateResult, updateEducation] = useMutation(UpdateEducationDocument);

  const isSubmitting = createResult.fetching || updateResult.fetching;

  const handleSubmit = async (data: EducationFormData) => {
    setError(null);

    try {
      if (mode === "create") {
        const result = await createEducation({
          userId,
          input: {
            institution: data.institution,
            degree: data.degree,
            field: data.field || null,
            startDate: data.startDate || null,
            endDate: data.endDate || null,
            isCurrent: data.isCurrent,
            description: data.description || null,
            gpa: data.gpa || null,
          },
        });

        if (result.error) {
          setError(result.error.message);
          return;
        }

        const response = result.data?.createEducation;
        if (response?.__typename === "EducationValidationError") {
          setError(response.message);
          return;
        }
      } else if (education) {
        const result = await updateEducation({
          id: education.id,
          input: {
            institution: data.institution,
            degree: data.degree,
            field: data.field || null,
            startDate: data.startDate || null,
            endDate: data.endDate || null,
            isCurrent: data.isCurrent,
            description: data.description || null,
            gpa: data.gpa || null,
          },
        });

        if (result.error) {
          setError(result.error.message);
          return;
        }

        const response = result.data?.updateEducation;
        if (response?.__typename === "EducationValidationError") {
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
          <DialogTitle>{mode === "create" ? "Add Education" : "Edit Education"}</DialogTitle>
          <DialogDescription>
            {mode === "create"
              ? "Add a new education entry to your profile."
              : "Update the details of this education entry."}
          </DialogDescription>
        </DialogHeader>

        {error && (
          <div className="rounded-md bg-destructive/15 p-3 text-sm text-destructive">{error}</div>
        )}

        <EducationForm
          initialData={
            education
              ? {
                  institution: education.institution,
                  degree: education.degree,
                  field: education.field ?? "",
                  startDate: education.startDate ?? "",
                  endDate: education.endDate ?? "",
                  isCurrent: education.isCurrent,
                  description: education.description ?? "",
                  gpa: education.gpa ?? "",
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
