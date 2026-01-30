"use client";

import { Skeleton } from "@/components/ui/skeleton";

export function ValidationPreviewSkeleton() {
  return (
    <div className="space-y-8">
      {/* Header */}
      <div className="space-y-2">
        <Skeleton className="h-9 w-64" />
        <Skeleton className="h-5 w-96" />
        <Skeleton className="h-20 w-full mt-4" />
      </div>

      {/* Corroborations Section */}
      <div className="space-y-4">
        <Skeleton className="h-7 w-80" />
        <div className="space-y-3">
          {[1, 2, 3].map((i) => (
            <Skeleton key={i} className="h-24 w-full" />
          ))}
        </div>
      </div>

      {/* Testimonials Section */}
      <div className="space-y-4">
        <Skeleton className="h-7 w-48" />
        <div className="space-y-3">
          {[1, 2].map((i) => (
            <Skeleton key={i} className="h-32 w-full" />
          ))}
        </div>
      </div>

      {/* Action buttons */}
      <div className="flex justify-between pt-4 border-t">
        <Skeleton className="h-10 w-20" />
        <div className="flex items-center gap-4">
          <Skeleton className="h-5 w-32" />
          <Skeleton className="h-10 w-28" />
        </div>
      </div>
    </div>
  );
}
