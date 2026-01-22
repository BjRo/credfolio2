package storage_test

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"
	"time"

	"backend/internal/config"
	"backend/internal/infrastructure/storage"
)

// testMinIOStorage returns a configured MinIO storage for testing.
// It skips the test if MinIO is not available and registers cleanup.
func testMinIOStorage(t *testing.T) *storage.MinIOStorage {
	t.Helper()

	cfg := config.MinIOConfig{
		Endpoint:  getEnvOrDefault("MINIO_ENDPOINT", "localhost:9000"),
		AccessKey: getEnvOrDefault("MINIO_ROOT_USER", "minioadmin"),
		SecretKey: getEnvOrDefault("MINIO_ROOT_PASSWORD", "minioadmin"),
		Bucket:    "test-bucket-" + time.Now().Format("20060102150405"),
		UseSSL:    false,
	}

	s, err := storage.NewMinIOStorage(cfg)
	if err != nil {
		t.Skipf("MinIO not available: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.EnsureBucket(ctx); err != nil {
		t.Skipf("MinIO not available: %v", err)
	}

	return s
}

func getEnvOrDefault(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

func TestMinIOStorage_Upload(t *testing.T) {
	s := testMinIOStorage(t)
	ctx := context.Background()

	content := []byte("test file content")
	key := "test-upload.txt"

	t.Cleanup(func() {
		if err := s.Delete(context.Background(), key); err != nil {
			t.Logf("cleanup failed: %v", err)
		}
	})

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
}

func TestMinIOStorage_Download(t *testing.T) {
	s := testMinIOStorage(t)
	ctx := context.Background()

	content := []byte("test download content")
	key := "test-download.txt"

	t.Cleanup(func() {
		if err := s.Delete(context.Background(), key); err != nil {
			t.Logf("cleanup failed: %v", err)
		}
	})

	// Upload first
	_, err := s.Upload(ctx, key, bytes.NewReader(content), int64(len(content)), "text/plain")
	if err != nil {
		t.Fatalf("Upload failed: %v", err)
	}

	// Download
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
		t.Fatalf("failed to read downloaded content: %v", err)
	}

	if !bytes.Equal(downloaded, content) {
		t.Errorf("downloaded content doesn't match: got %q, want %q", downloaded, content)
	}
}

func TestMinIOStorage_Download_NotFound(t *testing.T) {
	s := testMinIOStorage(t)
	ctx := context.Background()

	_, err := s.Download(ctx, "nonexistent-key")
	if err == nil {
		t.Error("expected error for nonexistent key, got nil")
	}
}

func TestMinIOStorage_Exists(t *testing.T) {
	s := testMinIOStorage(t)
	ctx := context.Background()

	key := "test-exists.txt"
	content := []byte("exists test")

	t.Cleanup(func() {
		if err := s.Delete(context.Background(), key); err != nil {
			t.Logf("cleanup failed: %v", err)
		}
	})

	// Check nonexistent
	exists, err := s.Exists(ctx, key)
	if err != nil {
		t.Fatalf("Exists failed: %v", err)
	}
	if exists {
		t.Error("expected false for nonexistent key")
	}

	// Upload
	_, err = s.Upload(ctx, key, bytes.NewReader(content), int64(len(content)), "text/plain")
	if err != nil {
		t.Fatalf("Upload failed: %v", err)
	}

	// Check exists
	exists, err = s.Exists(ctx, key)
	if err != nil {
		t.Fatalf("Exists failed: %v", err)
	}
	if !exists {
		t.Error("expected true for existing key")
	}
}

func TestMinIOStorage_Delete(t *testing.T) {
	s := testMinIOStorage(t)
	ctx := context.Background()

	key := "test-delete.txt"
	content := []byte("delete test")

	// Upload
	_, err := s.Upload(ctx, key, bytes.NewReader(content), int64(len(content)), "text/plain")
	if err != nil {
		t.Fatalf("Upload failed: %v", err)
	}

	// Delete
	err = s.Delete(ctx, key)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify deleted
	exists, err := s.Exists(ctx, key)
	if err != nil {
		t.Fatalf("Exists failed: %v", err)
	}
	if exists {
		t.Error("expected file to be deleted")
	}
}

func TestMinIOStorage_GetPresignedURL(t *testing.T) {
	s := testMinIOStorage(t)
	ctx := context.Background()

	key := "test-presigned.txt"
	content := []byte("presigned test")

	t.Cleanup(func() {
		if err := s.Delete(context.Background(), key); err != nil {
			t.Logf("cleanup failed: %v", err)
		}
	})

	// Upload
	_, err := s.Upload(ctx, key, bytes.NewReader(content), int64(len(content)), "text/plain")
	if err != nil {
		t.Fatalf("Upload failed: %v", err)
	}

	// Get presigned URL
	url, err := s.GetPresignedURL(ctx, key, time.Hour)
	if err != nil {
		t.Fatalf("GetPresignedURL failed: %v", err)
	}

	if url == "" {
		t.Error("expected non-empty URL")
	}
}
