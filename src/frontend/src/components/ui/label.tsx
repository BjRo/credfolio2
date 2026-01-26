"use client";

import { cva, type VariantProps } from "class-variance-authority";
import type * as React from "react";

import { cn } from "@/lib/utils";

const labelVariants = cva(
  "text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
);

function Label({
  className,
  ...props
}: React.ComponentProps<"label"> & VariantProps<typeof labelVariants>) {
  // biome-ignore lint/a11y/noLabelWithoutControl: htmlFor is passed via props spread
  return <label data-slot="label" className={cn(labelVariants(), className)} {...props} />;
}

export { Label };
