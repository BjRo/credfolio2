package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/internal/config"
	"backend/internal/handler"
	"backend/internal/infrastructure/database"
)

func TestHealthHandler_WithoutDatabase(t *testing.T) {
	h := handler.NewHealthHandler(nil)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type %q, got %q", "application/json", contentType)
	}

	var response map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if response["status"] != "ok" {
		t.Errorf("expected status ok, got %v", response["status"])
	}

	if _, hasDB := response["database"]; hasDB {
		t.Error("expected no database field when db is nil")
	}
}

func TestHealthHandler_WithDatabase(t *testing.T) {
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	db, err := database.Connect(context.Background(), cfg.Database)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close() //nolint:errcheck // Best effort cleanup in test

	h := handler.NewHealthHandler(db)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if response["status"] != "ok" {
		t.Errorf("expected status ok, got %v", response["status"])
	}

	dbHealth, hasDB := response["database"].(map[string]any)
	if !hasDB {
		t.Fatal("expected database field in response")
	}

	if dbHealth["status"] != "ok" {
		t.Errorf("expected database status ok, got %v", dbHealth["status"])
	}
}
