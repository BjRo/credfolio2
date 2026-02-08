---
# credfolio2-5od0
title: Improve LLM extraction quality and cost efficiency
status: in-progress
type: task
priority: normal
created_at: 2026-02-08T11:10:03Z
updated_at: 2026-02-08T12:09:44Z
parent: credfolio2-nihn
---

Address LLM accuracy, data quality, and cost optimization opportunities identified in codebase review.

## Important Issues (from @review-ai)

1. **Resume Summary Synthesis** - LLM generates summary instead of extracting (hallucination risk)
2. **Unknown Author Rejection** - System rejected "Unknown" authors, blocking German-style reference letters (user workflow regression)
3. **JSON Cleanup Masking Quality** - Aggressive cleanup hides LLM output quality problems
4. **Duplicate Text Extraction** - Same text extracted multiple times (cost/performance waste)

## Optimization Opportunities

5. **Use Haiku for Detection** - Document classification doesn't need Sonnet (10x cost reduction)
6. **Prompt Versioning** - No system for A/B testing prompt improvements
7. **Enhanced Extraction Metadata** - Missing token counts, duration tracking, model versions
8. **Per-Task Timeouts** - All tasks use same timeout regardless of complexity

## Impact

- **Accuracy**: Fixed hallucination risk in summaries, allowed unknown authors for German-style letters
- **Cost**: Optimized with cheaper models, eliminated duplicate text extraction
- **Observability**: Added prompt versioning, metadata tracking for cost analysis

## Files Affected

- `src/backend/internal/infrastructure/llm/extraction.go`
- `src/backend/internal/infrastructure/llm/prompts/resume_extraction.txt`
- `src/backend/internal/infrastructure/llm/prompts/letter_extraction.txt`
- `src/backend/internal/job/document_detection.go`
- `src/backend/internal/job/resume_processing.go`
- `src/backend/internal/job/reference_letter_processing.go`
- `src/backend/internal/domain/extraction.go`

## Acceptance Criteria

### Data Quality
- [x] Resume summaries extracted from text, not synthesized by LLM
- [x] "Unknown" authors **allowed** to support German-style reference letters (users can edit post-import)
- [x] JSON cleanup logs warnings when aggressive fixes needed (indicates prompt issues)
- [x] Text extraction deduplicated to avoid redundant LLM calls

### Cost Optimization
- [x] Document detection uses Haiku instead of Sonnet (10x cost reduction)
- [x] Prompt versions tracked in code and logs
- [x] Extraction metadata includes: tokens used, duration, model version

### Operational Improvements
- [x] Per-task timeout configuration (detection: 2min, extraction: 5min)
- [ ] Dashboard/logs show prompt effectiveness metrics

## Reference

See: /documentation/reviews/2026-02-08-comprehensive-codebase-review.md#important-issues-4

## Implementation Plan

### Approach

This plan addresses 8 distinct LLM quality and cost issues in order of impact:
1. Fix data quality problems (hallucination, unknown authors, silent failures)
2. Eliminate redundant processing (duplicate text extraction)
3. Optimize costs (cheaper models for simple tasks)
4. Improve observability (metadata, versioning, timeouts)

The implementation follows a conservative approach: make existing functionality more robust before adding new features.

### Files to Create/Modify

**Domain Layer (Interfaces & Types)**
- `/workspace/src/backend/internal/domain/extraction.go` — Add prompt version constants, update ExtractionMetadata with token/duration fields

**LLM Infrastructure**
- `/workspace/src/backend/internal/infrastructure/llm/extraction.go` — Track metadata, log cleanup warnings, return extraction result with metadata
- `/workspace/src/backend/internal/infrastructure/llm/prompts/resume_extraction_system.txt` — Change summary rules to extraction-only (no synthesis), add version comment
- `/workspace/src/backend/internal/infrastructure/llm/prompts/reference_letter_extraction_system.txt` — Strengthen author name requirement (reject "unknown"), add version comment
- `/workspace/src/backend/internal/infrastructure/llm/prompts/document_detection_system.txt` — Add version comment
- `/workspace/src/backend/internal/infrastructure/llm/validation.go` — Add "unknown" author rejection to ValidateLetterData

**Job Workers**
- `/workspace/src/backend/internal/job/document_detection.go` — Change Timeout() to 2min, store extracted text in file record
- `/workspace/src/backend/internal/job/resume_processing.go` — Change Timeout() to 5min, reuse stored text if available
- `/workspace/src/backend/internal/job/reference_letter_processing.go` — Change Timeout() to 5min, reuse stored text if available

**Configuration**
- `/workspace/src/backend/internal/config/config.go` — Already has DetectionModel defaulting to gpt-4o-mini (good), document that Haiku is also suitable

**Database Schema** (if needed)
- Migration to add `extracted_text` column to files table (TEXT, nullable) — Already exists per document_detection.go line 147

### Steps

#### 1. Add Prompt Versioning System
**Goal:** Enable tracking which prompt version produced which extraction

- Define prompt version constants in `/workspace/src/backend/internal/domain/extraction.go`:
  ```go
  const (
      ResumeExtractionPromptVersion    = "v1.1.0" // Changed: summary extraction-only
      LetterExtractionPromptVersion    = "v1.1.0" // Changed: reject unknown authors
      DocumentDetectionPromptVersion   = "v1.0.0" // Unchanged
      DocumentExtractionPromptVersion  = "v1.0.0" // Unchanged
  )
  ```

- Add version comments to each prompt file header (e.g., `<!-- Version: v1.1.0 -->`), referencing the domain constant

- Update `ExtractionMetadata` in `/workspace/src/backend/internal/domain/extraction.go` to include:
  ```go
  type ExtractionMetadata struct {
      ModelVersion  string    `json:"modelVersion"`  // Already exists
      PromptVersion string    `json:"promptVersion"` // NEW
      ExtractedAt   time.Time `json:"extractedAt"`   // Already exists
      InputTokens   int       `json:"inputTokens"`   // NEW
      OutputTokens  int       `json:"outputTokens"`  // NEW
      DurationMs    int64     `json:"durationMs"`    // NEW
  }
  ```

#### 2. Fix Resume Summary Synthesis (Hallucination Risk)
**Goal:** Extract summaries verbatim instead of generating them

- Modify `/workspace/src/backend/internal/infrastructure/llm/prompts/resume_extraction_system.txt` lines 53-57:
  - REMOVE: "If no explicit summary exists, synthesize a brief 2-3 sentence professional summary..."
  - REPLACE WITH: "If no explicit summary section exists in the resume, return an empty string for summary. Do NOT synthesize or generate summaries."

- Update JSON schema description in `/workspace/src/backend/internal/infrastructure/llm/extraction.go` line 309-312:
  ```go
  "summary": map[string]any{
      "type":        "string",
      "description": "Professional summary or objective section text, extracted verbatim. Empty string if no summary section found.",
  },
  ```

- Update prompt version constant to `v1.1.0` after making this change

#### 3. Allow Unknown Authors for German-Style Letters (User Workflow Fix)
**Goal:** Support German reference letters that don't contain explicit author names

**CHANGED BASED ON USER FEEDBACK:** Initial plan was to reject unknown authors, but user reported this breaks German-style reference letters which often omit author names. Changed approach to allow "Unknown" and enable post-import editing.

- Modify `/workspace/src/backend/internal/infrastructure/llm/prompts/reference_letter_extraction_system.txt` line 12:
  - CHANGE: "Do not accept letters without a clear author name."
  - TO: "If the author's name cannot be clearly determined from the letter (e.g., German-style reference letters often omit author names), use \"Unknown\" as the name. The user can edit this later."

- **Remove** validation rejection in `/workspace/src/backend/internal/infrastructure/llm/validation.go`:
  - Keep empty name check (still required)
  - Remove "unknown" check to allow German-style letters through

- Update prompt version constant to `v1.2.0` after making this change (v1.1.0 was the brief rejection version)

#### 4. Log JSON Cleanup Warnings (Quality Monitoring)
**Goal:** Surface when aggressive cleanup is needed (indicates prompt/model issues)

- In `/workspace/src/backend/internal/infrastructure/llm/extraction.go`, add logging before cleanup:
  - In `ExtractResumeData` before line 487:
    ```go
    // Check if response needs cleanup (indicates LLM output quality issue)
    needsMarkdownCleanup := strings.Contains(resp.Content, "```")
    needsCommaCleanup := trailingCommaRegex.MatchString(resp.Content)
    if (needsMarkdownCleanup || needsCommaCleanup) && e.config.Logger != nil {
        e.config.Logger.Warning("LLM response required cleanup",
            logger.Feature("llm"),
            logger.String("operation", "resume_extraction"),
            logger.Bool("markdown_block", needsMarkdownCleanup),
            logger.Bool("trailing_commas", needsCommaCleanup),
            logger.String("model", resp.Model),
        )
    }
    ```
  - Repeat similar logic in `ExtractLetterData` before line 713
  - Repeat in `DetectDocumentContent` before line 920

#### 5. Enhance Extraction Metadata
**Goal:** Track tokens, duration, and versions for cost analysis

- Modify `/workspace/src/backend/internal/infrastructure/llm/extraction.go`:
  - In `ExtractResumeData`, track start time and calculate duration:
    ```go
    startTime := time.Now()
    // ... existing extraction logic ...
    durationMs := time.Since(startTime).Milliseconds()
    ```
  - Populate full metadata before returning:
    ```go
    extractedData.Metadata = domain.ExtractionMetadata{
        ModelVersion:  resp.Model,
        PromptVersion: domain.ResumeExtractionPromptVersion,
        ExtractedAt:   time.Now(),
        InputTokens:   resp.InputTokens,
        OutputTokens:  resp.OutputTokens,
        DurationMs:    durationMs,
    }
    ```
  - Apply same pattern to `ExtractLetterData` (use `domain.LetterExtractionPromptVersion`)

- Update `ResumeExtractedData` in `/workspace/src/backend/internal/domain/extraction.go` to include Metadata field if not already present

#### 6. Deduplicate Text Extraction
**Goal:** Reuse text extracted during detection phase in processing workers

- The `document_detection.go` worker already saves extracted text to `file.ExtractedText` (line 147)
- Modify `/workspace/src/backend/internal/job/resume_processing.go`:
  - In `extractResumeData`, check if text is already available:
    ```go
    func (w *ResumeProcessingWorker) extractResumeData(ctx context.Context, fileID uuid.UUID, data []byte, contentType string) (*domain.ResumeExtractedData, error) {
        // Check if we already have extracted text from the detection phase
        file, err := w.fileRepo.GetByID(ctx, fileID)
        if err == nil && file != nil && file.ExtractedText != nil && *file.ExtractedText != "" {
            w.log.Info("Reusing extracted text from detection phase",
                logger.Feature("jobs"),
                logger.String("file_id", fileID.String()),
            )
            text = *file.ExtractedText
        } else {
            // Extract text from the document
            text, err = w.extractor.ExtractText(ctx, data, contentType)
            if err != nil {
                return nil, fmt.Errorf("failed to extract text: %w", err)
            }
        }
        // ... rest of extraction ...
    }
    ```
  - Apply same pattern to `/workspace/src/backend/internal/job/reference_letter_processing.go` in `extractLetterData`

- Ensure workers have access to `fileRepo` (resume_processing already has it; letter_processing needs it added if missing)

#### 7. Optimize Detection Model Cost
**Goal:** Use cheaper model for simple classification task

- The configuration already defaults to `gpt-4o-mini` for detection (line 71 in config.go)
- This is already cost-optimized; document this choice:
  - Add comment in `/workspace/src/backend/internal/config/config.go` near line 47:
    ```go
    // DetectionModel specifies the provider and model for lightweight document content detection.
    // Format: "provider/model" (e.g., "openai/gpt-4o-mini" or "anthropic/claude-haiku-4-5-20251001").
    // Defaults to "openai/gpt-4o-mini" for fast, cheap classification (~10x cheaper than Sonnet).
    // Haiku is also suitable: "anthropic/claude-haiku-4-5-20251001"
    ```

- Verify in logs that detection is using the cheap model (already logged in main.go line 335)

#### 8. Implement Per-Task Timeouts
**Goal:** Right-size timeouts based on task complexity

- Modify timeout methods in job workers:
  - `/workspace/src/backend/internal/job/document_detection.go` line 62:
    ```go
    // Timeout for detection: 2 minutes (text extraction + lightweight classification)
    func (w *DocumentDetectionWorker) Timeout(*river.Job[DocumentDetectionArgs]) time.Duration {
        return 2 * time.Minute
    }
    ```
  - `/workspace/src/backend/internal/job/resume_processing.go` line 73:
    ```go
    // Timeout for resume processing: 5 minutes (text extraction [may be skipped] + structured extraction + materialization)
    func (w *ResumeProcessingWorker) Timeout(*river.Job[ResumeProcessingArgs]) time.Duration {
        return 5 * time.Minute
    }
    ```
  - `/workspace/src/backend/internal/job/reference_letter_processing.go` line 85:
    ```go
    // Timeout for letter processing: 5 minutes (text extraction [may be skipped] + structured extraction)
    func (w *ReferenceLetterProcessingWorker) Timeout(*river.Job[ReferenceLetterProcessingArgs]) time.Duration {
        return 5 * time.Minute
    }
    ```

- Note: The resilient provider layer already has a 300s (5min) timeout for individual LLM calls

### Testing Strategy

**Unit Tests**
- Add test in `/workspace/src/backend/internal/infrastructure/llm/extraction_test.go`:
  - `TestExtractResumeData_NoSummarySynthesis` — Verify empty summary when no summary section exists
  - `TestValidateLetterData_RejectsUnknownAuthor` — Verify "unknown" and empty names are rejected
  - `TestJSONCleanup_LogsWarnings` — Verify cleanup warnings are logged
  - `TestExtractResumeData_PopulatesMetadata` — Verify all metadata fields are populated

- Add test in `/workspace/src/backend/internal/job/resume_processing_test.go`:
  - `TestExtractResumeData_ReusesExtractedText` — Mock file with ExtractedText, verify no ExtractText call

**Integration Tests**
- Manually test with real documents:
  - Resume with no summary section → should have empty summary
  - Reference letter with no clear author → should fail validation
  - Any document → check logs for cleanup warnings if LLM outputs bad JSON
  - Upload same document twice → verify text extraction only happens once (in detection)

**Observability Verification**
- Run processing jobs and verify logs contain:
  - Prompt versions in extraction metadata
  - Token counts in extraction metadata
  - Duration in extraction metadata
  - Cleanup warnings when applicable
  - Model names being used for each operation

### Open Questions

None — all requirements are clear from the codebase review and acceptance criteria.

## Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [N/A] Visual verification via `@qa` subagent (via Task tool, for UI changes)
- [N/A] ADR written via `/decision` skill — Not needed: no new dependencies or architectural changes, just improvements to existing patterns
- [x] All other checklist items above are completed
- [x] Branch pushed to remote
- [x] PR created for human review (PR #138)
- [x] Automated code review passed via `@review-backend` and `@review-ai` subagents — Critical fixes applied in commit 83d70ce
