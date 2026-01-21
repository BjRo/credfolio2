package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/uptrace/bun"
)

// HealthHandler handles health check requests.
type HealthHandler struct {
	db *bun.DB
}

// NewHealthHandler creates a new HealthHandler.
// If db is nil, the handler will not check database connectivity.
func NewHealthHandler(db *bun.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// healthResponse represents the health check response.
type healthResponse struct { //nolint:govet // Field ordering matches JSON convention (status first)
	Status   string           `json:"status"`
	Database *componentHealth `json:"database,omitempty"`
}

// componentHealth represents the health status of a component.
type componentHealth struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

// ServeHTTP implements the http.Handler interface.
func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := healthResponse{
		Status: "ok",
	}

	// Check database if configured
	if h.db != nil {
		dbHealth := h.checkDatabase(r.Context())
		response.Database = dbHealth
		if dbHealth.Status != "ok" {
			response.Status = "degraded"
		}
	}

	json.NewEncoder(w).Encode(response) //nolint:errcheck,gosec // ResponseWriter errors are not actionable
}

// checkDatabase verifies database connectivity.
func (h *HealthHandler) checkDatabase(ctx context.Context) *componentHealth {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := h.db.PingContext(ctx); err != nil {
		return &componentHealth{
			Status: "unhealthy",
			Error:  err.Error(),
		}
	}

	return &componentHealth{
		Status: "ok",
	}
}
