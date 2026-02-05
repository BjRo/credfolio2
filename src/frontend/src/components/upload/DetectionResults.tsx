"use client";

import { useCallback, useState } from "react";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { Textarea } from "@/components/ui/textarea";
import { GRAPHQL_ENDPOINT } from "@/lib/urql/client";
import type { DetectionResultsProps } from "./types";
import { CONFIDENCE_THRESHOLD } from "./types";

type ManualSelection = "resume" | "reference-letter" | "both" | null;
type CorrectionChoice = "just-resume" | "just-letter" | null;

export function DetectionResults({
  detection,
  fileName,
  userId,
  onProceed,
  onCancel,
}: DetectionResultsProps) {
  const isLowConfidence =
    detection.confidence < CONFIDENCE_THRESHOLD ||
    (!detection.hasCareerInfo && !detection.hasTestimonial);

  const [extractCareerInfo, setExtractCareerInfo] = useState(detection.hasCareerInfo);
  const [extractTestimonial, setExtractTestimonial] = useState(detection.hasTestimonial);
  const [manualSelection, setManualSelection] = useState<ManualSelection>(null);
  const [correctionOpen, setCorrectionOpen] = useState(false);
  const [correctionChoice, setCorrectionChoice] = useState<CorrectionChoice>(null);
  const [feedbackText, setFeedbackText] = useState("");

  const handleCorrectionChange = (choice: CorrectionChoice) => {
    setCorrectionChoice(choice);
    if (choice === "just-resume") {
      setExtractCareerInfo(true);
      setExtractTestimonial(false);
    } else if (choice === "just-letter") {
      setExtractCareerInfo(false);
      setExtractTestimonial(true);
    }
  };

  const handleManualSelection = (selection: ManualSelection) => {
    setManualSelection(selection);
  };

  const getManualSelections = (): { career: boolean; testimonial: boolean } => {
    switch (manualSelection) {
      case "resume":
        return { career: true, testimonial: false };
      case "reference-letter":
        return { career: false, testimonial: true };
      case "both":
        return { career: true, testimonial: true };
      default:
        return { career: false, testimonial: false };
    }
  };

  const submitFeedback = useCallback(
    (message: string) => {
      fetch(GRAPHQL_ENDPOINT, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          query: `mutation ReportDocumentFeedback($userId: ID!, $input: DocumentFeedbackInput!) {
            reportDocumentFeedback(userId: $userId, input: $input) { success }
          }`,
          variables: {
            userId,
            input: {
              fileId: detection.fileId,
              feedbackType: "DETECTION_CORRECTION",
              message,
            },
          },
        }),
      }).catch(() => {
        // Fire-and-forget: don't block the user flow on feedback failures
      });
    },
    [userId, detection.fileId]
  );

  const handleProceed = () => {
    // Submit correction feedback if the user made a correction
    if (correctionChoice || feedbackText.trim()) {
      const parts: string[] = [];
      if (correctionChoice) parts.push(`Correction: ${correctionChoice}`);
      if (feedbackText.trim()) parts.push(feedbackText.trim());
      submitFeedback(parts.join(" — "));
    }

    if (isLowConfidence) {
      const selections = getManualSelections();
      onProceed(selections.career, selections.testimonial);
    } else {
      onProceed(extractCareerInfo, extractTestimonial);
    }
  };

  const isProceedDisabled = isLowConfidence
    ? manualSelection === null
    : !extractCareerInfo && !extractTestimonial;

  // Determine which checkboxes to show based on corrections
  const showCareerInfo =
    correctionChoice === "just-letter"
      ? false
      : detection.hasCareerInfo || correctionChoice === "just-resume";
  const showTestimonial =
    correctionChoice === "just-resume"
      ? false
      : detection.hasTestimonial || correctionChoice === "just-letter";

  return (
    <div className="w-full max-w-xl mx-auto space-y-6">
      {/* File name */}
      <div className="text-center">
        <p className="text-sm text-muted-foreground">
          <span className="font-medium">{fileName}</span>
        </p>
      </div>

      {/* Summary */}
      {detection.summary && (
        <p className="text-sm text-muted-foreground text-center">{detection.summary}</p>
      )}

      {isLowConfidence ? (
        <LowConfidenceUI
          manualSelection={manualSelection}
          onSelectionChange={handleManualSelection}
        />
      ) : (
        <>
          {/* Content type checkboxes */}
          <div className="space-y-3">
            <h3 className="text-sm font-medium">We found the following content:</h3>

            {showCareerInfo && (
              <div className="flex items-start gap-3 p-3 rounded-lg border bg-card">
                <Checkbox
                  checked={extractCareerInfo}
                  onCheckedChange={(checked) => setExtractCareerInfo(checked)}
                  aria-label="Career Information"
                />
                <div>
                  <p className="text-sm font-medium">Career Information</p>
                  <p className="text-xs text-muted-foreground">
                    Resume/CV content including work experience, education, and skills
                  </p>
                </div>
              </div>
            )}

            {showTestimonial && (
              <div className="flex items-start gap-3 p-3 rounded-lg border bg-card">
                <Checkbox
                  checked={extractTestimonial}
                  onCheckedChange={(checked) => setExtractTestimonial(checked)}
                  aria-label="Testimonial"
                />
                <div>
                  <p className="text-sm font-medium">
                    Testimonial
                    {detection.testimonialAuthor && (
                      <span className="font-normal text-muted-foreground">
                        {" "}
                        from {detection.testimonialAuthor}
                      </span>
                    )}
                  </p>
                  <p className="text-xs text-muted-foreground">
                    Reference letter or recommendation content
                  </p>
                </div>
              </div>
            )}
          </div>

          {/* Correction UI */}
          <CorrectionUI
            isOpen={correctionOpen}
            onToggle={() => setCorrectionOpen(!correctionOpen)}
            correctionChoice={correctionChoice}
            onCorrectionChange={handleCorrectionChange}
            feedbackText={feedbackText}
            onFeedbackChange={setFeedbackText}
          />
        </>
      )}

      {/* Action buttons */}
      <div className="flex flex-col gap-3 pt-2">
        <Button onClick={handleProceed} disabled={isProceedDisabled}>
          Extract selected content
        </Button>
        <Button variant="ghost" onClick={onCancel}>
          Upload different document
        </Button>
      </div>
    </div>
  );
}

function LowConfidenceUI({
  manualSelection,
  onSelectionChange,
}: {
  manualSelection: ManualSelection;
  onSelectionChange: (selection: ManualSelection) => void;
}) {
  return (
    <div className="space-y-4">
      <div className="p-4 bg-warning/10 border border-warning/30 rounded-lg">
        <p className="text-sm font-medium text-warning-foreground">
          We're not sure what this document contains. Please tell us:
        </p>
      </div>

      <fieldset className="space-y-2">
        <legend className="sr-only">Document type</legend>

        <RadioOption
          name="manual-selection"
          value="resume"
          label="Resume / CV"
          checked={manualSelection === "resume"}
          onChange={() => onSelectionChange("resume")}
        />
        <RadioOption
          name="manual-selection"
          value="reference-letter"
          label="Reference letter"
          checked={manualSelection === "reference-letter"}
          onChange={() => onSelectionChange("reference-letter")}
        />
        <RadioOption
          name="manual-selection"
          value="both"
          label="Both"
          checked={manualSelection === "both"}
          onChange={() => onSelectionChange("both")}
        />
      </fieldset>
    </div>
  );
}

function CorrectionUI({
  isOpen,
  onToggle,
  correctionChoice,
  onCorrectionChange,
  feedbackText,
  onFeedbackChange,
}: {
  isOpen: boolean;
  onToggle: () => void;
  correctionChoice: CorrectionChoice;
  onCorrectionChange: (choice: CorrectionChoice) => void;
  feedbackText: string;
  onFeedbackChange: (text: string) => void;
}) {
  return (
    <div className="border-t pt-4">
      <button
        type="button"
        className="text-sm text-muted-foreground hover:text-foreground transition-colors"
        onClick={onToggle}
      >
        Not what you expected?
      </button>

      {isOpen && (
        <div className="mt-3 space-y-3">
          <fieldset className="space-y-2">
            <legend className="sr-only">Correction</legend>

            <RadioOption
              name="correction"
              value="just-resume"
              label="This is just a resume"
              checked={correctionChoice === "just-resume"}
              onChange={() => onCorrectionChange("just-resume")}
            />
            <RadioOption
              name="correction"
              value="just-letter"
              label="This is just a reference letter"
              checked={correctionChoice === "just-letter"}
              onChange={() => onCorrectionChange("just-letter")}
            />
          </fieldset>

          <Textarea
            placeholder="Tell us more about this document..."
            value={feedbackText}
            onChange={(e) => onFeedbackChange(e.target.value)}
            className="text-sm"
            rows={2}
          />
        </div>
      )}
    </div>
  );
}

function RadioOption({
  name,
  value,
  label,
  checked,
  onChange,
}: {
  name: string;
  value: string;
  label: string;
  checked: boolean;
  onChange: () => void;
}) {
  return (
    <label className="flex items-center gap-3 p-3 rounded-lg border bg-card cursor-pointer hover:border-primary/50 transition-colors">
      <input
        type="radio"
        name={name}
        value={value}
        checked={checked}
        onChange={onChange}
        className="h-4 w-4 text-primary"
        aria-label={label}
      />
      <span className="text-sm">{label}</span>
    </label>
  );
}
