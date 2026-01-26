"use client";

import { Briefcase, ChevronDown, ChevronUp, MapPin, Pencil, Plus, Trash2 } from "lucide-react";
import { useState } from "react";
import { formatDate } from "@/lib/utils";
import { DeleteExperienceDialog } from "./DeleteExperienceDialog";
import type { ProfileExperience, WorkExperience } from "./types";
import { WorkExperienceFormDialog } from "./WorkExperienceFormDialog";

const DESCRIPTION_COLLAPSE_THRESHOLD = 150;

// Unified experience type that can hold both resume-extracted and profile data
interface ExperienceItem {
  id?: string; // Only profile experiences have IDs
  company: string;
  title: string;
  location?: string | null;
  startDate?: string | null;
  endDate?: string | null;
  isCurrent: boolean;
  description?: string | null;
  highlights?: string[];
}

interface ExperienceCardProps {
  experience: ExperienceItem;
  isFirst: boolean;
  onEdit?: () => void;
  onDelete?: () => void;
}

function ExperienceCard({ experience, isFirst, onEdit, onDelete }: ExperienceCardProps) {
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

      {/* Edit/Delete buttons in top right */}
      {(onEdit || onDelete) && (
        <div className="absolute right-0 top-0 flex gap-1">
          {onEdit && (
            <button
              type="button"
              onClick={onEdit}
              className="p-1.5 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded transition-colors"
              aria-label="Edit experience"
            >
              <Pencil className="h-4 w-4" />
            </button>
          )}
          {onDelete && (
            <button
              type="button"
              onClick={onDelete}
              className="p-1.5 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded transition-colors"
              aria-label="Delete experience"
            >
              <Trash2 className="h-4 w-4" />
            </button>
          )}
        </div>
      )}

      <div className="flex flex-col sm:flex-row sm:justify-between sm:items-start gap-1 sm:gap-4 pr-16">
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
  // Resume-extracted experiences (read-only source)
  experiences?: WorkExperience[];
  // Profile experiences (editable, takes precedence when available)
  profileExperiences?: ProfileExperience[];
  // User ID for mutations (required for editing)
  userId?: string;
  // Callback after successful mutation
  onMutationSuccess?: () => void;
}

export function WorkExperienceSection({
  experiences = [],
  profileExperiences = [],
  userId,
  onMutationSuccess,
}: WorkExperienceSectionProps) {
  const [formDialogOpen, setFormDialogOpen] = useState(false);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [selectedExperience, setSelectedExperience] = useState<ExperienceItem | null>(null);

  // Use profile experiences if available, otherwise fall back to resume-extracted
  const hasProfileExperiences = profileExperiences.length > 0;
  const displayExperiences: ExperienceItem[] = hasProfileExperiences
    ? profileExperiences.map((exp) => ({
        id: exp.id,
        company: exp.company,
        title: exp.title,
        location: exp.location,
        startDate: exp.startDate,
        endDate: exp.endDate,
        isCurrent: exp.isCurrent,
        description: exp.description,
        highlights: exp.highlights,
      }))
    : experiences.map((exp) => ({
        company: exp.company,
        title: exp.title,
        location: exp.location,
        startDate: exp.startDate,
        endDate: exp.endDate,
        isCurrent: exp.isCurrent,
        description: exp.description,
      }));

  const isEditable = !!userId;

  const handleEdit = (experience: ExperienceItem) => {
    setSelectedExperience(experience);
    setFormDialogOpen(true);
  };

  const handleDelete = (experience: ExperienceItem) => {
    setSelectedExperience(experience);
    setDeleteDialogOpen(true);
  };

  const handleAddNew = () => {
    setSelectedExperience(null);
    setFormDialogOpen(true);
  };

  const handleFormDialogClose = (open: boolean) => {
    setFormDialogOpen(open);
    if (!open) {
      setSelectedExperience(null);
    }
  };

  const handleDeleteDialogClose = (open: boolean) => {
    setDeleteDialogOpen(open);
    if (!open) {
      setSelectedExperience(null);
    }
  };

  const handleSuccess = () => {
    onMutationSuccess?.();
  };

  if (displayExperiences.length === 0 && !isEditable) {
    return null;
  }

  return (
    <div className="bg-white shadow rounded-lg p-6 sm:p-8">
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-xl font-bold text-gray-900">Work Experience</h2>
        {isEditable && (
          <button
            type="button"
            onClick={handleAddNew}
            className="p-1.5 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded transition-colors"
            aria-label="Add work experience"
          >
            <Plus className="h-5 w-5" />
          </button>
        )}
      </div>

      {displayExperiences.length === 0 ? (
        <p className="text-gray-500 text-center py-8">
          No work experience yet. Click the + button to add your first position.
        </p>
      ) : (
        <div className="space-y-6 relative">
          {displayExperiences.map((exp, index) => (
            <ExperienceCard
              key={exp.id ?? `${exp.company}-${exp.title}-${index}`}
              experience={exp}
              isFirst={index === 0}
              onEdit={isEditable ? () => handleEdit(exp) : undefined}
              onDelete={isEditable && exp.id ? () => handleDelete(exp) : undefined}
            />
          ))}
        </div>
      )}

      {isEditable && userId && (
        <>
          <WorkExperienceFormDialog
            key={selectedExperience?.id ?? "new"}
            open={formDialogOpen}
            onOpenChange={handleFormDialogClose}
            userId={userId}
            experience={
              selectedExperience
                ? {
                    id: selectedExperience.id ?? "",
                    company: selectedExperience.company,
                    title: selectedExperience.title,
                    location: selectedExperience.location,
                    startDate: selectedExperience.startDate,
                    endDate: selectedExperience.endDate,
                    isCurrent: selectedExperience.isCurrent,
                    description: selectedExperience.description,
                    highlights: selectedExperience.highlights ?? [],
                  }
                : undefined
            }
            // If editing a resume-extracted experience (no ID), treat as create
            mode={selectedExperience && !selectedExperience.id ? "create" : undefined}
            onSuccess={handleSuccess}
          />

          {selectedExperience?.id && (
            <DeleteExperienceDialog
              open={deleteDialogOpen}
              onOpenChange={handleDeleteDialogClose}
              experienceId={selectedExperience.id}
              experienceTitle={selectedExperience.title}
              companyName={selectedExperience.company}
              onSuccess={handleSuccess}
            />
          )}
        </>
      )}
    </div>
  );
}

// Keep for backwards compatibility but mark as deprecated
/** @deprecated Use WorkExperienceSection with profileExperiences prop instead */
export function EditableWorkExperienceSection({
  experiences,
  userId,
  onMutationSuccess,
}: {
  experiences: ProfileExperience[];
  userId: string;
  onMutationSuccess?: () => void;
}) {
  return (
    <WorkExperienceSection
      profileExperiences={experiences}
      userId={userId}
      onMutationSuccess={onMutationSuccess}
    />
  );
}
