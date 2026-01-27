"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { useQuery } from "urql";
import { ResumeUpload } from "@/components";
import { GetUserResumesDocument, ResumeStatus } from "@/graphql/generated/graphql";

const DEMO_USER_ID = "00000000-0000-0000-0000-000000000001";

export default function Home() {
  const router = useRouter();
  const [result] = useQuery({
    query: GetUserResumesDocument,
    variables: { userId: DEMO_USER_ID },
  });

  const { fetching, data } = result;

  const completedResume = data?.resumes?.find((r) => r.status === ResumeStatus.Completed);

  useEffect(() => {
    if (completedResume) {
      router.push(`/profile/${completedResume.id}`);
    }
  }, [completedResume, router]);

  if (fetching) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-gray-50">
        <output className="animate-spin h-8 w-8 border-4 border-blue-500 border-t-transparent rounded-full block" />
      </div>
    );
  }

  if (completedResume) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-gray-50">
        <p className="text-gray-600">Redirecting to your profile...</p>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-3xl mx-auto">
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-gray-900">Upload Your Resume</h1>
          <p className="mt-2 text-sm text-gray-600">
            Upload your resume to extract your professional profile automatically
          </p>
        </div>

        <div className="bg-white shadow rounded-lg p-6">
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
          <p className="text-sm text-gray-500">
            Your resume will be processed by AI to extract your work experience, education, and
            skills.
          </p>
        </div>
      </div>
    </div>
  );
}
