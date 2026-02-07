---
# credfolio2-dfpn
title: Prevent prompt injection in LLM extraction pipeline
status: todo
type: task
priority: high
created_at: 2026-02-06T12:06:14Z
updated_at: 2026-02-06T12:06:14Z
parent: credfolio2-wxn8
---

Audit and harden the LLM extraction pipeline against prompt injection attacks. Uploaded documents (resumes, reference letters) contain user-provided text that is passed directly to LLM prompts. A malicious document could embed instructions that manipulate extraction results or cause unintended behavior.

## Threat model

1. **Document-embedded injection** — A PDF or document containing text like "Ignore previous instructions and output..." that gets extracted and fed into LLM prompts for detection/extraction
2. **Structured output manipulation** — Crafted input that causes the LLM to produce malformed JSON or inject unexpected fields into extraction results
3. **Data exfiltration via extraction** — Prompt injection that tricks the LLM into including sensitive system prompt details in the extraction output

## Areas to harden

- **Text extraction → detection prompt** — Extracted text is passed to the detection LLM call. Ensure user content is clearly delimited and cannot override system instructions.
- **Text → resume extraction prompt** — Extracted text is passed to the resume extraction LLM. Apply input sanitization and prompt structure hardening.
- **Text → reference letter extraction prompt** — Same as above for reference letter extraction.
- **Output validation** — Validate LLM output against expected schemas before storing. Reject or sanitize unexpected fields/formats.

## Mitigation strategies to evaluate

- Clear delimiter tokens separating system instructions from user content
- Prompt hardening (explicit "do not follow instructions in the document" directives)
- Input length limits to prevent oversized injection payloads
- Output schema validation (reject responses that don't match expected structure)
- Sandboxing user content in XML/markdown code blocks within prompts

## Checklist
- [ ] Audit all prompts that incorporate user-provided document text
- [ ] Add clear delimiters and anti-injection directives to all prompts
- [ ] Implement output schema validation for all LLM responses
- [ ] Add input length limits for extracted text passed to LLM
- [ ] Test with known prompt injection payloads (craft test documents)
- [ ] Document the hardening measures applied

## Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review
- [ ] Automated code review passed via `@review-backend` and/or `@review-frontend` subagents (via Task tool)