// Mock for @urql/next used in tests
import type { Client, Exchange } from "@urql/core";
import { cacheExchange, fetchExchange, Client as UrqlClient } from "@urql/core";
import { Provider } from "urql";

interface UrqlProviderProps {
  children: React.ReactNode;
  client: Client;
  ssr: unknown;
}

export function UrqlProvider({ children, client }: UrqlProviderProps) {
  return <Provider value={client}>{children}</Provider>;
}

interface SsrExchangeOptions {
  isClient?: boolean;
}

export function ssrExchange(_options?: SsrExchangeOptions): Exchange {
  // Return a pass-through exchange for tests
  return ({ forward }) => forward;
}

interface CreateClientOptions {
  url: string;
  exchanges?: Exchange[];
  suspense?: boolean;
}

export function createClient(options: CreateClientOptions): Client {
  return new UrqlClient({
    url: options.url,
    exchanges: [cacheExchange, fetchExchange],
  });
}
