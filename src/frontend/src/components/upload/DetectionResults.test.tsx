import { fireEvent, render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import { DetectionResults } from "./DetectionResults";
import type { DocumentDetectionResult } from "./types";

const HIGH_CONFIDENCE_HYBRID: DocumentDetectionResult = {
  hasCareerInfo: true,
  hasTestimonial: true,
  testimonialAuthor: "Jane Doe",
  confidence: 0.92,
  summary: "A reference letter with career details for a software engineer.",
  documentTypeHint: "HYBRID",
  fileId: "file-123",
};

const HIGH_CONFIDENCE_RESUME: DocumentDetectionResult = {
  hasCareerInfo: true,
  hasTestimonial: false,
  testimonialAuthor: null,
  confidence: 0.95,
  summary: "A software engineer resume with work experience and education.",
  documentTypeHint: "RESUME",
  fileId: "file-456",
};

const HIGH_CONFIDENCE_LETTER: DocumentDetectionResult = {
  hasCareerInfo: false,
  hasTestimonial: true,
  testimonialAuthor: "John Smith",
  confidence: 0.88,
  summary: "A reference letter from a former manager.",
  documentTypeHint: "REFERENCE_LETTER",
  fileId: "file-789",
};

const LOW_CONFIDENCE: DocumentDetectionResult = {
  hasCareerInfo: true,
  hasTestimonial: false,
  testimonialAuthor: null,
  confidence: 0.45,
  summary: "The document content is unclear.",
  documentTypeHint: "UNKNOWN",
  fileId: "file-low",
};

const NOTHING_DETECTED: DocumentDetectionResult = {
  hasCareerInfo: false,
  hasTestimonial: false,
  testimonialAuthor: null,
  confidence: 0.3,
  summary: "Could not identify document content.",
  documentTypeHint: "UNKNOWN",
  fileId: "file-empty",
};

const defaultProps = {
  fileName: "test-document.pdf",
  userId: "user-123",
  onProceed: vi.fn(),
  onCancel: vi.fn(),
};

describe("DetectionResults", () => {
  describe("High confidence detection", () => {
    it("shows document summary", () => {
      render(<DetectionResults detection={HIGH_CONFIDENCE_HYBRID} {...defaultProps} />);
      expect(
        screen.getByText("A reference letter with career details for a software engineer.")
      ).toBeInTheDocument();
    });

    it("shows file name", () => {
      render(<DetectionResults detection={HIGH_CONFIDENCE_HYBRID} {...defaultProps} />);
      expect(screen.getByText("test-document.pdf")).toBeInTheDocument();
    });

    it("shows career info checkbox pre-selected when detected", () => {
      render(<DetectionResults detection={HIGH_CONFIDENCE_HYBRID} {...defaultProps} />);
      const checkbox = screen.getByRole("checkbox", { name: /career information/i });
      expect(checkbox).toHaveAttribute("aria-checked", "true");
    });

    it("shows testimonial checkbox pre-selected with author name", () => {
      render(<DetectionResults detection={HIGH_CONFIDENCE_HYBRID} {...defaultProps} />);
      const checkbox = screen.getByRole("checkbox", { name: /testimonial/i });
      expect(checkbox).toHaveAttribute("aria-checked", "true");
      expect(screen.getByText(/Jane Doe/)).toBeInTheDocument();
    });

    it("does not show career info checkbox when not detected", () => {
      render(<DetectionResults detection={HIGH_CONFIDENCE_LETTER} {...defaultProps} />);
      expect(
        screen.queryByRole("checkbox", { name: /career information/i })
      ).not.toBeInTheDocument();
    });

    it("does not show testimonial checkbox when not detected", () => {
      render(<DetectionResults detection={HIGH_CONFIDENCE_RESUME} {...defaultProps} />);
      expect(screen.queryByRole("checkbox", { name: /testimonial/i })).not.toBeInTheDocument();
    });

    it("does not show confidence warning for high confidence", () => {
      render(<DetectionResults detection={HIGH_CONFIDENCE_HYBRID} {...defaultProps} />);
      expect(screen.queryByText(/not sure what this document contains/i)).not.toBeInTheDocument();
    });

    it("allows unchecking a content type", () => {
      render(<DetectionResults detection={HIGH_CONFIDENCE_HYBRID} {...defaultProps} />);
      const careerCheckbox = screen.getByRole("checkbox", { name: /career information/i });
      fireEvent.click(careerCheckbox);
      expect(careerCheckbox).toHaveAttribute("aria-checked", "false");
    });
  });

  describe("Proceed button", () => {
    it("is enabled when at least one content type is selected", () => {
      render(<DetectionResults detection={HIGH_CONFIDENCE_HYBRID} {...defaultProps} />);
      const proceedButton = screen.getByRole("button", { name: /extract selected content/i });
      expect(proceedButton).toBeEnabled();
    });

    it("is disabled when nothing is selected", () => {
      render(<DetectionResults detection={HIGH_CONFIDENCE_HYBRID} {...defaultProps} />);
      // Uncheck both
      fireEvent.click(screen.getByRole("checkbox", { name: /career information/i }));
      fireEvent.click(screen.getByRole("checkbox", { name: /testimonial/i }));
      const proceedButton = screen.getByRole("button", { name: /extract selected content/i });
      expect(proceedButton).toBeDisabled();
    });

    it("calls onProceed with correct selections when clicked", () => {
      const onProceed = vi.fn();
      render(
        <DetectionResults
          detection={HIGH_CONFIDENCE_HYBRID}
          {...defaultProps}
          onProceed={onProceed}
        />
      );
      fireEvent.click(screen.getByRole("button", { name: /extract selected content/i }));
      expect(onProceed).toHaveBeenCalledWith(true, true);
    });

    it("calls onProceed with only career info when testimonial unchecked", () => {
      const onProceed = vi.fn();
      render(
        <DetectionResults
          detection={HIGH_CONFIDENCE_HYBRID}
          {...defaultProps}
          onProceed={onProceed}
        />
      );
      fireEvent.click(screen.getByRole("checkbox", { name: /testimonial/i }));
      fireEvent.click(screen.getByRole("button", { name: /extract selected content/i }));
      expect(onProceed).toHaveBeenCalledWith(true, false);
    });

    it("shows description of what will be extracted", () => {
      render(<DetectionResults detection={HIGH_CONFIDENCE_HYBRID} {...defaultProps} />);
      expect(screen.getByText(/career information/i)).toBeInTheDocument();
      expect(screen.getByText(/testimonial/i)).toBeInTheDocument();
    });
  });

  describe("Cancel action", () => {
    it("shows upload different document option", () => {
      render(<DetectionResults detection={HIGH_CONFIDENCE_HYBRID} {...defaultProps} />);
      expect(
        screen.getByRole("button", { name: /upload different document/i })
      ).toBeInTheDocument();
    });

    it("calls onCancel when clicked", () => {
      const onCancel = vi.fn();
      render(
        <DetectionResults
          detection={HIGH_CONFIDENCE_HYBRID}
          {...defaultProps}
          onCancel={onCancel}
        />
      );
      fireEvent.click(screen.getByRole("button", { name: /upload different document/i }));
      expect(onCancel).toHaveBeenCalled();
    });
  });

  describe("Low confidence detection", () => {
    it("shows warning message for low confidence", () => {
      render(<DetectionResults detection={LOW_CONFIDENCE} {...defaultProps} />);
      expect(screen.getByText(/not sure what this document contains/i)).toBeInTheDocument();
    });

    it("shows manual selection radio options", () => {
      render(<DetectionResults detection={LOW_CONFIDENCE} {...defaultProps} />);
      expect(screen.getByRole("radio", { name: /resume/i })).toBeInTheDocument();
      expect(screen.getByRole("radio", { name: /reference letter/i })).toBeInTheDocument();
      expect(screen.getByRole("radio", { name: /both/i })).toBeInTheDocument();
    });

    it("disables proceed button until user selects a type", () => {
      render(<DetectionResults detection={LOW_CONFIDENCE} {...defaultProps} />);
      const proceedButton = screen.getByRole("button", { name: /extract selected content/i });
      expect(proceedButton).toBeDisabled();
    });

    it("enables proceed button after selecting resume", () => {
      render(<DetectionResults detection={LOW_CONFIDENCE} {...defaultProps} />);
      fireEvent.click(screen.getByRole("radio", { name: /resume/i }));
      const proceedButton = screen.getByRole("button", { name: /extract selected content/i });
      expect(proceedButton).toBeEnabled();
    });

    it("calls onProceed with career info only when resume selected", () => {
      const onProceed = vi.fn();
      render(
        <DetectionResults detection={LOW_CONFIDENCE} {...defaultProps} onProceed={onProceed} />
      );
      fireEvent.click(screen.getByRole("radio", { name: /resume/i }));
      fireEvent.click(screen.getByRole("button", { name: /extract selected content/i }));
      expect(onProceed).toHaveBeenCalledWith(true, false);
    });

    it("calls onProceed with testimonial only when reference letter selected", () => {
      const onProceed = vi.fn();
      render(
        <DetectionResults detection={LOW_CONFIDENCE} {...defaultProps} onProceed={onProceed} />
      );
      fireEvent.click(screen.getByRole("radio", { name: /reference letter/i }));
      fireEvent.click(screen.getByRole("button", { name: /extract selected content/i }));
      expect(onProceed).toHaveBeenCalledWith(false, true);
    });

    it("calls onProceed with both when both selected", () => {
      const onProceed = vi.fn();
      render(
        <DetectionResults detection={LOW_CONFIDENCE} {...defaultProps} onProceed={onProceed} />
      );
      fireEvent.click(screen.getByRole("radio", { name: /both/i }));
      fireEvent.click(screen.getByRole("button", { name: /extract selected content/i }));
      expect(onProceed).toHaveBeenCalledWith(true, true);
    });
  });

  describe("Nothing detected / empty detection", () => {
    it("shows manual selection when nothing detected", () => {
      render(<DetectionResults detection={NOTHING_DETECTED} {...defaultProps} />);
      expect(screen.getByText(/not sure what this document contains/i)).toBeInTheDocument();
      expect(screen.getByRole("radio", { name: /resume/i })).toBeInTheDocument();
    });
  });

  describe("Correction UI", () => {
    it("shows 'Not what you expected?' toggle for high confidence", () => {
      render(<DetectionResults detection={HIGH_CONFIDENCE_HYBRID} {...defaultProps} />);
      expect(screen.getByText(/not what you expected/i)).toBeInTheDocument();
    });

    it("expands correction options when clicked", () => {
      render(<DetectionResults detection={HIGH_CONFIDENCE_HYBRID} {...defaultProps} />);
      fireEvent.click(screen.getByText(/not what you expected/i));
      expect(screen.getByRole("radio", { name: /just a resume/i })).toBeInTheDocument();
      expect(screen.getByRole("radio", { name: /just a reference letter/i })).toBeInTheDocument();
    });

    it("applies 'just a resume' correction", () => {
      render(<DetectionResults detection={HIGH_CONFIDENCE_HYBRID} {...defaultProps} />);
      fireEvent.click(screen.getByText(/not what you expected/i));
      fireEvent.click(screen.getByRole("radio", { name: /just a resume/i }));

      // Should now only have career info selected
      const careerCheckbox = screen.getByRole("checkbox", { name: /career information/i });
      expect(careerCheckbox).toHaveAttribute("aria-checked", "true");
      expect(screen.queryByRole("checkbox", { name: /testimonial/i })).not.toBeInTheDocument();
    });

    it("applies 'just a reference letter' correction", () => {
      render(<DetectionResults detection={HIGH_CONFIDENCE_HYBRID} {...defaultProps} />);
      fireEvent.click(screen.getByText(/not what you expected/i));
      fireEvent.click(screen.getByRole("radio", { name: /just a reference letter/i }));

      // Should now only have testimonial selected
      expect(
        screen.queryByRole("checkbox", { name: /career information/i })
      ).not.toBeInTheDocument();
      const testimonialCheckbox = screen.getByRole("checkbox", { name: /testimonial/i });
      expect(testimonialCheckbox).toHaveAttribute("aria-checked", "true");
    });

    it("shows free-text feedback field", () => {
      render(<DetectionResults detection={HIGH_CONFIDENCE_HYBRID} {...defaultProps} />);
      fireEvent.click(screen.getByText(/not what you expected/i));
      expect(screen.getByPlaceholderText(/tell us more about this document/i)).toBeInTheDocument();
    });
  });

  describe("Error state (unreadable document)", () => {
    it("shows error message when detection indicates unreadable content", () => {
      const unreadable: DocumentDetectionResult = {
        hasCareerInfo: false,
        hasTestimonial: false,
        testimonialAuthor: null,
        confidence: 0,
        summary: "",
        documentTypeHint: "UNKNOWN",
        fileId: "file-bad",
      };
      render(<DetectionResults detection={unreadable} {...defaultProps} />);
      // Should show the low confidence / manual selection flow since nothing was detected
      expect(screen.getByText(/not sure what this document contains/i)).toBeInTheDocument();
    });
  });
});
