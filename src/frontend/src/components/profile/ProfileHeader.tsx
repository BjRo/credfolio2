"use client";

import { ChevronDown, ChevronUp, Mail, MapPin, Phone, User } from "lucide-react";
import { useState } from "react";
import type { ProfileData } from "./types";

interface ProfileHeaderProps {
  data: ProfileData;
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
      <h2 className="text-sm font-semibold text-muted-foreground uppercase tracking-wide mb-2">
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

export function ProfileHeader({ data }: ProfileHeaderProps) {
  return (
    <div className="bg-card shadow rounded-lg p-6 sm:p-8">
      <div className="flex flex-col sm:flex-row sm:items-start gap-4 sm:gap-6">
        <AvatarPlaceholder name={data.name} />
        <div className="flex-1 min-w-0">
          <div className="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-2">
            <div>
              <h1 className="text-2xl sm:text-3xl font-bold text-foreground truncate">
                {data.name}
              </h1>
              <div className="mt-2 space-y-1 text-muted-foreground">
                {data.email && (
                  <p className="flex items-center gap-2 text-sm sm:text-base">
                    <Mail className="w-4 h-4 flex-shrink-0" aria-hidden="true" />
                    <a href={`mailto:${data.email}`} className="hover:text-primary truncate">
                      {data.email}
                    </a>
                  </p>
                )}
                {data.phone && (
                  <p className="flex items-center gap-2 text-sm sm:text-base">
                    <Phone className="w-4 h-4 flex-shrink-0" aria-hidden="true" />
                    <a href={`tel:${data.phone}`} className="hover:text-primary">
                      {data.phone}
                    </a>
                  </p>
                )}
                {data.location && (
                  <p className="flex items-center gap-2 text-sm sm:text-base">
                    <MapPin className="w-4 h-4 flex-shrink-0" aria-hidden="true" />
                    <span>{data.location}</span>
                  </p>
                )}
              </div>
            </div>
            <div className="text-sm text-muted-foreground sm:text-right">
              <p>Confidence: {Math.round(data.confidence * 100)}%</p>
            </div>
          </div>
        </div>
      </div>

      {data.summary && <ProfileSummary summary={data.summary} />}
    </div>
  );
}
