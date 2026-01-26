"use client";

import { cacheExchange, fetchExchange } from "@urql/core";
import { createClient, ssrExchange, UrqlProvider as UrqlNextProvider } from "@urql/next";
import { useMemo } from "react";
import { GRAPHQL_ENDPOINT } from "./client";

interface UrqlProviderProps {
  children: React.ReactNode;
}

export function UrqlProvider({ children }: UrqlProviderProps) {
  const [client, ssr] = useMemo(() => {
    const ssr = ssrExchange({ isClient: typeof window !== "undefined" });
    const client = createClient({
      url: GRAPHQL_ENDPOINT,
      exchanges: [cacheExchange, ssr, fetchExchange],
      suspense: true,
      // Disable GET method to avoid URL length issues
      preferGetMethod: false,
      fetchOptions: {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
      },
    });
    return [client, ssr];
  }, []);

  return (
    <UrqlNextProvider client={client} ssr={ssr}>
      {children}
    </UrqlNextProvider>
  );
}
