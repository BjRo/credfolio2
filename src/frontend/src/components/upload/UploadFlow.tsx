"use client";

import { useRouter } from "next/navigation";
import { useState } from "react";
import { DetectionResults } from "./DetectionResults";
import { DocumentUpload } from "./DocumentUpload";
import { ExtractionProgress } from "./ExtractionProgress";
import { ExtractionReview } from "./ExtractionReview";
import { StepIndicator } from "./StepIndicator";
import type {
  DocumentDetectionResult,
  ExtractionResults,
  FlowStep,
  ProcessDocumentIds,
} from "./types";
import { FLOW_STEPS } from "./types";

const DEMO_USER_ID = "00000000-0000-0000-0000-000000000001";

export function UploadFlow() {
  const router = useRouter();
  const [currentStep, setCurrentStep] = useState<FlowStep>("upload");
  const [detection, setDetection] = useState<DocumentDetectionResult | null>(null);
  const [fileName, setFileName] = useState<string | null>(null);
  const [extractCareerInfo, setExtractCareerInfo] = useState(false);
  const [extractTestimonial, setExtractTestimonial] = useState(false);
  const [extractionResults, setExtractionResults] = useState<ExtractionResults | null>(null);
  const [processIds, setProcessIds] = useState<ProcessDocumentIds | null>(null);

  const handleDetectionComplete = (result: DocumentDetectionResult, name: string) => {
    setDetection(result);
    setFileName(name);
    setCurrentStep("review-detection");
  };

  const handleProceed = (careerInfo: boolean, testimonial: boolean) => {
    setExtractCareerInfo(careerInfo);
    setExtractTestimonial(testimonial);
    setCurrentStep("extract");
  };

  const handleCancel = () => {
    setDetection(null);
    setFileName(null);
    setCurrentStep("upload");
  };

  const handleExtractionComplete = (results: ExtractionResults, ids: ProcessDocumentIds) => {
    setExtractionResults(results);
    setProcessIds(ids);
    setCurrentStep("review-results");
  };

  const handleExtractionError = (_error: string) => {
    // Error is displayed by ExtractionProgress component
  };

  const handleImportComplete = (_profileId: string) => {
    setCurrentStep("import");
    router.push(`/profile/${DEMO_USER_ID}`);
  };

  const handleBackFromReview = () => {
    setExtractionResults(null);
    setProcessIds(null);
    setCurrentStep("review-detection");
  };

  return (
    <div className="w-full max-w-3xl mx-auto">
      <StepIndicator steps={FLOW_STEPS} currentStep={currentStep} />

      {currentStep === "upload" && (
        <DocumentUpload userId={DEMO_USER_ID} onDetectionComplete={handleDetectionComplete} />
      )}

      {currentStep === "review-detection" && detection && fileName && (
        <DetectionResults
          detection={detection}
          fileName={fileName}
          userId={DEMO_USER_ID}
          onProceed={handleProceed}
          onCancel={handleCancel}
        />
      )}

      {currentStep === "extract" && detection && (
        <ExtractionProgress
          userId={DEMO_USER_ID}
          fileId={detection.fileId}
          extractCareerInfo={extractCareerInfo}
          extractTestimonial={extractTestimonial}
          onComplete={handleExtractionComplete}
          onError={handleExtractionError}
        />
      )}

      {currentStep === "review-results" && extractionResults && detection && processIds && (
        <ExtractionReview
          userId={DEMO_USER_ID}
          fileId={detection.fileId}
          results={extractionResults}
          processDocumentIds={processIds}
          onImportComplete={handleImportComplete}
          onBack={handleBackFromReview}
        />
      )}

      {currentStep === "import" && (
        <div className="text-center py-12">
          <p className="text-sm text-muted-foreground">Redirecting to your profile...</p>
        </div>
      )}
    </div>
  );
}
