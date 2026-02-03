"use client";

import { useRouter } from "next/navigation";
import { type ChangeEvent, type DragEvent, useCallback, useEffect, useState } from "react";
import { GRAPHQL_ENDPOINT } from "@/lib/urql/client";

const ALLOWED_TYPES = {
  "application/pdf": ".pdf",
  "application/vnd.openxmlformats-officedocument.wordprocessingml.document": ".docx",
  "text/plain": ".txt",
};

const ALLOWED_EXTENSIONS = Object.values(ALLOWED_TYPES).join(", ");
const MAX_SIZE_BYTES = 10 * 1024 * 1024; // 10MB
const POLL_INTERVAL_MS = 2000; // Poll every 2 seconds

type UploadStatus = "idle" | "uploading" | "processing" | "success" | "error" | "duplicate";
type ResumeStatus = "PENDING" | "PROCESSING" | "COMPLETED" | "FAILED";

interface ResumeUploadResult {
  __typename: "UploadResumeResult";
  file: {
    id: string;
    filename: string;
  };
  resume: {
    id: string;
    status: ResumeStatus;
  };
}

interface ValidationError {
  __typename: "FileValidationError";
  message: string;
  field: string;
}

interface DuplicateFileDetected {
  __typename: "DuplicateFileDetected";
  existingFile: {
    id: string;
    filename: string;
    createdAt: string;
  };
  existingResume?: {
    id: string;
    status: ResumeStatus;
  };
  message: string;
}

interface ResumeUploadProps {
  userId: string;
  onUploadComplete?: (result: ResumeUploadResult) => void;
  onProcessingComplete?: (resumeId: string) => void;
  onError?: (error: string) => void;
}

export function ResumeUpload({
  userId,
  onUploadComplete,
  onProcessingComplete,
  onError,
}: ResumeUploadProps) {
  const router = useRouter();
  const [isDragOver, setIsDragOver] = useState(false);
  const [status, setStatus] = useState<UploadStatus>("idle");
  const [progress, setProgress] = useState(0);
  const [error, setError] = useState<string | null>(null);
  const [uploadedResume, setUploadedResume] = useState<ResumeUploadResult | null>(null);
  const [duplicateInfo, setDuplicateInfo] = useState<DuplicateFileDetected | null>(null);
  const [pendingFile, setPendingFile] = useState<File | null>(null);

  const validateFile = useCallback((file: File): string | null => {
    if (!Object.keys(ALLOWED_TYPES).includes(file.type)) {
      return `Invalid file type. Allowed types: ${ALLOWED_EXTENSIONS}`;
    }
    if (file.size > MAX_SIZE_BYTES) {
      return `File too large. Maximum size is ${MAX_SIZE_BYTES / (1024 * 1024)}MB`;
    }
    return null;
  }, []);

  // Poll for resume status
  useEffect(() => {
    if (status !== "processing" || !uploadedResume) return;

    const pollStatus = async () => {
      try {
        const response = await fetch(GRAPHQL_ENDPOINT, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            query: `
              query GetResumeStatus($id: ID!) {
                resume(id: $id) {
                  id
                  status
                }
              }
            `,
            variables: { id: uploadedResume.resume.id },
          }),
        });

        const result = await response.json();
        const resumeStatus = result.data?.resume?.status as ResumeStatus | undefined;

        if (resumeStatus === "COMPLETED") {
          setStatus("success");
          onProcessingComplete?.(uploadedResume.resume.id);
          // Auto-redirect to profile page
          router.push(`/profile/${uploadedResume.resume.id}`);
        } else if (resumeStatus === "FAILED") {
          setStatus("error");
          setError("Resume processing failed. Please try again.");
          onError?.("Resume processing failed");
        }
      } catch (err) {
        console.error("Failed to poll resume status:", err);
      }
    };

    const intervalId = setInterval(pollStatus, POLL_INTERVAL_MS);
    return () => clearInterval(intervalId);
  }, [status, uploadedResume, router, onProcessingComplete, onError]);

  const uploadFile = useCallback(
    async (file: File, forceReimport = false) => {
      const validationError = validateFile(file);
      if (validationError) {
        setError(validationError);
        setStatus("error");
        onError?.(validationError);
        return;
      }

      setStatus("uploading");
      setProgress(0);
      setError(null);
      setDuplicateInfo(null);

      const operations = JSON.stringify({
        query: `
          mutation UploadResume($userId: ID!, $file: Upload!, $forceReimport: Boolean) {
            uploadResume(userId: $userId, file: $file, forceReimport: $forceReimport) {
              ... on UploadResumeResult {
                __typename
                file {
                  id
                  filename
                  contentType
                  sizeBytes
                }
                resume {
                  id
                  status
                }
              }
              ... on FileValidationError {
                __typename
                message
                field
              }
              ... on DuplicateFileDetected {
                __typename
                existingFile {
                  id
                  filename
                  createdAt
                }
                existingResume {
                  id
                  status
                }
                message
              }
            }
          }
        `,
        variables: {
          userId,
          file: null,
          forceReimport: forceReimport || null,
        },
      });

      const map = JSON.stringify({
        "0": ["variables.file"],
      });

      const formData = new FormData();
      formData.append("operations", operations);
      formData.append("map", map);
      formData.append("0", file);

      try {
        const result = await new Promise<ResumeUploadResult | DuplicateFileDetected>(
          (resolve, reject) => {
            const xhr = new XMLHttpRequest();

            xhr.upload.addEventListener("progress", (event) => {
              if (event.lengthComputable) {
                const percentComplete = Math.round((event.loaded / event.total) * 100);
                setProgress(percentComplete);
              }
            });

            xhr.addEventListener("load", () => {
              if (xhr.status >= 200 && xhr.status < 300) {
                try {
                  const response = JSON.parse(xhr.responseText);
                  if (response.errors?.length) {
                    reject(new Error(response.errors[0].message));
                    return;
                  }
                  const data = response.data?.uploadResume;
                  if (!data) {
                    reject(new Error("No data returned from upload"));
                    return;
                  }
                  // Check for validation error union type
                  if (data.__typename === "FileValidationError") {
                    const validationErr = data as ValidationError;
                    reject(new Error(validationErr.message));
                    return;
                  }
                  // Handle both success and duplicate cases
                  resolve(data as ResumeUploadResult | DuplicateFileDetected);
                } catch (_parseError) {
                  reject(new Error("Failed to parse response"));
                }
              } else {
                reject(new Error(`Upload failed with status ${xhr.status}`));
              }
            });

            xhr.addEventListener("error", () => {
              reject(new Error("Network error during upload"));
            });

            xhr.addEventListener("abort", () => {
              reject(new Error("Upload was cancelled"));
            });

            xhr.open("POST", GRAPHQL_ENDPOINT);
            xhr.send(formData);
          }
        );

        // Handle duplicate detection
        if (result.__typename === "DuplicateFileDetected") {
          setStatus("duplicate");
          setDuplicateInfo(result);
          setPendingFile(file);
          return;
        }

        setStatus("processing");
        setUploadedResume(result);
        onUploadComplete?.(result);
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : "Upload failed";
        setError(errorMessage);
        setStatus("error");
        onError?.(errorMessage);
      }
    },
    [userId, validateFile, onUploadComplete, onError]
  );

  const handleForceReimport = useCallback(() => {
    if (pendingFile) {
      uploadFile(pendingFile, true);
    }
  }, [pendingFile, uploadFile]);

  const handleViewExisting = useCallback(() => {
    if (duplicateInfo?.existingResume) {
      router.push(`/profile/${duplicateInfo.existingResume.id}`);
    }
  }, [duplicateInfo, router]);

  const handleDragOver = useCallback((e: DragEvent<HTMLLabelElement>) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragOver(true);
  }, []);

  const handleDragLeave = useCallback((e: DragEvent<HTMLLabelElement>) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragOver(false);
  }, []);

  const handleDrop = useCallback(
    (e: DragEvent<HTMLLabelElement>) => {
      e.preventDefault();
      e.stopPropagation();
      setIsDragOver(false);

      const files = e.dataTransfer.files;
      if (files.length > 0) {
        uploadFile(files[0]);
      }
    },
    [uploadFile]
  );

  const handleFileSelect = useCallback(
    (e: ChangeEvent<HTMLInputElement>) => {
      const files = e.target.files;
      if (files && files.length > 0) {
        uploadFile(files[0]);
      }
      e.target.value = "";
    },
    [uploadFile]
  );

  const handleReset = useCallback(() => {
    setStatus("idle");
    setProgress(0);
    setError(null);
    setUploadedResume(null);
    setDuplicateInfo(null);
    setPendingFile(null);
  }, []);

  return (
    <div className="w-full max-w-xl mx-auto">
      {status === "duplicate" && duplicateInfo ? (
        <div className="p-6 border-2 border-warning bg-warning/10 rounded-lg">
          <div className="flex items-center gap-3 mb-4">
            <svg
              role="img"
              aria-label="Warning"
              className="w-8 h-8 text-warning"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
              />
            </svg>
            <h3 className="text-lg font-medium text-warning-foreground">Duplicate File Detected</h3>
          </div>
          <p className="text-sm text-warning-foreground/80 mb-4">{duplicateInfo.message}</p>
          {duplicateInfo.existingResume && (
            <p className="text-sm text-warning-foreground/80 mb-4">
              Previous extraction status:{" "}
              <span className="font-medium">{duplicateInfo.existingResume.status}</span>
            </p>
          )}
          <div className="flex gap-3 mt-4">
            {duplicateInfo.existingResume && (
              <button
                type="button"
                onClick={handleViewExisting}
                className="px-4 py-2 text-sm font-medium text-primary bg-primary/10 border border-primary/30 rounded-md hover:bg-primary/20 transition-colors"
              >
                View Existing Profile
              </button>
            )}
            <button
              type="button"
              onClick={handleForceReimport}
              className="px-4 py-2 text-sm font-medium text-warning-foreground bg-warning/20 border border-warning/50 rounded-md hover:bg-warning/30 transition-colors"
            >
              Re-import Anyway
            </button>
            <button
              type="button"
              onClick={handleReset}
              className="px-4 py-2 text-sm font-medium text-muted-foreground bg-muted border border-border rounded-md hover:bg-muted/80 transition-colors"
            >
              Cancel
            </button>
          </div>
        </div>
      ) : status === "processing" && uploadedResume ? (
        <div className="p-6 border-2 border-warning bg-warning/10 rounded-lg">
          <div className="flex items-center gap-3 mb-4">
            <svg
              role="img"
              aria-label="Processing"
              className="w-8 h-8 text-warning animate-spin"
              fill="none"
              viewBox="0 0 24 24"
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
            <h3 className="text-lg font-medium text-warning-foreground">Processing Resume</h3>
          </div>
          <p className="text-sm text-warning-foreground/80 mb-2">
            File &quot;{uploadedResume.file.filename}&quot; uploaded successfully.
          </p>
          <p className="text-sm text-warning-foreground/80">
            Extracting profile information... You&apos;ll be redirected automatically when complete.
          </p>
        </div>
      ) : status === "success" ? (
        <div className="p-6 border-2 border-success bg-success/10 rounded-lg">
          <div className="flex items-center gap-3 mb-4">
            <svg
              role="img"
              aria-label="Success"
              className="w-8 h-8 text-success"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M5 13l4 4L19 7"
              />
            </svg>
            <h3 className="text-lg font-medium text-success">Processing Complete</h3>
          </div>
          <p className="text-sm text-success/80 mb-4">Redirecting to your profile...</p>
        </div>
      ) : (
        <label
          onDragOver={handleDragOver}
          onDragLeave={handleDragLeave}
          onDrop={handleDrop}
          className={`
            relative p-8 border-2 border-dashed rounded-lg transition-colors cursor-pointer block
            ${
              isDragOver
                ? "border-primary bg-primary/10"
                : status === "error"
                  ? "border-destructive/50 bg-destructive/10"
                  : "border-border hover:border-primary/50 bg-card"
            }
          `}
        >
          <input
            type="file"
            accept={Object.keys(ALLOWED_TYPES).join(",")}
            onChange={handleFileSelect}
            className="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
            disabled={status === "uploading"}
          />

          <div className="text-center">
            {status === "uploading" ? (
              <>
                <svg
                  role="img"
                  aria-label="Uploading"
                  className="w-12 h-12 mx-auto text-warning animate-spin"
                  fill="none"
                  viewBox="0 0 24 24"
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
                <p className="mt-4 text-sm font-medium text-foreground">Uploading... {progress}%</p>
                <div className="w-full bg-muted rounded-full h-2 mt-2">
                  <div
                    className="bg-warning h-2 rounded-full transition-all duration-300"
                    style={{ width: `${progress}%` }}
                  />
                </div>
              </>
            ) : (
              <>
                <svg
                  role="img"
                  aria-label="Upload"
                  className="w-12 h-12 mx-auto text-muted-foreground"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"
                  />
                </svg>
                <p className="mt-4 text-sm font-medium text-foreground">
                  {isDragOver
                    ? "Drop your resume here"
                    : "Drag and drop your resume, or click to browse"}
                </p>
                <p className="mt-2 text-xs text-muted-foreground">
                  Supported formats: PDF, DOCX, TXT (max 10MB)
                </p>
              </>
            )}
          </div>

          {status === "error" && error && (
            <div className="mt-4 p-3 bg-destructive/10 border border-destructive/30 rounded-md">
              <div className="flex items-start gap-2">
                <svg
                  role="img"
                  aria-label="Error"
                  className="w-5 h-5 text-destructive flex-shrink-0 mt-0.5"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                  />
                </svg>
                <div>
                  <p className="text-sm font-medium text-destructive">Upload failed</p>
                  <p className="text-sm text-destructive/80">{error}</p>
                </div>
              </div>
              <button
                type="button"
                onClick={handleReset}
                className="mt-2 text-sm text-destructive hover:text-destructive/80 underline"
              >
                Try again
              </button>
            </div>
          )}
        </label>
      )}
    </div>
  );
}
