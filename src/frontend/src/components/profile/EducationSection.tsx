"use client";

import { GraduationCap, MoreVertical, Pencil, Plus, Trash2 } from "lucide-react";
import { useState } from "react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { formatDate } from "@/lib/utils";
import { DeleteEducationDialog } from "./DeleteEducationDialog";
import { EducationFormDialog } from "./EducationFormDialog";
import type { ProfileEducation } from "./types";

interface ActionMenuProps {
  onEdit?: () => void;
  onDelete?: () => void;
}

function ActionMenu({ onEdit, onDelete }: ActionMenuProps) {
  if (!onEdit && !onDelete) return null;

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <button
          type="button"
          className="p-1.5 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded transition-colors"
          aria-label="More actions"
        >
          <MoreVertical className="h-4 w-4" />
        </button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        {onEdit && (
          <DropdownMenuItem onClick={onEdit}>
            <Pencil className="h-4 w-4" />
            Edit
          </DropdownMenuItem>
        )}
        {onDelete && (
          <DropdownMenuItem
            onClick={onDelete}
            className="text-red-600 focus:text-red-600 focus:bg-red-50"
          >
            <Trash2 className="h-4 w-4" />
            Delete
          </DropdownMenuItem>
        )}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}

interface EducationCardProps {
  education: ProfileEducation;
  isFirst: boolean;
  onEdit?: () => void;
  onDelete?: () => void;
}

function EducationCard({ education, isFirst, onEdit, onDelete }: EducationCardProps) {
  const startDate = formatDate(education.startDate);
  const endDate = education.isCurrent ? "Present" : formatDate(education.endDate) || "Present";
  const dateRange = startDate ? `${startDate} - ${endDate}` : null;

  const degreeField = [education.degree, education.field].filter(Boolean).join(" in ");

  // Validate GPA - should be numeric, not a date or garbage
  const isValidGpa = education.gpa && /^[\d./]+$/.test(education.gpa.trim());

  return (
    <div className={`relative ${!isFirst ? "pt-6 border-t border-gray-200" : ""}`}>
      {/* Mobile: kebab menu positioned top-right */}
      {(onEdit || onDelete) && (
        <div className={`absolute right-0 sm:hidden ${isFirst ? "top-0" : "top-6"}`}>
          <ActionMenu onEdit={onEdit} onDelete={onDelete} />
        </div>
      )}

      <div className="flex sm:justify-between sm:items-start gap-1 sm:gap-4 pr-8 sm:pr-0">
        <div className="flex gap-3">
          <div className="hidden sm:flex w-10 h-10 rounded-lg bg-gray-100 items-center justify-center flex-shrink-0">
            <GraduationCap className="w-5 h-5 text-gray-500" aria-hidden="true" />
          </div>
          <div className="min-w-0">
            <h3 className="text-lg font-semibold text-gray-900 flex items-center gap-2">
              {education.institution}
              {education.isCurrent && (
                <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-green-100 text-green-800">
                  Current
                </span>
              )}
            </h3>
            {degreeField && <p className="text-gray-700">{degreeField}</p>}
            {isValidGpa && (
              <p className="text-sm text-gray-500">
                <span className="font-medium">GPA:</span> {education.gpa}
              </p>
            )}
            {dateRange && <p className="text-sm text-gray-500">{dateRange}</p>}
          </div>
        </div>
        {/* Desktop: kebab on right side */}
        <div className="hidden sm:flex items-center gap-1 flex-shrink-0">
          <ActionMenu onEdit={onEdit} onDelete={onDelete} />
        </div>
      </div>
      {education.description && (
        <p className="mt-3 text-gray-600 sm:ml-13">{education.description}</p>
      )}
    </div>
  );
}

interface EducationSectionProps {
  profileEducations?: ProfileEducation[];
  userId?: string;
  onMutationSuccess?: () => void;
}

export function EducationSection({
  profileEducations = [],
  userId,
  onMutationSuccess,
}: EducationSectionProps) {
  const [formDialogOpen, setFormDialogOpen] = useState(false);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [selectedEducation, setSelectedEducation] = useState<ProfileEducation | null>(null);

  const isEditable = !!userId;

  const handleEdit = (edu: ProfileEducation) => {
    setSelectedEducation(edu);
    setFormDialogOpen(true);
  };

  const handleDelete = (edu: ProfileEducation) => {
    setSelectedEducation(edu);
    setDeleteDialogOpen(true);
  };

  const handleAddNew = () => {
    setSelectedEducation(null);
    setFormDialogOpen(true);
  };

  const handleFormDialogClose = (open: boolean) => {
    setFormDialogOpen(open);
    if (!open) {
      setSelectedEducation(null);
    }
  };

  const handleDeleteDialogClose = (open: boolean) => {
    setDeleteDialogOpen(open);
    if (!open) {
      setSelectedEducation(null);
    }
  };

  const handleSuccess = () => {
    onMutationSuccess?.();
  };

  if (profileEducations.length === 0 && !isEditable) {
    return null;
  }

  return (
    <div className="bg-white shadow rounded-lg p-6 sm:p-8">
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-xl font-bold text-gray-900">Education</h2>
        {isEditable && (
          <button
            type="button"
            onClick={handleAddNew}
            className="p-1.5 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded transition-colors"
            aria-label="Add education"
          >
            <Plus className="h-5 w-5" />
          </button>
        )}
      </div>

      {profileEducations.length === 0 ? (
        <p className="text-gray-500 text-center py-8">
          No education entries yet. Click the + button to add your first entry.
        </p>
      ) : (
        <div className="space-y-6">
          {profileEducations.map((edu, index) => (
            <EducationCard
              key={edu.id}
              education={edu}
              isFirst={index === 0}
              onEdit={isEditable ? () => handleEdit(edu) : undefined}
              onDelete={isEditable ? () => handleDelete(edu) : undefined}
            />
          ))}
        </div>
      )}

      {isEditable && userId && (
        <>
          <EducationFormDialog
            key={selectedEducation?.id ?? "new"}
            open={formDialogOpen}
            onOpenChange={handleFormDialogClose}
            userId={userId}
            education={
              selectedEducation
                ? {
                    id: selectedEducation.id,
                    institution: selectedEducation.institution,
                    degree: selectedEducation.degree,
                    field: selectedEducation.field,
                    startDate: selectedEducation.startDate,
                    endDate: selectedEducation.endDate,
                    isCurrent: selectedEducation.isCurrent,
                    description: selectedEducation.description,
                    gpa: selectedEducation.gpa,
                  }
                : undefined
            }
            onSuccess={handleSuccess}
          />

          {selectedEducation && (
            <DeleteEducationDialog
              open={deleteDialogOpen}
              onOpenChange={handleDeleteDialogClose}
              educationId={selectedEducation.id}
              degree={selectedEducation.degree ?? "Education"}
              institutionName={selectedEducation.institution}
              onSuccess={handleSuccess}
            />
          )}
        </>
      )}
    </div>
  );
}
