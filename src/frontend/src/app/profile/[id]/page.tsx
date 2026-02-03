"use client";

import { useParams, useRouter } from "next/navigation";
import { useCallback, useState } from "react";
import { useQuery } from "urql";
import {
  EducationSection,
  ProfileActions,
  ProfileHeader,
  ProfileSkeleton,
  ReferenceLetterUploadModal,
  SkillsSection,
  TestimonialsSection,
  WorkExperienceSection,
} from "@/components/profile";
import { Button } from "@/components/ui/button";
import {
  GetProfileDocument,
  GetResumeDocument,
  GetTestimonialsDocument,
  ResumeStatus,
} from "@/graphql/generated/graphql";

export default function ProfilePage() {
  const params = useParams();
  const router = useRouter();
  const resumeId = params.id as string;
  const [isReferenceModalOpen, setIsReferenceModalOpen] = useState(false);

  const [resumeResult, _reexecuteResumeQuery] = useQuery({
    query: GetResumeDocument,
    variables: { id: resumeId },
    requestPolicy: "network-only", // Always fetch fresh data on page load
  });

  // Get user ID from resume to fetch their profile
  const userId = resumeResult.data?.resume?.user?.id;

  const [profileResult, reexecuteProfileQuery] = useQuery({
    query: GetProfileDocument,
    variables: { userId: userId || "" },
    pause: !userId, // Don't run until we have userId
    requestPolicy: "network-only", // Always fetch fresh data on page load
  });

  const { data, fetching, error } = resumeResult;
  const profile = profileResult.data?.profile;
  const profileId = profile?.id;

  // Fetch testimonials for the profile
  const [testimonialsResult, reexecuteTestimonialsQuery] = useQuery({
    query: GetTestimonialsDocument,
    variables: { profileId: profileId || "" },
    pause: !profileId, // Don't run until we have profileId
    requestPolicy: "network-only",
  });

  const testimonials = testimonialsResult.data?.testimonials ?? [];
  const testimonialsLoading = testimonialsResult.fetching;

  // Refetch profile and testimonials when mutations succeed - memoized to prevent unnecessary re-renders
  const handleMutationSuccess = useCallback(() => {
    reexecuteProfileQuery({ requestPolicy: "network-only" });
    reexecuteTestimonialsQuery({ requestPolicy: "network-only" });
  }, [reexecuteProfileQuery, reexecuteTestimonialsQuery]);

  // Handle reference letter upload success - navigate to validation preview
  const handleReferenceUploadSuccess = useCallback(
    (referenceLetterld: string) => {
      // Navigate to the validation preview page
      router.push(`/profile/${resumeId}/reference-letters/${referenceLetterld}/preview`);
    },
    [router, resumeId]
  );

  if (fetching) {
    return (
      <div className="min-h-screen bg-background py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-4xl mx-auto">
          <ProfileSkeleton />
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-background py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-4xl mx-auto text-center">
          <h1 className="text-2xl font-bold text-destructive mb-4">Error Loading Profile</h1>
          <p className="text-muted-foreground mb-6">{error.message}</p>
          <Button onClick={() => window.location.reload()}>Try Again</Button>
        </div>
      </div>
    );
  }

  const resume = data?.resume;

  if (!resume) {
    return (
      <div className="min-h-screen bg-background py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-4xl mx-auto text-center">
          <h1 className="text-2xl font-bold text-foreground mb-4">Profile Not Found</h1>
          <p className="text-muted-foreground mb-6">
            The resume you&apos;re looking for doesn&apos;t exist or has been removed.
          </p>
          <Button onClick={() => router.push("/upload-resume")}>Upload a Resume</Button>
        </div>
      </div>
    );
  }

  if (resume.status === ResumeStatus.Pending || resume.status === ResumeStatus.Processing) {
    return (
      <div className="min-h-screen bg-background py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-4xl mx-auto text-center">
          <h1 className="text-2xl font-bold text-foreground mb-4">Profile Processing</h1>
          <p className="text-muted-foreground mb-6">
            Your resume is still being processed. Please wait...
          </p>
          <div className="animate-spin h-8 w-8 border-4 border-primary border-t-transparent rounded-full mx-auto" />
        </div>
      </div>
    );
  }

  if (resume.status === ResumeStatus.Failed) {
    return (
      <div className="min-h-screen bg-background py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-4xl mx-auto text-center">
          <h1 className="text-2xl font-bold text-destructive mb-4">Processing Failed</h1>
          <p className="text-muted-foreground mb-6">
            {resume.errorMessage || "Failed to process resume"}
          </p>
          <Button onClick={() => router.push("/upload-resume")}>Try Again</Button>
        </div>
      </div>
    );
  }

  const extractedData = resume.extractedData;

  if (!extractedData) {
    return (
      <div className="min-h-screen bg-background py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-4xl mx-auto text-center">
          <h1 className="text-2xl font-bold text-foreground mb-4">No Profile Data</h1>
          <p className="text-muted-foreground mb-6">
            The resume was processed but no data could be extracted.
          </p>
          <Button onClick={() => router.push("/upload-resume")}>Upload a Different Resume</Button>
        </div>
      </div>
    );
  }

  const handleAddReference = () => {
    setIsReferenceModalOpen(true);
  };

  const handleExport = () => {
    // TODO: Implement PDF export
    alert("PDF export coming soon!");
  };

  const handleUploadAnother = () => {
    router.push("/upload-resume");
  };

  return (
    <div className="min-h-screen bg-background py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-4xl mx-auto space-y-6">
        <ProfileHeader
          data={extractedData}
          profileOverrides={profile ?? undefined}
          userId={userId}
          onMutationSuccess={handleMutationSuccess}
        />

        <WorkExperienceSection
          profileExperiences={profile?.experiences ?? []}
          userId={userId}
          onMutationSuccess={handleMutationSuccess}
        />

        <EducationSection
          profileEducations={profile?.educations ?? []}
          userId={userId}
          onMutationSuccess={handleMutationSuccess}
        />
        <SkillsSection
          profileSkills={profile?.skills ?? []}
          userId={userId}
          onMutationSuccess={handleMutationSuccess}
        />

        <TestimonialsSection
          testimonials={testimonials}
          isLoading={testimonialsLoading}
          onAddReference={handleAddReference}
          onTestimonialDeleted={handleMutationSuccess}
          onAuthorUpdated={handleMutationSuccess}
        />

        <ProfileActions
          onAddReference={handleAddReference}
          onExport={handleExport}
          onUploadAnother={handleUploadAnother}
        />

        {userId && (
          <ReferenceLetterUploadModal
            open={isReferenceModalOpen}
            onOpenChange={setIsReferenceModalOpen}
            userId={userId}
            onSuccess={handleReferenceUploadSuccess}
          />
        )}
      </div>
    </div>
  );
}
