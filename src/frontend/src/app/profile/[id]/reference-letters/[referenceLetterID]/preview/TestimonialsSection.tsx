"use client";

import { Quote } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Checkbox } from "@/components/ui/checkbox";
import type { TestimonialItem } from "./page";
import { SelectionControls } from "./SelectionControls";

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
          // biome-ignore lint/a11y/useSemanticElements: Custom styled checkbox card with inner Checkbox component
          <div
            key={`testimonial-${testimonial.quote.slice(0, 50)}-${index}`}
            role="checkbox"
            aria-checked={selectedTestimonials.has(index)}
            tabIndex={disabled ? -1 : 0}
            onClick={() => !disabled && onToggle(index)}
            onKeyDown={(e) => {
              if (!disabled && (e.key === " " || e.key === "Enter")) {
                e.preventDefault();
                onToggle(index);
              }
            }}
            className={`flex items-start gap-3 p-4 rounded-lg border transition-colors w-full text-left ${
              disabled
                ? "bg-muted cursor-not-allowed"
                : selectedTestimonials.has(index)
                  ? "bg-primary/5 border-primary/30 cursor-pointer"
                  : "bg-card border-border hover:bg-muted/50 cursor-pointer"
            }`}
          >
            <Checkbox
              checked={selectedTestimonials.has(index)}
              onCheckedChange={() => onToggle(index)}
              onClick={(e) => e.stopPropagation()}
              disabled={disabled}
              className="mt-1"
              tabIndex={-1}
              aria-hidden="true"
            />
            <div className="flex-1 space-y-3">
              <blockquote className="text-foreground">&ldquo;{testimonial.quote}&rdquo;</blockquote>
              <div className="flex flex-col gap-2">
                <p className="text-sm text-muted-foreground">
                  &mdash; {authorName}
                  {authorAttribution && `, ${authorAttribution}`}
                </p>
                {testimonial.skillsMentioned.length > 0 && (
                  <div className="flex flex-wrap gap-1">
                    <span className="text-xs text-muted-foreground mr-1">Skills validated:</span>
                    {testimonial.skillsMentioned.map((skill) => (
                      <Badge key={skill} variant="secondary" className="text-xs">
                        {skill}
                      </Badge>
                    ))}
                  </div>
                )}
              </div>
            </div>
          </div>
        ))}
      </div>
    </section>
  );
}
