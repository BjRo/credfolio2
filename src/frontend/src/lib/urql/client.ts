import { Client, cacheExchange, fetchExchange } from "@urql/core";

// For SSR, we need an absolute URL since the server doesn't know the host
// For client-side, we use the proxied endpoint to avoid port forwarding issues
const isServer = typeof window === "undefined";
export const GRAPHQL_ENDPOINT =
  process.env.NEXT_PUBLIC_GRAPHQL_URL ||
  (isServer ? "http://localhost:8080/graphql" : "/api/graphql");

export function createUrqlClient(url: string = GRAPHQL_ENDPOINT): Client {
  return new Client({
    url,
    exchanges: [cacheExchange, fetchExchange],
    // Disable GET method to avoid URL length issues
    preferGetMethod: false,
    // Don't set Content-Type header - urql will set it automatically
    // (application/json for regular requests, multipart/form-data for file uploads)
  });
}
