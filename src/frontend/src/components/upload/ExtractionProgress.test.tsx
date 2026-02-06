import { act, render, screen } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { ExtractionProgress } from "./ExtractionProgress";
import type { ExtractionResults, ProcessDocumentIds } from "./types";

const PROCESS_RESULT = {
  data: {
    processDocument: {
      __typename: "ProcessDocumentResult",
      resumeId: "resume-1",
      referenceLetterID: "letter-1",
    },
  },
};

const PROCESS_RESULT_RESUME_ONLY = {
  data: {
    processDocument: {
      __typename: "ProcessDocumentResult",
      resumeId: "resume-1",
      referenceLetterID: null,
    },
  },
};

const PROCESS_ERROR = {
  data: {
    processDocument: {
      __typename: "ProcessDocumentError",
      message: "File not found",
      field: "fileId",
    },
  },
};

const PROCESSING_STATUS = {
  data: {
    documentProcessingStatus: {
      allComplete: false,
      resume: { id: "resume-1", status: "PROCESSING", extractedData: null, errorMessage: null },
      referenceLetter: { id: "letter-1", status: "PROCESSING", extractedData: null },
    },
  },
};

const COMPLETED_STATUS = {
  data: {
    documentProcessingStatus: {
      allComplete: true,
      resume: {
        id: "resume-1",
        status: "COMPLETED",
        extractedData: {
          name: "John Doe",
          email: "john@example.com",
          phone: null,
          location: "Berlin",
          summary: "A developer",
          extractedAt: "2025-01-01T00:00:00Z",
          confidence: 0.95,
        },
        errorMessage: null,
      },
      referenceLetter: {
        id: "letter-1",
        status: "COMPLETED",
        extractedData: {
          author: { name: "Jane Smith", title: "CTO", company: "TechCo", relationship: "Manager" },
          testimonials: [{ quote: "Great engineer", skillsMentioned: ["Go"] }],
          skillMentions: [],
          experienceMentions: [],
          discoveredSkills: [],
          metadata: { extractedAt: "2025-01-01", modelVersion: "v1", processingTimeMs: 5000 },
        },
      },
    },
  },
};

const FAILED_STATUS = {
  data: {
    documentProcessingStatus: {
      allComplete: true,
      resume: {
        id: "resume-1",
        status: "FAILED",
        extractedData: null,
        errorMessage: "Extraction failed",
      },
      referenceLetter: null,
    },
  },
};

describe("ExtractionProgress", () => {
  let onComplete: ReturnType<typeof vi.fn>;
  let onError: ReturnType<typeof vi.fn>;
  let fetchSpy: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    vi.useFakeTimers();
    onComplete = vi.fn();
    onError = vi.fn();
    fetchSpy = vi.fn();
    vi.stubGlobal("fetch", fetchSpy);
  });

  afterEach(() => {
    vi.useRealTimers();
    vi.restoreAllMocks();
  });

  const makeProps = (overrides?: Record<string, unknown>) => ({
    userId: "user-1",
    fileId: "file-1",
    extractCareerInfo: true,
    extractTestimonial: true,
    onComplete,
    onError,
    ...overrides,
  });

  function mockFetchSequence(...responses: unknown[]) {
    for (const response of responses) {
      fetchSpy.mockResolvedValueOnce({
        json: () => Promise.resolve(response),
      });
    }
  }

  it("calls processDocument mutation on mount", async () => {
    mockFetchSequence(PROCESS_RESULT);

    await act(async () => {
      render(<ExtractionProgress {...makeProps()} />);
    });
    await act(async () => {});

    expect(fetchSpy).toHaveBeenCalledTimes(1);
    const callBody = JSON.parse(fetchSpy.mock.calls[0][1].body);
    expect(callBody.variables.userId).toBe("user-1");
    expect(callBody.variables.input.fileId).toBe("file-1");
    expect(callBody.variables.input.extractCareerInfo).toBe(true);
    expect(callBody.variables.input.extractTestimonial).toBe(true);
  });

  it("shows spinner while processing", async () => {
    mockFetchSequence(PROCESS_RESULT);

    await act(async () => {
      render(<ExtractionProgress {...makeProps()} />);
    });
    await act(async () => {});

    expect(screen.getByRole("img", { name: /loading/i })).toBeInTheDocument();
  });

  it("shows extracting messages for hybrid documents", async () => {
    mockFetchSequence(PROCESS_RESULT);

    await act(async () => {
      render(<ExtractionProgress {...makeProps()} />);
    });
    await act(async () => {});

    expect(
      screen.getByText("Extracting career information and testimonials...")
    ).toBeInTheDocument();
  });

  it("shows career info extraction message when extracting resume only", async () => {
    mockFetchSequence(PROCESS_RESULT_RESUME_ONLY);

    await act(async () => {
      render(<ExtractionProgress {...makeProps({ extractTestimonial: false })} />);
    });
    await act(async () => {});

    expect(screen.getByText("Extracting career information...")).toBeInTheDocument();
  });

  it("polls every 2 seconds", async () => {
    mockFetchSequence(PROCESS_RESULT, PROCESSING_STATUS, PROCESSING_STATUS, COMPLETED_STATUS);

    await act(async () => {
      render(<ExtractionProgress {...makeProps()} />);
    });
    await act(async () => {});

    // processDocument = 1 call
    expect(fetchSpy).toHaveBeenCalledTimes(1);

    // First poll at 2s
    await act(async () => {
      await vi.advanceTimersByTimeAsync(2000);
    });
    expect(fetchSpy).toHaveBeenCalledTimes(2);

    // Second poll at 4s
    await act(async () => {
      await vi.advanceTimersByTimeAsync(2000);
    });
    expect(fetchSpy).toHaveBeenCalledTimes(3);
  });

  it("calls onComplete with extraction results when allComplete is true", async () => {
    mockFetchSequence(PROCESS_RESULT, COMPLETED_STATUS);

    await act(async () => {
      render(<ExtractionProgress {...makeProps()} />);
    });
    await act(async () => {});

    await act(async () => {
      await vi.advanceTimersByTimeAsync(2000);
    });

    expect(onComplete).toHaveBeenCalledTimes(1);
    const [results, processIds] = onComplete.mock.calls[0] as [
      ExtractionResults,
      ProcessDocumentIds,
    ];
    expect(results.resume?.id).toBe("resume-1");
    expect(results.resume?.status).toBe("COMPLETED");
    expect(results.resume?.extractedData?.name).toBe("John Doe");
    expect(results.referenceLetter?.id).toBe("letter-1");
    expect(processIds.resumeId).toBe("resume-1");
    expect(processIds.referenceLetterID).toBe("letter-1");
  });

  it("calls onError when processDocument returns ProcessDocumentError", async () => {
    mockFetchSequence(PROCESS_ERROR);

    await act(async () => {
      render(<ExtractionProgress {...makeProps()} />);
    });
    await act(async () => {});

    expect(onError).toHaveBeenCalledWith("File not found");
  });

  it("shows error message when processDocument fails", async () => {
    mockFetchSequence(PROCESS_ERROR);

    await act(async () => {
      render(<ExtractionProgress {...makeProps()} />);
    });
    await act(async () => {});

    expect(screen.getByText("File not found")).toBeInTheDocument();
  });

  it("calls onError when extraction fails during polling", async () => {
    mockFetchSequence(PROCESS_RESULT_RESUME_ONLY, FAILED_STATUS);

    await act(async () => {
      render(<ExtractionProgress {...makeProps({ extractTestimonial: false })} />);
    });
    await act(async () => {});

    await act(async () => {
      await vi.advanceTimersByTimeAsync(2000);
    });

    expect(onError).toHaveBeenCalledWith("Extraction failed");
  });

  it("cleans up polling interval on unmount", async () => {
    mockFetchSequence(PROCESS_RESULT, PROCESSING_STATUS, PROCESSING_STATUS);

    let unmount: () => void;
    await act(async () => {
      const result = render(<ExtractionProgress {...makeProps()} />);
      unmount = result.unmount;
    });
    await act(async () => {});

    // First poll
    await act(async () => {
      await vi.advanceTimersByTimeAsync(2000);
    });
    expect(fetchSpy).toHaveBeenCalledTimes(2);

    // Unmount
    act(() => {
      unmount();
    });

    // Advance time — should NOT trigger more fetches
    await act(async () => {
      await vi.advanceTimersByTimeAsync(4000);
    });

    expect(fetchSpy).toHaveBeenCalledTimes(2);
  });

  it("does not call processDocument mutation again if effect re-fires", async () => {
    // Provide two responses in case the mutation fires twice
    const onError1 = vi.fn();
    const onError2 = vi.fn();
    mockFetchSequence(PROCESS_RESULT, PROCESS_RESULT);

    const { rerender } = render(<ExtractionProgress {...makeProps({ onError: onError1 })} />);
    await act(async () => {});

    expect(fetchSpy).toHaveBeenCalledTimes(1);

    // Re-render with a new onError reference causes startProcessing to get
    // a new reference, which re-fires the mount effect. The isStartedRef
    // guard should prevent a second mutation call.
    await act(async () => {
      rerender(<ExtractionProgress {...makeProps({ onError: onError2 })} />);
    });
    await act(async () => {});

    expect(fetchSpy).toHaveBeenCalledTimes(1);
  });
});
