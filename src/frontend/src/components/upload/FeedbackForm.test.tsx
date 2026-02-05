import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { FeedbackForm } from "./FeedbackForm";

describe("FeedbackForm", () => {
  const defaultProps = {
    userId: "user-1",
    fileId: "file-1",
  };

  let fetchSpy: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    fetchSpy = vi.fn().mockResolvedValue({
      json: () => Promise.resolve({ data: { reportDocumentFeedback: { success: true } } }),
    });
    vi.stubGlobal("fetch", fetchSpy);
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it("renders predefined feedback options", () => {
    render(<FeedbackForm {...defaultProps} />);

    expect(screen.getByLabelText("Missing information")).toBeInTheDocument();
    expect(screen.getByLabelText("Incorrect data")).toBeInTheDocument();
    expect(screen.getByLabelText("Wrong person")).toBeInTheDocument();
    expect(screen.getByLabelText("Other")).toBeInTheDocument();
  });

  it("renders free-text description textarea", () => {
    render(<FeedbackForm {...defaultProps} />);

    expect(screen.getByPlaceholderText("Tell us more...")).toBeInTheDocument();
  });

  it("disables submit button until a category is selected", () => {
    render(<FeedbackForm {...defaultProps} />);

    const submitButton = screen.getByRole("button", { name: "Send feedback" });
    expect(submitButton).toBeDisabled();
  });

  it("enables submit button after selecting a category", () => {
    render(<FeedbackForm {...defaultProps} />);

    fireEvent.click(screen.getByLabelText("Missing information"));

    const submitButton = screen.getByRole("button", { name: "Send feedback" });
    expect(submitButton).toBeEnabled();
  });

  it("calls reportDocumentFeedback mutation on submit", async () => {
    render(<FeedbackForm {...defaultProps} />);

    fireEvent.click(screen.getByLabelText("Incorrect data"));
    fireEvent.click(screen.getByRole("button", { name: "Send feedback" }));

    await waitFor(() => {
      expect(fetchSpy).toHaveBeenCalledTimes(1);
    });

    const callBody = JSON.parse(fetchSpy.mock.calls[0][1].body);
    expect(callBody.variables.userId).toBe("user-1");
    expect(callBody.variables.input.fileId).toBe("file-1");
    expect(callBody.variables.input.feedbackType).toBe("EXTRACTION_QUALITY");
    expect(callBody.variables.input.message).toBe("Incorrect data");
  });

  it("shows thank you confirmation after submission", async () => {
    render(<FeedbackForm {...defaultProps} />);

    fireEvent.click(screen.getByLabelText("Other"));
    fireEvent.click(screen.getByRole("button", { name: "Send feedback" }));

    await waitFor(() => {
      expect(screen.getByText("Thank you for your feedback.")).toBeInTheDocument();
    });
  });

  it("calls onSubmitted callback after submission", async () => {
    const onSubmitted = vi.fn();
    render(<FeedbackForm {...defaultProps} onSubmitted={onSubmitted} />);

    fireEvent.click(screen.getByLabelText("Wrong person"));
    fireEvent.click(screen.getByRole("button", { name: "Send feedback" }));

    await waitFor(() => {
      expect(onSubmitted).toHaveBeenCalledTimes(1);
    });
  });

  it("shows thank you even if mutation fails (fire-and-forget)", async () => {
    fetchSpy.mockRejectedValueOnce(new Error("Network error"));
    render(<FeedbackForm {...defaultProps} />);

    fireEvent.click(screen.getByLabelText("Missing information"));
    fireEvent.click(screen.getByRole("button", { name: "Send feedback" }));

    await waitFor(() => {
      expect(screen.getByText("Thank you for your feedback.")).toBeInTheDocument();
    });
  });

  it("composes message from category and free-text description", async () => {
    render(<FeedbackForm {...defaultProps} />);

    fireEvent.click(screen.getByLabelText("Missing information"));
    fireEvent.change(screen.getByPlaceholderText("Tell us more..."), {
      target: { value: "Skills were missed" },
    });
    fireEvent.click(screen.getByRole("button", { name: "Send feedback" }));

    await waitFor(() => {
      const callBody = JSON.parse(fetchSpy.mock.calls[0][1].body);
      expect(callBody.variables.input.message).toBe("Missing information: Skills were missed");
    });
  });
});
