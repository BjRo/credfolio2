import { describe, expect, it } from "vitest";
import { createUrqlClient, GRAPHQL_ENDPOINT } from "./client";

describe("URQL Client", () => {
  describe("createUrqlClient", () => {
    it("should create a client instance", () => {
      const client = createUrqlClient();
      expect(client).toBeDefined();
      expect(typeof client.query).toBe("function");
      expect(typeof client.mutation).toBe("function");
    });

    it("should create a client with exchanges configured", () => {
      const client = createUrqlClient();
      // Client should have query/mutation methods from exchanges
      expect(client.query).toBeDefined();
    });
  });

  describe("GRAPHQL_ENDPOINT", () => {
    it("should point to backend GraphQL endpoint", () => {
      expect(GRAPHQL_ENDPOINT).toBe("http://localhost:8080/graphql");
    });
  });
});
