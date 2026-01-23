package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"backend/internal/domain"
	"backend/internal/infrastructure/llm"
	"backend/internal/logger"
)

const maxExtractFileSize = 20 << 20 // 20MB

// supportedMediaTypes maps Content-Type to domain.ImageMediaType.
var supportedMediaTypes = map[string]domain.ImageMediaType{
	"image/jpeg":      domain.ImageMediaTypeJPEG,
	"image/png":       domain.ImageMediaTypePNG,
	"image/gif":       domain.ImageMediaTypeGIF,
	"image/webp":      domain.ImageMediaTypeWebP,
	"application/pdf": domain.ImageMediaTypePDF,
}

// ExtractHandler handles document text extraction requests.
type ExtractHandler struct {
	extractor *llm.DocumentExtractor
	log       logger.Logger
}

// NewExtractHandler creates a new ExtractHandler.
func NewExtractHandler(extractor *llm.DocumentExtractor, log logger.Logger) *ExtractHandler {
	return &ExtractHandler{
		extractor: extractor,
		log:       log,
	}
}

// extractResponse represents a successful extraction response.
type extractResponse struct {
	Text         string `json:"text"`
	InputTokens  int    `json:"inputTokens"`
	OutputTokens int    `json:"outputTokens"`
}

// extractErrorResponse represents an error response.
type extractErrorResponse struct {
	Error string `json:"error"`
}

// ServeHTTP implements the http.Handler interface.
func (h *ExtractHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Only accept POST
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		h.writeError(w, "Method not allowed")
		return
	}

	// Limit request body size
	r.Body = http.MaxBytesReader(w, r.Body, maxExtractFileSize)

	// Parse multipart form
	if err := r.ParseMultipartForm(maxExtractFileSize); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.writeError(w, "Failed to parse form: "+err.Error())
		return
	}

	// Get file from form
	file, header, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.writeError(w, "Missing file in request")
		return
	}
	defer file.Close() //nolint:errcheck // Best effort cleanup

	// Validate content type
	contentType := header.Header.Get("Content-Type")
	mediaType, ok := supportedMediaTypes[contentType]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		h.writeError(w, "Unsupported file type: "+contentType+". Supported: image/jpeg, image/png, image/gif, image/webp, application/pdf")
		return
	}

	// Read file data
	data, err := io.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.writeError(w, "Failed to read file: "+err.Error())
		return
	}

	h.log.Debug("Extracting text from document",
		logger.Feature("extraction"),
		logger.String("content_type", contentType),
		logger.Int("size_bytes", len(data)),
	)

	// Extract text
	result, err := h.extractor.ExtractText(r.Context(), llm.ExtractionRequest{
		Document:  data,
		MediaType: mediaType,
	})
	if err != nil {
		h.log.Error("Extraction failed",
			logger.Feature("extraction"),
			logger.Err(err),
		)
		w.WriteHeader(http.StatusInternalServerError)
		h.writeError(w, "Extraction failed: "+err.Error())
		return
	}

	h.log.Info("Extraction completed",
		logger.Feature("extraction"),
		logger.Int("input_tokens", result.InputTokens),
		logger.Int("output_tokens", result.OutputTokens),
	)

	// Return success response
	json.NewEncoder(w).Encode(extractResponse{ //nolint:errcheck,gosec // ResponseWriter errors are not actionable
		Text:         result.Text,
		InputTokens:  result.InputTokens,
		OutputTokens: result.OutputTokens,
	})
}

func (h *ExtractHandler) writeError(w http.ResponseWriter, msg string) {
	json.NewEncoder(w).Encode(extractErrorResponse{Error: msg}) //nolint:errcheck,gosec // ResponseWriter errors are not actionable
}

// ExtractUnavailableHandler returns 503 when extraction is not configured.
type ExtractUnavailableHandler struct{}

// NewExtractUnavailableHandler creates a new ExtractUnavailableHandler.
func NewExtractUnavailableHandler() *ExtractUnavailableHandler {
	return &ExtractUnavailableHandler{}
}

// ServeHTTP implements the http.Handler interface.
func (h *ExtractUnavailableHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusServiceUnavailable)
	json.NewEncoder(w).Encode(extractErrorResponse{ //nolint:errcheck,gosec // ResponseWriter errors are not actionable
		Error: "Document extraction is not available. ANTHROPIC_API_KEY not configured.",
	})
}
