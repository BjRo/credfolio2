import { redirect } from "next/navigation";
import { GetProfileDocument } from "@/graphql/generated/graphql";
import { createUrqlClient, GRAPHQL_ENDPOINT } from "@/lib/urql/client";

const DEMO_USER_ID = "00000000-0000-0000-0000-000000000001";

export default async function Home() {
  // Server-side GraphQL query
  const client = createUrqlClient(GRAPHQL_ENDPOINT);
  const result = await client.query(GetProfileDocument, { userId: DEMO_USER_ID }).toPromise();

  // Handle errors gracefully (backend down, network issues, etc.)
  if (result.error) {
    redirect("/upload");
  }

  const profile = result.data?.profileByUserId;

  // Server-side redirect (no loading spinner, instant navigation)
  if (profile) {
    redirect(`/profile/${profile.id}`);
  }
  redirect("/upload");
}
