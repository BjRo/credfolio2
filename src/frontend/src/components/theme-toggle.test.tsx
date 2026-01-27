import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";

const mockSetTheme = vi.fn();
let mockTheme = "system";
let mockResolvedTheme = "light";
let mockSystemTheme: string | undefined = "light";

vi.mock("next-themes", () => ({
  useTheme: () => ({
    theme: mockTheme,
    resolvedTheme: mockResolvedTheme,
    systemTheme: mockSystemTheme,
    setTheme: mockSetTheme,
  }),
}));

import { ThemeToggle } from "./theme-toggle";

describe("ThemeToggle", () => {
  beforeEach(() => {
    mockTheme = "system";
    mockResolvedTheme = "light";
    mockSystemTheme = "light";
    mockSetTheme.mockClear();
  });

  it("renders a button with accessible label", () => {
    render(<ThemeToggle />);
    expect(screen.getByRole("button", { name: /toggle theme/i })).toBeInTheDocument();
  });

  it("cycles from light to dark on click", async () => {
    mockTheme = "light";
    mockResolvedTheme = "light";
    const user = userEvent.setup();
    render(<ThemeToggle />);

    await user.click(screen.getByRole("button", { name: /toggle theme/i }));
    expect(mockSetTheme).toHaveBeenCalledWith("dark");
  });

  it("cycles from dark to system on click", async () => {
    mockTheme = "dark";
    mockResolvedTheme = "dark";
    const user = userEvent.setup();
    render(<ThemeToggle />);

    await user.click(screen.getByRole("button", { name: /toggle theme/i }));
    expect(mockSetTheme).toHaveBeenCalledWith("system");
  });

  it("skips redundant theme when system matches next in cycle", async () => {
    mockTheme = "system";
    mockResolvedTheme = "light";
    const user = userEvent.setup();
    render(<ThemeToggle />);

    // "system" → would be "light", but resolvedTheme is already "light", so skip to "dark"
    await user.click(screen.getByRole("button", { name: /toggle theme/i }));
    expect(mockSetTheme).toHaveBeenCalledWith("dark");
  });

  it("cycles from system to light when system resolves to dark", async () => {
    mockTheme = "system";
    mockResolvedTheme = "dark";
    mockSystemTheme = "dark";
    const user = userEvent.setup();
    render(<ThemeToggle />);

    // "system" → "light", resolvedTheme is "dark" so no skip needed
    await user.click(screen.getByRole("button", { name: /toggle theme/i }));
    expect(mockSetTheme).toHaveBeenCalledWith("light");
  });

  it("skips system when cycling from dark and OS also prefers dark", async () => {
    mockTheme = "dark";
    mockResolvedTheme = "dark";
    mockSystemTheme = "dark";
    const user = userEvent.setup();
    render(<ThemeToggle />);

    // "dark" → would be "system", but system resolves to "dark" too, so skip to "light"
    await user.click(screen.getByRole("button", { name: /toggle theme/i }));
    expect(mockSetTheme).toHaveBeenCalledWith("light");
  });

  it("skips system when cycling from light and OS also prefers light", async () => {
    mockTheme = "light";
    mockResolvedTheme = "light";
    mockSystemTheme = "light";
    const user = userEvent.setup();
    render(<ThemeToggle />);

    // "light" → would be "dark", no skip needed (dark !== light)
    await user.click(screen.getByRole("button", { name: /toggle theme/i }));
    expect(mockSetTheme).toHaveBeenCalledWith("dark");
  });
});
