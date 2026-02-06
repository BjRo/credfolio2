"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import { GRAPHQL_ENDPOINT } from "@/lib/urql/client";
import type { ExtractionProgressProps, ExtractionResults, ProcessDocumentIds } from "./types";

const POLL_INTERVAL_MS = 2000;

const PROCESS_DOCUMENT_MUTATION = `
  mutation ProcessDocument($userId: ID!, $input: ProcessDocumentInput!) {
    processDocument(userId: $userId, input: $input) {
      ... on ProcessDocumentResult {
        __typename
        resumeId
        referenceLetterID
      }
      ... on ProcessDocumentError {
        __typename
        message
        field
      }
    }
  }
`;

const DOCUMENT_PROCESSING_STATUS_QUERY = `
  query GetDocumentProcessingStatus($resumeId: ID, $referenceLetterID: ID) {
    documentProcessingStatus(resumeId: $resumeId, referenceLetterID: $referenceLetterID) {
      allComplete
      resume {
        id
        status
        extractedData {
          name
          email
          phone
          location
          summary
          extractedAt
          confidence
        }
        errorMessage
      }
      referenceLetter {
        id
        status
        extractedData {
          author { name title company relationship }
          testimonials { quote skillsMentioned }
          skillMentions { skill quote context }
          experienceMentions { company role quote }
          discoveredSkills { skill quote context }
          metadata { extractedAt modelVersion processingTimeMs }
        }
      }
    }
  }
`;

type ProgressStatus = "starting" | "processing" | "failed";

export function ExtractionProgress({
  userId,
  fileId,
  extractCareerInfo,
  extractTestimonial,
  onComplete,
  onError,
}: ExtractionProgressProps) {
  const [status, setStatus] = useState<ProgressStatus>("starting");
  const [error, setError] = useState<string | null>(null);
  const [processIds, setProcessIds] = useState<ProcessDocumentIds | null>(null);
  const pollIntervalRef = useRef<ReturnType<typeof setInterval> | null>(null);
  const mountedRef = useRef(true);
  const isStartedRef = useRef(false);

  const stopPolling = useCallback(() => {
    if (pollIntervalRef.current) {
      clearInterval(pollIntervalRef.current);
      pollIntervalRef.current = null;
    }
  }, []);

  // Step 1: Call processDocument mutation
  const startProcessing = useCallback(async () => {
    try {
      const response = await fetch(GRAPHQL_ENDPOINT, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          query: PROCESS_DOCUMENT_MUTATION,
          variables: {
            userId,
            input: { fileId, extractCareerInfo, extractTestimonial },
          },
        }),
      });
      const result = await response.json();
      const data = result.data?.processDocument;

      if (!mountedRef.current) return;

      if (data?.__typename === "ProcessDocumentError") {
        setStatus("failed");
        setError(data.message);
        onError(data.message);
        return;
      }

      const ids: ProcessDocumentIds = {
        resumeId: data.resumeId ?? null,
        referenceLetterID: data.referenceLetterID ?? null,
      };
      setProcessIds(ids);
      setStatus("processing");
    } catch {
      if (!mountedRef.current) return;
      setStatus("failed");
      setError("Failed to start extraction. Please try again.");
      onError("Failed to start extraction. Please try again.");
    }
  }, [userId, fileId, extractCareerInfo, extractTestimonial, onError]);

  // Step 2: Poll for completion
  useEffect(() => {
    if (status !== "processing" || !processIds) return;

    const poll = async () => {
      try {
        const response = await fetch(GRAPHQL_ENDPOINT, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            query: DOCUMENT_PROCESSING_STATUS_QUERY,
            variables: {
              resumeId: processIds.resumeId,
              referenceLetterID: processIds.referenceLetterID,
            },
          }),
        });
        const result = await response.json();
        const statusData = result.data?.documentProcessingStatus;

        if (!mountedRef.current) return;

        if (!statusData) return;

        // Check for failures
        const resumeFailed = statusData.resume?.status === "FAILED";
        const letterFailed = statusData.referenceLetter?.status === "FAILED";
        const allFailed =
          (processIds.resumeId ? resumeFailed : true) &&
          (processIds.referenceLetterID ? letterFailed : true);

        if (allFailed) {
          stopPolling();
          const errorMsg = statusData.resume?.errorMessage || "Extraction failed";
          setStatus("failed");
          setError(errorMsg);
          onError(errorMsg);
          return;
        }

        if (statusData.allComplete) {
          stopPolling();
          const results: ExtractionResults = {
            resume: statusData.resume ?? null,
            referenceLetter: statusData.referenceLetter ?? null,
          };
          onComplete(results, processIds);
        }
      } catch {
        // Network errors during polling are transient — keep polling
      }
    };

    pollIntervalRef.current = setInterval(poll, POLL_INTERVAL_MS);

    return () => {
      stopPolling();
    };
  }, [status, processIds, onComplete, onError, stopPolling]);

  // Mount trigger + cleanup
  useEffect(() => {
    mountedRef.current = true;
    if (!isStartedRef.current) {
      isStartedRef.current = true;
      startProcessing();
    }

    return () => {
      mountedRef.current = false;
    };
  }, [startProcessing]);

  const progressMessage = () => {
    if (status === "starting") return "Starting extraction...";
    if (extractCareerInfo && extractTestimonial) {
      return "Extracting career information and testimonials...";
    }
    if (extractCareerInfo) return "Extracting career information...";
    return "Extracting testimonials...";
  };

  if (status === "failed" && error) {
    return (
      <div className="text-center py-12 space-y-4">
        <div className="mx-auto w-12 h-12 rounded-full bg-destructive/10 flex items-center justify-center">
          <span className="text-destructive text-xl">!</span>
        </div>
        <p className="text-sm text-destructive">{error}</p>
      </div>
    );
  }

  return (
    <div className="text-center py-12 space-y-4">
      <div className="mx-auto w-12 h-12 flex items-center justify-center">
        <svg
          className="animate-spin h-8 w-8 text-primary"
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          role="img"
          aria-label="Loading spinner"
        >
          <circle
            className="opacity-25"
            cx="12"
            cy="12"
            r="10"
            stroke="currentColor"
            strokeWidth="4"
          />
          <path
            className="opacity-75"
            fill="currentColor"
            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
          />
        </svg>
      </div>
      <p className="text-sm text-muted-foreground">{progressMessage()}</p>
    </div>
  );
}
