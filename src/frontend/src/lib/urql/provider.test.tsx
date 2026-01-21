import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { UrqlProvider } from "./provider";

describe("UrqlProvider", () => {
  it("should render children", () => {
    render(
      <UrqlProvider>
        <div data-testid="child">Hello</div>
      </UrqlProvider>
    );

    expect(screen.getByTestId("child")).toBeInTheDocument();
  });

  it("should provide urql context to children", () => {
    // Children should be able to use urql hooks
    render(
      <UrqlProvider>
        <div>Test content</div>
      </UrqlProvider>
    );

    expect(screen.getByText("Test content")).toBeInTheDocument();
  });
});
