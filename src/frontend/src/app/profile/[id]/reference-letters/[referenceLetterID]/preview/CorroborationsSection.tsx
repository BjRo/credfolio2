"use client";

import { CheckCircle2 } from "lucide-react";
import { CheckboxCard } from "@/components/ui/checkbox-card";
import { SelectionControls } from "@/components/ui/selection-controls";
import type { ExperienceCorroboration, SkillCorroboration } from "./page";

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
              <CheckboxCard
                key={`${corr.profileSkillId}-${index}`}
                checked={selectedSkillCorroborations.has(corr.profileSkillId)}
                onToggle={() => onSkillToggle(corr.profileSkillId)}
                disabled={disabled}
                selectedClassName="bg-success/5 border-success/30"
              >
                <p className="font-medium text-foreground">{corr.skillName}</p>
                <blockquote className="mt-2 pl-3 border-l-2 border-muted-foreground/20 text-sm text-muted-foreground italic">
                  &ldquo;{corr.quote}&rdquo;
                </blockquote>
              </CheckboxCard>
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
              <CheckboxCard
                key={`${corr.profileExperienceId}-${index}`}
                checked={selectedExperienceCorroborations.has(corr.profileExperienceId)}
                onToggle={() => onExperienceToggle(corr.profileExperienceId)}
                disabled={disabled}
                selectedClassName="bg-success/5 border-success/30"
              >
                <p className="font-medium text-foreground">
                  {corr.role} at {corr.company}
                </p>
                <blockquote className="mt-2 pl-3 border-l-2 border-muted-foreground/20 text-sm text-muted-foreground italic">
                  &ldquo;{corr.quote}&rdquo;
                </blockquote>
              </CheckboxCard>
            ))}
          </div>
        </div>
      )}
    </section>
  );
}
