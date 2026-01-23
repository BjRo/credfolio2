"use client";

import { type ChangeEvent, type DragEvent, useCallback, useState } from "react";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

const SUPPORTED_TYPES = ["image/jpeg", "image/png", "image/gif", "image/webp", "application/pdf"];

const MAX_FILE_SIZE = 20 * 1024 * 1024; // 20MB

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

interface ExtractionResult {
  text: string;
  inputTokens: number;
  outputTokens: number;
}

function validateFile(file: File): string | null {
  if (!SUPPORTED_TYPES.includes(file.type)) {
    return `Unsupported file type: ${file.type}. Supported: JPEG, PNG, GIF, WebP, PDF`;
  }
  if (file.size > MAX_FILE_SIZE) {
    return `File too large (${(file.size / 1024 / 1024).toFixed(1)}MB). Maximum size is 20MB.`;
  }
  return null;
}

export default function ExtractTestPage() {
  const [file, setFile] = useState<File | null>(null);
  const [isDragging, setIsDragging] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [result, setResult] = useState<ExtractionResult | null>(null);
  const [error, setError] = useState<string | null>(null);

  const handleDragOver = useCallback((e: DragEvent) => {
    e.preventDefault();
    setIsDragging(true);
  }, []);

  const handleDragLeave = useCallback((e: DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
  }, []);

  const handleDrop = useCallback((e: DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
    setError(null);
    setResult(null);

    const droppedFile = e.dataTransfer.files[0];
    if (droppedFile) {
      const validationError = validateFile(droppedFile);
      if (validationError) {
        setError(validationError);
        return;
      }
      setFile(droppedFile);
    }
  }, []);

  const handleFileChange = useCallback((e: ChangeEvent<HTMLInputElement>) => {
    setError(null);
    setResult(null);

    const selectedFile = e.target.files?.[0];
    if (selectedFile) {
      const validationError = validateFile(selectedFile);
      if (validationError) {
        setError(validationError);
        return;
      }
      setFile(selectedFile);
    }
  }, []);

  const handleExtract = async () => {
    if (!file) return;

    setIsLoading(true);
    setError(null);
    setResult(null);

    try {
      const formData = new FormData();
      formData.append("file", file);

      const response = await fetch(`${API_URL}/api/extract`, {
        method: "POST",
        body: formData,
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || "Extraction failed");
      }

      setResult(data as ExtractionResult);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Unknown error occurred");
    } finally {
      setIsLoading(false);
    }
  };

  const handleClear = () => {
    setFile(null);
    setResult(null);
    setError(null);
  };

  return (
    <div className="min-h-screen bg-gray-50 py-12 px-4">
      <div className="max-w-3xl mx-auto">
        <h1 className="text-3xl font-bold mb-2">Document Extraction Test</h1>
        <p className="text-gray-600 mb-8">
          Upload an image or PDF to extract text using Claude&apos;s vision capabilities.
        </p>

        {/* Drop Zone */}
        <label
          htmlFor="file-input"
          className={cn(
            "block border-2 border-dashed rounded-lg p-8 text-center transition-colors cursor-pointer mb-6",
            isDragging && "border-blue-500 bg-blue-50",
            file && !isDragging && "border-green-500 bg-green-50",
            !file && !isDragging && "border-gray-300 hover:border-gray-400"
          )}
          onDragOver={handleDragOver}
          onDragLeave={handleDragLeave}
          onDrop={handleDrop}
        >
          <input
            type="file"
            accept={SUPPORTED_TYPES.join(",")}
            onChange={handleFileChange}
            className="hidden"
            id="file-input"
          />
          {file ? (
            <div>
              <p className="font-medium text-green-700">{file.name}</p>
              <p className="text-sm text-gray-500 mt-1">
                {(file.size / 1024).toFixed(1)} KB - {file.type}
              </p>
            </div>
          ) : (
            <div>
              <p className="text-gray-600">Drop a file here or click to select</p>
              <p className="text-sm text-gray-400 mt-2">
                Supports: JPEG, PNG, GIF, WebP, PDF (max 20MB)
              </p>
            </div>
          )}
        </label>

        {/* Action Buttons */}
        <div className="flex gap-4 mb-6">
          <Button onClick={handleExtract} disabled={!file || isLoading}>
            {isLoading ? "Extracting..." : "Extract Text"}
          </Button>
          {(file || result || error) && (
            <Button variant="outline" onClick={handleClear}>
              Clear
            </Button>
          )}
        </div>

        {/* Loading State */}
        {isLoading && (
          <div className="flex items-center justify-center p-8 bg-white rounded-lg shadow">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600" />
            <span className="ml-3 text-gray-600">Extracting text from document...</span>
          </div>
        )}

        {/* Error Display */}
        {error && (
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mb-6">
            {error}
          </div>
        )}

        {/* Result Display */}
        {result && (
          <div className="bg-white shadow rounded-lg p-6">
            <h2 className="text-lg font-medium mb-4">Extracted Text</h2>
            <div className="flex gap-4 text-sm text-gray-500 mb-4">
              <span>Input tokens: {result.inputTokens.toLocaleString()}</span>
              <span>Output tokens: {result.outputTokens.toLocaleString()}</span>
            </div>
            <pre className="whitespace-pre-wrap bg-gray-50 p-4 rounded text-sm overflow-auto max-h-96 font-mono">
              {result.text}
            </pre>
          </div>
        )}
      </div>
    </div>
  );
}
