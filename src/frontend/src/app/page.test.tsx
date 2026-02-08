import { redirect } from "next/navigation";
import { describe, expect, it, vi } from "vitest";
import { createUrqlClient } from "@/lib/urql/client";
import Home from "./page";

// Mock next/navigation redirect
vi.mock("next/navigation", () => ({
  redirect: vi.fn(),
}));

// Mock urql client
vi.mock("@/lib/urql/client", () => ({
  createUrqlClient: vi.fn(),
  GRAPHQL_ENDPOINT: "http://localhost:8080/graphql",
}));

const mockCreateUrqlClient = createUrqlClient as unknown as ReturnType<typeof vi.fn>;
const mockRedirect = redirect as unknown as ReturnType<typeof vi.fn>;

describe("Home Page (Server Component)", () => {
  it("redirects to upload page when user has no profile", async () => {
    // Mock urql client to return no profile
    const mockQuery = vi.fn().mockReturnValue({
      toPromise: vi.fn().mockResolvedValue({
        data: { profileByUserId: null },
      }),
    });
    mockCreateUrqlClient.mockReturnValue({ query: mockQuery } as never);

    // Call the server component
    await Home();

    // Verify redirect was called with /upload
    expect(mockRedirect).toHaveBeenCalledWith("/upload");
  });

  it("redirects to profile page when profile exists", async () => {
    // Mock urql client to return a profile
    const mockQuery = vi.fn().mockReturnValue({
      toPromise: vi.fn().mockResolvedValue({
        data: { profileByUserId: { id: "profile-123" } },
      }),
    });
    mockCreateUrqlClient.mockReturnValue({ query: mockQuery } as never);

    // Call the server component
    await Home();

    // Verify redirect was called with profile ID
    expect(mockRedirect).toHaveBeenCalledWith("/profile/profile-123");
  });

  it("redirects to upload page when query returns undefined data", async () => {
    // Mock urql client to return undefined
    const mockQuery = vi.fn().mockReturnValue({
      toPromise: vi.fn().mockResolvedValue({
        data: undefined,
      }),
    });
    mockCreateUrqlClient.mockReturnValue({ query: mockQuery } as never);

    // Call the server component
    await Home();

    // Verify redirect was called with /upload
    expect(mockRedirect).toHaveBeenCalledWith("/upload");
  });
});
