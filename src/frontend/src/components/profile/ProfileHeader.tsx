"use client";

import { ChevronDown, ChevronUp, Mail, MapPin, Pencil, Phone, User } from "lucide-react";
import { useState } from "react";
import { Button } from "@/components/ui/button";
import { ProfileHeaderFormDialog } from "./ProfileHeaderFormDialog";
import type { ProfileData } from "./types";

interface ProfileHeaderOverrides {
  name?: string | null;
  email?: string | null;
  phone?: string | null;
  location?: string | null;
  summary?: string | null;
}

interface ProfileHeaderProps {
  data: ProfileData;
  profileOverrides?: ProfileHeaderOverrides;
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

function AvatarPlaceholder({ name }: { name: string }) {
  const initials = name
    .split(" ")
    .map((n) => n[0])
    .join("")
    .toUpperCase()
    .slice(0, 2);

  return (
    <div
      className="w-20 h-20 rounded-full bg-primary flex items-center justify-center flex-shrink-0"
      role="img"
      aria-label={`Avatar for ${name}`}
    >
      {initials ? (
        <span className="text-2xl font-semibold text-primary-foreground">{initials}</span>
      ) : (
        <User className="w-10 h-10 text-primary-foreground" aria-hidden="true" />
      )}
    </div>
  );
}

export function ProfileHeader({
  data,
  profileOverrides,
  userId,
  onMutationSuccess,
}: ProfileHeaderProps) {
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);

  // Merge profile overrides with extracted data (profile overrides take precedence if set)
  const displayName = profileOverrides?.name || data.name;
  const displayEmail = profileOverrides?.email ?? data.email;
  const displayPhone = profileOverrides?.phone ?? data.phone;
  const displayLocation = profileOverrides?.location ?? data.location;
  const displaySummary = profileOverrides?.summary ?? data.summary;

  const handleEditSuccess = () => {
    onMutationSuccess?.();
  };

  return (
    <>
      <div className="bg-card border rounded-lg p-6 sm:p-8">
        <div className="flex flex-col sm:flex-row sm:items-start gap-4 sm:gap-6">
          <AvatarPlaceholder name={displayName} />
          <div className="flex-1 min-w-0">
            <div className="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-2">
              <div>
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
              {userId && (
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={() => setIsEditDialogOpen(true)}
                  aria-label="Edit profile"
                  className="flex-shrink-0"
                >
                  <Pencil className="w-4 h-4" />
                </Button>
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
          onSuccess={handleEditSuccess}
        />
      )}
    </>
  );
}
