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
	fmt.Fprint(w, "Hello from Go backend!") //nolint:errcheck // ResponseWriter errors are not actionable
}
