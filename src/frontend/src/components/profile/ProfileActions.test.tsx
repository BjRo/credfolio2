import { fireEvent, render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import { ProfileActions } from "./ProfileActions";

describe("ProfileActions", () => {
  it("renders Add Reference button when handler provided", () => {
    const onAddReference = vi.fn();
    render(<ProfileActions onAddReference={onAddReference} />);
    expect(screen.getByRole("button", { name: /add reference letter/i })).toBeInTheDocument();
  });

  it("renders Export PDF button when handler provided", () => {
    const onExport = vi.fn();
    render(<ProfileActions onExport={onExport} />);
    expect(screen.getByRole("button", { name: /export pdf/i })).toBeInTheDocument();
  });

  it("renders Upload Another button when handler provided", () => {
    const onUploadAnother = vi.fn();
    render(<ProfileActions onUploadAnother={onUploadAnother} />);
    expect(screen.getByRole("button", { name: /upload another resume/i })).toBeInTheDocument();
  });

  it("calls onAddReference when clicked", () => {
    const onAddReference = vi.fn();
    render(<ProfileActions onAddReference={onAddReference} />);
    fireEvent.click(screen.getByRole("button", { name: /add reference letter/i }));
    expect(onAddReference).toHaveBeenCalledOnce();
  });

  it("calls onExport when clicked", () => {
    const onExport = vi.fn();
    render(<ProfileActions onExport={onExport} />);
    fireEvent.click(screen.getByRole("button", { name: /export pdf/i }));
    expect(onExport).toHaveBeenCalledOnce();
  });

  it("calls onUploadAnother when clicked", () => {
    const onUploadAnother = vi.fn();
    render(<ProfileActions onUploadAnother={onUploadAnother} />);
    fireEvent.click(screen.getByRole("button", { name: /upload another resume/i }));
    expect(onUploadAnother).toHaveBeenCalledOnce();
  });

  it("does not render buttons when handlers not provided", () => {
    render(<ProfileActions />);
    expect(screen.queryByRole("button", { name: /add reference letter/i })).not.toBeInTheDocument();
    expect(screen.queryByRole("button", { name: /export pdf/i })).not.toBeInTheDocument();
    expect(
      screen.queryByRole("button", { name: /upload another resume/i })
    ).not.toBeInTheDocument();
  });
});
