"use client";

import {
  Briefcase,
  ChevronDown,
  ChevronUp,
  MapPin,
  MoreVertical,
  Pencil,
  Plus,
  Trash2,
} from "lucide-react";
import { useState } from "react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  calculateDurationMonths,
  calculateTotalTenure,
  formatDate,
  formatDuration,
} from "@/lib/utils";
import { DeleteExperienceDialog } from "./DeleteExperienceDialog";
import type { ProfileExperience } from "./types";
import { WorkExperienceFormDialog } from "./WorkExperienceFormDialog";

const DESCRIPTION_COLLAPSE_THRESHOLD = 150;

// Company group containing one or more roles
interface CompanyGroup {
  company: string;
  location?: string | null;
  totalTenureMonths: number | null;
  hasCurrentRole: boolean;
  roles: ProfileExperience[];
}

/**
 * Group experiences by company name (case-insensitive).
 * Roles within each group are sorted by start date (newest first).
 */
function groupExperiencesByCompany(experiences: ProfileExperience[]): CompanyGroup[] {
  const groupMap = new Map<string, ProfileExperience[]>();

  for (const exp of experiences) {
    const key = exp.company.toLowerCase().trim();
    const existing = groupMap.get(key) || [];
    existing.push(exp);
    groupMap.set(key, existing);
  }

  const groups: CompanyGroup[] = [];

  for (const roles of groupMap.values()) {
    // Sort by start date descending (newest first)
    // Roles without dates go to the end
    roles.sort((a, b) => {
      if (!a.startDate && !b.startDate) return 0;
      if (!a.startDate) return 1;
      if (!b.startDate) return -1;
      return b.startDate.localeCompare(a.startDate);
    });

    const company = roles[0].company; // Use original casing from first role
    const location = roles[0].location; // Use location from most recent role
    const hasCurrentRole = roles.some((r) => r.isCurrent);
    const totalTenureMonths = calculateTotalTenure(roles);

    groups.push({
      company,
      location,
      totalTenureMonths,
      hasCurrentRole,
      roles,
    });
  }

  // Sort groups by most recent role's start date
  groups.sort((a, b) => {
    const aDate = a.roles[0]?.startDate;
    const bDate = b.roles[0]?.startDate;
    if (!aDate && !bDate) return 0;
    if (!aDate) return 1;
    if (!bDate) return -1;
    return bDate.localeCompare(aDate);
  });

  return groups;
}

// Reusable action menu for edit/delete
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
          className="p-1.5 text-muted-foreground hover:text-foreground hover:bg-muted rounded transition-colors"
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
            className="text-destructive focus:text-destructive focus:bg-destructive/10"
          >
            <Trash2 className="h-4 w-4" />
            Delete
          </DropdownMenuItem>
        )}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}

// Flat experience card for single roles at a company (original style)
interface ExperienceCardProps {
  experience: ProfileExperience;
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

  const durationMonths = calculateDurationMonths(
    experience.startDate,
    experience.endDate,
    experience.isCurrent
  );
  const duration = durationMonths !== null ? formatDuration(durationMonths) : null;

  return (
    <div className={`group/card relative ${!isFirst ? "pt-6 border-t border-border" : ""}`}>
      {/* Green current indicator - vertically centered with title row */}
      {experience.isCurrent && (
        <span
          className={`absolute -left-3 w-1.5 h-1.5 bg-green-500 rounded-full ${
            isFirst ? "top-3" : "top-9"
          }`}
        />
      )}

      {/* Mobile: kebab menu positioned top-right - hidden until hover/focus */}
      {(onEdit || onDelete) && (
        <div
          className={`absolute right-0 sm:hidden ${isFirst ? "top-0" : "top-6"} opacity-0 group-hover/card:opacity-100 group-focus-within/card:opacity-100 focus-within:opacity-100 transition-opacity`}
        >
          <ActionMenu onEdit={onEdit} onDelete={onDelete} />
        </div>
      )}

      <div className="flex sm:justify-between sm:items-start gap-1 sm:gap-4 pr-8 sm:pr-0">
        <div className="flex gap-3">
          <div className="hidden sm:flex w-10 h-10 rounded-lg bg-muted dark:border dark:border-border items-center justify-center flex-shrink-0">
            <Briefcase className="w-5 h-5 text-muted-foreground" aria-hidden="true" />
          </div>
          <div className="min-w-0">
            <h3 className="text-lg font-semibold text-foreground flex items-center gap-2">
              {experience.title}
              {experience.isCurrent && (
                <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400">
                  Current
                </span>
              )}
            </h3>
            <p className="text-foreground">{experience.company}</p>
            {/* Date and duration - below company */}
            {dateRange && (
              <p className="text-sm text-muted-foreground">
                {dateRange}
                {duration && <span className="hidden sm:inline"> · {duration}</span>}
              </p>
            )}
            {/* Duration on mobile - separate line */}
            {duration && dateRange && (
              <p className="text-sm text-muted-foreground sm:hidden">{duration}</p>
            )}
            {experience.location && (
              <p className="text-sm text-muted-foreground flex items-center gap-1">
                <MapPin className="w-3 h-3" aria-hidden="true" />
                {experience.location}
              </p>
            )}
          </div>
        </div>
        {/* Desktop: kebab on right side - hidden until hover/focus */}
        <div className="hidden sm:flex items-center gap-1 flex-shrink-0 opacity-0 group-hover/card:opacity-100 group-focus-within/card:opacity-100 focus-within:opacity-100 transition-opacity">
          <ActionMenu onEdit={onEdit} onDelete={onDelete} />
        </div>
      </div>

      {experience.description && (
        <div className="mt-3 sm:ml-13">
          <p
            className={`text-foreground whitespace-pre-line ${
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
      )}
    </div>
  );
}

// Role card for displaying individual roles within a company group (multi-role)
interface RoleCardProps {
  role: ProfileExperience;
  isFirst: boolean;
  isLast: boolean;
  onEdit?: () => void;
  onDelete?: () => void;
}

function RoleCard({ role, isFirst, isLast, onEdit, onDelete }: RoleCardProps) {
  const [isExpanded, setIsExpanded] = useState(false);
  const hasLongDescription =
    role.description && role.description.length > DESCRIPTION_COLLAPSE_THRESHOLD;

  const startDate = formatDate(role.startDate);
  const endDate = role.isCurrent ? "Present" : formatDate(role.endDate) || "N/A";
  const dateRange = startDate ? `${startDate} - ${endDate}` : null;

  const durationMonths = calculateDurationMonths(role.startDate, role.endDate, role.isCurrent);
  const duration = durationMonths !== null ? formatDuration(durationMonths) : null;

  return (
    <div className={`group/card relative ${!isFirst ? "mt-4" : ""}`}>
      {/* Timeline connector - always shown for multi-role groups */}
      {/* Dot */}
      <span
        className={`absolute left-0 w-2 h-2 rounded-full ${
          role.isCurrent ? "bg-green-500" : "bg-muted-foreground/30"
        }`}
        style={{ top: "6px" }}
      />
      {/* Vertical line */}
      {!isLast && (
        <span
          className="absolute left-[3px] w-0.5 bg-border"
          style={{ top: "14px", bottom: "-16px" }}
        />
      )}

      <div className="pl-5">
        {/* Mobile: kebab menu positioned top-right - hidden until hover/focus */}
        {(onEdit || onDelete) && (
          <div className="absolute right-0 top-0 sm:hidden opacity-0 group-hover/card:opacity-100 group-focus-within/card:opacity-100 focus-within:opacity-100 transition-opacity">
            <ActionMenu onEdit={onEdit} onDelete={onDelete} />
          </div>
        )}

        <div className="flex sm:justify-between sm:items-start gap-1 sm:gap-4 pr-8 sm:pr-0">
          <div className="min-w-0 flex-1">
            <h4 className="text-base font-semibold text-foreground flex items-center gap-2 flex-wrap">
              {role.title}
              {role.isCurrent && (
                <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400">
                  Current
                </span>
              )}
            </h4>
            {/* Date and duration */}
            {dateRange && (
              <p className="text-sm text-muted-foreground">
                {dateRange}
                {duration && <span className="hidden sm:inline"> · {duration}</span>}
              </p>
            )}
            {/* Duration on mobile - separate line */}
            {duration && dateRange && (
              <p className="text-sm text-muted-foreground sm:hidden">{duration}</p>
            )}
          </div>
          {/* Desktop: kebab on right side - hidden until hover/focus */}
          <div className="hidden sm:flex items-center gap-1 flex-shrink-0 opacity-0 group-hover/card:opacity-100 group-focus-within/card:opacity-100 focus-within:opacity-100 transition-opacity">
            <ActionMenu onEdit={onEdit} onDelete={onDelete} />
          </div>
        </div>

        {role.description && (
          <div className="mt-2">
            <p
              className={`text-foreground text-sm whitespace-pre-line ${
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
              {role.description}
            </p>
            {hasLongDescription && (
              <button
                type="button"
                onClick={() => setIsExpanded(!isExpanded)}
                className="mt-1 text-sm text-primary hover:text-primary/80 flex items-center gap-1"
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
    </div>
  );
}

// Company group component showing company header and all roles
interface CompanyExperienceGroupProps {
  group: CompanyGroup;
  isFirst: boolean;
  isEditable: boolean;
  onEdit: (experience: ProfileExperience) => void;
  onDelete: (experience: ProfileExperience) => void;
}

function CompanyExperienceGroup({
  group,
  isFirst,
  isEditable,
  onEdit,
  onDelete,
}: CompanyExperienceGroupProps) {
  const totalDuration =
    group.totalTenureMonths !== null ? formatDuration(group.totalTenureMonths) : null;

  return (
    <div className={`relative ${!isFirst ? "pt-6 border-t border-border" : ""}`}>
      {/* Company header */}
      <div className="flex gap-3 mb-4">
        <div className="hidden sm:flex w-10 h-10 rounded-lg bg-muted dark:border dark:border-border items-center justify-center flex-shrink-0">
          <Briefcase className="w-5 h-5 text-muted-foreground" aria-hidden="true" />
        </div>
        <div className="min-w-0 flex-1">
          <h3 className="text-lg font-semibold text-foreground">{group.company}</h3>
          {totalDuration && <p className="text-sm text-muted-foreground">{totalDuration}</p>}
          {group.location && (
            <p className="text-sm text-muted-foreground flex items-center gap-1">
              <MapPin className="w-3 h-3" aria-hidden="true" />
              {group.location}
            </p>
          )}
        </div>
      </div>

      {/* Roles list with timeline */}
      <div className="sm:ml-13">
        {group.roles.map((role, index) => (
          <RoleCard
            key={role.id}
            role={role}
            isFirst={index === 0}
            isLast={index === group.roles.length - 1}
            onEdit={isEditable ? () => onEdit(role) : undefined}
            onDelete={isEditable ? () => onDelete(role) : undefined}
          />
        ))}
      </div>
    </div>
  );
}

interface WorkExperienceSectionProps {
  profileExperiences?: ProfileExperience[];
  userId?: string;
  onMutationSuccess?: () => void;
}

export function WorkExperienceSection({
  profileExperiences = [],
  userId,
  onMutationSuccess,
}: WorkExperienceSectionProps) {
  const [formDialogOpen, setFormDialogOpen] = useState(false);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [selectedExperience, setSelectedExperience] = useState<ProfileExperience | null>(null);

  const companyGroups = groupExperiencesByCompany(profileExperiences);

  const isEditable = !!userId;

  const handleEdit = (experience: ProfileExperience) => {
    setSelectedExperience(experience);
    setFormDialogOpen(true);
  };

  const handleDelete = (experience: ProfileExperience) => {
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

  if (companyGroups.length === 0 && !isEditable) {
    return null;
  }

  return (
    <div className="bg-card border rounded-lg p-6 sm:p-8">
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-xl font-bold text-foreground">Work Experience</h2>
        {isEditable && (
          <button
            type="button"
            onClick={handleAddNew}
            className="p-1.5 text-muted-foreground hover:text-foreground hover:bg-muted rounded transition-colors"
            aria-label="Add work experience"
          >
            <Plus className="h-5 w-5" />
          </button>
        )}
      </div>

      {companyGroups.length === 0 ? (
        <p className="text-muted-foreground text-center py-8">
          No work experience yet. Click the + button to add your first position.
        </p>
      ) : (
        <div className="space-y-6 relative">
          {companyGroups.map((group, index) =>
            group.roles.length > 1 ? (
              // Multi-role: show grouped view with company header and timeline
              <CompanyExperienceGroup
                key={`${group.company}-${index}`}
                group={group}
                isFirst={index === 0}
                isEditable={isEditable}
                onEdit={handleEdit}
                onDelete={handleDelete}
              />
            ) : (
              // Single role: show flat card
              <ExperienceCard
                key={group.roles[0].id}
                experience={group.roles[0]}
                isFirst={index === 0}
                onEdit={isEditable ? () => handleEdit(group.roles[0]) : undefined}
                onDelete={isEditable ? () => handleDelete(group.roles[0]) : undefined}
              />
            )
          )}
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
                    id: selectedExperience.id,
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
            onSuccess={handleSuccess}
          />

          {selectedExperience && (
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
