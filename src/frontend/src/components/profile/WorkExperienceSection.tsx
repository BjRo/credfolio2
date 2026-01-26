"use client";

import { Briefcase, ChevronDown, ChevronUp, MapPin } from "lucide-react";
import { useState } from "react";
import { formatDate } from "@/lib/utils";
import type { WorkExperience } from "./types";

interface ExperienceCardProps {
  experience: WorkExperience;
  isFirst: boolean;
}

const DESCRIPTION_COLLAPSE_THRESHOLD = 150;

function ExperienceCard({ experience, isFirst }: ExperienceCardProps) {
  const [isExpanded, setIsExpanded] = useState(false);
  const hasLongDescription =
    experience.description && experience.description.length > DESCRIPTION_COLLAPSE_THRESHOLD;

  const startDate = formatDate(experience.startDate);
  const endDate = experience.isCurrent ? "Present" : formatDate(experience.endDate) || "N/A";
  const dateRange = startDate ? `${startDate} - ${endDate}` : null;

  return (
    <div
      className={`${!isFirst ? "pt-6 border-t border-gray-200" : ""} ${
        experience.isCurrent ? "relative" : ""
      }`}
    >
      {experience.isCurrent && (
        <span className="absolute -left-3 top-6 w-1.5 h-1.5 bg-green-500 rounded-full" />
      )}
      <div className="flex flex-col sm:flex-row sm:justify-between sm:items-start gap-1 sm:gap-4">
        <div className="flex gap-3">
          <div className="hidden sm:flex w-10 h-10 rounded-lg bg-gray-100 items-center justify-center flex-shrink-0">
            <Briefcase className="w-5 h-5 text-gray-500" aria-hidden="true" />
          </div>
          <div className="min-w-0">
            <h3 className="text-lg font-semibold text-gray-900 flex items-center gap-2">
              {experience.title}
              {experience.isCurrent && (
                <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-green-100 text-green-800">
                  Current
                </span>
              )}
            </h3>
            <p className="text-gray-700">{experience.company}</p>
            {experience.location && (
              <p className="text-sm text-gray-500 flex items-center gap-1">
                <MapPin className="w-3 h-3" aria-hidden="true" />
                {experience.location}
              </p>
            )}
          </div>
        </div>
        <div className="text-sm text-gray-500 sm:text-right flex-shrink-0">
          {dateRange && <p>{dateRange}</p>}
        </div>
      </div>
      {experience.description && (
        <div className="mt-3 sm:ml-13">
          <p
            className={`text-gray-600 whitespace-pre-line ${
              !isExpanded && hasLongDescription ? "line-clamp-3" : ""
            }`}
            style={
              !isExpanded && hasLongDescription
                ? {
                    display: "-webkit-box",
                    WebkitLineClamp: 3,
                    WebkitBoxOrient: "vertical",
                    overflow: "hidden",
                  }
                : undefined
            }
          >
            {experience.description}
          </p>
          {hasLongDescription && (
            <button
              type="button"
              onClick={() => setIsExpanded(!isExpanded)}
              className="mt-2 text-sm text-blue-600 hover:text-blue-800 flex items-center gap-1"
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
      )}
    </div>
  );
}

interface WorkExperienceSectionProps {
  experience: WorkExperience[];
}

export function WorkExperienceSection({ experience }: WorkExperienceSectionProps) {
  if (experience.length === 0) {
    return null;
  }

  return (
    <div className="bg-white shadow rounded-lg p-6 sm:p-8">
      <h2 className="text-xl font-bold text-gray-900 mb-6">Work Experience</h2>
      <div className="space-y-6 relative">
        {experience.map((exp, index) => (
          <ExperienceCard
            key={`${exp.company}-${exp.title}-${index}`}
            experience={exp}
            isFirst={index === 0}
          />
        ))}
      </div>
    </div>
  );
}
