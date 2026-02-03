import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { useMutation } from "urql";
import { beforeEach, describe, expect, it, type Mock, vi } from "vitest";
import { AuthorEditModal } from "./AuthorEditModal";

// Mock urql
vi.mock("urql", () => ({
  useMutation: vi.fn(),
}));

// Mock next/image
vi.mock("next/image", () => ({
  default: ({ src, alt, ...props }: { src: string; alt: string }) => (
    // biome-ignore lint/performance/noImgElement: this is a mock for testing
    <img src={src} alt={alt} {...props} />
  ),
}));

// Mock URL.createObjectURL
const mockCreateObjectURL = vi.fn(() => "blob:http://localhost/mock-image");
global.URL.createObjectURL = mockCreateObjectURL;

describe("AuthorEditModal", () => {
  const mockOnOpenChange = vi.fn();
  const mockOnSuccess = vi.fn();
  const mockUpdateAuthor = vi.fn();

  const defaultAuthor = {
    id: "author-1",
    name: "John Manager",
    title: "Engineering Manager",
    company: "Acme Corp",
    linkedInUrl: "https://linkedin.com/in/johnmanager",
    imageUrl: null,
  };

  beforeEach(() => {
    vi.clearAllMocks();
    (useMutation as Mock).mockReturnValue([{ fetching: false }, mockUpdateAuthor]);
    mockUpdateAuthor.mockResolvedValue({ data: { updateAuthor: { id: "author-1" } } });
  });

  describe("Rendering", () => {
    it("renders dialog with correct title for known author", () => {
      render(
        <AuthorEditModal
          open={true}
          onOpenChange={mockOnOpenChange}
          author={defaultAuthor}
          onSuccess={mockOnSuccess}
        />
      );

      expect(screen.getByRole("dialog")).toBeInTheDocument();
      expect(screen.getByText("Edit Author")).toBeInTheDocument();
      expect(screen.getByText("Update the author's information.")).toBeInTheDocument();
    });

    it("renders dialog with correct title for unknown author", () => {
      const unknownAuthor = { ...defaultAuthor, name: "unknown" };
      render(
        <AuthorEditModal
          open={true}
          onOpenChange={mockOnOpenChange}
          author={unknownAuthor}
          onSuccess={mockOnSuccess}
        />
      );

      expect(screen.getByText("Add Author Details")).toBeInTheDocument();
      expect(
        screen.getByText("The author of this testimonial wasn't detected. Add their details below.")
      ).toBeInTheDocument();
    });

    it("renders all form fields", () => {
      render(
        <AuthorEditModal
          open={true}
          onOpenChange={mockOnOpenChange}
          author={defaultAuthor}
          onSuccess={mockOnSuccess}
        />
      );

      expect(screen.getByLabelText(/name/i)).toBeInTheDocument();
      expect(screen.getByLabelText(/title \/ position/i)).toBeInTheDocument();
      expect(screen.getByLabelText(/company \/ organization/i)).toBeInTheDocument();
      expect(screen.getByLabelText(/linkedin profile/i)).toBeInTheDocument();
      expect(screen.getByRole("button", { name: /upload author photo/i })).toBeInTheDocument();
    });

    it("populates fields with author data", () => {
      render(
        <AuthorEditModal
          open={true}
          onOpenChange={mockOnOpenChange}
          author={defaultAuthor}
          onSuccess={mockOnSuccess}
        />
      );

      expect(screen.getByLabelText(/name/i)).toHaveValue("John Manager");
      expect(screen.getByLabelText(/title \/ position/i)).toHaveValue("Engineering Manager");
      expect(screen.getByLabelText(/company \/ organization/i)).toHaveValue("Acme Corp");
      expect(screen.getByLabelText(/linkedin profile/i)).toHaveValue(
        "https://linkedin.com/in/johnmanager"
      );
    });

    it("clears name field for unknown authors", () => {
      const unknownAuthor = { ...defaultAuthor, name: "unknown" };
      render(
        <AuthorEditModal
          open={true}
          onOpenChange={mockOnOpenChange}
          author={unknownAuthor}
          onSuccess={mockOnSuccess}
        />
      );

      expect(screen.getByLabelText(/name/i)).toHaveValue("");
    });

    it("renders save and cancel buttons", () => {
      render(
        <AuthorEditModal
          open={true}
          onOpenChange={mockOnOpenChange}
          author={defaultAuthor}
          onSuccess={mockOnSuccess}
        />
      );

      expect(screen.getByRole("button", { name: /save/i })).toBeInTheDocument();
      expect(screen.getByRole("button", { name: /cancel/i })).toBeInTheDocument();
    });

    it("shows required indicator on name field", () => {
      render(
        <AuthorEditModal
          open={true}
          onOpenChange={mockOnOpenChange}
          author={defaultAuthor}
          onSuccess={mockOnSuccess}
        />
      );

      // The asterisk should be in the label
      const nameLabel = screen.getByText(/name/i).closest("label");
      expect(nameLabel?.textContent).toContain("*");
    });
  });

  describe("Validation", () => {
    it("shows error when name is empty", async () => {
      const user = userEvent.setup();
      render(
        <AuthorEditModal
          open={true}
          onOpenChange={mockOnOpenChange}
          author={{ ...defaultAuthor, name: "" }}
          onSuccess={mockOnSuccess}
        />
      );

      const saveButton = screen.getByRole("button", { name: /save/i });
      await user.click(saveButton);

      expect(screen.getByText(/name is required/i)).toBeInTheDocument();
      expect(mockUpdateAuthor).not.toHaveBeenCalled();
    });

    it("shows error for invalid LinkedIn URL", async () => {
      const user = userEvent.setup();
      render(
        <AuthorEditModal
          open={true}
          onOpenChange={mockOnOpenChange}
          author={{ ...defaultAuthor, linkedInUrl: null }}
          onSuccess={mockOnSuccess}
        />
      );

      const linkedInInput = screen.getByLabelText(/linkedin profile/i);
      await user.clear(linkedInInput);
      await user.type(linkedInInput, "https://twitter.com/user");

      const saveButton = screen.getByRole("button", { name: /save/i });
      await user.click(saveButton);

      expect(
        screen.getByText(
          /linkedin url must be in the format: https:\/\/linkedin\.com\/in\/username/i
        )
      ).toBeInTheDocument();
      expect(mockUpdateAuthor).not.toHaveBeenCalled();
    });

    it("accepts valid LinkedIn URL formats", async () => {
      const user = userEvent.setup();
      const validUrls = [
        "https://linkedin.com/in/user",
        "https://www.linkedin.com/in/user",
        "https://linkedin.com/in/user-name",
        "http://linkedin.com/in/user123",
      ];

      for (const url of validUrls) {
        vi.clearAllMocks();
        mockUpdateAuthor.mockResolvedValue({ data: { updateAuthor: { id: "author-1" } } });

        const { unmount } = render(
          <AuthorEditModal
            open={true}
            onOpenChange={mockOnOpenChange}
            author={{ ...defaultAuthor, linkedInUrl: null }}
            onSuccess={mockOnSuccess}
          />
        );

        const linkedInInput = screen.getByLabelText(/linkedin profile/i);
        await user.clear(linkedInInput);
        await user.type(linkedInInput, url);

        const saveButton = screen.getByRole("button", { name: /save/i });
        await user.click(saveButton);

        await waitFor(() => {
          expect(mockUpdateAuthor).toHaveBeenCalled();
        });

        unmount();
      }
    });

    it("allows empty LinkedIn URL", async () => {
      const user = userEvent.setup();
      render(
        <AuthorEditModal
          open={true}
          onOpenChange={mockOnOpenChange}
          author={{ ...defaultAuthor, linkedInUrl: null }}
          onSuccess={mockOnSuccess}
        />
      );

      // Leave LinkedIn field empty
      const saveButton = screen.getByRole("button", { name: /save/i });
      await user.click(saveButton);

      await waitFor(() => {
        expect(mockUpdateAuthor).toHaveBeenCalled();
      });
    });

    it("clears error when user starts typing", async () => {
      const user = userEvent.setup();
      render(
        <AuthorEditModal
          open={true}
          onOpenChange={mockOnOpenChange}
          author={{ ...defaultAuthor, name: "" }}
          onSuccess={mockOnSuccess}
        />
      );

      // Trigger validation error
      const saveButton = screen.getByRole("button", { name: /save/i });
      await user.click(saveButton);
      expect(screen.getByText(/name is required/i)).toBeInTheDocument();

      // Start typing
      const nameInput = screen.getByLabelText(/name/i);
      await user.type(nameInput, "J");

      expect(screen.queryByText(/name is required/i)).not.toBeInTheDocument();
    });
  });

  describe("Form Submission", () => {
    it("calls updateAuthor mutation with correct data", async () => {
      const user = userEvent.setup();
      render(
        <AuthorEditModal
          open={true}
          onOpenChange={mockOnOpenChange}
          author={defaultAuthor}
          onSuccess={mockOnSuccess}
        />
      );

      // Modify the name
      const nameInput = screen.getByLabelText(/name/i);
      await user.clear(nameInput);
      await user.type(nameInput, "John A. Manager");

      const saveButton = screen.getByRole("button", { name: /save/i });
      await user.click(saveButton);

      await waitFor(() => {
        expect(mockUpdateAuthor).toHaveBeenCalledWith({
          id: "author-1",
          input: {
            name: "John A. Manager",
            title: "Engineering Manager",
            company: "Acme Corp",
            linkedInUrl: "https://linkedin.com/in/johnmanager",
          },
        });
      });
    });

    it("trims whitespace from form values", async () => {
      const user = userEvent.setup();
      render(
        <AuthorEditModal
          open={true}
          onOpenChange={mockOnOpenChange}
          author={{ ...defaultAuthor, name: "", title: null, company: null, linkedInUrl: null }}
          onSuccess={mockOnSuccess}
        />
      );

      await user.type(screen.getByLabelText(/name/i), "  John Doe  ");
      await user.type(screen.getByLabelText(/title \/ position/i), "  Manager  ");

      const saveButton = screen.getByRole("button", { name: /save/i });
      await user.click(saveButton);

      await waitFor(() => {
        expect(mockUpdateAuthor).toHaveBeenCalledWith({
          id: "author-1",
          input: expect.objectContaining({
            name: "John Doe",
            title: "Manager",
          }),
        });
      });
    });

    it("converts empty optional fields to null", async () => {
      const user = userEvent.setup();
      render(
        <AuthorEditModal
          open={true}
          onOpenChange={mockOnOpenChange}
          author={defaultAuthor}
          onSuccess={mockOnSuccess}
        />
      );

      // Clear optional fields
      await user.clear(screen.getByLabelText(/title \/ position/i));
      await user.clear(screen.getByLabelText(/company \/ organization/i));
      await user.clear(screen.getByLabelText(/linkedin profile/i));

      const saveButton = screen.getByRole("button", { name: /save/i });
      await user.click(saveButton);

      await waitFor(() => {
        expect(mockUpdateAuthor).toHaveBeenCalledWith({
          id: "author-1",
          input: expect.objectContaining({
            title: null,
            company: null,
            linkedInUrl: null,
          }),
        });
      });
    });

    it("calls onSuccess and closes modal after successful submission", async () => {
      const user = userEvent.setup();
      render(
        <AuthorEditModal
          open={true}
          onOpenChange={mockOnOpenChange}
          author={defaultAuthor}
          onSuccess={mockOnSuccess}
        />
      );

      const saveButton = screen.getByRole("button", { name: /save/i });
      await user.click(saveButton);

      await waitFor(() => {
        expect(mockOnSuccess).toHaveBeenCalled();
        expect(mockOnOpenChange).toHaveBeenCalledWith(false);
      });
    });

    it("shows error message on mutation failure", async () => {
      const user = userEvent.setup();
      mockUpdateAuthor.mockResolvedValue({ error: { message: "Server error" } });

      render(
        <AuthorEditModal
          open={true}
          onOpenChange={mockOnOpenChange}
          author={defaultAuthor}
          onSuccess={mockOnSuccess}
        />
      );

      const saveButton = screen.getByRole("button", { name: /save/i });
      await user.click(saveButton);

      await waitFor(() => {
        expect(screen.getByText("Server error")).toBeInTheDocument();
      });

      expect(mockOnSuccess).not.toHaveBeenCalled();
      expect(mockOnOpenChange).not.toHaveBeenCalledWith(false);
    });

    it("closes modal when cancel button is clicked", async () => {
      const user = userEvent.setup();
      render(
        <AuthorEditModal
          open={true}
          onOpenChange={mockOnOpenChange}
          author={defaultAuthor}
          onSuccess={mockOnSuccess}
        />
      );

      const cancelButton = screen.getByRole("button", { name: /cancel/i });
      await user.click(cancelButton);

      expect(mockOnOpenChange).toHaveBeenCalledWith(false);
    });
  });

  describe("Loading State", () => {
    it("disables all inputs when submitting", () => {
      (useMutation as Mock).mockReturnValue([{ fetching: true }, mockUpdateAuthor]);

      render(
        <AuthorEditModal
          open={true}
          onOpenChange={mockOnOpenChange}
          author={defaultAuthor}
          onSuccess={mockOnSuccess}
        />
      );

      expect(screen.getByLabelText(/name/i)).toBeDisabled();
      expect(screen.getByLabelText(/title \/ position/i)).toBeDisabled();
      expect(screen.getByLabelText(/company \/ organization/i)).toBeDisabled();
      expect(screen.getByLabelText(/linkedin profile/i)).toBeDisabled();
    });

    it("shows saving text when submitting", () => {
      (useMutation as Mock).mockReturnValue([{ fetching: true }, mockUpdateAuthor]);

      render(
        <AuthorEditModal
          open={true}
          onOpenChange={mockOnOpenChange}
          author={defaultAuthor}
          onSuccess={mockOnSuccess}
        />
      );

      expect(screen.getByRole("button", { name: /saving/i })).toBeInTheDocument();
    });

    it("disables buttons when submitting", () => {
      (useMutation as Mock).mockReturnValue([{ fetching: true }, mockUpdateAuthor]);

      render(
        <AuthorEditModal
          open={true}
          onOpenChange={mockOnOpenChange}
          author={defaultAuthor}
          onSuccess={mockOnSuccess}
        />
      );

      expect(screen.getByRole("button", { name: /saving/i })).toBeDisabled();
      expect(screen.getByRole("button", { name: /cancel/i })).toBeDisabled();
    });
  });

  describe("Image Upload", () => {
    it("renders image upload button", () => {
      render(
        <AuthorEditModal
          open={true}
          onOpenChange={mockOnOpenChange}
          author={defaultAuthor}
          onSuccess={mockOnSuccess}
        />
      );

      expect(screen.getByRole("button", { name: /upload author photo/i })).toBeInTheDocument();
      expect(screen.getByText(/add photo/i)).toBeInTheDocument();
    });

    it("shows existing image when author has imageUrl", () => {
      const authorWithImage = {
        ...defaultAuthor,
        imageUrl: "https://example.com/photo.jpg",
      };

      render(
        <AuthorEditModal
          open={true}
          onOpenChange={mockOnOpenChange}
          author={authorWithImage}
          onSuccess={mockOnSuccess}
        />
      );

      const image = screen.getByAltText(/photo of john manager/i);
      expect(image).toHaveAttribute("src", "https://example.com/photo.jpg");
    });

    it("shows remove button when image is present", () => {
      const authorWithImage = {
        ...defaultAuthor,
        imageUrl: "https://example.com/photo.jpg",
      };

      render(
        <AuthorEditModal
          open={true}
          onOpenChange={mockOnOpenChange}
          author={authorWithImage}
          onSuccess={mockOnSuccess}
        />
      );

      expect(screen.getByRole("button", { name: /remove photo/i })).toBeInTheDocument();
    });

    it("previews selected image file", async () => {
      const user = userEvent.setup();
      render(
        <AuthorEditModal
          open={true}
          onOpenChange={mockOnOpenChange}
          author={defaultAuthor}
          onSuccess={mockOnSuccess}
        />
      );

      const file = new File(["test"], "photo.jpg", { type: "image/jpeg" });
      const fileInput = screen.getByLabelText(/upload author photo/i, {
        selector: 'input[type="file"]',
      });

      await user.upload(fileInput, file);

      // Should show the blob URL created by URL.createObjectURL
      await waitFor(() => {
        const image = screen.getByAltText(/photo of john manager/i);
        expect(image).toHaveAttribute("src", "blob:http://localhost/mock-image");
      });
    });

    it("removes image preview when remove button is clicked", async () => {
      const user = userEvent.setup();
      const authorWithImage = {
        ...defaultAuthor,
        imageUrl: "https://example.com/photo.jpg",
      };

      render(
        <AuthorEditModal
          open={true}
          onOpenChange={mockOnOpenChange}
          author={authorWithImage}
          onSuccess={mockOnSuccess}
        />
      );

      // Verify image is shown
      expect(screen.getByAltText(/photo of john manager/i)).toBeInTheDocument();

      // Click remove button
      const removeButton = screen.getByRole("button", { name: /remove photo/i });
      await user.click(removeButton);

      // Image should be replaced with "Add photo" text
      expect(screen.queryByAltText(/photo of john manager/i)).not.toBeInTheDocument();
      expect(screen.getByText(/add photo/i)).toBeInTheDocument();
    });

    it("hides remove button when submitting", () => {
      (useMutation as Mock).mockReturnValue([{ fetching: true }, mockUpdateAuthor]);

      const authorWithImage = {
        ...defaultAuthor,
        imageUrl: "https://example.com/photo.jpg",
      };

      render(
        <AuthorEditModal
          open={true}
          onOpenChange={mockOnOpenChange}
          author={authorWithImage}
          onSuccess={mockOnSuccess}
        />
      );

      expect(screen.queryByRole("button", { name: /remove photo/i })).not.toBeInTheDocument();
    });
  });
});
