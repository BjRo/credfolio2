---
paths:
  - "src/frontend/**"
---

# Visual Verification with Fixture Resume

> **Note:** For routine visual verification during development, use the `@qa` subagent (via Task tool) instead of running these commands manually. The QA subagent handles dev server management, browser automation, and error checking automatically. The commands below are reference documentation for the underlying `agent-browser` CLI.

A fixture resume is available at `fixtures/CV_TEMPLATE_0004.pdf` for testing the resume upload and profile display flow.

## How to upload via agent-browser

```bash
# 1. Start dev servers (ensure ports are free first)
pnpm dev &

# 2. Navigate to the upload page
agent-browser open http://localhost:3000/upload-resume

# 3. Upload the fixture resume using CSS selector for the hidden file input
agent-browser upload 'input[type="file"]' /workspace/fixtures/CV_TEMPLATE_0004.pdf

# 4. Wait for extraction to complete (redirects to profile page)
agent-browser wait --load networkidle
# If not auto-redirected, wait ~30s and re-snapshot:
sleep 30 && agent-browser snapshot -c

# 5. Verify the profile page renders correctly
agent-browser screenshot --full /path/to/screenshot.png
agent-browser errors  # Check for JS errors
```

## Key points

- The upload page's file input is hidden (opacity: 0). Use the CSS selector `'input[type="file"]'` with `agent-browser upload` — do **not** try to click it or use a ref.
- The profile page URL is `/profile/{resumeId}` — it takes a **resume ID** (not user ID). After upload, the page auto-redirects.
- Resume extraction takes ~15-30 seconds (LLM processing). Wait before checking the profile.
- The demo user ID is `00000000-0000-0000-0000-000000000001` (seeded automatically on server start).
- You can also create test data directly via GraphQL mutations (`createEducation`, `createExperience`) without uploading a resume.
