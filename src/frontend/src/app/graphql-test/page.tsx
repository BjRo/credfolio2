"use client";

import { useQuery } from "urql";
import { graphql } from "@/graphql/generated";

// Use a valid UUID format for the test query
const TEST_USER_ID = "00000000-0000-0000-0000-000000000000";

const TestQuery = graphql(`
  query TestConnection($userId: ID!) {
    referenceLetters(userId: $userId) {
      id
      title
      status
      authorName
      createdAt
    }
  }
`);

export default function GraphQLTestPage() {
  const [result] = useQuery({
    query: TestQuery,
    variables: { userId: TEST_USER_ID },
  });

  const { data, fetching, error } = result;

  return (
    <div style={{ padding: "2rem", fontFamily: "system-ui" }}>
      <h1>GraphQL Connection Test</h1>

      <section style={{ marginTop: "1rem" }}>
        <h2>Status</h2>
        <ul>
          <li>
            <strong>Fetching:</strong> {fetching ? "Yes" : "No"}
          </li>
          <li>
            <strong>Error:</strong> {error ? error.message : "None"}
          </li>
          <li>
            <strong>Data received:</strong> {data ? "Yes" : "No"}
          </li>
        </ul>
      </section>

      {error && (
        <section style={{ marginTop: "1rem", color: "red" }}>
          <h2>Error Details</h2>
          <pre style={{ background: "#fee", padding: "1rem", overflow: "auto" }}>
            {JSON.stringify(error, null, 2)}
          </pre>
        </section>
      )}

      {data && (
        <section style={{ marginTop: "1rem", color: "green" }}>
          <h2>✓ Connection Successful!</h2>
          <p>Reference Letters found: {data.referenceLetters.length}</p>
          {data.referenceLetters.length === 0 ? (
            <p style={{ color: "#666" }}>(Empty result is expected - no data for test user UUID)</p>
          ) : (
            <ul>
              {data.referenceLetters.map((letter) => (
                <li key={letter.id}>
                  <strong>{letter.title || "Untitled"}</strong> - {letter.status}
                  {letter.authorName && ` (from ${letter.authorName})`}
                </li>
              ))}
            </ul>
          )}
        </section>
      )}

      <section style={{ marginTop: "2rem", color: "#666" }}>
        <h2>Instructions</h2>
        <ol>
          <li>
            Start the backend: <code>cd src/backend && pnpm dev</code>
          </li>
          <li>
            Ensure PostgreSQL is running: <code>docker-compose up -d</code>
          </li>
          <li>If you see &quot;Connection Successful&quot; above, everything works!</li>
          <li>If you see a network error, check that the backend is running on port 8080</li>
        </ol>
      </section>
    </div>
  );
}
