package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"backend/internal/domain"
)

// MockStorage is an in-memory implementation of domain.Storage for testing.
type MockStorage struct { //nolint:govet // Field ordering prioritizes readability over memory alignment
	mu      sync.RWMutex
	objects map[string]*mockObject
}

type mockObject struct { //nolint:govet // Field ordering prioritizes readability over memory alignment
	data        []byte
	contentType string
}

// NewMockStorage creates a new in-memory storage for testing.
func NewMockStorage() *MockStorage {
	return &MockStorage{
		objects: make(map[string]*mockObject),
	}
}

// Upload stores data with the given key.
func (s *MockStorage) Upload(_ context.Context, key string, reader io.Reader, _ int64, contentType string) (*domain.StorageObject, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read data: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.objects[key] = &mockObject{
		data:        data,
		contentType: contentType,
	}

	return &domain.StorageObject{
		Key:         key,
		Size:        int64(len(data)),
		ContentType: contentType,
		ETag:        fmt.Sprintf("mock-etag-%d", len(data)),
	}, nil
}

// Download retrieves data by key.
func (s *MockStorage) Download(_ context.Context, key string) (io.ReadCloser, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	obj, ok := s.objects[key]
	if !ok {
		return nil, fmt.Errorf("object not found: %s", key)
	}

	return io.NopCloser(bytes.NewReader(obj.data)), nil
}

// Delete removes an object by key.
func (s *MockStorage) Delete(_ context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.objects, key)
	return nil
}

// GetPresignedURL generates a mock URL for testing.
func (s *MockStorage) GetPresignedURL(_ context.Context, key string, expiry time.Duration) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, ok := s.objects[key]; !ok {
		return "", fmt.Errorf("object not found: %s", key)
	}

	return fmt.Sprintf("http://mock-storage/%s?expiry=%s", key, expiry), nil
}

// Exists checks if an object exists.
func (s *MockStorage) Exists(_ context.Context, key string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.objects[key]
	return ok, nil
}

// Verify MockStorage implements domain.Storage.
var _ domain.Storage = (*MockStorage)(nil)
