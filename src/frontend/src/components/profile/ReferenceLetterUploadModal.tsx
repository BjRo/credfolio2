"use client";

import { Upload } from "lucide-react";
import { type ChangeEvent, type DragEvent, useCallback, useEffect, useState } from "react";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
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
type ReferenceLetterStatus = "PENDING" | "PROCESSING" | "COMPLETED" | "FAILED";

interface UploadFileResult {
  __typename: "UploadFileResult";
  file: {
    id: string;
    filename: string;
  };
  referenceLetter: {
    id: string;
    status: ReferenceLetterStatus;
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
  existingReferenceLetter?: {
    id: string;
    status: ReferenceLetterStatus;
  };
  message: string;
}

interface ReferenceLetterUploadModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  userId: string;
  onSuccess?: (referenceLetterld: string) => void;
}

export function ReferenceLetterUploadModal({
  open,
  onOpenChange,
  userId,
  onSuccess,
}: ReferenceLetterUploadModalProps) {
  const [isDragOver, setIsDragOver] = useState(false);
  const [status, setStatus] = useState<UploadStatus>("idle");
  const [progress, setProgress] = useState(0);
  const [error, setError] = useState<string | null>(null);
  const [uploadedLetter, setUploadedLetter] = useState<UploadFileResult | null>(null);
  const [duplicateInfo, setDuplicateInfo] = useState<DuplicateFileDetected | null>(null);
  const [pendingFile, setPendingFile] = useState<File | null>(null);

  // Reset state when modal opens/closes
  useEffect(() => {
    if (!open) {
      setStatus("idle");
      setProgress(0);
      setError(null);
      setUploadedLetter(null);
      setDuplicateInfo(null);
      setPendingFile(null);
      setIsDragOver(false);
    }
  }, [open]);

  const validateFile = useCallback((file: File): string | null => {
    if (!Object.keys(ALLOWED_TYPES).includes(file.type)) {
      return `Invalid file type. Allowed types: ${ALLOWED_EXTENSIONS}`;
    }
    if (file.size > MAX_SIZE_BYTES) {
      return `File too large. Maximum size is ${MAX_SIZE_BYTES / (1024 * 1024)}MB`;
    }
    return null;
  }, []);

  // Poll for reference letter status
  useEffect(() => {
    if (status !== "processing" || !uploadedLetter) return;

    const pollStatus = async () => {
      try {
        const response = await fetch(GRAPHQL_ENDPOINT, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            query: `
              query GetReferenceLetterStatus($id: ID!) {
                referenceLetter(id: $id) {
                  id
                  status
                }
              }
            `,
            variables: { id: uploadedLetter.referenceLetter.id },
          }),
        });

        const result = await response.json();
        const letterStatus = result.data?.referenceLetter?.status as
          | ReferenceLetterStatus
          | undefined;

        if (letterStatus === "COMPLETED") {
          setStatus("success");
          onSuccess?.(uploadedLetter.referenceLetter.id);
          onOpenChange(false);
        } else if (letterStatus === "FAILED") {
          setStatus("error");
          setError("Reference letter processing failed. Please try again.");
        }
      } catch (err) {
        console.error("Failed to poll reference letter status:", err);
      }
    };

    const intervalId = setInterval(pollStatus, POLL_INTERVAL_MS);
    return () => clearInterval(intervalId);
  }, [status, uploadedLetter, onSuccess, onOpenChange]);

  const uploadFile = useCallback(
    async (file: File, forceReimport = false) => {
      const validationError = validateFile(file);
      if (validationError) {
        setError(validationError);
        setStatus("error");
        return;
      }

      setStatus("uploading");
      setProgress(0);
      setError(null);
      setDuplicateInfo(null);

      const operations = JSON.stringify({
        query: `
          mutation UploadReferenceLetter($userId: ID!, $file: Upload!, $forceReimport: Boolean) {
            uploadFile(userId: $userId, file: $file, forceReimport: $forceReimport) {
              ... on UploadFileResult {
                __typename
                file {
                  id
                  filename
                  contentType
                  sizeBytes
                }
                referenceLetter {
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
                existingReferenceLetter {
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
        // NOTE: Uses XMLHttpRequest for file uploads (see DocumentUpload.tsx for detailed explanation)
        // Follows GraphQL multipart request spec for file uploads with progress tracking.
        const result = await new Promise<UploadFileResult | DuplicateFileDetected>(
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
                  const data = response.data?.uploadFile;
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
                  resolve(data as UploadFileResult | DuplicateFileDetected);
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
        setUploadedLetter(result);
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : "Upload failed";
        setError(errorMessage);
        setStatus("error");
      }
    },
    [userId, validateFile]
  );

  const handleForceReimport = useCallback(() => {
    if (pendingFile) {
      uploadFile(pendingFile, true);
    }
  }, [pendingFile, uploadFile]);

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
    setUploadedLetter(null);
    setDuplicateInfo(null);
    setPendingFile(null);
  }, []);

  const handleCancel = useCallback(() => {
    handleReset();
    onOpenChange(false);
  }, [handleReset, onOpenChange]);

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>Add Reference Letter</DialogTitle>
          <DialogDescription>
            Upload a reference letter to validate your profile. The letter will be processed to
            extract skills, experiences, and testimonials.
          </DialogDescription>
        </DialogHeader>

        <div className="py-4">
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
                <h3 className="text-lg font-medium text-warning-foreground">
                  Duplicate File Detected
                </h3>
              </div>
              <p className="text-sm text-warning-foreground/80 mb-4">{duplicateInfo.message}</p>
              {duplicateInfo.existingReferenceLetter && (
                <p className="text-sm text-warning-foreground/80 mb-4">
                  Previous extraction status:{" "}
                  <span className="font-medium">
                    {duplicateInfo.existingReferenceLetter.status}
                  </span>
                </p>
              )}
              <div className="flex gap-3 mt-4">
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
          ) : status === "processing" && uploadedLetter ? (
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
                <h3 className="text-lg font-medium text-warning-foreground">
                  Processing Reference Letter
                </h3>
              </div>
              <p className="text-sm text-warning-foreground/80 mb-2">
                File &quot;{uploadedLetter.file.filename}&quot; uploaded successfully.
              </p>
              <p className="text-sm text-warning-foreground/80">
                Extracting information from your reference letter...
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
              <p className="text-sm text-success/80">
                Your reference letter has been processed successfully.
              </p>
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
                data-testid="file-input"
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
                    <p className="mt-4 text-sm font-medium text-foreground">
                      Uploading... {progress}%
                    </p>
                    <div className="w-full bg-muted rounded-full h-2 mt-2">
                      <div
                        className="bg-warning h-2 rounded-full transition-all duration-300"
                        style={{ width: `${progress}%` }}
                      />
                    </div>
                  </>
                ) : (
                  <>
                    <Upload className="w-12 h-12 mx-auto text-muted-foreground" />
                    <p className="mt-4 text-sm font-medium text-foreground">
                      {isDragOver
                        ? "Drop your reference letter here"
                        : "Drag and drop your reference letter, or click to browse"}
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

        <DialogFooter>
          <Button variant="outline" onClick={handleCancel} disabled={status === "processing"}>
            Cancel
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
