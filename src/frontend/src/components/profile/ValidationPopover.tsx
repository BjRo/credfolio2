"use client";

import { CheckCircle2, FileText, Loader2 } from "lucide-react";
import type { ReactNode } from "react";
import { useQuery } from "urql";

import { HoverCard, HoverCardContent, HoverCardTrigger } from "@/components/ui/hover-card";
import {
  AuthorRelationship,
  GetExperienceValidationsDocument,
  GetSkillValidationsDocument,
} from "@/graphql/generated/graphql";

interface ValidationPopoverProps {
  /** The ID of the skill or experience to show validations for */
  itemId: string;
  /** Whether this is a skill or experience validation */
  type: "skill" | "experience";
  /** The item name (for display in the popover header) */
  itemName: string;
  /** The number of validations (for the header) */
  validationCount: number;
  /** The trigger element (usually the skill tag or experience card) */
  children: ReactNode;
}

const RELATIONSHIP_LABELS: Record<AuthorRelationship, string> = {
  [AuthorRelationship.Manager]: "Manager",
  [AuthorRelationship.Peer]: "Peer",
  [AuthorRelationship.DirectReport]: "Direct Report",
  [AuthorRelationship.Client]: "Client",
  [AuthorRelationship.Mentor]: "Mentor",
  [AuthorRelationship.Professor]: "Professor",
  [AuthorRelationship.Colleague]: "Colleague",
  [AuthorRelationship.Other]: "Reference",
};

/**
 * Displays a hover popover with validation details for a skill or experience.
 * Shows who provided the validation and the quote snippet.
 */
export function ValidationPopover({
  itemId,
  type,
  itemName,
  validationCount,
  children,
}: ValidationPopoverProps) {
  // Don't show popover if no validations
  if (validationCount === 0) {
    return <>{children}</>;
  }

  return (
    <HoverCard openDelay={300} closeDelay={100}>
      <HoverCardTrigger asChild>{children}</HoverCardTrigger>
      <HoverCardContent className="w-80" side="top" align="start">
        <ValidationPopoverContent
          itemId={itemId}
          type={type}
          itemName={itemName}
          validationCount={validationCount}
        />
      </HoverCardContent>
    </HoverCard>
  );
}

interface ValidationPopoverContentProps {
  itemId: string;
  type: "skill" | "experience";
  itemName: string;
  validationCount: number;
}

function ValidationPopoverContent({
  itemId,
  type,
  itemName,
  validationCount,
}: ValidationPopoverContentProps) {
  // Fetch skill validations
  const [skillResult] = useQuery({
    query: GetSkillValidationsDocument,
    variables: { skillId: itemId },
    pause: type !== "skill",
  });

  // Fetch experience validations
  const [experienceResult] = useQuery({
    query: GetExperienceValidationsDocument,
    variables: { experienceId: itemId },
    pause: type !== "experience",
  });

  const result = type === "skill" ? skillResult : experienceResult;
  const validations =
    type === "skill"
      ? skillResult.data?.skillValidations
      : experienceResult.data?.experienceValidations;

  const isLoading = result.fetching;
  const error = result.error;

  return (
    <div className="space-y-3">
      {/* Header */}
      <div className="flex items-center justify-between border-b border-border/50 pb-2">
        <span className="font-medium text-sm">{itemName}</span>
        <span className="text-xs text-muted-foreground">
          {validationCount} validation{validationCount === 1 ? "" : "s"}
        </span>
      </div>

      {/* Loading state */}
      {isLoading && (
        <div className="flex items-center justify-center py-4">
          <Loader2 className="h-4 w-4 animate-spin text-muted-foreground" />
        </div>
      )}

      {/* Error state */}
      {error && (
        <div className="text-sm text-destructive py-2">Failed to load validation details</div>
      )}

      {/* Validations list */}
      {!isLoading && !error && validations && (
        <div className="space-y-3">
          {/* Resume source indicator */}
          <div className="flex items-start gap-2 text-sm">
            <FileText className="h-4 w-4 text-muted-foreground mt-0.5 shrink-0" />
            <div>
              <div className="font-medium text-foreground/80">Resume</div>
              <div className="text-xs text-muted-foreground">Listed in your profile</div>
            </div>
          </div>

          {/* Reference letter validations */}
          {validations.map((validation) => {
            const author = validation.referenceLetter?.extractedData?.author;
            const authorName = author?.name ?? "Reference";
            const authorTitle = author?.title;
            const authorCompany = author?.company;
            const relationship = author?.relationship
              ? RELATIONSHIP_LABELS[author.relationship]
              : "Reference";

            return (
              <div key={validation.id} className="flex items-start gap-2 text-sm">
                <CheckCircle2 className="h-4 w-4 text-green-500 mt-0.5 shrink-0" />
                <div className="min-w-0">
                  <div className="font-medium text-foreground/80">{authorName}</div>
                  {(authorTitle || authorCompany) && (
                    <div className="text-xs text-muted-foreground">
                      {[authorTitle, authorCompany].filter(Boolean).join(", ")}
                    </div>
                  )}
                  <div className="text-xs text-muted-foreground">{relationship}</div>
                  {validation.quoteSnippet && (
                    <blockquote className="mt-1.5 text-xs text-foreground/70 italic border-l-2 border-green-500/30 pl-2">
                      &ldquo;{validation.quoteSnippet}&rdquo;
                    </blockquote>
                  )}
                </div>
              </div>
            );
          })}

          {/* Link to testimonials section */}
          <a
            href="#testimonials"
            className="block text-xs text-primary hover:underline pt-1 border-t border-border/50"
          >
            View full testimonials →
          </a>
        </div>
      )}
    </div>
  );
}
