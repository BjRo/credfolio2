"use client";

import { Quote } from "lucide-react";
import { CheckboxCard } from "@/components/ui/checkbox-card";
import { SelectionControls } from "@/components/ui/selection-controls";
import type { TestimonialItem } from "./page";

interface TestimonialsSectionProps {
  testimonials: TestimonialItem[];
  selectedTestimonials: Set<number>;
  onToggle: (index: number) => void;
  onSelectAll: () => void;
  onDeselectAll: () => void;
  authorName: string;
  authorAttribution: string;
  disabled?: boolean;
}

export function TestimonialsSection({
  testimonials,
  selectedTestimonials,
  onToggle,
  onSelectAll,
  onDeselectAll,
  authorName,
  authorAttribution,
  disabled = false,
}: TestimonialsSectionProps) {
  return (
    <section className="space-y-4">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <Quote className="h-5 w-5 text-primary" />
          <h2 className="text-xl font-semibold text-foreground">Testimonials to Add</h2>
        </div>
        <SelectionControls
          selectedCount={selectedTestimonials.size}
          totalCount={testimonials.length}
          onSelectAll={onSelectAll}
          onDeselectAll={onDeselectAll}
          disabled={disabled}
        />
      </div>

      {/* biome-ignore lint/a11y/useSemanticElements: Using role="group" for checkbox group semantics */}
      <div className="space-y-4" role="group" aria-label="Testimonials">
        {testimonials.map((testimonial, index) => (
          <CheckboxCard
            key={`testimonial-${testimonial.quote.slice(0, 50)}-${index}`}
            checked={selectedTestimonials.has(index)}
            onToggle={() => onToggle(index)}
            disabled={disabled}
            selectedClassName="bg-primary/5 border-primary/30"
          >
            <div className="space-y-3">
              <blockquote className="text-foreground">&ldquo;{testimonial.quote}&rdquo;</blockquote>
              <div className="flex flex-col gap-2">
                <p className="text-sm text-muted-foreground">
                  &mdash; {authorName}
                  {authorAttribution && `, ${authorAttribution}`}
                </p>
              </div>
            </div>
          </CheckboxCard>
        ))}
      </div>
    </section>
  );
}
