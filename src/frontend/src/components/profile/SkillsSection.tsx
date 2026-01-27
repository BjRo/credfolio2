"use client";

import { MoreVertical, Pencil, Plus, Trash2 } from "lucide-react";
import { useState } from "react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { SkillCategory } from "@/graphql/generated/graphql";
import { DeleteSkillDialog } from "./DeleteSkillDialog";
import { SkillFormDialog } from "./SkillFormDialog";
import type { ProfileSkill } from "./types";

const CATEGORY_LABELS: Record<SkillCategory, string> = {
  [SkillCategory.Technical]: "Technical",
  [SkillCategory.Soft]: "Soft Skills",
  [SkillCategory.Domain]: "Domain",
};

const CATEGORY_ORDER: SkillCategory[] = [
  SkillCategory.Technical,
  SkillCategory.Soft,
  SkillCategory.Domain,
];

interface EditableSkillTagProps {
  skill: ProfileSkill;
  onEdit?: () => void;
  onDelete?: () => void;
}

function EditableSkillTag({ skill, onEdit, onDelete }: EditableSkillTagProps) {
  if (!onEdit && !onDelete) {
    return (
      <span className="inline-flex items-center px-3 py-1.5 bg-primary/10 text-primary rounded-full text-sm font-medium border border-primary/20">
        {skill.name}
      </span>
    );
  }

  return (
    <span className="inline-flex items-center gap-1 pl-3 pr-1 py-1 bg-primary/10 text-primary rounded-full text-sm font-medium border border-primary/20">
      {skill.name}
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <button
            type="button"
            className="p-0.5 text-primary/60 hover:text-primary hover:bg-primary/10 rounded-full transition-colors"
            aria-label={`Actions for ${skill.name}`}
          >
            <MoreVertical className="h-3.5 w-3.5" />
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
              className="text-destructive focus:text-destructive focus:bg-destructive/10"
            >
              <Trash2 className="h-4 w-4" />
              Delete
            </DropdownMenuItem>
          )}
        </DropdownMenuContent>
      </DropdownMenu>
    </span>
  );
}

interface ReadOnlySkillTagProps {
  skill: string;
}

function ReadOnlySkillTag({ skill }: ReadOnlySkillTagProps) {
  return (
    <span className="inline-flex items-center px-3 py-1.5 bg-primary/10 text-primary rounded-full text-sm font-medium border border-primary/20">
      {skill}
    </span>
  );
}

interface SkillsSectionProps {
  profileSkills?: ProfileSkill[];
  extractedSkills?: string[];
  userId?: string;
  onMutationSuccess?: () => void;
}

export function SkillsSection({
  profileSkills = [],
  extractedSkills = [],
  userId,
  onMutationSuccess,
}: SkillsSectionProps) {
  const [formDialogOpen, setFormDialogOpen] = useState(false);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [selectedSkill, setSelectedSkill] = useState<ProfileSkill | null>(null);
  const [defaultCategory, setDefaultCategory] = useState<SkillCategory>(SkillCategory.Technical);

  const isEditable = !!userId;
  const hasProfileSkills = profileSkills.length > 0;
  const hasExtractedSkills = extractedSkills.length > 0;

  const handleEdit = (skill: ProfileSkill) => {
    setSelectedSkill(skill);
    setFormDialogOpen(true);
  };

  const handleDelete = (skill: ProfileSkill) => {
    setSelectedSkill(skill);
    setDeleteDialogOpen(true);
  };

  const handleAddNew = (category?: SkillCategory) => {
    setSelectedSkill(null);
    setDefaultCategory(category ?? SkillCategory.Technical);
    setFormDialogOpen(true);
  };

  const handleFormDialogClose = (open: boolean) => {
    setFormDialogOpen(open);
    if (!open) {
      setSelectedSkill(null);
    }
  };

  const handleDeleteDialogClose = (open: boolean) => {
    setDeleteDialogOpen(open);
    if (!open) {
      setSelectedSkill(null);
    }
  };

  const handleSuccess = () => {
    onMutationSuccess?.();
  };

  // Group profile skills by category
  const groupedSkills = CATEGORY_ORDER.reduce(
    (acc, category) => {
      const skills = profileSkills.filter((s) => s.category === category);
      if (skills.length > 0) {
        acc[category] = skills;
      }
      return acc;
    },
    {} as Record<SkillCategory, ProfileSkill[]>
  );

  // If there are no skills at all and not editable, don't render
  if (!hasProfileSkills && !hasExtractedSkills && !isEditable) {
    return null;
  }

  return (
    <div className="bg-card shadow rounded-lg p-6 sm:p-8">
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-xl font-bold text-foreground">Skills</h2>
        {isEditable && (
          <button
            type="button"
            onClick={() => handleAddNew()}
            className="p-1.5 text-muted-foreground hover:text-foreground hover:bg-muted rounded transition-colors"
            aria-label="Add skill"
          >
            <Plus className="h-5 w-5" />
          </button>
        )}
      </div>

      {/* Profile skills grouped by category */}
      {hasProfileSkills && (
        <div className="space-y-4">
          {CATEGORY_ORDER.map((category) => {
            const skills = groupedSkills[category];
            if (!skills) return null;

            return (
              <div key={category}>
                <h3 className="text-sm font-medium text-muted-foreground mb-2">
                  {CATEGORY_LABELS[category]}
                </h3>
                <ul
                  className="flex flex-wrap gap-2 list-none"
                  aria-label={`${CATEGORY_LABELS[category]} skills`}
                >
                  {skills.map((skill) => (
                    <li key={skill.id}>
                      <EditableSkillTag
                        skill={skill}
                        onEdit={isEditable ? () => handleEdit(skill) : undefined}
                        onDelete={isEditable ? () => handleDelete(skill) : undefined}
                      />
                    </li>
                  ))}
                </ul>
              </div>
            );
          })}
        </div>
      )}

      {/* Extracted skills (read-only, shown when no profile skills exist) */}
      {!hasProfileSkills && hasExtractedSkills && (
        <ul className="flex flex-wrap gap-2 list-none" aria-label="Skills list">
          {extractedSkills.map((skill) => (
            <li key={skill}>
              <ReadOnlySkillTag skill={skill} />
            </li>
          ))}
        </ul>
      )}

      {/* Empty state for editable with no skills */}
      {!hasProfileSkills && !hasExtractedSkills && isEditable && (
        <p className="text-muted-foreground text-center py-8">
          No skills yet. Click the + button to add your first skill.
        </p>
      )}

      {isEditable && userId && (
        <>
          <SkillFormDialog
            key={selectedSkill?.id ?? "new"}
            open={formDialogOpen}
            onOpenChange={handleFormDialogClose}
            userId={userId}
            skill={
              selectedSkill
                ? {
                    id: selectedSkill.id,
                    name: selectedSkill.name,
                    category: selectedSkill.category,
                  }
                : undefined
            }
            defaultCategory={defaultCategory}
            onSuccess={handleSuccess}
          />

          {selectedSkill && (
            <DeleteSkillDialog
              open={deleteDialogOpen}
              onOpenChange={handleDeleteDialogClose}
              skillId={selectedSkill.id}
              skillName={selectedSkill.name}
              onSuccess={handleSuccess}
            />
          )}
        </>
      )}
    </div>
  );
}
