package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/internal/handler"
)

func TestRootHandler(t *testing.T) {
	h := handler.NewRootHandler()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	expected := "Hello from Go backend!"
	if w.Body.String() != expected {
		t.Errorf("expected body %q, got %q", expected, w.Body.String())
	}
}
