"use client";

import { useState } from "react";
import { FileUpload } from "@/components";

// Demo user ID for testing
const DEMO_USER_ID = "00000000-0000-0000-0000-000000000001";

export default function UploadPage() {
  const [uploads, setUploads] = useState<Array<{ id: string; filename: string; status: string }>>(
    []
  );

  return (
    <div className="min-h-screen bg-muted/50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-3xl mx-auto">
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-foreground">Upload Reference Letter</h1>
          <p className="mt-2 text-sm text-muted-foreground">
            Upload your reference letter document for AI-powered analysis
          </p>
        </div>

        <div className="bg-card shadow rounded-lg p-6 mb-8">
          <FileUpload
            userId={DEMO_USER_ID}
            onUploadComplete={(result) => {
              setUploads((prev) => [
                {
                  id: result.referenceLetter.id,
                  filename: result.file.filename,
                  status: result.referenceLetter.status,
                },
                ...prev,
              ]);
            }}
            onError={(error) => {
              console.error("Upload error:", error);
            }}
          />
        </div>

        {uploads.length > 0 && (
          <div className="bg-card shadow rounded-lg p-6">
            <h2 className="text-lg font-medium text-foreground mb-4">Recent Uploads</h2>
            <ul className="divide-y divide-border">
              {uploads.map((upload) => (
                <li key={upload.id} className="py-3 flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-foreground">{upload.filename}</p>
                    <p className="text-xs text-muted-foreground">ID: {upload.id}</p>
                  </div>
                  <span
                    className={`px-2 py-1 text-xs font-medium rounded-full ${
                      upload.status === "COMPLETED"
                        ? "bg-primary/10 text-primary"
                        : upload.status === "PROCESSING"
                          ? "bg-primary/20 text-primary"
                          : upload.status === "FAILED"
                            ? "bg-destructive/10 text-destructive"
                            : "bg-muted text-muted-foreground"
                    }`}
                  >
                    {upload.status}
                  </span>
                </li>
              ))}
            </ul>
          </div>
        )}
      </div>
    </div>
  );
}
