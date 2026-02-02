"use client";

import { CheckCircle2, FileText, MessageSquareQuote, Plus, User } from "lucide-react";
import { useCallback } from "react";
import { type GetTestimonialsQuery, TestimonialRelationship } from "@/graphql/generated/graphql";

const RELATIONSHIP_LABELS: Record<TestimonialRelationship, string> = {
  [TestimonialRelationship.Manager]: "Manager",
  [TestimonialRelationship.Peer]: "Peer",
  [TestimonialRelationship.DirectReport]: "Direct Report",
  [TestimonialRelationship.Client]: "Client",
  [TestimonialRelationship.Other]: "Other",
};

type Testimonial = GetTestimonialsQuery["testimonials"][number];

interface TestimonialCardProps {
  testimonial: Testimonial;
  onSkillClick?: (skillId: string) => void;
}

function TestimonialCard({ testimonial, onSkillClick }: TestimonialCardProps) {
  const handleSkillClick = useCallback(
    (skillId: string) => {
      if (onSkillClick) {
        onSkillClick(skillId);
      } else {
        // Default behavior: scroll to the skill element
        const element = document.getElementById(`skill-${skillId}`);
        if (element) {
          element.scrollIntoView({ behavior: "smooth", block: "center" });
          // Briefly highlight the element
          element.classList.add("ring-2", "ring-primary", "ring-offset-2");
          setTimeout(() => {
            element.classList.remove("ring-2", "ring-primary", "ring-offset-2");
          }, 2000);
        }
      }
    },
    [onSkillClick]
  );

  // Get the source PDF URL from the reference letter
  const sourceUrl = testimonial.referenceLetter?.file?.url;

  return (
    <div className="bg-muted/30 rounded-lg p-6 border border-border/50 relative">
      {/* Source Badge - shows link to original reference letter PDF */}
      {sourceUrl && (
        <a
          href={sourceUrl}
          target="_blank"
          rel="noopener noreferrer"
          className="absolute top-4 right-4 flex items-center gap-1 text-xs text-muted-foreground hover:text-primary transition-colors"
          title="View source document"
        >
          <FileText className="h-4 w-4" />
          <span className="sr-only">View source</span>
        </a>
      )}

      {/* Quote */}
      <blockquote className="relative">
        <span className="absolute -top-2 -left-1 text-4xl text-primary/20 font-serif">&ldquo;</span>
        <p className="text-foreground pl-4 pr-2 italic leading-relaxed">{testimonial.quote}</p>
        <span className="absolute -bottom-4 right-0 text-4xl text-primary/20 font-serif">
          &rdquo;
        </span>
      </blockquote>

      {/* Attribution */}
      <div className="mt-6 pt-4 border-t border-border/50">
        <div className="flex items-start gap-3">
          <div className="flex-shrink-0 w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center">
            <User className="h-5 w-5 text-primary" />
          </div>
          <div className="flex-1 min-w-0">
            <p className="font-semibold text-foreground">{testimonial.authorName}</p>
            {(testimonial.authorTitle || testimonial.authorCompany) && (
              <p className="text-sm text-muted-foreground">
                {testimonial.authorTitle}
                {testimonial.authorTitle && testimonial.authorCompany && " at "}
                {testimonial.authorCompany}
              </p>
            )}
            <span className="inline-flex items-center mt-1 px-2 py-0.5 text-xs font-medium bg-secondary text-secondary-foreground rounded-full">
              {RELATIONSHIP_LABELS[testimonial.relationship]}
            </span>
          </div>
        </div>

        {/* Validated Skills */}
        {testimonial.validatedSkills && testimonial.validatedSkills.length > 0 && (
          <div className="mt-4 pt-3 border-t border-border/30">
            <div className="flex items-center gap-2 flex-wrap">
              <span className="text-xs text-muted-foreground flex items-center gap-1">
                <CheckCircle2 className="h-3 w-3" />
                Validates:
              </span>
              {testimonial.validatedSkills.map((skill, index) => (
                <span key={skill.id} className="inline-flex items-center">
                  <button
                    type="button"
                    onClick={() => handleSkillClick(skill.id)}
                    className="text-xs text-primary hover:text-primary/80 hover:underline transition-colors"
                  >
                    {skill.name}
                  </button>
                  {index < testimonial.validatedSkills.length - 1 && (
                    <span className="text-muted-foreground mx-1">·</span>
                  )}
                </span>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

interface TestimonialsSectionProps {
  testimonials: Testimonial[];
  isLoading?: boolean;
  onAddReference?: () => void;
  onSkillClick?: (skillId: string) => void;
}

export function TestimonialsSection({
  testimonials,
  isLoading = false,
  onAddReference,
  onSkillClick,
}: TestimonialsSectionProps) {
  // Don't render if no testimonials and no way to add one
  if (testimonials.length === 0 && !onAddReference) {
    return null;
  }

  return (
    <div id="testimonials" className="bg-card border rounded-lg p-6 sm:p-8">
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-2">
          <MessageSquareQuote className="h-5 w-5 text-primary" />
          <h2 className="text-xl font-bold text-foreground">What Others Say</h2>
        </div>
        {onAddReference && (
          <button
            type="button"
            onClick={onAddReference}
            className="p-1.5 text-muted-foreground hover:text-foreground hover:bg-muted rounded transition-colors"
            aria-label="Add reference letter"
          >
            <Plus className="h-5 w-5" />
          </button>
        )}
      </div>

      {isLoading ? (
        <div className="space-y-4">
          {[1, 2].map((i) => (
            <div
              key={i}
              className="bg-muted/30 rounded-lg p-6 border border-border/50 animate-pulse"
            >
              <div className="h-20 bg-muted rounded" />
              <div className="mt-6 pt-4 border-t border-border/50">
                <div className="flex items-start gap-3">
                  <div className="w-10 h-10 rounded-full bg-muted" />
                  <div className="flex-1 space-y-2">
                    <div className="h-4 bg-muted rounded w-32" />
                    <div className="h-3 bg-muted rounded w-48" />
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>
      ) : testimonials.length > 0 ? (
        <div className="space-y-4">
          {testimonials.map((testimonial) => (
            <TestimonialCard
              key={testimonial.id}
              testimonial={testimonial}
              onSkillClick={onSkillClick}
            />
          ))}
        </div>
      ) : (
        <div className="text-center py-8">
          <MessageSquareQuote className="h-12 w-12 text-muted-foreground/50 mx-auto mb-4" />
          <p className="text-muted-foreground mb-4">No testimonials yet.</p>
          <p className="text-sm text-muted-foreground mb-6">
            Add a reference letter to include testimonials from people who&apos;ve worked with you.
          </p>
          {onAddReference && (
            <button
              type="button"
              onClick={onAddReference}
              className="inline-flex items-center gap-2 px-4 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90 transition-colors"
            >
              <Plus className="h-4 w-4" />
              Add Reference Letter
            </button>
          )}
        </div>
      )}
    </div>
  );
}
