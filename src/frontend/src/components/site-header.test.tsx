import { render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";

vi.mock("next-themes", () => ({
  useTheme: () => ({
    theme: "system",
    setTheme: vi.fn(),
  }),
}));

import { SiteHeader } from "./site-header";

describe("SiteHeader", () => {
  it("renders the app name as a link to home", () => {
    render(<SiteHeader />);
    const link = screen.getByRole("link", { name: "Credfolio" });
    expect(link).toHaveAttribute("href", "/");
  });

  it("renders the theme toggle button", () => {
    render(<SiteHeader />);
    expect(screen.getByRole("button", { name: /toggle theme/i })).toBeInTheDocument();
  });

  it("renders as a banner landmark", () => {
    render(<SiteHeader />);
    expect(screen.getByRole("banner")).toBeInTheDocument();
  });

  it("renders upload navigation link", () => {
    render(<SiteHeader />);
    const link = screen.getByRole("link", { name: "Upload" });
    expect(link).toHaveAttribute("href", "/upload");
  });
});
