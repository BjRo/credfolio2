"use client";

import { useCallback, useState } from "react";
import { Button } from "@/components/ui/button";
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
  const canImport = !!(hasResumeData || hasLetterData);

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
  }, [userId, processDocumentIds, onImportComplete]);

  return (
    <div className="space-y-6 w-full">
      {anyFailed && <PartialFailureWarning results={results} />}

      {hasResumeData && <CareerInfoSection data={results.resume as ResumeExtractionData} />}

      {hasLetterData && (
        <TestimonialSection data={results.referenceLetter as ReferenceLetterExtractionData} />
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

      <div className="flex flex-col gap-3 pt-2">
        <Button onClick={handleImport} disabled={!canImport || importStatus === "importing"}>
          {importStatus === "importing" ? "Importing..." : "Import to profile"}
        </Button>
        <Button variant="ghost" onClick={onBack}>
          Back
        </Button>
      </div>
    </div>
  );
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

function TestimonialSection({ data }: { data: ReferenceLetterExtractionData }) {
  const { extractedData } = data;
  if (!extractedData) return null;

  const { author, testimonials, skillMentions, experienceMentions, discoveredSkills } =
    extractedData;

  return (
    <div className="space-y-3">
      <h3 className="text-lg font-semibold">Testimonial</h3>

      <div className="rounded-lg border p-4 space-y-1">
        <p className="font-medium">{author.name}</p>
        <p className="text-sm text-muted-foreground">
          {[author.title, author.company].filter(Boolean).join(", ")}
          {author.relationship && ` \u00B7 ${author.relationship}`}
        </p>
      </div>

      {testimonials.length > 0 && (
        <div className="space-y-2">
          <h4 className="text-sm font-medium">Quotes</h4>
          {testimonials.map((t) => (
            <blockquote
              key={t.quote}
              className="border-l-2 border-primary/30 pl-3 text-sm italic text-muted-foreground"
            >
              &ldquo;{t.quote}&rdquo;
              {t.skillsMentioned && t.skillsMentioned.length > 0 && (
                <span className="block mt-1 text-xs not-italic">
                  Skills: {t.skillsMentioned.join(", ")}
                </span>
              )}
            </blockquote>
          ))}
        </div>
      )}

      {skillMentions.length > 0 && (
        <div className="space-y-2">
          <h4 className="text-sm font-medium">Skill Mentions</h4>
          {skillMentions.map((s) => (
            <div key={s.skill} className="text-sm">
              <span className="font-medium">{s.skill}</span>
              <span className="text-muted-foreground"> &mdash; &ldquo;{s.quote}&rdquo;</span>
            </div>
          ))}
        </div>
      )}

      {experienceMentions.length > 0 && (
        <div className="space-y-2">
          <h4 className="text-sm font-medium">Experience Mentions</h4>
          {experienceMentions.map((e) => (
            <div key={`${e.company}-${e.role}`} className="text-sm">
              <span className="font-medium">
                {e.role} at {e.company}
              </span>
              <span className="text-muted-foreground"> &mdash; &ldquo;{e.quote}&rdquo;</span>
            </div>
          ))}
        </div>
      )}

      {discoveredSkills.length > 0 && (
        <div className="space-y-2">
          <h4 className="text-sm font-medium">Discovered Skills</h4>
          <div className="flex flex-wrap gap-2">
            {discoveredSkills.map((s) => (
              <span
                key={s.skill}
                className="inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-medium"
              >
                {s.skill}
              </span>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
