import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";
import { TestimonialRelationship } from "@/graphql/generated/graphql";
import { TestimonialsSection } from "./TestimonialsSection";

const mockTestimonials = [
  {
    __typename: "Testimonial" as const,
    id: "1",
    quote: "Great team player with excellent leadership skills.",
    authorName: "John Manager",
    authorTitle: "Engineering Manager",
    authorCompany: "Acme Corp",
    relationship: TestimonialRelationship.Manager,
    createdAt: "2024-01-01T00:00:00Z",
  },
  {
    __typename: "Testimonial" as const,
    id: "2",
    quote: "A brilliant collaborator who consistently delivers high-quality work.",
    authorName: "Sarah Peer",
    authorTitle: "Senior Engineer",
    authorCompany: "Acme Corp",
    relationship: TestimonialRelationship.Peer,
    createdAt: "2024-01-02T00:00:00Z",
  },
];

describe("TestimonialsSection", () => {
  it("renders the section heading", () => {
    render(<TestimonialsSection testimonials={mockTestimonials} />);
    expect(screen.getByRole("heading", { level: 2, name: "What Others Say" })).toBeInTheDocument();
  });

  it("returns null when no testimonials and no onAddReference", () => {
    const { container } = render(<TestimonialsSection testimonials={[]} />);
    expect(container.firstChild).toBeNull();
  });

  it("renders testimonial quotes", () => {
    render(<TestimonialsSection testimonials={mockTestimonials} />);
    expect(
      screen.getByText("Great team player with excellent leadership skills.")
    ).toBeInTheDocument();
    expect(
      screen.getByText("A brilliant collaborator who consistently delivers high-quality work.")
    ).toBeInTheDocument();
  });

  it("renders author names", () => {
    render(<TestimonialsSection testimonials={mockTestimonials} />);
    expect(screen.getByText("John Manager")).toBeInTheDocument();
    expect(screen.getByText("Sarah Peer")).toBeInTheDocument();
  });

  it("renders author title and company", () => {
    render(<TestimonialsSection testimonials={mockTestimonials} />);
    expect(screen.getByText("Engineering Manager at Acme Corp")).toBeInTheDocument();
    expect(screen.getByText("Senior Engineer at Acme Corp")).toBeInTheDocument();
  });

  it("renders relationship badges", () => {
    render(<TestimonialsSection testimonials={mockTestimonials} />);
    expect(screen.getByText("Manager")).toBeInTheDocument();
    expect(screen.getByText("Peer")).toBeInTheDocument();
  });

  describe("Empty state", () => {
    it("shows empty state when onAddReference is provided", () => {
      render(<TestimonialsSection testimonials={[]} onAddReference={() => {}} />);
      expect(screen.getByText("No testimonials yet.")).toBeInTheDocument();
      expect(screen.getByText(/Add a reference letter/)).toBeInTheDocument();
    });

    it("shows add button in empty state", () => {
      render(<TestimonialsSection testimonials={[]} onAddReference={() => {}} />);
      expect(screen.getByRole("button", { name: "Add Reference Letter" })).toBeInTheDocument();
    });

    it("calls onAddReference when add button is clicked", async () => {
      const user = userEvent.setup();
      const onAddReference = vi.fn();
      render(<TestimonialsSection testimonials={[]} onAddReference={onAddReference} />);

      await user.click(screen.getByRole("button", { name: "Add Reference Letter" }));
      expect(onAddReference).toHaveBeenCalled();
    });
  });

  describe("Loading state", () => {
    it("shows loading skeleton when isLoading is true", () => {
      render(<TestimonialsSection testimonials={[]} isLoading={true} onAddReference={() => {}} />);
      // Loading state shows animated skeleton
      const skeletons = document.querySelectorAll(".animate-pulse");
      expect(skeletons.length).toBeGreaterThan(0);
    });
  });

  describe("Header add button", () => {
    it("shows add button in header when onAddReference is provided", () => {
      render(<TestimonialsSection testimonials={mockTestimonials} onAddReference={() => {}} />);
      expect(screen.getByRole("button", { name: "Add reference letter" })).toBeInTheDocument();
    });

    it("calls onAddReference when header add button is clicked", async () => {
      const user = userEvent.setup();
      const onAddReference = vi.fn();
      render(
        <TestimonialsSection testimonials={mockTestimonials} onAddReference={onAddReference} />
      );

      await user.click(screen.getByRole("button", { name: "Add reference letter" }));
      expect(onAddReference).toHaveBeenCalled();
    });

    it("hides add button in header when onAddReference is not provided", () => {
      render(<TestimonialsSection testimonials={mockTestimonials} />);
      expect(
        screen.queryByRole("button", { name: "Add reference letter" })
      ).not.toBeInTheDocument();
    });
  });
});
