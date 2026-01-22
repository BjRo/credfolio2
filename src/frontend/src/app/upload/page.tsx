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
    <div className="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-3xl mx-auto">
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-gray-900">Upload Reference Letter</h1>
          <p className="mt-2 text-sm text-gray-600">
            Upload your reference letter document for AI-powered analysis
          </p>
        </div>

        <div className="bg-white shadow rounded-lg p-6 mb-8">
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
          <div className="bg-white shadow rounded-lg p-6">
            <h2 className="text-lg font-medium text-gray-900 mb-4">Recent Uploads</h2>
            <ul className="divide-y divide-gray-200">
              {uploads.map((upload) => (
                <li key={upload.id} className="py-3 flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-gray-900">{upload.filename}</p>
                    <p className="text-xs text-gray-500">ID: {upload.id}</p>
                  </div>
                  <span
                    className={`px-2 py-1 text-xs font-medium rounded-full ${
                      upload.status === "COMPLETED"
                        ? "bg-green-100 text-green-800"
                        : upload.status === "PROCESSING"
                          ? "bg-blue-100 text-blue-800"
                          : upload.status === "FAILED"
                            ? "bg-red-100 text-red-800"
                            : "bg-yellow-100 text-yellow-800"
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
