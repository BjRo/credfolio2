package job

import (
	"context"
	"testing"
	"time"

	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"

	"backend/internal/domain"
)

// resumeMockExtractor implements domain.DocumentExtractor for span testing.
type resumeMockExtractor struct {
	extractTextResult  string
	extractResumeData  *domain.ResumeExtractedData
	extractTextError   error
	extractResumeError error
}

func (e *resumeMockExtractor) ExtractText(_ context.Context, _ []byte, _ string) (string, error) {
	return e.extractTextResult, e.extractTextError
}

func (e *resumeMockExtractor) ExtractResumeData(_ context.Context, _ string) (*domain.ResumeExtractedData, error) {
	return e.extractResumeData, e.extractResumeError
}

func (e *resumeMockExtractor) ExtractLetterData(_ context.Context, _ string, _ []domain.ProfileSkillContext) (*domain.ExtractedLetterData, error) {
	return nil, nil
}

func (e *resumeMockExtractor) DetectDocumentContent(_ context.Context, _ string) (*domain.DocumentDetectionResult, error) {
	return nil, nil
}

func setupTestTracingForJobs(t *testing.T) *tracetest.InMemoryExporter {
	t.Helper()
	exporter := tracetest.NewInMemoryExporter()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))
	prev := otel.GetTracerProvider()
	otel.SetTracerProvider(tp)
	t.Cleanup(func() {
		otel.SetTracerProvider(prev)
		_ = tp.Shutdown(context.Background()) //nolint:errcheck // best effort cleanup in test
	})
	return exporter
}

func TestExtractResumeData_CreatesParentSpan(t *testing.T) {
	exporter := setupTestTracingForJobs(t)

	extractor := &resumeMockExtractor{
		extractTextResult: "Some resume text",
		extractResumeData: &domain.ResumeExtractedData{
			Name:       "Test User",
			Skills:     []string{},
			Experience: []domain.WorkExperience{},
			Education:  []domain.Education{},
		},
	}

	worker := &ResumeProcessingWorker{
		extractor: extractor,
	}

	_, err := worker.extractResumeData(context.Background(), []byte("pdf data"), "application/pdf")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	spans := exporter.GetSpans()
	var found bool
	for _, s := range spans {
		if s.Name == "resume_extraction" {
			found = true
			break
		}
	}
	if !found {
		names := make([]string, len(spans))
		for i, s := range spans {
			names[i] = s.Name
		}
		t.Errorf("expected span named 'resume_extraction', got spans: %v", names)
	}
}

func TestResumeProcessingWorker_Timeout_TenMinutes(t *testing.T) {
	worker := &ResumeProcessingWorker{}
	timeout := worker.Timeout(nil)
	if timeout != 10*time.Minute {
		t.Errorf("Timeout = %v, want 10m (safety net)", timeout)
	}
}

func TestResumeProcessingArgs_InsertOpts_MaxAttempts(t *testing.T) {
	args := ResumeProcessingArgs{}
	opts := args.InsertOpts()
	if opts.MaxAttempts != 2 {
		t.Errorf("MaxAttempts = %d, want 2", opts.MaxAttempts)
	}
}
