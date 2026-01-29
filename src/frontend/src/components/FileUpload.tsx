"use client";

import { type ChangeEvent, type DragEvent, useCallback, useState } from "react";
import { GRAPHQL_ENDPOINT } from "@/lib/urql/client";

const ALLOWED_TYPES = {
  "application/pdf": ".pdf",
  "application/vnd.openxmlformats-officedocument.wordprocessingml.document": ".docx",
  "text/plain": ".txt",
};

const ALLOWED_EXTENSIONS = Object.values(ALLOWED_TYPES).join(", ");
const MAX_SIZE_BYTES = 10 * 1024 * 1024; // 10MB

type UploadStatus = "idle" | "uploading" | "success" | "error";

interface UploadResult {
  file: {
    id: string;
    filename: string;
  };
  referenceLetter: {
    id: string;
    status: string;
  };
}

interface ValidationError {
  message: string;
  field: string;
}

interface FileUploadProps {
  userId: string;
  onUploadComplete?: (result: UploadResult) => void;
  onError?: (error: string) => void;
}

export function FileUpload({ userId, onUploadComplete, onError }: FileUploadProps) {
  const [isDragOver, setIsDragOver] = useState(false);
  const [status, setStatus] = useState<UploadStatus>("idle");
  const [progress, setProgress] = useState(0);
  const [error, setError] = useState<string | null>(null);
  const [uploadedFile, setUploadedFile] = useState<UploadResult | null>(null);

  const validateFile = useCallback((file: File): string | null => {
    if (!Object.keys(ALLOWED_TYPES).includes(file.type)) {
      return `Invalid file type. Allowed types: ${ALLOWED_EXTENSIONS}`;
    }
    if (file.size > MAX_SIZE_BYTES) {
      return `File too large. Maximum size is ${MAX_SIZE_BYTES / (1024 * 1024)}MB`;
    }
    return null;
  }, []);

  const uploadFile = useCallback(
    async (file: File) => {
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

      const operations = JSON.stringify({
        query: `
        mutation UploadFile($userId: ID!, $file: Upload!) {
          uploadFile(userId: $userId, file: $file) {
            ... on UploadFileResult {
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
              message
              field
            }
          }
        }
      `,
        variables: {
          userId,
          file: null,
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
        // Use XMLHttpRequest for progress tracking
        const result = await new Promise<UploadResult>((resolve, reject) => {
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
                if ("message" in data && "field" in data) {
                  const validationErr = data as ValidationError;
                  reject(new Error(validationErr.message));
                  return;
                }
                resolve(data as UploadResult);
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
        });

        setStatus("success");
        setUploadedFile(result);
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
      // Reset input so same file can be selected again
      e.target.value = "";
    },
    [uploadFile]
  );

  const handleReset = useCallback(() => {
    setStatus("idle");
    setProgress(0);
    setError(null);
    setUploadedFile(null);
  }, []);

  return (
    <div className="w-full max-w-xl mx-auto">
      {status === "success" && uploadedFile ? (
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
            <h3 className="text-lg font-medium text-success">Upload Complete</h3>
          </div>
          <p className="text-sm text-success/80 mb-4">
            File "{uploadedFile.file.filename}" uploaded successfully and is now being processed.
          </p>
          <button
            type="button"
            onClick={handleReset}
            className="px-4 py-2 text-sm font-medium text-success bg-card border border-success rounded-md hover:bg-success/10 focus:outline-none focus:ring-2 focus:ring-success focus:ring-offset-2"
          >
            Upload Another File
          </button>
        </div>
      ) : (
        <label
          onDragOver={handleDragOver}
          onDragLeave={handleDragLeave}
          onDrop={handleDrop}
          className={`
            relative p-8 border-2 border-dashed rounded-lg transition-colors cursor-pointer
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
                  {isDragOver ? "Drop file here" : "Drag and drop a file here, or click to browse"}
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
