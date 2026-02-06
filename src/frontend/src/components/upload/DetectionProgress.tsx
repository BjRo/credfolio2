"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import { GRAPHQL_ENDPOINT } from "@/lib/urql/client";
import type { DetectionProgressProps, DocumentDetectionResult } from "./types";

const POLL_INTERVAL_MS = 2000;

const DETECTION_STATUS_QUERY = `
  query GetDocumentDetectionStatus($fileId: ID!) {
    documentDetectionStatus(fileId: $fileId) {
      fileId
      status
      detection {
        hasCareerInfo
        hasTestimonial
        testimonialAuthor
        confidence
        summary
        documentTypeHint
        fileId
      }
      error
    }
  }
`;

type ProgressStatus = "polling" | "failed";

export function DetectionProgress({
  fileId,
  onDetectionComplete,
  onError,
}: DetectionProgressProps) {
  const [status, setStatus] = useState<ProgressStatus>("polling");
  const [error, setError] = useState<string | null>(null);
  const pollIntervalRef = useRef<ReturnType<typeof setInterval> | null>(null);
  const mountedRef = useRef(true);

  const stopPolling = useCallback(() => {
    if (pollIntervalRef.current) {
      clearInterval(pollIntervalRef.current);
      pollIntervalRef.current = null;
    }
  }, []);

  useEffect(() => {
    mountedRef.current = true;

    const poll = async () => {
      try {
        const response = await fetch(GRAPHQL_ENDPOINT, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            query: DETECTION_STATUS_QUERY,
            variables: { fileId },
          }),
        });
        const result = await response.json();
        const statusData = result.data?.documentDetectionStatus;

        if (!mountedRef.current) return;

        if (!statusData) return;

        if (statusData.status === "FAILED") {
          stopPolling();
          const errMsg = statusData.error || "Detection failed";
          setStatus("failed");
          setError(errMsg);
          onError(errMsg);
          return;
        }

        if (statusData.status === "COMPLETED" && statusData.detection) {
          stopPolling();
          const detection: DocumentDetectionResult = statusData.detection;
          onDetectionComplete(detection);
        }
      } catch {
        // Network errors during polling are transient — keep polling
      }
    };

    pollIntervalRef.current = setInterval(poll, POLL_INTERVAL_MS);

    return () => {
      mountedRef.current = false;
      stopPolling();
    };
  }, [fileId, onDetectionComplete, onError, stopPolling]);

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
      <p className="text-sm text-muted-foreground">Analyzing document...</p>
    </div>
  );
}
