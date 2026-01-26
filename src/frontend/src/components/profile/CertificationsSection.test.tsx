import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { type Certification, CertificationsSection } from "./CertificationsSection";

const mockCertifications: Certification[] = [
  {
    name: "AWS Solutions Architect",
    issuer: "Amazon Web Services",
    date: "2023",
    url: "https://aws.amazon.com/verify/123",
  },
  {
    name: "Google Cloud Professional",
    issuer: "Google",
    date: "2022",
    url: null,
  },
  {
    name: "Kubernetes Administrator",
    issuer: "CNCF",
    date: null,
    url: null,
  },
];

describe("CertificationsSection", () => {
  it("renders the section heading", () => {
    render(<CertificationsSection certifications={mockCertifications} />);
    expect(screen.getByRole("heading", { level: 2, name: "Certifications" })).toBeInTheDocument();
  });

  it("renders certification names", () => {
    render(<CertificationsSection certifications={mockCertifications} />);
    expect(screen.getByText("AWS Solutions Architect")).toBeInTheDocument();
    expect(screen.getByText("Google Cloud Professional")).toBeInTheDocument();
    expect(screen.getByText("Kubernetes Administrator")).toBeInTheDocument();
  });

  it("renders issuers", () => {
    render(<CertificationsSection certifications={mockCertifications} />);
    expect(screen.getByText("Amazon Web Services")).toBeInTheDocument();
    expect(screen.getByText("Google")).toBeInTheDocument();
    expect(screen.getByText("CNCF")).toBeInTheDocument();
  });

  it("renders dates when provided", () => {
    render(<CertificationsSection certifications={mockCertifications} />);
    expect(screen.getByText("2023")).toBeInTheDocument();
    expect(screen.getByText("2022")).toBeInTheDocument();
  });

  it("renders credential link when URL is provided", () => {
    render(<CertificationsSection certifications={mockCertifications} />);
    const link = screen.getByRole("link", { name: /AWS Solutions Architect/i });
    expect(link).toHaveAttribute("href", "https://aws.amazon.com/verify/123");
    expect(link).toHaveAttribute("target", "_blank");
    expect(link).toHaveAttribute("rel", "noopener noreferrer");
  });

  it("does not render link when URL is not provided", () => {
    render(<CertificationsSection certifications={mockCertifications} />);
    expect(
      screen.queryByRole("link", { name: /Google Cloud Professional/i })
    ).not.toBeInTheDocument();
    expect(screen.getByText("Google Cloud Professional")).toBeInTheDocument();
  });

  it("returns null when certifications array is empty", () => {
    const { container } = render(<CertificationsSection certifications={[]} />);
    expect(container.firstChild).toBeNull();
  });
});
