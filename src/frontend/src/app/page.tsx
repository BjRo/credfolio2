"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { useQuery } from "urql";
import { GetProfileDocument } from "@/graphql/generated/graphql";

const DEMO_USER_ID = "00000000-0000-0000-0000-000000000001";

export default function Home() {
  const router = useRouter();
  const [result] = useQuery({
    query: GetProfileDocument,
    variables: { userId: DEMO_USER_ID },
  });

  const { fetching, data } = result;

  const profile = data?.profileByUserId;

  useEffect(() => {
    if (fetching) return;
    if (profile) {
      router.push(`/profile/${profile.id}`);
    } else {
      router.push("/upload");
    }
  }, [fetching, profile, router]);

  if (fetching) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-muted/50">
        <output className="animate-spin h-8 w-8 border-4 border-primary border-t-transparent rounded-full block" />
      </div>
    );
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-muted/50">
      <p className="text-muted-foreground">Redirecting...</p>
    </div>
  );
}
