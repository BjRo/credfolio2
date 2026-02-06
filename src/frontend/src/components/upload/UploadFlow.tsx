"use client";

import { useRouter } from "next/navigation";
import { useState } from "react";
import { DetectionProgress } from "./DetectionProgress";
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
  const [fileId, setFileId] = useState<string | null>(null);
  const [fileName, setFileName] = useState<string | null>(null);
  const [detection, setDetection] = useState<DocumentDetectionResult | null>(null);
  const [extractCareerInfo, setExtractCareerInfo] = useState(false);
  const [extractTestimonial, setExtractTestimonial] = useState(false);
  const [extractionResults, setExtractionResults] = useState<ExtractionResults | null>(null);
  const [processIds, setProcessIds] = useState<ProcessDocumentIds | null>(null);

  const handleUploadComplete = (id: string, name: string) => {
    setFileId(id);
    setFileName(name);
    setCurrentStep("detect");
  };

  const handleDetectionComplete = (result: DocumentDetectionResult) => {
    setDetection(result);
    setCurrentStep("review-detection");
  };

  const handleDetectionError = (_error: string) => {
    // Error is displayed by DetectionProgress component
  };

  const handleProceed = (careerInfo: boolean, testimonial: boolean) => {
    setExtractCareerInfo(careerInfo);
    setExtractTestimonial(testimonial);
    setCurrentStep("extract");
  };

  const handleCancel = () => {
    setDetection(null);
    setFileId(null);
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

  const handleImportComplete = (profileId: string) => {
    setCurrentStep("import");
    router.push(`/profile/${profileId}`);
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
        <DocumentUpload userId={DEMO_USER_ID} onUploadComplete={handleUploadComplete} />
      )}

      {currentStep === "detect" && fileId && (
        <DetectionProgress
          fileId={fileId}
          onDetectionComplete={handleDetectionComplete}
          onError={handleDetectionError}
        />
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

      {currentStep === "extract" && fileId && (
        <ExtractionProgress
          userId={DEMO_USER_ID}
          fileId={fileId}
          extractCareerInfo={extractCareerInfo}
          extractTestimonial={extractTestimonial}
          onComplete={handleExtractionComplete}
          onError={handleExtractionError}
        />
      )}

      {currentStep === "review-results" && extractionResults && fileId && processIds && (
        <ExtractionReview
          userId={DEMO_USER_ID}
          fileId={fileId}
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
