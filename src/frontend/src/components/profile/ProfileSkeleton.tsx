function SkeletonBox({ className }: { className?: string }) {
  return <div className={`animate-pulse bg-muted rounded ${className}`} />;
}

export function ProfileSkeleton() {
  return (
    <div className="space-y-6">
      {/* Header skeleton */}
      <div className="bg-card border rounded-lg p-8">
        <div className="flex items-start justify-between">
          <div className="space-y-3">
            <SkeletonBox className="h-8 w-48" />
            <SkeletonBox className="h-4 w-40" />
            <SkeletonBox className="h-4 w-36" />
            <SkeletonBox className="h-4 w-32" />
          </div>
          <SkeletonBox className="h-4 w-24" />
        </div>
        <div className="mt-6 space-y-2">
          <SkeletonBox className="h-4 w-20" />
          <SkeletonBox className="h-4 w-full" />
          <SkeletonBox className="h-4 w-3/4" />
        </div>
      </div>

      {/* Experience skeleton */}
      <div className="bg-card border rounded-lg p-8">
        <SkeletonBox className="h-6 w-36 mb-6" />
        <div className="space-y-6">
          {[1, 2].map((i) => (
            <div key={i} className={i > 1 ? "pt-6 border-t border-border" : ""}>
              <div className="flex justify-between items-start">
                <div className="space-y-2">
                  <SkeletonBox className="h-5 w-40" />
                  <SkeletonBox className="h-4 w-32" />
                  <SkeletonBox className="h-3 w-28" />
                </div>
                <SkeletonBox className="h-3 w-32" />
              </div>
              <div className="mt-2 space-y-2">
                <SkeletonBox className="h-4 w-full" />
                <SkeletonBox className="h-4 w-2/3" />
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Education skeleton */}
      <div className="bg-card border rounded-lg p-8">
        <SkeletonBox className="h-6 w-28 mb-6" />
        <div className="flex justify-between items-start">
          <div className="space-y-2">
            <SkeletonBox className="h-5 w-44" />
            <SkeletonBox className="h-4 w-56" />
          </div>
          <SkeletonBox className="h-3 w-32" />
        </div>
      </div>

      {/* Skills skeleton */}
      <div className="bg-card border rounded-lg p-8">
        <SkeletonBox className="h-6 w-20 mb-4" />
        <div className="flex flex-wrap gap-2">
          {[1, 2, 3, 4, 5, 6].map((i) => (
            <SkeletonBox key={i} className="h-7 w-20 rounded-full" />
          ))}
        </div>
      </div>
    </div>
  );
}
