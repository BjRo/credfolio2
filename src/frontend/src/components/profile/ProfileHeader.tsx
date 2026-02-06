"use client";

import {
  Camera,
  ChevronDown,
  ChevronUp,
  Loader2,
  Mail,
  MapPin,
  Pencil,
  Phone,
  User,
  X,
} from "lucide-react";
import Image from "next/image";
import { useCallback, useRef, useState } from "react";
import { useMutation } from "urql";
import { Button } from "@/components/ui/button";
import { DeleteProfilePhotoDocument } from "@/graphql/generated/graphql";
import { GRAPHQL_UPLOAD_ENDPOINT } from "@/lib/urql/client";
import { ProfileHeaderFormDialog } from "./ProfileHeaderFormDialog";

interface ProfileHeaderProps {
  profile: {
    name?: string | null;
    email?: string | null;
    phone?: string | null;
    location?: string | null;
    summary?: string | null;
    profilePhotoUrl?: string | null;
  };
  userId?: string;
  onMutationSuccess?: () => void;
}

interface ProfileSummaryProps {
  summary: string;
  collapsedLines?: number;
}

const SUMMARY_COLLAPSE_THRESHOLD = 200;

function ProfileSummary({ summary, collapsedLines = 3 }: ProfileSummaryProps) {
  const [isExpanded, setIsExpanded] = useState(false);
  const isLongSummary = summary.length > SUMMARY_COLLAPSE_THRESHOLD;

  return (
    <div className="mt-6">
      <h2 className="text-sm font-semibold text-foreground uppercase tracking-wide mb-2">
        Summary
      </h2>
      <div className="relative">
        <p
          className={`text-foreground ${
            !isExpanded && isLongSummary ? `line-clamp-${collapsedLines}` : ""
          }`}
          style={
            !isExpanded && isLongSummary
              ? {
                  display: "-webkit-box",
                  WebkitLineClamp: collapsedLines,
                  WebkitBoxOrient: "vertical",
                  overflow: "hidden",
                }
              : undefined
          }
        >
          {summary}
        </p>
        {isLongSummary && (
          <button
            type="button"
            onClick={() => setIsExpanded(!isExpanded)}
            className="mt-2 text-sm text-primary hover:text-primary/80 flex items-center gap-1"
            aria-expanded={isExpanded}
          >
            {isExpanded ? (
              <>
                Show less <ChevronUp className="w-4 h-4" aria-hidden="true" />
              </>
            ) : (
              <>
                Show more <ChevronDown className="w-4 h-4" aria-hidden="true" />
              </>
            )}
          </button>
        )}
      </div>
    </div>
  );
}

interface ProfileAvatarProps {
  name: string;
  photoUrl?: string | null;
  userId?: string;
  onUploadSuccess?: () => void;
}

function ProfileAvatar({ name, photoUrl, userId, onUploadSuccess }: ProfileAvatarProps) {
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [isHovered, setIsHovered] = useState(false);
  const [isUploading, setIsUploading] = useState(false);
  const [deleteResult, deletePhoto] = useMutation(DeleteProfilePhotoDocument);

  const isLoading = isUploading || deleteResult.fetching;
  const canEdit = !!userId;

  const initials = name
    .split(" ")
    .map((n) => n[0])
    .join("")
    .toUpperCase()
    .slice(0, 2);

  // Use XHR for file upload following GraphQL multipart request spec
  const uploadFile = useCallback(
    async (file: File) => {
      if (!userId) return;

      setIsUploading(true);

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

      try {
        const response = await new Promise<{ success: boolean; error?: string }>(
          (resolve, reject) => {
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
          }
        );

        if (response.success) {
          onUploadSuccess?.();
        }
      } catch (err) {
        console.error("Profile photo upload failed:", err);
      } finally {
        setIsUploading(false);
      }
    },
    [userId, onUploadSuccess]
  );

  const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file || !userId) return;

    // Reset input so the same file can be selected again
    e.target.value = "";

    await uploadFile(file);
  };

  const handleDeletePhoto = async () => {
    if (!userId) return;

    const result = await deletePhoto({ userId });

    if (result.data?.deleteProfilePhoto?.__typename === "DeleteProfilePhotoResult") {
      onUploadSuccess?.();
    }
  };

  const handleClick = () => {
    if (canEdit && !isLoading) {
      fileInputRef.current?.click();
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if ((e.key === "Enter" || e.key === " ") && canEdit && !isLoading) {
      e.preventDefault();
      fileInputRef.current?.click();
    }
  };

  return (
    <div className="relative w-20 h-20 flex-shrink-0">
      {/* Avatar Display - using button for keyboard accessibility */}
      <button
        type="button"
        className={`w-20 h-20 rounded-full overflow-hidden flex items-center justify-center border-0 ${
          photoUrl ? "bg-muted" : "bg-primary"
        } ${canEdit ? "cursor-pointer focus:outline-none focus:ring-2 focus:ring-primary focus:ring-offset-2" : ""}`}
        aria-label={canEdit ? `Change profile photo for ${name}` : `Avatar for ${name}`}
        onClick={handleClick}
        onKeyDown={handleKeyDown}
        onMouseEnter={() => setIsHovered(true)}
        onMouseLeave={() => setIsHovered(false)}
        onFocus={() => setIsHovered(true)}
        onBlur={() => setIsHovered(false)}
        disabled={!canEdit || isLoading}
      >
        {isLoading ? (
          <Loader2 className="w-8 h-8 text-primary-foreground animate-spin" aria-hidden="true" />
        ) : photoUrl ? (
          <Image
            src={photoUrl}
            alt={`Profile photo of ${name}`}
            width={80}
            height={80}
            className="w-full h-full object-cover"
            unoptimized
          />
        ) : initials ? (
          <span className="text-2xl font-semibold text-primary-foreground">{initials}</span>
        ) : (
          <User className="w-10 h-10 text-primary-foreground" aria-hidden="true" />
        )}
      </button>

      {/* Hover Overlay */}
      {canEdit && isHovered && !isLoading && (
        <div
          className="absolute inset-0 rounded-full bg-black/50 flex items-center justify-center pointer-events-none"
          aria-hidden="true"
        >
          <Camera className="w-6 h-6 text-white" />
        </div>
      )}

      {/* Delete Button */}
      {canEdit && photoUrl && isHovered && !isLoading && (
        <button
          type="button"
          onClick={(e) => {
            e.stopPropagation();
            handleDeletePhoto();
          }}
          className="absolute -top-1 -right-1 w-6 h-6 rounded-full bg-destructive text-destructive-foreground flex items-center justify-center hover:bg-destructive/90"
          aria-label="Remove photo"
        >
          <X className="w-3 h-3" />
        </button>
      )}

      {/* Hidden File Input */}
      <input
        ref={fileInputRef}
        type="file"
        accept="image/jpeg,image/png,image/gif,image/webp"
        className="hidden"
        onChange={handleFileChange}
        aria-label="Upload profile photo"
      />
    </div>
  );
}

export function ProfileHeader({ profile, userId, onMutationSuccess }: ProfileHeaderProps) {
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);

  const displayName = profile.name || "Unknown";
  const displayEmail = profile.email;
  const displayPhone = profile.phone;
  const displayLocation = profile.location;
  const displaySummary = profile.summary;
  const displayPhotoUrl = profile.profilePhotoUrl;

  const handleEditSuccess = () => {
    onMutationSuccess?.();
  };

  return (
    <>
      <div className="bg-card border rounded-lg p-6 sm:p-8 relative">
        {/* Edit button - positioned absolutely in top-right */}
        {userId && (
          <Button
            variant="ghost"
            size="icon"
            onClick={() => setIsEditDialogOpen(true)}
            aria-label="Edit profile"
            className="absolute top-4 right-4 sm:top-6 sm:right-6"
          >
            <Pencil className="w-4 h-4" />
          </Button>
        )}

        <div className="flex flex-col sm:flex-row sm:items-start gap-4 sm:gap-6">
          <ProfileAvatar
            name={displayName}
            photoUrl={displayPhotoUrl}
            userId={userId}
            onUploadSuccess={onMutationSuccess}
          />
          <div className="flex-1 min-w-0 pr-8">
            <h1 className="text-2xl sm:text-3xl font-bold text-foreground truncate">
              {displayName}
            </h1>
            <div className="mt-2 space-y-1 text-muted-foreground">
              {displayEmail && (
                <p className="flex items-center gap-2 text-sm sm:text-base">
                  <Mail className="w-4 h-4 flex-shrink-0" aria-hidden="true" />
                  <a href={`mailto:${displayEmail}`} className="hover:text-primary truncate">
                    {displayEmail}
                  </a>
                </p>
              )}
              {displayPhone && (
                <p className="flex items-center gap-2 text-sm sm:text-base">
                  <Phone className="w-4 h-4 flex-shrink-0" aria-hidden="true" />
                  <a href={`tel:${displayPhone}`} className="hover:text-primary">
                    {displayPhone}
                  </a>
                </p>
              )}
              {displayLocation && (
                <p className="flex items-center gap-2 text-sm sm:text-base">
                  <MapPin className="w-4 h-4 flex-shrink-0" aria-hidden="true" />
                  <span>{displayLocation}</span>
                </p>
              )}
            </div>
          </div>
        </div>

        {displaySummary && <ProfileSummary summary={displaySummary} />}
      </div>

      {userId && (
        <ProfileHeaderFormDialog
          open={isEditDialogOpen}
          onOpenChange={setIsEditDialogOpen}
          userId={userId}
          headerData={{
            name: displayName,
            email: displayEmail,
            phone: displayPhone,
            location: displayLocation,
            summary: displaySummary,
          }}
          photoUrl={displayPhotoUrl}
          onSuccess={handleEditSuccess}
        />
      )}
    </>
  );
}
