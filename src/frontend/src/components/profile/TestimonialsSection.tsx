"use client";

import {
  CheckCircle2,
  ChevronDown,
  ChevronUp,
  FileText,
  Linkedin,
  MessageSquareQuote,
  Plus,
  User,
} from "lucide-react";
import { useCallback, useMemo, useState } from "react";
import { type GetTestimonialsQuery, TestimonialRelationship } from "@/graphql/generated/graphql";

const RELATIONSHIP_LABELS: Record<TestimonialRelationship, string> = {
  [TestimonialRelationship.Manager]: "Manager",
  [TestimonialRelationship.Peer]: "Peer",
  [TestimonialRelationship.DirectReport]: "Direct Report",
  [TestimonialRelationship.Client]: "Client",
  [TestimonialRelationship.Other]: "Other",
};

// Number of quotes to show before collapsing
const COLLAPSE_THRESHOLD = 2;

type Testimonial = GetTestimonialsQuery["testimonials"][number];

// Author info extracted from testimonial (using author entity or legacy fields)
interface AuthorInfo {
  id: string;
  name: string;
  title: string | null | undefined;
  company: string | null | undefined;
  linkedInUrl: string | null | undefined;
}

// A group of testimonials from the same author
interface TestimonialGroup {
  author: AuthorInfo;
  relationship: TestimonialRelationship;
  testimonials: Testimonial[];
}

// Get author info from testimonial, preferring author entity over legacy fields
function getAuthorInfo(testimonial: Testimonial): AuthorInfo {
  if (testimonial.author) {
    return {
      id: testimonial.author.id,
      name: testimonial.author.name,
      title: testimonial.author.title,
      company: testimonial.author.company,
      linkedInUrl: testimonial.author.linkedInUrl,
    };
  }
  // Fallback to legacy fields
  return {
    id: `legacy-${testimonial.authorName}`,
    name: testimonial.authorName,
    title: testimonial.authorTitle,
    company: testimonial.authorCompany,
    linkedInUrl: null,
  };
}

// Group testimonials by author
function groupTestimonialsByAuthor(testimonials: Testimonial[]): TestimonialGroup[] {
  const groups = new Map<string, TestimonialGroup>();

  for (const testimonial of testimonials) {
    const author = getAuthorInfo(testimonial);
    const existingGroup = groups.get(author.id);

    if (existingGroup) {
      existingGroup.testimonials.push(testimonial);
    } else {
      groups.set(author.id, {
        author,
        relationship: testimonial.relationship,
        testimonials: [testimonial],
      });
    }
  }

  return Array.from(groups.values());
}

interface QuoteItemProps {
  testimonial: Testimonial;
  onSkillClick?: (skillId: string) => void;
}

function QuoteItem({ testimonial, onSkillClick }: QuoteItemProps) {
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
    <div className="pl-4 border-l-2 border-primary/20">
      {/* Quote */}
      <blockquote className="relative">
        <span className="absolute -top-1 -left-1 text-2xl text-primary/20 font-serif">&ldquo;</span>
        <p className="text-foreground pl-4 pr-2 italic leading-relaxed text-sm">
          {testimonial.quote}
          <span className="text-2xl text-primary/20 font-serif leading-none align-bottom">
            &rdquo;
          </span>
        </p>
      </blockquote>

      {/* Source badge and validated skills */}
      <div className="mt-2 flex items-center gap-3 flex-wrap">
        {sourceUrl && (
          <a
            href={sourceUrl}
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center gap-1 px-2 py-0.5 text-xs text-muted-foreground hover:text-primary transition-colors"
            title="View original reference letter"
          >
            <FileText className="h-3 w-3" />
            <span>Source</span>
          </a>
        )}

        {/* Validated Skills */}
        {testimonial.validatedSkills && testimonial.validatedSkills.length > 0 && (
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
        )}
      </div>
    </div>
  );
}

interface TestimonialGroupCardProps {
  group: TestimonialGroup;
  onSkillClick?: (skillId: string) => void;
}

function TestimonialGroupCard({ group, onSkillClick }: TestimonialGroupCardProps) {
  const { author, relationship, testimonials } = group;
  const [isExpanded, setIsExpanded] = useState(testimonials.length <= COLLAPSE_THRESHOLD);

  const visibleTestimonials = isExpanded ? testimonials : testimonials.slice(0, COLLAPSE_THRESHOLD);
  const hiddenCount = testimonials.length - COLLAPSE_THRESHOLD;

  return (
    <div className="bg-muted/30 rounded-lg p-6 border border-border/50">
      {/* Author header */}
      <div className="flex items-start gap-3 mb-4">
        <div className="flex-shrink-0 w-10 h-10 rounded-full bg-primary/10 flex items-center justify-center">
          <User className="h-5 w-5 text-primary" />
        </div>
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2">
            <p className="font-semibold text-foreground">{author.name}</p>
            {author.linkedInUrl && (
              <a
                href={author.linkedInUrl}
                target="_blank"
                rel="noopener noreferrer"
                className="text-muted-foreground hover:text-primary transition-colors"
                aria-label="LinkedIn profile"
              >
                <Linkedin className="h-4 w-4" />
              </a>
            )}
          </div>
          {(author.title || author.company) && (
            <p className="text-sm text-muted-foreground">
              {author.title}
              {author.title && author.company && " at "}
              {author.company}
            </p>
          )}
          <span className="inline-flex items-center px-2 py-0.5 text-xs font-medium bg-secondary text-secondary-foreground rounded-full mt-1">
            {RELATIONSHIP_LABELS[relationship]}
          </span>
        </div>
      </div>

      {/* Quotes */}
      <div className="space-y-4">
        {visibleTestimonials.map((testimonial) => (
          <QuoteItem key={testimonial.id} testimonial={testimonial} onSkillClick={onSkillClick} />
        ))}
      </div>

      {/* Expand/collapse button */}
      {testimonials.length > COLLAPSE_THRESHOLD && (
        <button
          type="button"
          onClick={() => setIsExpanded(!isExpanded)}
          className="mt-4 flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground transition-colors"
        >
          {isExpanded ? (
            <>
              <ChevronUp className="h-4 w-4" />
              Show less
            </>
          ) : (
            <>
              <ChevronDown className="h-4 w-4" />
              Show {hiddenCount} more
            </>
          )}
        </button>
      )}
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
  // Group testimonials by author (must be called before any early returns)
  const groups = useMemo(() => groupTestimonialsByAuthor(testimonials), [testimonials]);

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
          {groups.map((group) => (
            <TestimonialGroupCard key={group.author.id} group={group} onSkillClick={onSkillClick} />
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
