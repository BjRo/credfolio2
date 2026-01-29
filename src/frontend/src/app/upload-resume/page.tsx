"use client";

import { ResumeUpload } from "@/components";

// Demo user ID for testing
const DEMO_USER_ID = "00000000-0000-0000-0000-000000000001";

export default function UploadResumePage() {
  return (
    <div className="min-h-screen bg-muted/50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-3xl mx-auto">
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-foreground">Upload Your Resume</h1>
          <p className="mt-2 text-sm text-muted-foreground">
            Upload your resume to extract your professional profile automatically
          </p>
        </div>

        <div className="bg-card shadow rounded-lg p-6">
          <ResumeUpload
            userId={DEMO_USER_ID}
            onUploadComplete={(result) => {
              console.log("Upload complete:", result);
            }}
            onProcessingComplete={(resumeId) => {
              console.log("Processing complete, redirecting to:", resumeId);
            }}
            onError={(error) => {
              console.error("Upload error:", error);
            }}
          />
        </div>

        <div className="mt-8 text-center">
          <p className="text-sm text-muted-foreground">
            Your resume will be processed by AI to extract your work experience, education, and
            skills.
          </p>
        </div>
      </div>
    </div>
  );
}
