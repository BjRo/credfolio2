package handler

import (
	"fmt"
	"net/http"
)

// RootHandler handles requests to the root path.
type RootHandler struct{}

// NewRootHandler creates a new RootHandler.
func NewRootHandler() *RootHandler {
	return &RootHandler{}
}

// ServeHTTP implements the http.Handler interface.
func (h *RootHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintf(w, "Hello from Go backend!")
}
