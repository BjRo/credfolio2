"use client";

import { UploadFlow } from "@/components/upload";

export default function UploadPage() {
  return (
    <div className="min-h-screen bg-muted/50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-3xl mx-auto">
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-foreground">Upload Document</h1>
          <p className="mt-2 text-sm text-muted-foreground">
            Upload a resume, reference letter, or any career document to get started
          </p>
        </div>

        <UploadFlow />
      </div>
    </div>
  );
}
