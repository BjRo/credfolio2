import { render } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import Home from "./page";

describe("Home Page", () => {
  it("should render successfully", () => {
    render(<Home />);
    expect(document.body).toBeTruthy();
  });
});
