"use client";

import { Camera, Loader2, X } from "lucide-react";
import Image from "next/image";
import { useCallback, useRef, useState } from "react";
import { useMutation } from "urql";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { UpdateAuthorDocument } from "@/graphql/generated/graphql";
import { GRAPHQL_UPLOAD_ENDPOINT } from "@/lib/urql/client";

interface AuthorEditModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  author: {
    id: string;
    name: string;
    title?: string | null;
    company?: string | null;
    linkedInUrl?: string | null;
    imageUrl?: string | null;
  };
  onSuccess?: () => void;
}

interface FormData {
  name: string;
  title: string;
  company: string;
  linkedInUrl: string;
}

const LINKEDIN_URL_PATTERN = /^https?:\/\/(www\.)?linkedin\.com\/in\/[\w-]+\/?$/i;

function validateLinkedInUrl(url: string): string | null {
  if (!url) return null;
  if (!LINKEDIN_URL_PATTERN.test(url)) {
    return "LinkedIn URL must be in the format: https://linkedin.com/in/username";
  }
  return null;
}

// Upload author image using XHR following GraphQL multipart request spec
async function uploadAuthorImageXhr(
  authorId: string,
  file: File
): Promise<{ success: boolean; error?: string }> {
  const operations = JSON.stringify({
    query: `
      mutation UploadAuthorImage($authorId: ID!, $file: Upload!) {
        uploadAuthorImage(authorId: $authorId, file: $file) {
          ... on UploadAuthorImageResult {
            __typename
            file {
              id
            }
            author {
              id
              imageUrl
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
      authorId,
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
          const data = result.data?.uploadAuthorImage;
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

export function AuthorEditModal({ open, onOpenChange, author, onSuccess }: AuthorEditModalProps) {
  const isUnknownAuthor = author.name === "unknown" || !author.name;

  const [formData, setFormData] = useState<FormData>({
    name: isUnknownAuthor ? "" : author.name,
    title: author.title ?? "",
    company: author.company ?? "",
    linkedInUrl: author.linkedInUrl ?? "",
  });
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [previewImageUrl, setPreviewImageUrl] = useState<string | null>(author.imageUrl ?? null);
  const [pendingImageFile, setPendingImageFile] = useState<File | null>(null);
  const [imageRemoved, setImageRemoved] = useState(false);
  const [isUploading, setIsUploading] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const [updateResult, updateAuthor] = useMutation(UpdateAuthorDocument);

  const isSubmitting = updateResult.fetching || isUploading;

  const handleInputChange = useCallback(
    (field: keyof FormData, value: string) => {
      setFormData((prev) => ({ ...prev, [field]: value }));
      // Clear error when user types
      if (errors[field]) {
        setErrors((prev) => {
          const next = { ...prev };
          delete next[field];
          return next;
        });
      }
    },
    [errors]
  );

  const handleFileChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    // Reset input so the same file can be selected again
    e.target.value = "";

    // Store the file for upload on submit
    setPendingImageFile(file);
    setImageRemoved(false);

    // Create a preview URL immediately
    const localPreviewUrl = URL.createObjectURL(file);
    setPreviewImageUrl(localPreviewUrl);
  }, []);

  const handleRemoveImage = useCallback(() => {
    setPreviewImageUrl(null);
    setPendingImageFile(null);
    setImageRemoved(true);
  }, []);

  const validate = (): boolean => {
    const newErrors: Record<string, string> = {};

    if (!formData.name.trim()) {
      newErrors.name = "Name is required";
    }

    const linkedInError = validateLinkedInUrl(formData.linkedInUrl);
    if (linkedInError) {
      newErrors.linkedInUrl = linkedInError;
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validate()) return;

    try {
      // Step 1: Upload image if there's a pending file
      if (pendingImageFile) {
        setIsUploading(true);
        try {
          await uploadAuthorImageXhr(author.id, pendingImageFile);
        } catch (err) {
          setErrors({ submit: err instanceof Error ? err.message : "Failed to upload image" });
          setIsUploading(false);
          return;
        }
        setIsUploading(false);
      }

      // Step 2: Update author details (name, title, company, linkedInUrl)
      // If image was removed and no new image uploaded, clear the imageId
      const result = await updateAuthor({
        id: author.id,
        input: {
          name: formData.name.trim(),
          title: formData.title.trim() || null,
          company: formData.company.trim() || null,
          linkedInUrl: formData.linkedInUrl.trim() || null,
          // Clear imageId if user removed the image and didn't upload a new one
          ...(imageRemoved && !pendingImageFile ? { imageId: "" } : {}),
        },
      });

      if (result.error) {
        setErrors({ submit: result.error.message });
        return;
      }

      onSuccess?.();
      onOpenChange(false);
    } catch (err) {
      setErrors({
        submit: err instanceof Error ? err.message : "An error occurred",
      });
    }
  };

  const handleCancel = () => {
    setErrors({});
    onOpenChange(false);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>{isUnknownAuthor ? "Add Author Details" : "Edit Author"}</DialogTitle>
          <DialogDescription>
            {isUnknownAuthor
              ? "The author of this testimonial wasn't detected. Add their details below."
              : "Update the author's information."}
          </DialogDescription>
        </DialogHeader>

        {errors.submit && (
          <div className="rounded-md bg-destructive/15 p-3 text-sm text-destructive">
            {errors.submit}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          {/* Author Image */}
          <div className="flex flex-col items-center gap-3">
            <div className="relative w-20 h-20">
              <button
                type="button"
                className={`w-20 h-20 rounded-full overflow-hidden flex items-center justify-center border-2 border-dashed border-muted-foreground/30 ${
                  previewImageUrl ? "bg-muted" : "bg-muted/50"
                } cursor-pointer hover:border-primary/50 transition-colors focus:outline-none focus:ring-2 focus:ring-primary focus:ring-offset-2`}
                onClick={() => fileInputRef.current?.click()}
                aria-label="Upload author photo"
                disabled={isSubmitting}
              >
                {previewImageUrl ? (
                  <Image
                    src={previewImageUrl}
                    alt={`Photo of ${formData.name || "author"}`}
                    width={80}
                    height={80}
                    className="w-full h-full object-cover"
                    unoptimized
                  />
                ) : (
                  <div className="flex flex-col items-center gap-1 text-muted-foreground">
                    <Camera className="w-6 h-6" />
                    <span className="text-xs">Add photo</span>
                  </div>
                )}
              </button>
              {previewImageUrl && !isSubmitting && (
                <button
                  type="button"
                  onClick={handleRemoveImage}
                  className="absolute -top-1 -right-1 w-6 h-6 rounded-full bg-destructive text-destructive-foreground flex items-center justify-center hover:bg-destructive/90"
                  aria-label="Remove photo"
                >
                  <X className="w-3 h-3" />
                </button>
              )}
            </div>
            <input
              ref={fileInputRef}
              type="file"
              accept="image/jpeg,image/png,image/gif,image/webp"
              className="hidden"
              onChange={handleFileChange}
              aria-label="Upload author photo"
              disabled={isSubmitting}
            />
            <p className="text-xs text-muted-foreground">
              Click to upload a profile photo (optional)
            </p>
          </div>

          {/* Name */}
          <div className="space-y-2">
            <Label htmlFor="name">
              Name <span className="text-destructive">*</span>
            </Label>
            <Input
              id="name"
              value={formData.name}
              onChange={(e) => handleInputChange("name", e.target.value)}
              placeholder="e.g., John Smith"
              className={errors.name ? "border-destructive" : ""}
              disabled={isSubmitting}
            />
            {errors.name && <p className="text-sm text-destructive">{errors.name}</p>}
          </div>

          {/* Title */}
          <div className="space-y-2">
            <Label htmlFor="title">Title / Position</Label>
            <Input
              id="title"
              value={formData.title}
              onChange={(e) => handleInputChange("title", e.target.value)}
              placeholder="e.g., Engineering Manager"
              disabled={isSubmitting}
            />
          </div>

          {/* Company */}
          <div className="space-y-2">
            <Label htmlFor="company">Company / Organization</Label>
            <Input
              id="company"
              value={formData.company}
              onChange={(e) => handleInputChange("company", e.target.value)}
              placeholder="e.g., Acme Corp"
              disabled={isSubmitting}
            />
          </div>

          {/* LinkedIn URL */}
          <div className="space-y-2">
            <Label htmlFor="linkedInUrl">LinkedIn Profile</Label>
            <Input
              id="linkedInUrl"
              type="url"
              value={formData.linkedInUrl}
              onChange={(e) => handleInputChange("linkedInUrl", e.target.value)}
              placeholder="https://linkedin.com/in/username"
              className={errors.linkedInUrl ? "border-destructive" : ""}
              disabled={isSubmitting}
            />
            {errors.linkedInUrl && <p className="text-sm text-destructive">{errors.linkedInUrl}</p>}
          </div>

          {/* Actions */}
          <div className="flex justify-end gap-3 pt-4">
            <Button type="button" variant="outline" onClick={handleCancel} disabled={isSubmitting}>
              Cancel
            </Button>
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting ? (
                <>
                  <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                  Saving...
                </>
              ) : (
                "Save"
              )}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}
