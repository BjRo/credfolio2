export interface DocumentDetectionResult {
  hasCareerInfo: boolean;
  hasTestimonial: boolean;
  testimonialAuthor: string | null;
  confidence: number;
  summary: string;
  documentTypeHint: "RESUME" | "REFERENCE_LETTER" | "HYBRID" | "UNKNOWN";
  fileId: string;
}

export type FlowStep =
  | "upload"
  | "detect"
  | "review-detection"
  | "extract"
  | "review-results"
  | "import";

export const FLOW_STEPS: { key: FlowStep; label: string }[] = [
  { key: "upload", label: "Upload" },
  { key: "detect", label: "Analyze" },
  { key: "review-detection", label: "Review Detection" },
  { key: "extract", label: "Extract" },
  { key: "review-results", label: "Review Results" },
  { key: "import", label: "Import" },
];

export interface DocumentUploadProps {
  userId: string;
  onUploadComplete: (fileId: string, fileName: string) => void;
  onError?: (error: string) => void;
}

export interface DetectionProgressProps {
  fileId: string;
  onDetectionComplete: (detection: DocumentDetectionResult) => void;
  onError: (error: string) => void;
}

export interface StepIndicatorProps {
  steps: { key: string; label: string }[];
  currentStep: string;
}

export interface DetectionResultsProps {
  detection: DocumentDetectionResult;
  fileName: string;
  userId: string;
  onProceed: (extractCareerInfo: boolean, extractTestimonial: boolean) => void;
  onCancel: () => void;
  onFeedbackSubmitted?: () => void;
}

export const CONFIDENCE_THRESHOLD = 0.7;

export interface ProcessDocumentIds {
  resumeId: string | null;
  referenceLetterID: string | null;
}

export interface ResumeExtractionData {
  id: string;
  status: "PENDING" | "PROCESSING" | "COMPLETED" | "FAILED";
  extractedData: {
    name: string;
    email: string | null;
    phone: string | null;
    location: string | null;
    summary: string | null;
    extractedAt: string;
    confidence: number;
  } | null;
  errorMessage: string | null;
}

export interface ReferenceLetterExtractionData {
  id: string;
  status: "PENDING" | "PROCESSING" | "COMPLETED" | "FAILED" | "APPLIED";
  extractedData: {
    author: {
      name: string;
      title: string | null;
      company: string | null;
      relationship: string;
    };
    testimonials: Array<{
      quote: string;
      skillsMentioned: string[] | null;
    }>;
    skillMentions: Array<{
      skill: string;
      quote: string;
      context: string | null;
    }>;
    experienceMentions: Array<{
      company: string;
      role: string;
      quote: string;
    }>;
    discoveredSkills: Array<{
      skill: string;
      quote: string;
      context: string | null;
    }>;
    metadata: {
      extractedAt: string;
      modelVersion: string;
      processingTimeMs: number | null;
    };
  } | null;
}

export interface ExtractionResults {
  resume: ResumeExtractionData | null;
  referenceLetter: ReferenceLetterExtractionData | null;
}

export interface ExtractionProgressProps {
  userId: string;
  fileId: string;
  extractCareerInfo: boolean;
  extractTestimonial: boolean;
  onComplete: (results: ExtractionResults, processIds: ProcessDocumentIds) => void;
  onError: (error: string) => void;
}

export interface ExtractionReviewProps {
  userId: string;
  fileId: string;
  results: ExtractionResults;
  processDocumentIds: ProcessDocumentIds;
  onImportComplete: (profileId: string) => void;
  onBack: () => void;
}

export interface FeedbackFormProps {
  userId: string;
  fileId: string;
  onSubmitted?: () => void;
}
