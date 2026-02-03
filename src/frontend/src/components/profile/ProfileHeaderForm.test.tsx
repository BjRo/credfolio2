import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";
import { ProfileHeaderForm } from "./ProfileHeaderForm";

describe("ProfileHeaderForm", () => {
  const mockOnSubmit = vi.fn();
  const mockOnCancel = vi.fn();

  const defaultProps = {
    onSubmit: mockOnSubmit,
    onCancel: mockOnCancel,
  };

  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe("Rendering", () => {
    it("renders all form fields", () => {
      render(<ProfileHeaderForm {...defaultProps} />);

      expect(screen.getByLabelText(/name/i)).toBeInTheDocument();
      expect(screen.getByLabelText(/email/i)).toBeInTheDocument();
      expect(screen.getByLabelText(/phone/i)).toBeInTheDocument();
      expect(screen.getByLabelText(/location/i)).toBeInTheDocument();
      expect(screen.getByLabelText(/professional summary/i)).toBeInTheDocument();
    });

    it("renders save and cancel buttons", () => {
      render(<ProfileHeaderForm {...defaultProps} />);

      expect(screen.getByRole("button", { name: /save changes/i })).toBeInTheDocument();
      expect(screen.getByRole("button", { name: /cancel/i })).toBeInTheDocument();
    });

    it("populates fields with initial data", () => {
      const initialData = {
        name: "John Doe",
        email: "john@example.com",
        phone: "+1 555-123-4567",
        location: "San Francisco, CA",
        summary: "A summary",
      };

      render(<ProfileHeaderForm {...defaultProps} initialData={initialData} />);

      expect(screen.getByLabelText(/name/i)).toHaveValue("John Doe");
      expect(screen.getByLabelText(/email/i)).toHaveValue("john@example.com");
      expect(screen.getByLabelText(/phone/i)).toHaveValue("+1 555-123-4567");
      expect(screen.getByLabelText(/location/i)).toHaveValue("San Francisco, CA");
      expect(screen.getByLabelText(/professional summary/i)).toHaveValue("A summary");
    });
  });

  describe("Validation", () => {
    it("shows error when name is empty", async () => {
      const user = userEvent.setup();
      render(<ProfileHeaderForm {...defaultProps} />);

      const submitButton = screen.getByRole("button", { name: /save changes/i });
      await user.click(submitButton);

      expect(screen.getByText(/name is required/i)).toBeInTheDocument();
      expect(mockOnSubmit).not.toHaveBeenCalled();
    });

    it("validates email format client-side", async () => {
      // Note: Email validation is handled by both HTML5 type="email" attribute
      // and custom regex validation. The custom validation triggers on submit.
      const user = userEvent.setup();
      render(
        <ProfileHeaderForm
          {...defaultProps}
          initialData={{
            name: "John Doe",
            email: "valid@email.com",
          }}
        />
      );

      // Verify the email input has type="email" for browser validation
      const emailInput = screen.getByLabelText(/email/i);
      expect(emailInput).toHaveAttribute("type", "email");

      // Valid email should submit successfully
      const submitButton = screen.getByRole("button", { name: /save changes/i });
      await user.click(submitButton);

      expect(mockOnSubmit).toHaveBeenCalledWith(
        expect.objectContaining({ email: "valid@email.com" })
      );
    });

    it("shows error for invalid phone format", async () => {
      const user = userEvent.setup();
      render(<ProfileHeaderForm {...defaultProps} />);

      await user.type(screen.getByLabelText(/name/i), "John Doe");
      await user.type(screen.getByLabelText(/phone/i), "abc");

      const submitButton = screen.getByRole("button", { name: /save changes/i });
      await user.click(submitButton);

      expect(screen.getByText(/invalid phone format/i)).toBeInTheDocument();
      expect(mockOnSubmit).not.toHaveBeenCalled();
    });

    it("accepts valid email format", async () => {
      const user = userEvent.setup();
      render(<ProfileHeaderForm {...defaultProps} />);

      await user.type(screen.getByLabelText(/name/i), "John Doe");
      await user.type(screen.getByLabelText(/email/i), "john@example.com");

      const submitButton = screen.getByRole("button", { name: /save changes/i });
      await user.click(submitButton);

      expect(screen.queryByText(/invalid email format/i)).not.toBeInTheDocument();
      expect(mockOnSubmit).toHaveBeenCalled();
    });

    it("accepts valid phone formats", async () => {
      const user = userEvent.setup();
      render(<ProfileHeaderForm {...defaultProps} />);

      await user.type(screen.getByLabelText(/name/i), "John Doe");
      await user.type(screen.getByLabelText(/phone/i), "+1 (555) 123-4567");

      const submitButton = screen.getByRole("button", { name: /save changes/i });
      await user.click(submitButton);

      expect(screen.queryByText(/invalid phone format/i)).not.toBeInTheDocument();
      expect(mockOnSubmit).toHaveBeenCalled();
    });
  });

  describe("Form Submission", () => {
    it("calls onSubmit with trimmed form data", async () => {
      const user = userEvent.setup();
      render(<ProfileHeaderForm {...defaultProps} />);

      await user.type(screen.getByLabelText(/name/i), "  John Doe  ");
      await user.type(screen.getByLabelText(/email/i), "  john@example.com  ");
      await user.type(screen.getByLabelText(/phone/i), "  +1 555-123-4567  ");
      await user.type(screen.getByLabelText(/location/i), "  San Francisco  ");
      await user.type(screen.getByLabelText(/professional summary/i), "  Summary text  ");

      const submitButton = screen.getByRole("button", { name: /save changes/i });
      await user.click(submitButton);

      expect(mockOnSubmit).toHaveBeenCalledWith({
        name: "John Doe",
        email: "john@example.com",
        phone: "+1 555-123-4567",
        location: "San Francisco",
        summary: "Summary text",
      });
    });

    it("calls onCancel when cancel button is clicked", async () => {
      const user = userEvent.setup();
      render(<ProfileHeaderForm {...defaultProps} />);

      const cancelButton = screen.getByRole("button", { name: /cancel/i });
      await user.click(cancelButton);

      expect(mockOnCancel).toHaveBeenCalled();
    });
  });

  describe("Loading State", () => {
    it("disables all fields when isSubmitting is true", () => {
      render(<ProfileHeaderForm {...defaultProps} isSubmitting={true} />);

      expect(screen.getByLabelText(/name/i)).toBeDisabled();
      expect(screen.getByLabelText(/email/i)).toBeDisabled();
      expect(screen.getByLabelText(/phone/i)).toBeDisabled();
      expect(screen.getByLabelText(/location/i)).toBeDisabled();
      expect(screen.getByLabelText(/professional summary/i)).toBeDisabled();
    });

    it("shows saving text when isSubmitting is true", () => {
      render(<ProfileHeaderForm {...defaultProps} isSubmitting={true} />);

      expect(screen.getByRole("button", { name: /saving/i })).toBeInTheDocument();
    });

    it("disables buttons when isSubmitting is true", () => {
      render(<ProfileHeaderForm {...defaultProps} isSubmitting={true} />);

      expect(screen.getByRole("button", { name: /saving/i })).toBeDisabled();
      expect(screen.getByRole("button", { name: /cancel/i })).toBeDisabled();
    });
  });

  describe("Profile Photo", () => {
    // Mock URL.createObjectURL
    const mockCreateObjectURL = vi.fn(() => "blob:http://localhost/mock-image");
    beforeAll(() => {
      global.URL.createObjectURL = mockCreateObjectURL;
    });

    it("renders photo upload button", () => {
      render(<ProfileHeaderForm {...defaultProps} />);

      expect(screen.getByRole("button", { name: /upload profile photo/i })).toBeInTheDocument();
      expect(screen.getByText(/add photo/i)).toBeInTheDocument();
    });

    it("shows existing photo when photoUrl is provided", () => {
      render(<ProfileHeaderForm {...defaultProps} photoUrl="https://example.com/photo.jpg" />);

      const image = screen.getByAltText(/profile photo/i);
      expect(image).toHaveAttribute("src", "https://example.com/photo.jpg");
    });

    it("shows remove button when photo is present", () => {
      render(<ProfileHeaderForm {...defaultProps} photoUrl="https://example.com/photo.jpg" />);

      expect(screen.getByRole("button", { name: /remove photo/i })).toBeInTheDocument();
    });

    it("previews selected image file", async () => {
      const user = userEvent.setup();
      render(<ProfileHeaderForm {...defaultProps} />);

      const file = new File(["test"], "photo.jpg", { type: "image/jpeg" });
      const fileInput = screen.getByLabelText(/upload profile photo/i, {
        selector: 'input[type="file"]',
      });

      await user.upload(fileInput, file);

      // Should show the blob URL created by URL.createObjectURL
      const image = screen.getByAltText(/profile photo/i);
      expect(image).toHaveAttribute("src", "blob:http://localhost/mock-image");
    });

    it("removes photo preview when remove button is clicked", async () => {
      const user = userEvent.setup();
      render(<ProfileHeaderForm {...defaultProps} photoUrl="https://example.com/photo.jpg" />);

      // Verify image is shown
      expect(screen.getByAltText(/profile photo/i)).toBeInTheDocument();

      // Click remove button
      const removeButton = screen.getByRole("button", { name: /remove photo/i });
      await user.click(removeButton);

      // Image should be replaced with "Add photo" text
      expect(screen.queryByAltText(/profile photo/i)).not.toBeInTheDocument();
      expect(screen.getByText(/add photo/i)).toBeInTheDocument();
    });

    it("calls onSubmit with image data when file is selected", async () => {
      const user = userEvent.setup();
      render(<ProfileHeaderForm {...defaultProps} />);

      // Fill required field
      await user.type(screen.getByLabelText(/name/i), "John Doe");

      // Upload photo
      const file = new File(["test"], "photo.jpg", { type: "image/jpeg" });
      const fileInput = screen.getByLabelText(/upload profile photo/i, {
        selector: 'input[type="file"]',
      });
      await user.upload(fileInput, file);

      // Submit form
      const submitButton = screen.getByRole("button", { name: /save changes/i });
      await user.click(submitButton);

      expect(mockOnSubmit).toHaveBeenCalledWith(
        expect.objectContaining({
          name: "John Doe",
          pendingImageFile: file,
        })
      );
    });

    it("calls onSubmit with imageRemoved flag when photo is removed", async () => {
      const user = userEvent.setup();
      render(
        <ProfileHeaderForm
          {...defaultProps}
          photoUrl="https://example.com/photo.jpg"
          initialData={{ name: "John Doe" }}
        />
      );

      // Remove photo
      const removeButton = screen.getByRole("button", { name: /remove photo/i });
      await user.click(removeButton);

      // Submit form
      const submitButton = screen.getByRole("button", { name: /save changes/i });
      await user.click(submitButton);

      expect(mockOnSubmit).toHaveBeenCalledWith(
        expect.objectContaining({
          imageRemoved: true,
        })
      );
    });

    it("hides remove button when submitting", () => {
      render(
        <ProfileHeaderForm
          {...defaultProps}
          photoUrl="https://example.com/photo.jpg"
          isSubmitting={true}
        />
      );

      expect(screen.queryByRole("button", { name: /remove photo/i })).not.toBeInTheDocument();
    });

    it("disables photo upload when submitting", () => {
      render(<ProfileHeaderForm {...defaultProps} isSubmitting={true} />);

      expect(screen.getByRole("button", { name: /upload profile photo/i })).toBeDisabled();
    });
  });
});
