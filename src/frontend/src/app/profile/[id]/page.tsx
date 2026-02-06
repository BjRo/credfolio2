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
import { GetProfileByIdDocument, GetTestimonialsDocument } from "@/graphql/generated/graphql";

export default function ProfilePage() {
  const params = useParams();
  const router = useRouter();
  const profileId = params.id as string;
  const [isReferenceModalOpen, setIsReferenceModalOpen] = useState(false);

  const [profileResult, reexecuteProfileQuery] = useQuery({
    query: GetProfileByIdDocument,
    variables: { id: profileId },
    requestPolicy: "network-only",
  });

  const profile = profileResult.data?.profile;
  const userId = profile?.user?.id;

  // Fetch testimonials for the profile
  const [testimonialsResult, reexecuteTestimonialsQuery] = useQuery({
    query: GetTestimonialsDocument,
    variables: { profileId },
    pause: !profile,
    requestPolicy: "network-only",
  });

  const testimonials = testimonialsResult.data?.testimonials ?? [];
  const testimonialsLoading = testimonialsResult.fetching;

  const handleMutationSuccess = useCallback(() => {
    reexecuteProfileQuery({ requestPolicy: "network-only" });
    reexecuteTestimonialsQuery({ requestPolicy: "network-only" });
  }, [reexecuteProfileQuery, reexecuteTestimonialsQuery]);

  const handleReferenceUploadSuccess = useCallback(
    (referenceLetterld: string) => {
      router.push(`/profile/${profileId}/reference-letters/${referenceLetterld}/preview`);
    },
    [router, profileId]
  );

  if (profileResult.fetching) {
    return (
      <div className="min-h-screen bg-background py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-4xl mx-auto">
          <ProfileSkeleton />
        </div>
      </div>
    );
  }

  if (profileResult.error) {
    return (
      <div className="min-h-screen bg-background py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-4xl mx-auto text-center">
          <h1 className="text-2xl font-bold text-destructive mb-4">Error Loading Profile</h1>
          <p className="text-muted-foreground mb-6">{profileResult.error.message}</p>
          <Button onClick={() => window.location.reload()}>Try Again</Button>
        </div>
      </div>
    );
  }

  if (!profile) {
    return (
      <div className="min-h-screen bg-background py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-4xl mx-auto text-center">
          <h1 className="text-2xl font-bold text-foreground mb-4">Profile Not Found</h1>
          <p className="text-muted-foreground mb-6">
            The profile you&apos;re looking for doesn&apos;t exist or has been removed.
          </p>
          <Button onClick={() => router.push("/upload")}>Upload a Document</Button>
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
    router.push("/upload");
  };

  return (
    <div className="min-h-screen bg-background py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-4xl mx-auto space-y-6">
        <ProfileHeader
          profile={profile}
          userId={userId}
          onMutationSuccess={handleMutationSuccess}
        />

        <WorkExperienceSection
          profileExperiences={profile.experiences ?? []}
          userId={userId}
          onMutationSuccess={handleMutationSuccess}
        />

        <EducationSection
          profileEducations={profile.educations ?? []}
          userId={userId}
          onMutationSuccess={handleMutationSuccess}
        />
        <SkillsSection
          profileSkills={profile.skills ?? []}
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
