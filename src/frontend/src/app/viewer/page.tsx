"use client";

import { ArrowLeft, Info, X } from "lucide-react";
import dynamic from "next/dynamic";
import { useRouter, useSearchParams } from "next/navigation";
import { Suspense, useCallback, useEffect, useState } from "react";
import { useQuery } from "urql";
import { Button } from "@/components/ui/button";
import { GetReferenceLetterForViewerDocument } from "@/graphql/generated/graphql";

const PDFViewer = dynamic(
  () => import("@/components/viewer/PDFViewer").then((mod) => mod.PDFViewer),
  { ssr: false }
);

const UUID_REGEX = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;

function InfoBanner({ onDismiss }: { onDismiss: () => void }) {
  return (
    <div className="flex items-center gap-2 bg-amber-50 dark:bg-amber-950/20 border-b border-amber-200 dark:border-amber-800 px-4 py-2 text-sm">
      <Info className="size-4 text-amber-600 dark:text-amber-400 shrink-0" />
      <span className="text-amber-800 dark:text-amber-200 flex-1">
        Could not locate exact quote — showing full document
      </span>
      <Button
        variant="ghost"
        size="icon-sm"
        onClick={onDismiss}
        aria-label="Dismiss banner"
        className="text-amber-600 dark:text-amber-400 hover:bg-amber-100 dark:hover:bg-amber-900/30"
      >
        <X className="size-4" />
      </Button>
    </div>
  );
}

function ErrorPage({
  title,
  description,
  onBack,
}: {
  title: string;
  description: string;
  onBack: () => void;
}) {
  return (
    <div className="min-h-screen bg-background flex items-center justify-center px-4">
      <div className="text-center">
        <h1 className="text-2xl font-bold text-destructive mb-4">{title}</h1>
        <p className="text-muted-foreground mb-6">{description}</p>
        <Button onClick={onBack}>Go back</Button>
      </div>
    </div>
  );
}

function LoadingSkeleton() {
  return (
    <output className="min-h-screen bg-background flex items-center justify-center">
      <div className="animate-spin h-8 w-8 border-4 border-primary border-t-transparent rounded-full" />
    </output>
  );
}

export default function ViewerPage() {
  return (
    <Suspense fallback={<LoadingSkeleton />}>
      <ViewerPageContent />
    </Suspense>
  );
}

function ViewerPageContent() {
  const searchParams = useSearchParams();
  const router = useRouter();
  const letterId = searchParams.get("letterId");
  const highlight = searchParams.get("highlight");

  const [highlightNotFound, setHighlightNotFound] = useState(false);
  const [bannerDismissed, setBannerDismissed] = useState(false);

  // biome-ignore lint/correctness/useExhaustiveDependencies: intentionally reset state when highlight param changes
  useEffect(() => {
    setHighlightNotFound(false);
    setBannerDismissed(false);
  }, [highlight]);

  const isValidId = letterId !== null && UUID_REGEX.test(letterId);

  const [result] = useQuery({
    query: GetReferenceLetterForViewerDocument,
    variables: { id: letterId ?? "" },
    pause: !isValidId,
  });

  const handleHighlightResult = useCallback((found: boolean) => {
    if (!found) {
      setHighlightNotFound(true);
    }
  }, []);

  const handleBack = useCallback(() => {
    if (window.history.length > 1) {
      router.back();
    } else {
      router.push("/");
    }
  }, [router]);

  // Invalid or missing letterId
  if (!isValidId) {
    return (
      <ErrorPage
        title="Document not found"
        description="The document link is invalid or missing a document ID."
        onBack={handleBack}
      />
    );
  }

  // Loading
  if (result.fetching) {
    return <LoadingSkeleton />;
  }

  // GraphQL error
  if (result.error) {
    return (
      <ErrorPage
        title="Failed to load document"
        description="An error occurred while loading the document. Please try again."
        onBack={handleBack}
      />
    );
  }

  const letter = result.data?.referenceLetter;

  // Not found
  if (!letter) {
    return (
      <ErrorPage
        title="Document not found"
        description="The requested document could not be found."
        onBack={handleBack}
      />
    );
  }

  // No file
  if (!letter.file?.url) {
    return (
      <ErrorPage
        title="Document file unavailable"
        description="The document file is not available. It may have been removed."
        onBack={handleBack}
      />
    );
  }

  const documentTitle =
    letter.title ||
    (letter.authorName ? `Reference from ${letter.authorName}` : "Reference Letter");

  const documentSubtitle =
    letter.authorTitle && letter.organization
      ? `${letter.authorTitle}, ${letter.organization}`
      : letter.authorTitle || letter.organization || null;

  const toolbarLeft = (
    <>
      <Button variant="ghost" size="icon-sm" onClick={handleBack} aria-label="Go back">
        <ArrowLeft className="size-4" />
      </Button>
      <div className="flex-1 min-w-0">
        <h1 className="text-sm font-medium truncate">{documentTitle}</h1>
        {documentSubtitle && (
          <p className="text-xs text-muted-foreground truncate">{documentSubtitle}</p>
        )}
      </div>
    </>
  );

  return (
    <div
      data-testid="viewer-container"
      className="flex flex-col h-[calc(100dvh-var(--header-height))]"
    >
      {/* Info banner for highlight not found */}
      {highlight && highlightNotFound && !bannerDismissed && (
        <InfoBanner onDismiss={() => setBannerDismissed(true)} />
      )}

      {/* PDF Viewer */}
      <div className="flex-1 min-h-0">
        <PDFViewer
          fileUrl={letter.file.url}
          highlightText={highlight || undefined}
          onHighlightResult={handleHighlightResult}
          toolbarLeft={toolbarLeft}
        />
      </div>
    </div>
  );
}
