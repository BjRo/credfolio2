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
import {
  DeleteProfilePhotoDocument,
  UpdateProfileHeaderDocument,
} from "@/graphql/generated/graphql";
import { GRAPHQL_UPLOAD_ENDPOINT } from "@/lib/urql/client";
import { ProfileHeaderForm, type ProfileHeaderFormData } from "./ProfileHeaderForm";

// Upload profile photo using XHR following GraphQL multipart request spec
async function uploadProfilePhotoXhr(
  userId: string,
  file: File
): Promise<{ success: boolean; error?: string }> {
  const operations = JSON.stringify({
    query: `
      mutation UploadProfilePhoto($userId: ID!, $file: Upload!) {
        uploadProfilePhoto(userId: $userId, file: $file) {
          ... on UploadProfilePhotoResult {
            __typename
            profile {
              profilePhotoUrl
            }
          }
          ... on FileValidationError {
            __typename
            message
            field
          }
        }
      }
    `,
    variables: {
      userId,
      file: null,
    },
  });

  const map = JSON.stringify({
    "0": ["variables.file"],
  });

  const formData = new FormData();
  formData.append("operations", operations);
  formData.append("map", map);
  formData.append("0", file);

  return new Promise((resolve, reject) => {
    const xhr = new XMLHttpRequest();

    xhr.addEventListener("load", () => {
      if (xhr.status >= 200 && xhr.status < 300) {
        try {
          const result = JSON.parse(xhr.responseText);
          if (result.errors?.length) {
            reject(new Error(result.errors[0].message));
            return;
          }
          const data = result.data?.uploadProfilePhoto;
          if (data?.__typename === "FileValidationError") {
            reject(new Error(data.message));
            return;
          }
          resolve({ success: true });
        } catch {
          reject(new Error("Failed to parse response"));
        }
      } else {
        reject(new Error(`Upload failed with status ${xhr.status}`));
      }
    });

    xhr.addEventListener("error", () => {
      reject(new Error("Network error during upload"));
    });

    xhr.open("POST", GRAPHQL_UPLOAD_ENDPOINT);
    xhr.send(formData);
  });
}

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
  photoUrl?: string | null;
  onSuccess?: () => void;
}

export function ProfileHeaderFormDialog({
  open,
  onOpenChange,
  userId,
  headerData,
  photoUrl,
  onSuccess,
}: ProfileHeaderFormDialogProps) {
  const [error, setError] = useState<string | null>(null);
  const [isUploading, setIsUploading] = useState(false);

  const [updateResult, updateProfileHeader] = useMutation(UpdateProfileHeaderDocument);
  const [, deletePhoto] = useMutation(DeleteProfilePhotoDocument);

  const isSubmitting = updateResult.fetching || isUploading;

  const handleSubmit = async (data: ProfileHeaderFormData) => {
    setError(null);

    try {
      // Step 1: Handle image upload if there's a pending file
      if (data.pendingImageFile) {
        setIsUploading(true);
        try {
          await uploadProfilePhotoXhr(userId, data.pendingImageFile);
        } catch (err) {
          setError(err instanceof Error ? err.message : "Failed to upload image");
          setIsUploading(false);
          return;
        }
        setIsUploading(false);
      }

      // Step 2: Handle image deletion if image was removed
      if (data.imageRemoved && !data.pendingImageFile) {
        const deleteResult = await deletePhoto({ userId });
        if (deleteResult.error) {
          setError(deleteResult.error.message);
          return;
        }
      }

      // Step 3: Update profile header data
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
          photoUrl={photoUrl}
          onSubmit={handleSubmit}
          onCancel={handleCancel}
          isSubmitting={isSubmitting}
        />
      </DialogContent>
    </Dialog>
  );
}
