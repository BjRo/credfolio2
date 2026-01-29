import { render } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { ProfileSkeleton } from "./ProfileSkeleton";

describe("ProfileSkeleton", () => {
  it("renders skeleton structure", () => {
    const { container } = render(<ProfileSkeleton />);

    // Should have multiple skeleton boxes with animate-pulse
    const skeletonBoxes = container.querySelectorAll(".animate-pulse");
    expect(skeletonBoxes.length).toBeGreaterThan(10);
  });

  it("renders multiple card sections", () => {
    const { container } = render(<ProfileSkeleton />);

    // Should have multiple card sections
    const cards = container.querySelectorAll(".bg-card.border.rounded-lg");
    expect(cards.length).toBe(4); // Header, Experience, Education, Skills
  });
});
