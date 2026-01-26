"use client";

import { useParams, useRouter } from "next/navigation";
import { useQuery } from "urql";
import {
  EditableWorkExperienceSection,
  EducationSection,
  ProfileActions,
  ProfileHeader,
  ProfileSkeleton,
  SkillsSection,
  WorkExperienceSection,
} from "@/components/profile";
import { Button } from "@/components/ui/button";
import { GetProfileDocument, GetResumeDocument, ResumeStatus } from "@/graphql/generated/graphql";

export default function ProfilePage() {
  const params = useParams();
  const router = useRouter();
  const resumeId = params.id as string;

  const [resumeResult, _reexecuteResumeQuery] = useQuery({
    query: GetResumeDocument,
    variables: { id: resumeId },
  });

  // Get user ID from resume to fetch their profile
  const userId = resumeResult.data?.resume?.user?.id;

  const [profileResult, reexecuteProfileQuery] = useQuery({
    query: GetProfileDocument,
    variables: { userId: userId || "" },
    pause: !userId, // Don't run until we have userId
  });

  const { data, fetching, error } = resumeResult;
  const profile = profileResult.data?.profile;

  // Refetch profile when mutations succeed
  const handleMutationSuccess = () => {
    reexecuteProfileQuery({ requestPolicy: "network-only" });
  };

  if (fetching) {
    return (
      <div className="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-4xl mx-auto">
          <ProfileSkeleton />
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-4xl mx-auto text-center">
          <h1 className="text-2xl font-bold text-red-600 mb-4">Error Loading Profile</h1>
          <p className="text-gray-600 mb-6">{error.message}</p>
          <Button onClick={() => window.location.reload()}>Try Again</Button>
        </div>
      </div>
    );
  }

  const resume = data?.resume;

  if (!resume) {
    return (
      <div className="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-4xl mx-auto text-center">
          <h1 className="text-2xl font-bold text-gray-900 mb-4">Profile Not Found</h1>
          <p className="text-gray-600 mb-6">
            The resume you&apos;re looking for doesn&apos;t exist or has been removed.
          </p>
          <Button onClick={() => router.push("/upload-resume")}>Upload a Resume</Button>
        </div>
      </div>
    );
  }

  if (resume.status === ResumeStatus.Pending || resume.status === ResumeStatus.Processing) {
    return (
      <div className="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-4xl mx-auto text-center">
          <h1 className="text-2xl font-bold text-gray-900 mb-4">Profile Processing</h1>
          <p className="text-gray-600 mb-6">Your resume is still being processed. Please wait...</p>
          <div className="animate-spin h-8 w-8 border-4 border-blue-500 border-t-transparent rounded-full mx-auto" />
        </div>
      </div>
    );
  }

  if (resume.status === ResumeStatus.Failed) {
    return (
      <div className="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-4xl mx-auto text-center">
          <h1 className="text-2xl font-bold text-red-600 mb-4">Processing Failed</h1>
          <p className="text-gray-600 mb-6">{resume.errorMessage || "Failed to process resume"}</p>
          <Button onClick={() => router.push("/upload-resume")}>Try Again</Button>
        </div>
      </div>
    );
  }

  const extractedData = resume.extractedData;

  if (!extractedData) {
    return (
      <div className="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-4xl mx-auto text-center">
          <h1 className="text-2xl font-bold text-gray-900 mb-4">No Profile Data</h1>
          <p className="text-gray-600 mb-6">
            The resume was processed but no data could be extracted.
          </p>
          <Button onClick={() => router.push("/upload-resume")}>Upload a Different Resume</Button>
        </div>
      </div>
    );
  }

  const handleAddReference = () => {
    router.push("/upload");
  };

  const handleExport = () => {
    // TODO: Implement PDF export
    alert("PDF export coming soon!");
  };

  const handleUploadAnother = () => {
    router.push("/upload-resume");
  };

  // Use profile experiences if available (editable), otherwise fall back to extracted data
  const hasProfileExperiences = profile && profile.experiences.length > 0;

  return (
    <div className="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-4xl mx-auto space-y-6">
        <ProfileHeader data={extractedData} />

        {/* Show editable section for manual profile management */}
        {userId && (
          <EditableWorkExperienceSection
            experiences={profile?.experiences ?? []}
            userId={userId}
            onMutationSuccess={handleMutationSuccess}
          />
        )}

        {/* Also show extracted resume data if not yet migrated to profile */}
        {!hasProfileExperiences && extractedData.experience.length > 0 && (
          <div className="opacity-60">
            <p className="text-sm text-gray-500 mb-4 text-center">
              Below is data extracted from your resume. Use the section above to manually manage
              your work experience.
            </p>
            <WorkExperienceSection experience={extractedData.experience} />
          </div>
        )}

        <EducationSection education={extractedData.education} />
        <SkillsSection skills={extractedData.skills} />
        <ProfileActions
          onAddReference={handleAddReference}
          onExport={handleExport}
          onUploadAnother={handleUploadAnother}
        />
      </div>
    </div>
  );
}
