import { fireEvent, render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import { ProfileActions, ProfileActionsBar } from "./ProfileActions";

describe("ProfileActions (mobile card)", () => {
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

  it("returns null when no handlers provided", () => {
    const { container } = render(<ProfileActions />);
    expect(container.innerHTML).toBe("");
  });
});

describe("ProfileActionsBar (desktop icon bar)", () => {
  it("renders icon buttons with aria-labels when handlers provided", () => {
    render(
      <ProfileActionsBar onAddReference={vi.fn()} onExport={vi.fn()} onUploadAnother={vi.fn()} />
    );
    expect(screen.getByRole("button", { name: "Add Reference Letter" })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Export PDF" })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Upload Another Resume" })).toBeInTheDocument();
  });

  it("renders tooltips with action labels", () => {
    render(<ProfileActionsBar onAddReference={vi.fn()} onExport={vi.fn()} />);
    const tooltips = screen.getAllByRole("tooltip");
    expect(tooltips).toHaveLength(2);
    expect(tooltips[0]).toHaveTextContent("Add Reference Letter");
    expect(tooltips[1]).toHaveTextContent("Export PDF");
  });

  it("calls handlers when icon buttons are clicked", () => {
    const onAddReference = vi.fn();
    const onExport = vi.fn();
    const onUploadAnother = vi.fn();
    render(
      <ProfileActionsBar
        onAddReference={onAddReference}
        onExport={onExport}
        onUploadAnother={onUploadAnother}
      />
    );

    fireEvent.click(screen.getByRole("button", { name: "Add Reference Letter" }));
    expect(onAddReference).toHaveBeenCalledOnce();

    fireEvent.click(screen.getByRole("button", { name: "Export PDF" }));
    expect(onExport).toHaveBeenCalledOnce();

    fireEvent.click(screen.getByRole("button", { name: "Upload Another Resume" }));
    expect(onUploadAnother).toHaveBeenCalledOnce();
  });

  it("returns null when no handlers provided", () => {
    const { container } = render(<ProfileActionsBar />);
    expect(container.innerHTML).toBe("");
  });

  it("only renders buttons for provided handlers", () => {
    render(<ProfileActionsBar onExport={vi.fn()} />);
    expect(screen.getByRole("button", { name: "Export PDF" })).toBeInTheDocument();
    expect(screen.queryByRole("button", { name: "Add Reference Letter" })).not.toBeInTheDocument();
    expect(screen.queryByRole("button", { name: "Upload Another Resume" })).not.toBeInTheDocument();
  });
});
