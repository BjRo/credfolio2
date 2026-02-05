"use client";

import { useCallback, useState } from "react";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { GRAPHQL_ENDPOINT } from "@/lib/urql/client";
import type { FeedbackFormProps } from "./types";

const FEEDBACK_CATEGORIES = [
  "Missing information",
  "Incorrect data",
  "Wrong person",
  "Other",
] as const;

const REPORT_FEEDBACK_MUTATION = `
  mutation ReportDocumentFeedback($userId: ID!, $input: DocumentFeedbackInput!) {
    reportDocumentFeedback(userId: $userId, input: $input) { success }
  }
`;

export function FeedbackForm({ userId, fileId, onSubmitted }: FeedbackFormProps) {
  const [category, setCategory] = useState<string | null>(null);
  const [description, setDescription] = useState("");
  const [submitted, setSubmitted] = useState(false);

  const handleSubmit = useCallback(() => {
    const message = description.trim() ? `${category}: ${description.trim()}` : (category ?? "");

    fetch(GRAPHQL_ENDPOINT, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        query: REPORT_FEEDBACK_MUTATION,
        variables: {
          userId,
          input: { fileId, feedbackType: "EXTRACTION_QUALITY", message },
        },
      }),
    }).catch(() => {});

    setSubmitted(true);
    onSubmitted?.();
  }, [userId, fileId, category, description, onSubmitted]);

  if (submitted) {
    return <p className="text-sm text-muted-foreground">Thank you for your feedback.</p>;
  }

  return (
    <div className="space-y-3">
      <fieldset className="space-y-2">
        <legend className="text-sm font-medium">What went wrong?</legend>
        {FEEDBACK_CATEGORIES.map((cat) => (
          <label key={cat} className="flex items-center gap-3 cursor-pointer">
            <input
              type="radio"
              name="feedback-category"
              value={cat}
              checked={category === cat}
              onChange={() => setCategory(cat)}
              className="accent-primary"
            />
            <span className="text-sm">{cat}</span>
          </label>
        ))}
      </fieldset>
      <Textarea
        placeholder="Tell us more..."
        value={description}
        onChange={(e) => setDescription(e.target.value)}
        rows={3}
      />
      <Button onClick={handleSubmit} disabled={!category} size="sm" variant="outline">
        Send feedback
      </Button>
    </div>
  );
}
