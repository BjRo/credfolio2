export interface DocumentDetectionResult {
  hasCareerInfo: boolean;
  hasTestimonial: boolean;
  testimonialAuthor: string | null;
  confidence: number;
  summary: string;
  documentTypeHint: "RESUME" | "REFERENCE_LETTER" | "HYBRID" | "UNKNOWN";
  fileId: string;
}

export type FlowStep = "upload" | "review-detection" | "extract" | "review-results" | "import";

export const FLOW_STEPS: { key: FlowStep; label: string }[] = [
  { key: "upload", label: "Upload" },
  { key: "review-detection", label: "Review Detection" },
  { key: "extract", label: "Extract" },
  { key: "review-results", label: "Review Results" },
  { key: "import", label: "Import" },
];

export interface DocumentUploadProps {
  userId: string;
  onDetectionComplete: (detection: DocumentDetectionResult, fileName: string) => void;
  onError?: (error: string) => void;
}

export interface StepIndicatorProps {
  steps: { key: string; label: string }[];
  currentStep: string;
}
