"use client";

import { useState } from "react";
import { DetectionResults } from "./DetectionResults";
import { DocumentUpload } from "./DocumentUpload";
import { StepIndicator } from "./StepIndicator";
import type { DocumentDetectionResult, FlowStep } from "./types";
import { FLOW_STEPS } from "./types";

const DEMO_USER_ID = "00000000-0000-0000-0000-000000000001";

export function UploadFlow() {
  const [currentStep, setCurrentStep] = useState<FlowStep>("upload");
  const [detection, setDetection] = useState<DocumentDetectionResult | null>(null);
  const [fileName, setFileName] = useState<string | null>(null);

  const handleDetectionComplete = (result: DocumentDetectionResult, name: string) => {
    setDetection(result);
    setFileName(name);
    setCurrentStep("review-detection");
  };

  const handleProceed = (_extractCareerInfo: boolean, _extractTestimonial: boolean) => {
    setCurrentStep("extract");
  };

  const handleCancel = () => {
    setDetection(null);
    setFileName(null);
    setCurrentStep("upload");
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

      {currentStep !== "upload" && currentStep !== "review-detection" && (
        <div className="text-center py-12 text-muted-foreground">
          <p className="text-sm mt-2">Next steps coming soon...</p>
        </div>
      )}
    </div>
  );
}
