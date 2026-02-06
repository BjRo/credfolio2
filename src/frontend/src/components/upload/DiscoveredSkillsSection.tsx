"use client";

import { Sparkles } from "lucide-react";
import { useMemo } from "react";
import { CheckboxCard } from "@/components/ui/checkbox-card";
import { SelectionControls } from "@/components/ui/selection-controls";
import { SkillCategory } from "@/graphql/generated/graphql";

export interface DiscoveredSkillItem {
  name: string;
  quote: string;
  category: SkillCategory;
}

const CATEGORY_LABELS: Record<SkillCategory, string> = {
  [SkillCategory.Technical]: "Technical",
  [SkillCategory.Soft]: "Soft Skills",
  [SkillCategory.Domain]: "Domain Knowledge",
};

const CATEGORY_ORDER: SkillCategory[] = [
  SkillCategory.Technical,
  SkillCategory.Soft,
  SkillCategory.Domain,
];

interface DiscoveredSkillsSectionProps {
  discoveredSkills: DiscoveredSkillItem[];
  selected: Map<string, { selected: boolean; category: SkillCategory }>;
  onToggle: (skillName: string) => void;
  onCategoryChange: (skillName: string, category: SkillCategory) => void;
  onSelectAll: () => void;
  onDeselectAll: () => void;
  disabled?: boolean;
  description?: string;
  unselectedClassName?: string;
}

export function DiscoveredSkillsSection({
  discoveredSkills,
  selected,
  onToggle,
  onCategoryChange,
  onSelectAll,
  onDeselectAll,
  disabled = false,
  description = "These skills were mentioned in the reference letter but aren\u2019t in your profile yet. Select any you want to add.",
  unselectedClassName,
}: DiscoveredSkillsSectionProps) {
  const selectedCount = [...selected.values()].filter((v) => v.selected).length;

  // Group skills by their current category (may be overridden by user)
  const grouped = useMemo(() => {
    const groups = new Map<SkillCategory, DiscoveredSkillItem[]>();
    for (const skill of discoveredSkills) {
      const category = selected.get(skill.name)?.category ?? skill.category;
      const group = groups.get(category) ?? [];
      group.push(skill);
      groups.set(category, group);
    }
    return groups;
  }, [discoveredSkills, selected]);

  return (
    <section className="space-y-4">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <Sparkles className="h-5 w-5 text-warning" />
          <h2 className="text-xl font-semibold text-foreground">Skills Your Reference Noticed</h2>
        </div>
        <SelectionControls
          selectedCount={selectedCount}
          totalCount={discoveredSkills.length}
          onSelectAll={onSelectAll}
          onDeselectAll={onDeselectAll}
          disabled={disabled}
        />
      </div>

      <p className="text-sm text-muted-foreground">{description}</p>

      {CATEGORY_ORDER.filter((cat) => grouped.has(cat)).map((cat) => (
        <div key={cat} className="space-y-3">
          <h3 className="text-sm font-medium text-muted-foreground uppercase tracking-wide">
            {CATEGORY_LABELS[cat]}
          </h3>
          {/* biome-ignore lint/a11y/useSemanticElements: Using role="group" for checkbox group semantics */}
          <div
            className="space-y-3"
            role="group"
            aria-label={`${CATEGORY_LABELS[cat]} discovered skills`}
          >
            {grouped.get(cat)?.map((skill) => {
              const entry = selected.get(skill.name);
              const isSelected = entry?.selected ?? false;
              const currentCategory = entry?.category ?? skill.category;
              return (
                <CheckboxCard
                  key={skill.name}
                  checked={isSelected}
                  onToggle={() => onToggle(skill.name)}
                  disabled={disabled}
                  selectedClassName="bg-warning/5 border-warning/50"
                  borderStyle="border-2 border-dashed"
                  unselectedClassName={unselectedClassName}
                >
                  <div className="flex items-start justify-between gap-3">
                    <div className="flex-1 min-w-0">
                      <p className="font-medium text-foreground">{skill.name}</p>
                      {skill.quote && (
                        <blockquote className="mt-2 pl-3 border-l-2 border-warning/30 text-sm text-muted-foreground italic">
                          &ldquo;{skill.quote}&rdquo;
                        </blockquote>
                      )}
                    </div>
                    <select
                      value={currentCategory}
                      onChange={(e) => {
                        e.stopPropagation();
                        onCategoryChange(skill.name, e.target.value as SkillCategory);
                      }}
                      onClick={(e) => e.stopPropagation()}
                      onKeyDown={(e) => e.stopPropagation()}
                      disabled={disabled}
                      className="text-xs px-2 py-1 rounded border border-border bg-background text-foreground shrink-0"
                      aria-label={`Category for ${skill.name}`}
                    >
                      {CATEGORY_ORDER.map((c) => (
                        <option key={c} value={c}>
                          {CATEGORY_LABELS[c]}
                        </option>
                      ))}
                    </select>
                  </div>
                </CheckboxCard>
              );
            })}
          </div>
        </div>
      ))}
    </section>
  );
}
