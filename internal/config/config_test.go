package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_WhoisDefaults(t *testing.T) {
	// Point the database at a temporary directory so Load does not create one
	// under the real home directory.
	t.Setenv("DATABASE_PATH", filepath.Join(t.TempDir(), "db.sqlite"))

	cfg, err := Load(filepath.Join(t.TempDir(), "does-not-exist.json"))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if !cfg.Whois.Enabled {
		t.Error("whois enrichment should be on by default")
	}
	if cfg.Whois.RDAPURL != DefaultRDAPURL {
		t.Errorf("RDAPURL = %q, want %q", cfg.Whois.RDAPURL, DefaultRDAPURL)
	}
	if cfg.Whois.TTLHours != DefaultWhoisTTLHours {
		t.Errorf("TTLHours = %d, want %d", cfg.Whois.TTLHours, DefaultWhoisTTLHours)
	}
}

func TestLoad_WhoisFromEnv(t *testing.T) {
	t.Setenv("DATABASE_PATH", filepath.Join(t.TempDir(), "db.sqlite"))
	t.Setenv("WHOIS_ENABLED", "false")
	// A URL without a trailing slash must still produce a usable endpoint.
	t.Setenv("WHOIS_RDAP_URL", "https://rdap.example.com/ip")
	t.Setenv("WHOIS_TTL_HOURS", "24")

	cfg, err := Load(filepath.Join(t.TempDir(), "does-not-exist.json"))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if cfg.Whois.Enabled {
		t.Error("WHOIS_ENABLED=false was ignored")
	}
	if cfg.Whois.RDAPURL != "https://rdap.example.com/ip/" {
		t.Errorf("RDAPURL = %q, want a trailing slash appended", cfg.Whois.RDAPURL)
	}
	if cfg.Whois.TTLHours != 24 {
		t.Errorf("TTLHours = %d, want 24", cfg.Whois.TTLHours)
	}
}

func TestLoad_WhoisFileOverridesEnv(t *testing.T) {
	t.Setenv("DATABASE_PATH", filepath.Join(t.TempDir(), "db.sqlite"))
	t.Setenv("WHOIS_ENABLED", "true")

	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(`{"whois":{"enabled":false,"ttl_hours":0}}`), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if cfg.Whois.Enabled {
		t.Error("the config file must win over the environment")
	}
	// A zero or negative TTL in the file falls back to the default rather
	// than making every lookup instantly stale.
	if cfg.Whois.TTLHours != DefaultWhoisTTLHours {
		t.Errorf("TTLHours = %d, want the default %d", cfg.Whois.TTLHours, DefaultWhoisTTLHours)
	}
}
