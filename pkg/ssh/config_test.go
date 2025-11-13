package ssh

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig_Success(t *testing.T) {
	cfgPath := filepath.Join("examples", "correct_config.yaml")
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatalf("failed to read %s: %v", cfgPath, err)
	}

	cfg, err := ParseYAML(data)
	if err != nil {
		t.Fatalf("unexpected error parsing %s: %v", cfgPath, err)
	}

	if cfg == nil {
		t.Fatal("config is nil after parsing")
	}

	if _, ok := cfg.Hosts["local"]; !ok {
		t.Errorf("expected host alias 'local' to exist")
	}

	// this can be named as wanted
	server, ok := cfg.Hosts["server1"]
	if !ok {
		t.Fatalf("expected host alias 'server' to exist")
	}

	if server.IP != "1.0.0.1" {
		t.Errorf("expected IP '1.0.0.1', got %q", server.IP)
	}

	if server.Username != "user" {
		t.Errorf("expected Username 'user', got %q", server.Username)
	}

	if server.Port != 1000 {
		t.Errorf("expected Port 1000, got %d", server.Port)
	}

	if server.KeyFile == "" {
		t.Errorf("expected KeyFile to be non-empty")
	}
}
