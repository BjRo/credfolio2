"use client";

import { Sparkles } from "lucide-react";
import { CheckboxCard } from "@/components/ui/checkbox-card";
import { SelectionControls } from "@/components/ui/selection-controls";
import type { DiscoveredSkill } from "./page";

interface DiscoveredSkillsSectionProps {
  discoveredSkills: DiscoveredSkill[];
  selectedSkills: Set<string>;
  onToggle: (skillName: string) => void;
  onSelectAll: () => void;
  onDeselectAll: () => void;
  disabled?: boolean;
}

export function DiscoveredSkillsSection({
  discoveredSkills,
  selectedSkills,
  onToggle,
  onSelectAll,
  onDeselectAll,
  disabled = false,
}: DiscoveredSkillsSectionProps) {
  return (
    <section className="space-y-4">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <Sparkles className="h-5 w-5 text-warning" />
          <h2 className="text-xl font-semibold text-foreground">Skills Your Reference Noticed</h2>
        </div>
        <SelectionControls
          selectedCount={selectedSkills.size}
          totalCount={discoveredSkills.length}
          onSelectAll={onSelectAll}
          onDeselectAll={onDeselectAll}
          disabled={disabled}
        />
      </div>

      <p className="text-sm text-muted-foreground">
        These skills were mentioned in the reference letter but aren&apos;t in your profile yet.
        Select any you want to add.
      </p>

      {/* biome-ignore lint/a11y/useSemanticElements: Using role="group" for checkbox group semantics */}
      <div className="space-y-3" role="group" aria-label="Discovered skills">
        {discoveredSkills.map((skill) => (
          <CheckboxCard
            key={skill.name}
            checked={selectedSkills.has(skill.name)}
            onToggle={() => onToggle(skill.name)}
            disabled={disabled}
            selectedClassName="bg-warning/5 border-warning/50"
            borderStyle="border-2 border-dashed"
            unselectedClassName="bg-card border-warning/20 hover:bg-warning/5"
          >
            <p className="font-medium text-foreground">{skill.name}</p>
            {skill.quote && (
              <blockquote className="mt-2 pl-3 border-l-2 border-warning/30 text-sm text-muted-foreground italic">
                &ldquo;{skill.quote}&rdquo;
              </blockquote>
            )}
          </CheckboxCard>
        ))}
      </div>
    </section>
  );
}
