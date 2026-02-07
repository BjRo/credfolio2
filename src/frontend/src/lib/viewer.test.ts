import { describe, expect, it } from "vitest";
import { buildViewerUrl } from "./viewer";

describe("buildViewerUrl", () => {
  it("builds URL with letterId only", () => {
    const url = buildViewerUrl("abc-123");
    expect(url).toBe("/viewer?letterId=abc-123");
  });

  it("builds URL with letterId and highlight text", () => {
    const url = buildViewerUrl("abc-123", "great engineer");
    expect(url).toBe("/viewer?letterId=abc-123&highlight=great+engineer");
  });

  it("encodes special characters in highlight text", () => {
    const url = buildViewerUrl("abc-123", "skills: React & TypeScript");
    expect(url).toContain("letterId=abc-123");
    expect(url).toContain("highlight=");
    // Verify round-trip decoding
    const params = new URLSearchParams(url.split("?")[1]);
    expect(params.get("highlight")).toBe("skills: React & TypeScript");
  });

  it("omits highlight param when highlight is empty string", () => {
    const url = buildViewerUrl("abc-123", "");
    expect(url).toBe("/viewer?letterId=abc-123");
  });

  it("truncates highlight text longer than 500 characters", () => {
    const longText = "a".repeat(600);
    const url = buildViewerUrl("abc-123", longText);
    const params = new URLSearchParams(url.split("?")[1]);
    expect(params.get("highlight")?.length).toBe(500);
  });

  it("does not truncate highlight text at exactly 500 characters", () => {
    const exactText = "b".repeat(500);
    const url = buildViewerUrl("abc-123", exactText);
    const params = new URLSearchParams(url.split("?")[1]);
    expect(params.get("highlight")?.length).toBe(500);
  });
});
