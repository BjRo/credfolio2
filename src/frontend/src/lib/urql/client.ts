import { Client, cacheExchange, fetchExchange } from "@urql/core";

export const GRAPHQL_ENDPOINT =
  process.env.NEXT_PUBLIC_GRAPHQL_URL || "http://localhost:8080/graphql";

export function createUrqlClient(url: string = GRAPHQL_ENDPOINT): Client {
  return new Client({
    url,
    exchanges: [cacheExchange, fetchExchange],
  });
}
