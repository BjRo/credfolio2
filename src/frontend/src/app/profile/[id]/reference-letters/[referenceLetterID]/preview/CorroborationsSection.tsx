"use client";

import { CheckCircle2 } from "lucide-react";
import { Checkbox } from "@/components/ui/checkbox";
import type { ExperienceCorroboration, SkillCorroboration } from "./page";
import { SelectionControls } from "./SelectionControls";

interface CorroborationsSectionProps {
  skillCorroborations: SkillCorroboration[];
  experienceCorroborations: ExperienceCorroboration[];
  selectedSkillCorroborations: Set<string>;
  selectedExperienceCorroborations: Set<string>;
  onSkillToggle: (profileSkillId: string) => void;
  onExperienceToggle: (profileExperienceId: string) => void;
  onSelectAllSkills: () => void;
  onDeselectAllSkills: () => void;
  onSelectAllExperiences: () => void;
  onDeselectAllExperiences: () => void;
  disabled?: boolean;
}

export function CorroborationsSection({
  skillCorroborations,
  experienceCorroborations,
  selectedSkillCorroborations,
  selectedExperienceCorroborations,
  onSkillToggle,
  onExperienceToggle,
  onSelectAllSkills,
  onDeselectAllSkills,
  onSelectAllExperiences,
  onDeselectAllExperiences,
  disabled = false,
}: CorroborationsSectionProps) {
  return (
    <section className="space-y-6">
      <div className="flex items-center gap-2">
        <CheckCircle2 className="h-5 w-5 text-success" />
        <h2 className="text-xl font-semibold text-foreground">
          Skills & Experiences Your Reference Validates
        </h2>
      </div>

      {/* Skills Corroborations */}
      {skillCorroborations.length > 0 && (
        <div className="space-y-4">
          <div className="flex justify-between items-center">
            <h3 className="text-lg font-medium text-foreground">Skills</h3>
            <SelectionControls
              selectedCount={selectedSkillCorroborations.size}
              totalCount={skillCorroborations.length}
              onSelectAll={onSelectAllSkills}
              onDeselectAll={onDeselectAllSkills}
              disabled={disabled}
            />
          </div>

          {/* biome-ignore lint/a11y/useSemanticElements: Using role="group" for checkbox group semantics */}
          <div className="space-y-3" role="group" aria-label="Skills corroborations">
            {skillCorroborations.map((corr, index) => (
              // biome-ignore lint/a11y/useSemanticElements: Custom styled checkbox card with inner Checkbox component
              <div
                key={`${corr.profileSkillId}-${index}`}
                role="checkbox"
                aria-checked={selectedSkillCorroborations.has(corr.profileSkillId)}
                tabIndex={disabled ? -1 : 0}
                onClick={() => !disabled && onSkillToggle(corr.profileSkillId)}
                onKeyDown={(e) => {
                  if (!disabled && (e.key === " " || e.key === "Enter")) {
                    e.preventDefault();
                    onSkillToggle(corr.profileSkillId);
                  }
                }}
                className={`flex items-start gap-3 p-4 rounded-lg border transition-colors ${
                  disabled
                    ? "bg-muted cursor-not-allowed"
                    : selectedSkillCorroborations.has(corr.profileSkillId)
                      ? "bg-success/5 border-success/30 cursor-pointer"
                      : "bg-card border-border hover:bg-muted/50 cursor-pointer"
                }`}
              >
                <Checkbox
                  checked={selectedSkillCorroborations.has(corr.profileSkillId)}
                  onCheckedChange={() => onSkillToggle(corr.profileSkillId)}
                  onClick={(e) => e.stopPropagation()}
                  disabled={disabled}
                  className="mt-1"
                  tabIndex={-1}
                  aria-hidden="true"
                />
                <div className="flex-1">
                  <p className="font-medium text-foreground">{corr.skillName}</p>
                  <blockquote className="mt-2 pl-3 border-l-2 border-muted-foreground/20 text-sm text-muted-foreground italic">
                    &ldquo;{corr.quote}&rdquo;
                  </blockquote>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Experience Corroborations */}
      {experienceCorroborations.length > 0 && (
        <div className="space-y-4">
          <div className="flex justify-between items-center">
            <h3 className="text-lg font-medium text-foreground">Experiences</h3>
            <SelectionControls
              selectedCount={selectedExperienceCorroborations.size}
              totalCount={experienceCorroborations.length}
              onSelectAll={onSelectAllExperiences}
              onDeselectAll={onDeselectAllExperiences}
              disabled={disabled}
            />
          </div>

          {/* biome-ignore lint/a11y/useSemanticElements: Using role="group" for checkbox group semantics */}
          <div className="space-y-3" role="group" aria-label="Experience corroborations">
            {experienceCorroborations.map((corr, index) => (
              // biome-ignore lint/a11y/useSemanticElements: Custom styled checkbox card with inner Checkbox component
              <div
                key={`${corr.profileExperienceId}-${index}`}
                role="checkbox"
                aria-checked={selectedExperienceCorroborations.has(corr.profileExperienceId)}
                tabIndex={disabled ? -1 : 0}
                onClick={() => !disabled && onExperienceToggle(corr.profileExperienceId)}
                onKeyDown={(e) => {
                  if (!disabled && (e.key === " " || e.key === "Enter")) {
                    e.preventDefault();
                    onExperienceToggle(corr.profileExperienceId);
                  }
                }}
                className={`flex items-start gap-3 p-4 rounded-lg border transition-colors ${
                  disabled
                    ? "bg-muted cursor-not-allowed"
                    : selectedExperienceCorroborations.has(corr.profileExperienceId)
                      ? "bg-success/5 border-success/30 cursor-pointer"
                      : "bg-card border-border hover:bg-muted/50 cursor-pointer"
                }`}
              >
                <Checkbox
                  checked={selectedExperienceCorroborations.has(corr.profileExperienceId)}
                  onCheckedChange={() => onExperienceToggle(corr.profileExperienceId)}
                  onClick={(e) => e.stopPropagation()}
                  disabled={disabled}
                  className="mt-1"
                  tabIndex={-1}
                  aria-hidden="true"
                />
                <div className="flex-1">
                  <p className="font-medium text-foreground">
                    {corr.role} at {corr.company}
                  </p>
                  <blockquote className="mt-2 pl-3 border-l-2 border-muted-foreground/20 text-sm text-muted-foreground italic">
                    &ldquo;{corr.quote}&rdquo;
                  </blockquote>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </section>
  );
}
