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
import { UpdateProfileHeaderDocument } from "@/graphql/generated/graphql";
import { ProfileHeaderForm, type ProfileHeaderFormData } from "./ProfileHeaderForm";

interface ProfileHeaderFormDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  userId: string;
  headerData?: {
    name: string;
    email?: string | null;
    phone?: string | null;
    location?: string | null;
    summary?: string | null;
  };
  onSuccess?: () => void;
}

export function ProfileHeaderFormDialog({
  open,
  onOpenChange,
  userId,
  headerData,
  onSuccess,
}: ProfileHeaderFormDialogProps) {
  const [error, setError] = useState<string | null>(null);

  const [updateResult, updateProfileHeader] = useMutation(UpdateProfileHeaderDocument);

  const isSubmitting = updateResult.fetching;

  const handleSubmit = async (data: ProfileHeaderFormData) => {
    setError(null);

    try {
      const result = await updateProfileHeader({
        userId,
        input: {
          name: data.name || null,
          email: data.email || null,
          phone: data.phone || null,
          location: data.location || null,
          summary: data.summary || null,
        },
      });

      if (result.error) {
        setError(result.error.message);
        return;
      }

      const response = result.data?.updateProfileHeader;
      if (response?.__typename === "ProfileHeaderValidationError") {
        setError(response.message);
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
      <DialogContent className="sm:max-w-[500px] max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>Edit Profile</DialogTitle>
          <DialogDescription>
            Update your personal information and professional summary.
          </DialogDescription>
        </DialogHeader>

        {error && (
          <div className="rounded-md bg-destructive/15 p-3 text-sm text-destructive">{error}</div>
        )}

        <ProfileHeaderForm
          initialData={
            headerData
              ? {
                  name: headerData.name,
                  email: headerData.email ?? "",
                  phone: headerData.phone ?? "",
                  location: headerData.location ?? "",
                  summary: headerData.summary ?? "",
                }
              : undefined
          }
          onSubmit={handleSubmit}
          onCancel={handleCancel}
          isSubmitting={isSubmitting}
        />
      </DialogContent>
    </Dialog>
  );
}
