package config

import (
	"testing"
	"time"
)

// Default values matching docker-compose.yml
const (
	defaultDBHost     = "localhost"
	defaultDBPort     = 5432
	defaultDBUser     = "credfolio"
	defaultDBPassword = "credfolio_dev" //nolint:gosec // test credentials
	defaultDBName     = "credfolio"
	defaultDBSSLMode  = "disable"

	defaultMinIOEndpoint  = "localhost:9000"
	defaultMinIOAccessKey = "minioadmin"
	defaultMinIOSecretKey = "minioadmin" //nolint:gosec // test credentials
	defaultMinIOBucket    = "credfolio"

	defaultServerPort = 8080
)

func TestLoad_DatabaseDefaults(t *testing.T) {
	clearEnv(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Host", cfg.Database.Host, defaultDBHost},
		{"Port", cfg.Database.Port, defaultDBPort},
		{"User", cfg.Database.User, defaultDBUser},
		{"Password", cfg.Database.Password, defaultDBPassword},
		{"Name", cfg.Database.Name, defaultDBName},
		{"SSLMode", cfg.Database.SSLMode, defaultDBSSLMode},
	}

	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("Database.%s = %v, want %v", tt.name, tt.got, tt.want)
		}
	}
}

func TestLoad_MinIODefaults(t *testing.T) {
	clearEnv(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Endpoint", cfg.MinIO.Endpoint, defaultMinIOEndpoint},
		{"AccessKey", cfg.MinIO.AccessKey, defaultMinIOAccessKey},
		{"SecretKey", cfg.MinIO.SecretKey, defaultMinIOSecretKey},
		{"UseSSL", cfg.MinIO.UseSSL, false},
		{"Bucket", cfg.MinIO.Bucket, defaultMinIOBucket},
	}

	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("MinIO.%s = %v, want %v", tt.name, tt.got, tt.want)
		}
	}
}

func TestLoad_ServerDefaults(t *testing.T) {
	clearEnv(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Port", cfg.Server.Port, defaultServerPort},
		{"ReadTimeout", cfg.Server.ReadTimeout, 15 * time.Second},
		{"WriteTimeout", cfg.Server.WriteTimeout, 15 * time.Second},
		{"IdleTimeout", cfg.Server.IdleTimeout, 60 * time.Second},
	}

	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("Server.%s = %v, want %v", tt.name, tt.got, tt.want)
		}
	}
}

func TestLoad_DatabaseOverrides(t *testing.T) {
	clearEnv(t)

	t.Setenv("POSTGRES_HOST", "db.example.com")
	t.Setenv("POSTGRES_PORT", "5433")
	t.Setenv("POSTGRES_USER", "myuser")
	t.Setenv("POSTGRES_PASSWORD", "mypassword")
	t.Setenv("POSTGRES_DB", "mydb")
	t.Setenv("POSTGRES_SSLMODE", "require")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Host", cfg.Database.Host, "db.example.com"},
		{"Port", cfg.Database.Port, 5433},
		{"User", cfg.Database.User, "myuser"},
		{"Password", cfg.Database.Password, "mypassword"},
		{"Name", cfg.Database.Name, "mydb"},
		{"SSLMode", cfg.Database.SSLMode, "require"},
	}

	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("Database.%s = %v, want %v", tt.name, tt.got, tt.want)
		}
	}
}

func TestLoad_MinIOOverrides(t *testing.T) {
	clearEnv(t)

	t.Setenv("MINIO_ENDPOINT", "minio.example.com:9000")
	t.Setenv("MINIO_ROOT_USER", "minio_user")
	t.Setenv("MINIO_ROOT_PASSWORD", "minio_pass")
	t.Setenv("MINIO_USE_SSL", "true")
	t.Setenv("MINIO_BUCKET", "mybucket")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"Endpoint", cfg.MinIO.Endpoint, "minio.example.com:9000"},
		{"AccessKey", cfg.MinIO.AccessKey, "minio_user"},
		{"SecretKey", cfg.MinIO.SecretKey, "minio_pass"},
		{"UseSSL", cfg.MinIO.UseSSL, true},
		{"Bucket", cfg.MinIO.Bucket, "mybucket"},
	}

	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("MinIO.%s = %v, want %v", tt.name, tt.got, tt.want)
		}
	}
}

func TestLoad_ServerOverrides(t *testing.T) {
	clearEnv(t)

	t.Setenv("SERVER_PORT", "3000")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Server.Port != 3000 {
		t.Errorf("Server.Port = %d, want %d", cfg.Server.Port, 3000)
	}
}

func TestLoad_DatabaseURL_Computed(t *testing.T) {
	clearEnv(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// #nosec G101 - test credentials only
	want := "postgres://credfolio:credfolio_dev@localhost:5432/credfolio?sslmode=disable"
	if cfg.Database.URL() != want {
		t.Errorf("Database.URL() = %q, want %q", cfg.Database.URL(), want)
	}
}

func TestLoad_DatabaseURL_WithOverrides(t *testing.T) {
	clearEnv(t)

	t.Setenv("POSTGRES_HOST", "db.example.com")
	t.Setenv("POSTGRES_PORT", "5433")
	t.Setenv("POSTGRES_USER", "myuser")
	t.Setenv("POSTGRES_PASSWORD", "mypassword")
	t.Setenv("POSTGRES_DB", "mydb")
	t.Setenv("POSTGRES_SSLMODE", "require")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	want := "postgres://myuser:mypassword@db.example.com:5433/mydb?sslmode=require"
	if cfg.Database.URL() != want {
		t.Errorf("Database.URL() = %q, want %q", cfg.Database.URL(), want)
	}
}

func TestLoad_DatabaseURL_FromEnvVar(t *testing.T) {
	clearEnv(t)

	// DATABASE_URL should take precedence over individual settings
	customURL := "postgres://override:pass@custom.host:5555/overridedb?sslmode=verify-full"
	t.Setenv("DATABASE_URL", customURL)
	t.Setenv("POSTGRES_HOST", "ignored.host")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Database.URL() != customURL {
		t.Errorf("Database.URL() = %q, want %q", cfg.Database.URL(), customURL)
	}
}

func TestLoad_InvalidPort(t *testing.T) {
	clearEnv(t)

	t.Setenv("POSTGRES_PORT", "not-a-number")

	_, err := Load()
	if err == nil {
		t.Error("Load() expected error for invalid port, got nil")
	}
}

func TestLoad_InvalidServerPort(t *testing.T) {
	clearEnv(t)

	t.Setenv("SERVER_PORT", "invalid")

	_, err := Load()
	if err == nil {
		t.Error("Load() expected error for invalid server port, got nil")
	}
}

// clearEnv clears all config-related environment variables for test isolation.
func clearEnv(t *testing.T) {
	t.Helper()

	vars := []string{
		"DATABASE_URL",
		"POSTGRES_HOST",
		"POSTGRES_PORT",
		"POSTGRES_USER",
		"POSTGRES_PASSWORD",
		"POSTGRES_DB",
		"POSTGRES_SSLMODE",
		"MINIO_ENDPOINT",
		"MINIO_ROOT_USER",
		"MINIO_ROOT_PASSWORD",
		"MINIO_USE_SSL",
		"MINIO_BUCKET",
		"SERVER_PORT",
	}

	for _, v := range vars {
		t.Setenv(v, "")
	}
}
