import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { StepIndicator } from "./StepIndicator";

const steps = [
  { key: "upload", label: "Upload" },
  { key: "review", label: "Review" },
  { key: "extract", label: "Extract" },
];

describe("StepIndicator", () => {
  it("renders all step labels", () => {
    render(<StepIndicator steps={steps} currentStep="upload" />);
    expect(screen.getByText("Upload")).toBeInTheDocument();
    expect(screen.getByText("Review")).toBeInTheDocument();
    expect(screen.getByText("Extract")).toBeInTheDocument();
  });

  it("renders as a navigation landmark", () => {
    render(<StepIndicator steps={steps} currentStep="upload" />);
    expect(screen.getByRole("navigation", { name: /progress/i })).toBeInTheDocument();
  });

  it("marks the current step with aria-current", () => {
    render(<StepIndicator steps={steps} currentStep="review" />);
    const currentItem = screen.getByText("Review").closest("li");
    expect(currentItem).toHaveAttribute("aria-current", "step");
  });

  it("does not mark non-current steps with aria-current", () => {
    render(<StepIndicator steps={steps} currentStep="review" />);
    const uploadItem = screen.getByText("Upload").closest("li");
    const extractItem = screen.getByText("Extract").closest("li");
    expect(uploadItem).not.toHaveAttribute("aria-current");
    expect(extractItem).not.toHaveAttribute("aria-current");
  });

  it("renders steps as an ordered list", () => {
    render(<StepIndicator steps={steps} currentStep="upload" />);
    const list = screen.getByRole("list");
    expect(list).toBeInTheDocument();
    const items = screen.getAllByRole("listitem");
    expect(items).toHaveLength(3);
  });
});
