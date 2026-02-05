"use client";

import { useState } from "react";
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

  return (
    <div className="w-full max-w-3xl mx-auto">
      <StepIndicator steps={FLOW_STEPS} currentStep={currentStep} />

      {currentStep === "upload" && (
        <DocumentUpload userId={DEMO_USER_ID} onDetectionComplete={handleDetectionComplete} />
      )}

      {currentStep !== "upload" && (
        <div className="text-center py-12 text-muted-foreground">
          <p>
            Detection results ready for: <span className="font-medium">{fileName}</span>
          </p>
          {detection && (
            <p className="text-sm mt-2">
              Found: {detection.hasCareerInfo && "Career info"}
              {detection.hasCareerInfo && detection.hasTestimonial && " & "}
              {detection.hasTestimonial &&
                `Testimonial${detection.testimonialAuthor ? ` from ${detection.testimonialAuthor}` : ""}`}
            </p>
          )}
          <p className="text-sm mt-2">Next steps coming soon...</p>
        </div>
      )}
    </div>
  );
}
