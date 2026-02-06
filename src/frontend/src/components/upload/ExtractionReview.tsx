"use client";

import { Briefcase, GraduationCap, Quote, Sparkles, Wrench } from "lucide-react";
import { useCallback, useMemo, useState } from "react";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { SelectionControls } from "@/components/ui/selection-controls";
import { GRAPHQL_ENDPOINT } from "@/lib/urql/client";
import { FeedbackForm } from "./FeedbackForm";
import type {
  ExtractionResults,
  ExtractionReviewProps,
  ReferenceLetterExtractionData,
  ResumeExtractionData,
} from "./types";

const IMPORT_MUTATION = `
  mutation ImportDocumentResults($userId: ID!, $input: ImportDocumentResultsInput!) {
    importDocumentResults(userId: $userId, input: $input) {
      ... on ImportDocumentResultsResult {
        __typename
        profile { id }
        importedCount { experiences educations skills testimonials }
      }
      ... on ImportDocumentResultsError {
        __typename
        message
        field
      }
    }
  }
`;

type ImportStatus = "idle" | "importing" | "success" | "error";

export function ExtractionReview({
  userId,
  fileId,
  results,
  processDocumentIds,
  onImportComplete,
  onBack,
}: ExtractionReviewProps) {
  const [importStatus, setImportStatus] = useState<ImportStatus>("idle");
  const [importError, setImportError] = useState<string | null>(null);
  const [feedbackOpen, setFeedbackOpen] = useState(false);

  const hasResumeData = results.resume?.status === "COMPLETED" && results.resume.extractedData;
  const hasLetterData =
    results.referenceLetter?.status === "COMPLETED" && results.referenceLetter.extractedData;
  const resumeFailed = results.resume?.status === "FAILED";
  const letterFailed = results.referenceLetter?.status === "FAILED";
  const anyFailed = (results.resume && resumeFailed) || (results.referenceLetter && letterFailed);

  // Extract data arrays for selection
  const experiences = results.resume?.extractedData?.experiences ?? [];
  const educations = results.resume?.extractedData?.educations ?? [];
  const skills = results.resume?.extractedData?.skills ?? [];
  const testimonials = results.referenceLetter?.extractedData?.testimonials ?? [];
  const discoveredSkills = results.referenceLetter?.extractedData?.discoveredSkills ?? [];

  // Selection state — all pre-selected except discovered skills
  const [selectedExperiences, setSelectedExperiences] = useState<Set<number>>(
    () => new Set(experiences.map((_, i) => i))
  );
  const [selectedEducation, setSelectedEducation] = useState<Set<number>>(
    () => new Set(educations.map((_, i) => i))
  );
  const [selectedSkills, setSelectedSkills] = useState<Set<string>>(() => new Set(skills));
  const [selectedTestimonials, setSelectedTestimonials] = useState<Set<number>>(
    () => new Set(testimonials.map((_, i) => i))
  );
  const [selectedDiscoveredSkills, setSelectedDiscoveredSkills] = useState<Set<string>>(
    () => new Set()
  );

  // Toggle handlers
  const toggleExperience = useCallback((index: number) => {
    setSelectedExperiences((prev) => {
      const next = new Set(prev);
      if (next.has(index)) next.delete(index);
      else next.add(index);
      return next;
    });
  }, []);

  const toggleEducation = useCallback((index: number) => {
    setSelectedEducation((prev) => {
      const next = new Set(prev);
      if (next.has(index)) next.delete(index);
      else next.add(index);
      return next;
    });
  }, []);

  const toggleSkill = useCallback((name: string) => {
    setSelectedSkills((prev) => {
      const next = new Set(prev);
      if (next.has(name)) next.delete(name);
      else next.add(name);
      return next;
    });
  }, []);

  const toggleTestimonial = useCallback((index: number) => {
    setSelectedTestimonials((prev) => {
      const next = new Set(prev);
      if (next.has(index)) next.delete(index);
      else next.add(index);
      return next;
    });
  }, []);

  const toggleDiscoveredSkill = useCallback((name: string) => {
    setSelectedDiscoveredSkills((prev) => {
      const next = new Set(prev);
      if (next.has(name)) next.delete(name);
      else next.add(name);
      return next;
    });
  }, []);

  // Total selected count
  const totalSelected = useMemo(
    () =>
      selectedExperiences.size +
      selectedEducation.size +
      selectedSkills.size +
      selectedTestimonials.size +
      selectedDiscoveredSkills.size,
    [
      selectedExperiences,
      selectedEducation,
      selectedSkills,
      selectedTestimonials,
      selectedDiscoveredSkills,
    ]
  );

  const canImport = !!(hasResumeData || hasLetterData) && totalSelected > 0;
  const isImporting = importStatus === "importing";

  const handleImport = useCallback(async () => {
    setImportStatus("importing");
    setImportError(null);
    try {
      const response = await fetch(GRAPHQL_ENDPOINT, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          query: IMPORT_MUTATION,
          variables: {
            userId,
            input: {
              resumeId: processDocumentIds.resumeId,
              referenceLetterID: processDocumentIds.referenceLetterID,
              selectedExperienceIndices: hasResumeData ? [...selectedExperiences] : null,
              selectedEducationIndices: hasResumeData ? [...selectedEducation] : null,
              selectedSkills: hasResumeData ? [...selectedSkills] : null,
              selectedTestimonialIndices: hasLetterData ? [...selectedTestimonials] : null,
              selectedDiscoveredSkills: hasLetterData ? [...selectedDiscoveredSkills] : null,
            },
          },
        }),
      });
      const result = await response.json();
      const data = result.data?.importDocumentResults;

      if (data?.__typename === "ImportDocumentResultsError") {
        setImportStatus("error");
        setImportError(data.message);
        return;
      }

      setImportStatus("success");
      onImportComplete(data.profile.id);
    } catch {
      setImportStatus("error");
      setImportError("Import failed. Please try again.");
    }
  }, [
    userId,
    processDocumentIds,
    onImportComplete,
    hasResumeData,
    hasLetterData,
    selectedExperiences,
    selectedEducation,
    selectedSkills,
    selectedTestimonials,
    selectedDiscoveredSkills,
  ]);

  return (
    <div className="space-y-6 w-full">
      {anyFailed && <PartialFailureWarning results={results} />}

      {hasResumeData && <CareerInfoSection data={results.resume as ResumeExtractionData} />}

      {experiences.length > 0 && (
        <ExperiencesSection
          experiences={experiences}
          selected={selectedExperiences}
          onToggle={toggleExperience}
          onSelectAll={() => setSelectedExperiences(new Set(experiences.map((_, i) => i)))}
          onDeselectAll={() => setSelectedExperiences(new Set())}
          disabled={isImporting}
        />
      )}

      {educations.length > 0 && (
        <EducationSection
          educations={educations}
          selected={selectedEducation}
          onToggle={toggleEducation}
          onSelectAll={() => setSelectedEducation(new Set(educations.map((_, i) => i)))}
          onDeselectAll={() => setSelectedEducation(new Set())}
          disabled={isImporting}
        />
      )}

      {skills.length > 0 && (
        <SkillsSection
          skills={skills}
          selected={selectedSkills}
          onToggle={toggleSkill}
          onSelectAll={() => setSelectedSkills(new Set(skills))}
          onDeselectAll={() => setSelectedSkills(new Set())}
          disabled={isImporting}
        />
      )}

      {hasLetterData && (
        <AuthorInfoSection data={results.referenceLetter as ReferenceLetterExtractionData} />
      )}

      {testimonials.length > 0 && (
        <TestimonialsSection
          testimonials={testimonials}
          selected={selectedTestimonials}
          onToggle={toggleTestimonial}
          onSelectAll={() => setSelectedTestimonials(new Set(testimonials.map((_, i) => i)))}
          onDeselectAll={() => setSelectedTestimonials(new Set())}
          authorName={results.referenceLetter?.extractedData?.author.name ?? ""}
          authorAttribution={formatAuthorAttribution(
            results.referenceLetter?.extractedData?.author
          )}
          disabled={isImporting}
        />
      )}

      {discoveredSkills.length > 0 && (
        <DiscoveredSkillsSection
          discoveredSkills={discoveredSkills}
          selected={selectedDiscoveredSkills}
          onToggle={toggleDiscoveredSkill}
          onSelectAll={() =>
            setSelectedDiscoveredSkills(new Set(discoveredSkills.map((s) => s.skill)))
          }
          onDeselectAll={() => setSelectedDiscoveredSkills(new Set())}
          disabled={isImporting}
        />
      )}

      <div className="border-t pt-4">
        <button
          type="button"
          onClick={() => setFeedbackOpen(!feedbackOpen)}
          className="text-sm text-muted-foreground hover:text-foreground underline-offset-4 hover:underline"
        >
          Something doesn&apos;t look right?
        </button>
        {feedbackOpen && (
          <div className="mt-3">
            <FeedbackForm userId={userId} fileId={fileId} />
          </div>
        )}
      </div>

      {importStatus === "error" && importError && (
        <div className="p-3 bg-destructive/10 border border-destructive/30 rounded-md">
          <p className="text-sm text-destructive">{importError}</p>
        </div>
      )}

      <div className="flex justify-between items-center pt-2">
        <Button variant="ghost" onClick={onBack}>
          Back
        </Button>
        <div className="flex items-center gap-4">
          <span className="text-sm text-muted-foreground">{totalSelected} item(s) selected</span>
          <Button onClick={handleImport} disabled={!canImport || isImporting}>
            {isImporting ? "Importing..." : "Import Selected"}
          </Button>
        </div>
      </div>
    </div>
  );
}

function formatAuthorAttribution(
  author: { title: string | null; company: string | null; relationship: string } | undefined
): string {
  if (!author) return "";
  const parts = [author.title, author.company].filter(Boolean);
  let attribution = parts.join(", ");
  if (author.relationship) {
    if (attribution) attribution += ` \u00B7 ${author.relationship}`;
    else attribution = author.relationship;
  }
  return attribution;
}

function PartialFailureWarning({ results }: { results: ExtractionResults }) {
  const resumeFailed = results.resume?.status === "FAILED";
  const letterFailed = results.referenceLetter?.status === "FAILED";

  let message = "Some extraction could not be completed.";
  if (resumeFailed && !letterFailed) {
    message = "Career information extraction could not be completed.";
  } else if (letterFailed && !resumeFailed) {
    message = "Testimonial extraction could not be completed.";
  }

  return (
    <div className="p-3 bg-warning/10 border border-warning/30 rounded-md">
      <p className="text-sm text-warning-foreground">{message}</p>
    </div>
  );
}

function CareerInfoSection({ data }: { data: ResumeExtractionData }) {
  const { extractedData } = data;
  if (!extractedData) return null;

  return (
    <div className="space-y-3">
      <h3 className="text-lg font-semibold">Career Information</h3>
      <div className="rounded-lg border p-4 space-y-2">
        <p className="font-medium text-base">{extractedData.name}</p>
        {extractedData.email && (
          <p className="text-sm text-muted-foreground">{extractedData.email}</p>
        )}
        {extractedData.phone && (
          <p className="text-sm text-muted-foreground">{extractedData.phone}</p>
        )}
        {extractedData.location && (
          <p className="text-sm text-muted-foreground">{extractedData.location}</p>
        )}
        {extractedData.summary && <p className="text-sm mt-2">{extractedData.summary}</p>}
      </div>
    </div>
  );
}

function AuthorInfoSection({ data }: { data: ReferenceLetterExtractionData }) {
  const { extractedData } = data;
  if (!extractedData) return null;

  const { author } = extractedData;
  return (
    <div className="space-y-3">
      <h3 className="text-lg font-semibold">Reference Letter</h3>
      <div className="rounded-lg border p-4 space-y-1">
        <p className="font-medium">{author.name}</p>
        <p className="text-sm text-muted-foreground">
          {[author.title, author.company].filter(Boolean).join(", ")}
          {author.relationship && ` \u00B7 ${author.relationship}`}
        </p>
      </div>
    </div>
  );
}

// --- Selectable section components ---

interface CheckboxCardProps {
  checked: boolean;
  onToggle: () => void;
  disabled: boolean;
  selectedClassName: string;
  borderStyle?: string;
  children: React.ReactNode;
}

function CheckboxCard({
  checked,
  onToggle,
  disabled,
  selectedClassName,
  borderStyle = "border",
  children,
}: CheckboxCardProps) {
  return (
    // biome-ignore lint/a11y/useSemanticElements: Custom styled checkbox card with inner Checkbox component
    <div
      role="checkbox"
      aria-checked={checked}
      tabIndex={disabled ? -1 : 0}
      onClick={() => !disabled && onToggle()}
      onKeyDown={(e) => {
        if (!disabled && (e.key === " " || e.key === "Enter")) {
          e.preventDefault();
          onToggle();
        }
      }}
      className={`flex items-start gap-3 p-4 rounded-lg ${borderStyle} transition-colors w-full text-left ${
        disabled
          ? "bg-muted cursor-not-allowed"
          : checked
            ? `${selectedClassName} cursor-pointer`
            : "bg-card border-border hover:bg-muted/50 cursor-pointer"
      }`}
    >
      <Checkbox
        checked={checked}
        onCheckedChange={() => onToggle()}
        onClick={(e) => e.stopPropagation()}
        disabled={disabled}
        className="mt-1"
        tabIndex={-1}
        aria-hidden="true"
      />
      <div className="flex-1">{children}</div>
    </div>
  );
}

function ExperiencesSection({
  experiences,
  selected,
  onToggle,
  onSelectAll,
  onDeselectAll,
  disabled,
}: {
  experiences: ResumeExtractionData["extractedData"] extends infer T
    ? T extends { experiences: infer E }
      ? E
      : never
    : never;
  selected: Set<number>;
  onToggle: (index: number) => void;
  onSelectAll: () => void;
  onDeselectAll: () => void;
  disabled: boolean;
}) {
  return (
    <section className="space-y-4">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <Briefcase className="h-5 w-5 text-success" />
          <h2 className="text-xl font-semibold text-foreground">Work Experience</h2>
        </div>
        <SelectionControls
          selectedCount={selected.size}
          totalCount={experiences.length}
          onSelectAll={onSelectAll}
          onDeselectAll={onDeselectAll}
          disabled={disabled}
        />
      </div>
      {/* biome-ignore lint/a11y/useSemanticElements: Using role="group" for checkbox group semantics */}
      <div className="space-y-3" role="group" aria-label="Work experiences">
        {experiences.map((exp, index) => (
          <CheckboxCard
            key={`exp-${exp.company}-${exp.title}-${index}`}
            checked={selected.has(index)}
            onToggle={() => onToggle(index)}
            disabled={disabled}
            selectedClassName="bg-success/5 border-success/30"
          >
            <p className="font-medium text-foreground">{exp.title}</p>
            <p className="text-sm text-muted-foreground">
              {exp.company}
              {exp.location && ` \u00B7 ${exp.location}`}
            </p>
            {(exp.startDate || exp.endDate) && (
              <p className="text-xs text-muted-foreground mt-1">
                {exp.startDate ?? ""}
                {exp.startDate && exp.endDate ? " – " : ""}
                {exp.isCurrent ? "Present" : (exp.endDate ?? "")}
              </p>
            )}
            {exp.description && (
              <p className="text-sm text-muted-foreground mt-2">{exp.description}</p>
            )}
          </CheckboxCard>
        ))}
      </div>
    </section>
  );
}

function EducationSection({
  educations,
  selected,
  onToggle,
  onSelectAll,
  onDeselectAll,
  disabled,
}: {
  educations: ResumeExtractionData["extractedData"] extends infer T
    ? T extends { educations: infer E }
      ? E
      : never
    : never;
  selected: Set<number>;
  onToggle: (index: number) => void;
  onSelectAll: () => void;
  onDeselectAll: () => void;
  disabled: boolean;
}) {
  return (
    <section className="space-y-4">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <GraduationCap className="h-5 w-5 text-success" />
          <h2 className="text-xl font-semibold text-foreground">Education</h2>
        </div>
        <SelectionControls
          selectedCount={selected.size}
          totalCount={educations.length}
          onSelectAll={onSelectAll}
          onDeselectAll={onDeselectAll}
          disabled={disabled}
        />
      </div>
      {/* biome-ignore lint/a11y/useSemanticElements: Using role="group" for checkbox group semantics */}
      <div className="space-y-3" role="group" aria-label="Education entries">
        {educations.map((edu, index) => (
          <CheckboxCard
            key={`edu-${edu.institution}-${index}`}
            checked={selected.has(index)}
            onToggle={() => onToggle(index)}
            disabled={disabled}
            selectedClassName="bg-success/5 border-success/30"
          >
            <p className="font-medium text-foreground">
              {edu.degree && `${edu.degree}`}
              {edu.degree && edu.field && " in "}
              {edu.field && `${edu.field}`}
              {!edu.degree && !edu.field && edu.institution}
            </p>
            {(edu.degree || edu.field) && (
              <p className="text-sm text-muted-foreground">{edu.institution}</p>
            )}
            {(edu.startDate || edu.endDate) && (
              <p className="text-xs text-muted-foreground mt-1">
                {edu.startDate ?? ""}
                {edu.startDate && edu.endDate ? " – " : ""}
                {edu.endDate ?? ""}
              </p>
            )}
          </CheckboxCard>
        ))}
      </div>
    </section>
  );
}

function SkillsSection({
  skills,
  selected,
  onToggle,
  onSelectAll,
  onDeselectAll,
  disabled,
}: {
  skills: string[];
  selected: Set<string>;
  onToggle: (name: string) => void;
  onSelectAll: () => void;
  onDeselectAll: () => void;
  disabled: boolean;
}) {
  return (
    <section className="space-y-4">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <Wrench className="h-5 w-5 text-success" />
          <h2 className="text-xl font-semibold text-foreground">Skills</h2>
        </div>
        <SelectionControls
          selectedCount={selected.size}
          totalCount={skills.length}
          onSelectAll={onSelectAll}
          onDeselectAll={onDeselectAll}
          disabled={disabled}
        />
      </div>
      {/* biome-ignore lint/a11y/useSemanticElements: Using role="group" for checkbox group semantics */}
      <div className="flex flex-wrap gap-2" role="group" aria-label="Skills">
        {skills.map((skill) => (
          // biome-ignore lint/a11y/useSemanticElements: Custom styled checkbox card with inner Checkbox component
          <div
            key={skill}
            role="checkbox"
            aria-checked={selected.has(skill)}
            tabIndex={disabled ? -1 : 0}
            onClick={() => !disabled && onToggle(skill)}
            onKeyDown={(e) => {
              if (!disabled && (e.key === " " || e.key === "Enter")) {
                e.preventDefault();
                onToggle(skill);
              }
            }}
            className={`inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full border transition-colors text-sm ${
              disabled
                ? "bg-muted cursor-not-allowed"
                : selected.has(skill)
                  ? "bg-success/5 border-success/30 cursor-pointer"
                  : "bg-card border-border hover:bg-muted/50 cursor-pointer"
            }`}
          >
            <Checkbox
              checked={selected.has(skill)}
              onCheckedChange={() => onToggle(skill)}
              onClick={(e) => e.stopPropagation()}
              disabled={disabled}
              className="h-3.5 w-3.5"
              tabIndex={-1}
              aria-hidden="true"
            />
            {skill}
          </div>
        ))}
      </div>
    </section>
  );
}

function TestimonialsSection({
  testimonials,
  selected,
  onToggle,
  onSelectAll,
  onDeselectAll,
  authorName,
  authorAttribution,
  disabled,
}: {
  testimonials: Array<{ quote: string; skillsMentioned: string[] | null }>;
  selected: Set<number>;
  onToggle: (index: number) => void;
  onSelectAll: () => void;
  onDeselectAll: () => void;
  authorName: string;
  authorAttribution: string;
  disabled: boolean;
}) {
  return (
    <section className="space-y-4">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <Quote className="h-5 w-5 text-primary" />
          <h2 className="text-xl font-semibold text-foreground">Testimonials</h2>
        </div>
        <SelectionControls
          selectedCount={selected.size}
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
            checked={selected.has(index)}
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
                {testimonial.skillsMentioned && testimonial.skillsMentioned.length > 0 && (
                  <div className="flex flex-wrap gap-1">
                    <span className="text-xs text-muted-foreground mr-1">Skills mentioned:</span>
                    {testimonial.skillsMentioned.map((skill) => (
                      <Badge key={skill} variant="secondary" className="text-xs">
                        {skill}
                      </Badge>
                    ))}
                  </div>
                )}
              </div>
            </div>
          </CheckboxCard>
        ))}
      </div>
    </section>
  );
}

function DiscoveredSkillsSection({
  discoveredSkills,
  selected,
  onToggle,
  onSelectAll,
  onDeselectAll,
  disabled,
}: {
  discoveredSkills: Array<{ skill: string; quote: string; context: string | null }>;
  selected: Set<string>;
  onToggle: (name: string) => void;
  onSelectAll: () => void;
  onDeselectAll: () => void;
  disabled: boolean;
}) {
  return (
    <section className="space-y-4">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <Sparkles className="h-5 w-5 text-warning" />
          <h2 className="text-xl font-semibold text-foreground">Skills Your Reference Noticed</h2>
        </div>
        <SelectionControls
          selectedCount={selected.size}
          totalCount={discoveredSkills.length}
          onSelectAll={onSelectAll}
          onDeselectAll={onDeselectAll}
          disabled={disabled}
        />
      </div>
      <p className="text-sm text-muted-foreground">
        These skills were mentioned in the reference letter but aren&apos;t in the resume. Select
        any you want to add.
      </p>
      {/* biome-ignore lint/a11y/useSemanticElements: Using role="group" for checkbox group semantics */}
      <div className="space-y-3" role="group" aria-label="Discovered skills">
        {discoveredSkills.map((skill) => (
          <CheckboxCard
            key={skill.skill}
            checked={selected.has(skill.skill)}
            onToggle={() => onToggle(skill.skill)}
            disabled={disabled}
            selectedClassName="bg-warning/5 border-warning/50"
            borderStyle="border-2 border-dashed"
          >
            <p className="font-medium text-foreground">{skill.skill}</p>
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
