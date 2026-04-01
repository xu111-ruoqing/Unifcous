package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadAppliesEnvOverrides(t *testing.T) {
	configPath := writeTestConfig(t)

	t.Setenv("DB_HOST", "db-from-env")
	t.Setenv("DB_PORT", "15432")
	t.Setenv("API_PORT", "18080")
	t.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://127.0.0.1:5173")
	t.Setenv("CORS_ALLOW_CREDENTIALS", "true")
	t.Setenv("SEED_COMPETITIONS_ON_START", "false")

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg.Database.Host != "db-from-env" {
		t.Fatalf("expected DB host from env, got %q", cfg.Database.Host)
	}
	if cfg.Database.Port != 15432 {
		t.Fatalf("expected DB port from env 15432, got %d", cfg.Database.Port)
	}
	if cfg.Server.Port != 18080 {
		t.Fatalf("expected server port from env 18080, got %d", cfg.Server.Port)
	}
	if len(cfg.CORS.AllowedOrigins) != 2 {
		t.Fatalf("expected 2 allowed origins, got %d", len(cfg.CORS.AllowedOrigins))
	}
	if !cfg.CORS.AllowCredentials {
		t.Fatal("expected CORS allow_credentials=true from env")
	}
	if cfg.Bootstrap.ShouldSeedCompetitions() {
		t.Fatal("expected bootstrap seed to be disabled by env")
	}
}

func TestLoadReturnsErrorForInvalidBooleanEnv(t *testing.T) {
	configPath := writeTestConfig(t)

	t.Setenv("CORS_ALLOW_CREDENTIALS", "not-a-bool")

	_, err := Load(configPath)
	if err == nil {
		t.Fatal("expected error for invalid boolean env value")
	}
	if !strings.Contains(err.Error(), "invalid boolean") {
		t.Fatalf("expected invalid boolean error, got %v", err)
	}
}

func writeTestConfig(t *testing.T) string {
	t.Helper()

	content := `server:
  port: 8080
  mode: debug

database:
  host: 127.0.0.1
  port: 5432
  user: unifocus
  password: test
  dbname: unifocus_dev

redis:
  host: 127.0.0.1
  port: 6379

jwt:
  secret: test-secret-key-with-minimum-length

nlp_service:
  url: http://127.0.0.1:8000

cors:
  allowed_origins:
    - http://localhost:3000
`

	dir := t.TempDir()
	path := filepath.Join(dir, "config.test.yaml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}
	return path
}
