package storage_test

import (
	"bytes"
	"context"
	"io"
	"testing"
	"time"

	"backend/internal/infrastructure/storage"
)

const testKey = "test-key"

func TestMockStorage_Upload(t *testing.T) {
	s := storage.NewMockStorage()
	ctx := context.Background()

	content := []byte("test content")
	key := testKey

	obj, err := s.Upload(ctx, key, bytes.NewReader(content), int64(len(content)), "text/plain")
	if err != nil {
		t.Fatalf("Upload failed: %v", err)
	}

	if obj.Key != key {
		t.Errorf("expected key %q, got %q", key, obj.Key)
	}
	if obj.Size != int64(len(content)) {
		t.Errorf("expected size %d, got %d", len(content), obj.Size)
	}
	if obj.ContentType != "text/plain" {
		t.Errorf("expected content type %q, got %q", "text/plain", obj.ContentType)
	}
}

func TestMockStorage_Download(t *testing.T) {
	s := storage.NewMockStorage()
	ctx := context.Background()

	content := []byte("test content")
	key := testKey

	_, err := s.Upload(ctx, key, bytes.NewReader(content), int64(len(content)), "text/plain")
	if err != nil {
		t.Fatalf("Upload failed: %v", err)
	}

	reader, err := s.Download(ctx, key)
	if err != nil {
		t.Fatalf("Download failed: %v", err)
	}
	defer func() {
		if closeErr := reader.Close(); closeErr != nil {
			t.Logf("failed to close reader: %v", closeErr)
		}
	}()

	downloaded, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("failed to read: %v", err)
	}

	if !bytes.Equal(downloaded, content) {
		t.Errorf("content mismatch: got %q, want %q", downloaded, content)
	}
}

func TestMockStorage_Download_NotFound(t *testing.T) {
	s := storage.NewMockStorage()
	ctx := context.Background()

	_, err := s.Download(ctx, "nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent key")
	}
}

func TestMockStorage_Exists(t *testing.T) {
	s := storage.NewMockStorage()
	ctx := context.Background()

	key := testKey

	// Should not exist
	exists, err := s.Exists(ctx, key)
	if err != nil {
		t.Fatalf("Exists failed: %v", err)
	}
	if exists {
		t.Error("expected false for nonexistent key")
	}

	// Upload
	_, err = s.Upload(ctx, key, bytes.NewReader([]byte("data")), 4, "text/plain")
	if err != nil {
		t.Fatalf("Upload failed: %v", err)
	}

	// Should exist now
	exists, err = s.Exists(ctx, key)
	if err != nil {
		t.Fatalf("Exists failed: %v", err)
	}
	if !exists {
		t.Error("expected true for existing key")
	}
}

func TestMockStorage_Delete(t *testing.T) {
	s := storage.NewMockStorage()
	ctx := context.Background()

	key := testKey

	_, err := s.Upload(ctx, key, bytes.NewReader([]byte("data")), 4, "text/plain")
	if err != nil {
		t.Fatalf("Upload failed: %v", err)
	}

	err = s.Delete(ctx, key)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	exists, err := s.Exists(ctx, key)
	if err != nil {
		t.Fatalf("Exists failed: %v", err)
	}
	if exists {
		t.Error("expected key to be deleted")
	}
}

func TestMockStorage_GetPresignedURL(t *testing.T) {
	s := storage.NewMockStorage()
	ctx := context.Background()

	key := testKey

	_, err := s.Upload(ctx, key, bytes.NewReader([]byte("data")), 4, "text/plain")
	if err != nil {
		t.Fatalf("Upload failed: %v", err)
	}

	url, err := s.GetPresignedURL(ctx, key, time.Hour)
	if err != nil {
		t.Fatalf("GetPresignedURL failed: %v", err)
	}

	if url == "" {
		t.Error("expected non-empty URL")
	}
}

func TestMockStorage_GetPresignedURL_NotFound(t *testing.T) {
	s := storage.NewMockStorage()
	ctx := context.Background()

	_, err := s.GetPresignedURL(ctx, "nonexistent", time.Hour)
	if err == nil {
		t.Error("expected error for nonexistent key")
	}
}
