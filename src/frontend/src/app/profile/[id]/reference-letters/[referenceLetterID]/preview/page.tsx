"use client";

import { useParams, useRouter } from "next/navigation";
import { useCallback, useMemo, useState } from "react";
import { useMutation, useQuery } from "urql";
import { Button } from "@/components/ui/button";
import {
  ApplyReferenceLetterValidationsDocument,
  GetProfileDocument,
  GetReferenceLetterDocument,
  ReferenceLetterStatus,
  SkillCategory,
} from "@/graphql/generated/graphql";
import { CorroborationsSection } from "./CorroborationsSection";
import { DiscoveredSkillsSection } from "./DiscoveredSkillsSection";
import { TestimonialsSection } from "./TestimonialsSection";
import { ValidationPreviewSkeleton } from "./ValidationPreviewSkeleton";

// Types for selection state
export interface SkillCorroboration {
  profileSkillId: string;
  skillName: string;
  quote: string;
}

export interface ExperienceCorroboration {
  profileExperienceId: string;
  company: string;
  role: string;
  quote: string;
}

export interface TestimonialItem {
  quote: string;
  skillsMentioned: string[];
  pageNumber: number | null | undefined;
}

export interface DiscoveredSkill {
  name: string;
  quote: string;
}

export default function ValidationPreviewPage() {
  const params = useParams();
  const router = useRouter();
  const resumeId = params.id as string;
  const referenceLetterID = params.referenceLetterID as string;

  // Selection state
  const [selectedSkillCorroborations, setSelectedSkillCorroborations] = useState<Set<string>>(
    new Set()
  );
  const [selectedExperienceCorroborations, setSelectedExperienceCorroborations] = useState<
    Set<string>
  >(new Set());
  const [selectedTestimonials, setSelectedTestimonials] = useState<Set<number>>(new Set());
  const [selectedDiscoveredSkills, setSelectedDiscoveredSkills] = useState<Set<string>>(new Set());
  const [isInitialized, setIsInitialized] = useState(false);

  // Query reference letter with extracted data
  const [referenceLetterResult] = useQuery({
    query: GetReferenceLetterDocument,
    variables: { id: referenceLetterID },
    requestPolicy: "network-only",
  });

  // Get user ID from reference letter
  const userId = referenceLetterResult.data?.referenceLetter?.user?.id;

  // Query user's profile to get current skills and experiences
  const [profileResult] = useQuery({
    query: GetProfileDocument,
    variables: { userId: userId || "" },
    pause: !userId,
    requestPolicy: "network-only",
  });

  // Apply validations mutation
  const [applyResult, applyValidations] = useMutation(ApplyReferenceLetterValidationsDocument);

  const referenceLetter = referenceLetterResult.data?.referenceLetter;
  const profile = profileResult.data?.profile;
  const extractedData = referenceLetter?.extractedData;

  // Compute skill corroborations (matching extracted skill mentions with profile skills)
  const skillCorroborations = useMemo((): SkillCorroboration[] => {
    if (!extractedData?.skillMentions || !profile?.skills) return [];

    const corroborations: SkillCorroboration[] = [];

    for (const skillMention of extractedData.skillMentions) {
      // Find matching profile skill (case-insensitive)
      const matchingSkill = profile.skills.find(
        (s) => s.normalizedName?.toLowerCase() === skillMention.skill.toLowerCase()
      );

      if (matchingSkill) {
        corroborations.push({
          profileSkillId: matchingSkill.id,
          skillName: matchingSkill.name,
          quote: skillMention.quote,
        });
      }
    }

    return corroborations;
  }, [extractedData?.skillMentions, profile?.skills]);

  // Compute experience corroborations (matching extracted experience mentions with profile experiences)
  const experienceCorroborations = useMemo((): ExperienceCorroboration[] => {
    if (!extractedData?.experienceMentions || !profile?.experiences) return [];

    const corroborations: ExperienceCorroboration[] = [];

    for (const expMention of extractedData.experienceMentions) {
      // Find matching profile experience (case-insensitive company match)
      const matchingExp = profile.experiences.find(
        (e) => e.company.toLowerCase() === expMention.company.toLowerCase()
      );

      if (matchingExp) {
        corroborations.push({
          profileExperienceId: matchingExp.id,
          company: matchingExp.company,
          role: expMention.role,
          quote: expMention.quote,
        });
      }
    }

    return corroborations;
  }, [extractedData?.experienceMentions, profile?.experiences]);

  // Testimonials from extracted data
  const testimonials = useMemo((): TestimonialItem[] => {
    if (!extractedData?.testimonials) return [];

    return extractedData.testimonials.map((t) => ({
      quote: t.quote,
      skillsMentioned: t.skillsMentioned || [],
      pageNumber: t.pageNumber,
    }));
  }, [extractedData?.testimonials]);

  // Discovered skills (skills mentioned but not in profile)
  const discoveredSkills = useMemo((): DiscoveredSkill[] => {
    if (!extractedData?.discoveredSkills) return [];

    return extractedData.discoveredSkills.map((ds) => ({
      name: ds.skill,
      quote: ds.quote || "",
    }));
  }, [extractedData?.discoveredSkills]);

  // Initialize selections when data loads (select all corroborations and testimonials by default)
  useMemo(() => {
    if (isInitialized || !extractedData || !profile) return;

    // Pre-select all skill corroborations
    const skillIds = new Set(skillCorroborations.map((c) => c.profileSkillId));
    setSelectedSkillCorroborations(skillIds);

    // Pre-select all experience corroborations
    const expIds = new Set(experienceCorroborations.map((c) => c.profileExperienceId));
    setSelectedExperienceCorroborations(expIds);

    // Pre-select all testimonials
    const testimonialIndices = new Set(testimonials.map((_, i) => i));
    setSelectedTestimonials(testimonialIndices);

    // Discovered skills are NOT pre-selected
    setSelectedDiscoveredSkills(new Set());

    setIsInitialized(true);
  }, [
    extractedData,
    profile,
    isInitialized,
    skillCorroborations,
    experienceCorroborations,
    testimonials,
  ]);

  // Selection handlers
  const handleSkillCorroborationToggle = useCallback((profileSkillId: string) => {
    setSelectedSkillCorroborations((prev) => {
      const newSet = new Set(prev);
      if (newSet.has(profileSkillId)) {
        newSet.delete(profileSkillId);
      } else {
        newSet.add(profileSkillId);
      }
      return newSet;
    });
  }, []);

  const handleExperienceCorroborationToggle = useCallback((profileExperienceId: string) => {
    setSelectedExperienceCorroborations((prev) => {
      const newSet = new Set(prev);
      if (newSet.has(profileExperienceId)) {
        newSet.delete(profileExperienceId);
      } else {
        newSet.add(profileExperienceId);
      }
      return newSet;
    });
  }, []);

  const handleTestimonialToggle = useCallback((index: number) => {
    setSelectedTestimonials((prev) => {
      const newSet = new Set(prev);
      if (newSet.has(index)) {
        newSet.delete(index);
      } else {
        newSet.add(index);
      }
      return newSet;
    });
  }, []);

  const handleDiscoveredSkillToggle = useCallback((skillName: string) => {
    setSelectedDiscoveredSkills((prev) => {
      const newSet = new Set(prev);
      if (newSet.has(skillName)) {
        newSet.delete(skillName);
      } else {
        newSet.add(skillName);
      }
      return newSet;
    });
  }, []);

  // Select/Deselect all handlers
  const handleSelectAllSkillCorroborations = useCallback(() => {
    setSelectedSkillCorroborations(new Set(skillCorroborations.map((c) => c.profileSkillId)));
  }, [skillCorroborations]);

  const handleDeselectAllSkillCorroborations = useCallback(() => {
    setSelectedSkillCorroborations(new Set());
  }, []);

  const handleSelectAllExperienceCorroborations = useCallback(() => {
    setSelectedExperienceCorroborations(
      new Set(experienceCorroborations.map((c) => c.profileExperienceId))
    );
  }, [experienceCorroborations]);

  const handleDeselectAllExperienceCorroborations = useCallback(() => {
    setSelectedExperienceCorroborations(new Set());
  }, []);

  const handleSelectAllTestimonials = useCallback(() => {
    setSelectedTestimonials(new Set(testimonials.map((_, i) => i)));
  }, [testimonials]);

  const handleDeselectAllTestimonials = useCallback(() => {
    setSelectedTestimonials(new Set());
  }, []);

  const handleSelectAllDiscoveredSkills = useCallback(() => {
    setSelectedDiscoveredSkills(new Set(discoveredSkills.map((d) => d.name)));
  }, [discoveredSkills]);

  const handleDeselectAllDiscoveredSkills = useCallback(() => {
    setSelectedDiscoveredSkills(new Set());
  }, []);

  // Apply selected validations
  const handleApplySelected = useCallback(async () => {
    if (!userId) return;

    // Build skill validations input
    const skillValidationsInput = skillCorroborations
      .filter((c) => selectedSkillCorroborations.has(c.profileSkillId))
      .map((c) => ({
        profileSkillID: c.profileSkillId,
        quoteSnippet: c.quote,
      }));

    // Build experience validations input
    const experienceValidationsInput = experienceCorroborations
      .filter((c) => selectedExperienceCorroborations.has(c.profileExperienceId))
      .map((c) => ({
        profileExperienceID: c.profileExperienceId,
        quoteSnippet: c.quote,
      }));

    // Build testimonials input
    const testimonialsInput = testimonials
      .filter((_, i) => selectedTestimonials.has(i))
      .map((t) => ({
        quote: t.quote,
        skillsMentioned: t.skillsMentioned,
        pageNumber: t.pageNumber,
      }));

    // Build new skills input (use SOFT as default category for discovered skills)
    const newSkillsInput = discoveredSkills
      .filter((d) => selectedDiscoveredSkills.has(d.name))
      .map((d) => ({
        name: d.name,
        category: SkillCategory.Soft, // Default to soft skills for discovered skills
        quoteContext: d.quote || undefined,
      }));

    const result = await applyValidations({
      userId,
      input: {
        referenceLetterID,
        skillValidations: skillValidationsInput,
        experienceValidations: experienceValidationsInput,
        testimonials: testimonialsInput,
        newSkills: newSkillsInput,
      },
    });

    if (result.data?.applyReferenceLetterValidations?.__typename === "ApplyValidationsResult") {
      // Navigate back to profile page
      router.push(`/profile/${resumeId}`);
    }
  }, [
    userId,
    referenceLetterID,
    resumeId,
    skillCorroborations,
    selectedSkillCorroborations,
    experienceCorroborations,
    selectedExperienceCorroborations,
    testimonials,
    selectedTestimonials,
    discoveredSkills,
    selectedDiscoveredSkills,
    applyValidations,
    router,
  ]);

  // Handle cancel
  const handleCancel = useCallback(() => {
    router.push(`/profile/${resumeId}`);
  }, [router, resumeId]);

  // Loading state
  if (referenceLetterResult.fetching || (userId && profileResult.fetching)) {
    return (
      <div className="min-h-screen bg-background py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-4xl mx-auto">
          <ValidationPreviewSkeleton />
        </div>
      </div>
    );
  }

  // Error state
  if (referenceLetterResult.error) {
    return (
      <div className="min-h-screen bg-background py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-4xl mx-auto text-center">
          <h1 className="text-2xl font-bold text-destructive mb-4">
            Error Loading Reference Letter
          </h1>
          <p className="text-muted-foreground mb-6">{referenceLetterResult.error.message}</p>
          <Button onClick={() => router.push(`/profile/${resumeId}`)}>Back to Profile</Button>
        </div>
      </div>
    );
  }

  // Not found state
  if (!referenceLetter) {
    return (
      <div className="min-h-screen bg-background py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-4xl mx-auto text-center">
          <h1 className="text-2xl font-bold text-foreground mb-4">Reference Letter Not Found</h1>
          <p className="text-muted-foreground mb-6">
            The reference letter you&apos;re looking for doesn&apos;t exist.
          </p>
          <Button onClick={() => router.push(`/profile/${resumeId}`)}>Back to Profile</Button>
        </div>
      </div>
    );
  }

  // Still processing state
  if (
    referenceLetter.status === ReferenceLetterStatus.Pending ||
    referenceLetter.status === ReferenceLetterStatus.Processing
  ) {
    return (
      <div className="min-h-screen bg-background py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-4xl mx-auto text-center">
          <h1 className="text-2xl font-bold text-foreground mb-4">Processing Reference Letter</h1>
          <p className="text-muted-foreground mb-6">
            Your reference letter is still being processed. Please wait...
          </p>
          <div className="animate-spin h-8 w-8 border-4 border-primary border-t-transparent rounded-full mx-auto" />
        </div>
      </div>
    );
  }

  // Failed state
  if (referenceLetter.status === ReferenceLetterStatus.Failed) {
    return (
      <div className="min-h-screen bg-background py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-4xl mx-auto text-center">
          <h1 className="text-2xl font-bold text-destructive mb-4">Processing Failed</h1>
          <p className="text-muted-foreground mb-6">
            Failed to extract information from your reference letter.
          </p>
          <Button onClick={() => router.push(`/profile/${resumeId}`)}>Back to Profile</Button>
        </div>
      </div>
    );
  }

  // No extracted data
  if (!extractedData) {
    return (
      <div className="min-h-screen bg-background py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-4xl mx-auto text-center">
          <h1 className="text-2xl font-bold text-foreground mb-4">No Data Extracted</h1>
          <p className="text-muted-foreground mb-6">
            No information could be extracted from your reference letter.
          </p>
          <Button onClick={() => router.push(`/profile/${resumeId}`)}>Back to Profile</Button>
        </div>
      </div>
    );
  }

  // Calculate total selected
  const totalSelected =
    selectedSkillCorroborations.size +
    selectedExperienceCorroborations.size +
    selectedTestimonials.size +
    selectedDiscoveredSkills.size;

  // Check if already applied
  const isAlreadyApplied = referenceLetter.status === ReferenceLetterStatus.Applied;

  // Get author info for display
  const authorName = extractedData.author?.name || referenceLetter.authorName || "Unknown Author";
  const authorTitle = extractedData.author?.title || referenceLetter.authorTitle || "";
  const authorCompany = extractedData.author?.company || referenceLetter.organization || "";
  const authorRelationship = extractedData.author?.relationship || "";

  const authorAttribution = [authorTitle, authorCompany, authorRelationship]
    .filter(Boolean)
    .join(" | ");

  return (
    <div className="min-h-screen bg-background py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-4xl mx-auto space-y-8">
        {/* Header */}
        <div>
          <h1 className="text-3xl font-bold text-foreground mb-2">Validation Preview</h1>
          <p className="text-muted-foreground">
            Review what this reference letter validates. Select the items you want to add to your
            profile.
          </p>
          <div className="mt-4 p-4 bg-muted rounded-lg">
            <p className="text-sm font-medium text-foreground">Reference from {authorName}</p>
            {authorAttribution && (
              <p className="text-sm text-muted-foreground">{authorAttribution}</p>
            )}
          </div>
        </div>

        {isAlreadyApplied && (
          <div className="p-4 bg-warning/10 border border-warning rounded-lg">
            <p className="text-warning-foreground font-medium">
              This reference letter has already been applied to your profile.
            </p>
          </div>
        )}

        {/* Corroborations Section */}
        {(skillCorroborations.length > 0 || experienceCorroborations.length > 0) && (
          <CorroborationsSection
            skillCorroborations={skillCorroborations}
            experienceCorroborations={experienceCorroborations}
            selectedSkillCorroborations={selectedSkillCorroborations}
            selectedExperienceCorroborations={selectedExperienceCorroborations}
            onSkillToggle={handleSkillCorroborationToggle}
            onExperienceToggle={handleExperienceCorroborationToggle}
            onSelectAllSkills={handleSelectAllSkillCorroborations}
            onDeselectAllSkills={handleDeselectAllSkillCorroborations}
            onSelectAllExperiences={handleSelectAllExperienceCorroborations}
            onDeselectAllExperiences={handleDeselectAllExperienceCorroborations}
            disabled={isAlreadyApplied}
          />
        )}

        {/* Testimonials Section */}
        {testimonials.length > 0 && (
          <TestimonialsSection
            testimonials={testimonials}
            selectedTestimonials={selectedTestimonials}
            onToggle={handleTestimonialToggle}
            onSelectAll={handleSelectAllTestimonials}
            onDeselectAll={handleDeselectAllTestimonials}
            authorName={authorName}
            authorAttribution={authorAttribution}
            disabled={isAlreadyApplied}
          />
        )}

        {/* Discovered Skills Section */}
        {discoveredSkills.length > 0 && (
          <DiscoveredSkillsSection
            discoveredSkills={discoveredSkills}
            selectedSkills={selectedDiscoveredSkills}
            onToggle={handleDiscoveredSkillToggle}
            onSelectAll={handleSelectAllDiscoveredSkills}
            onDeselectAll={handleDeselectAllDiscoveredSkills}
            disabled={isAlreadyApplied}
          />
        )}

        {/* Empty state */}
        {skillCorroborations.length === 0 &&
          experienceCorroborations.length === 0 &&
          testimonials.length === 0 &&
          discoveredSkills.length === 0 && (
            <div className="text-center py-12">
              <p className="text-muted-foreground">
                No validations could be found in this reference letter.
              </p>
            </div>
          )}

        {/* Error from mutation */}
        {applyResult.data?.applyReferenceLetterValidations?.__typename ===
          "ApplyValidationsError" && (
          <div className="p-4 bg-destructive/10 border border-destructive rounded-lg">
            <p className="text-destructive font-medium">
              {applyResult.data.applyReferenceLetterValidations.message}
            </p>
          </div>
        )}

        {/* Action buttons */}
        <div className="flex justify-between items-center pt-4 border-t">
          <Button variant="outline" onClick={handleCancel}>
            Cancel
          </Button>
          <div className="flex items-center gap-4">
            <span className="text-sm text-muted-foreground">{totalSelected} item(s) selected</span>
            <Button
              onClick={handleApplySelected}
              disabled={totalSelected === 0 || applyResult.fetching || isAlreadyApplied}
            >
              {applyResult.fetching ? "Applying..." : "Apply Selected"}
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}
