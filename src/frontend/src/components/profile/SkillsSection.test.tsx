import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { SkillsSection } from "./SkillsSection";

const mockSkills = ["JavaScript", "TypeScript", "React", "Node.js", "GraphQL"];

describe("SkillsSection", () => {
  it("renders the section heading", () => {
    render(<SkillsSection skills={mockSkills} />);
    expect(screen.getByRole("heading", { level: 2, name: "Skills" })).toBeInTheDocument();
  });

  it("renders all skills as badges", () => {
    render(<SkillsSection skills={mockSkills} />);
    expect(screen.getByText("JavaScript")).toBeInTheDocument();
    expect(screen.getByText("TypeScript")).toBeInTheDocument();
    expect(screen.getByText("React")).toBeInTheDocument();
    expect(screen.getByText("Node.js")).toBeInTheDocument();
    expect(screen.getByText("GraphQL")).toBeInTheDocument();
  });

  it("returns null when skills array is empty", () => {
    const { container } = render(<SkillsSection skills={[]} />);
    expect(container.firstChild).toBeNull();
  });
});
