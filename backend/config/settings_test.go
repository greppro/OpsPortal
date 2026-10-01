package config

import (
	"reflect"
	"testing"
	"time"
)

var settingsEnvKeys = []string{
	"HOST", "PORT", "DB_PATH", "JWT_SECRET", "ADMIN_PASSWORD", "CORS_ALLOW_ORIGINS", "METRICS_TOKEN",
	"PROBE_ENABLED", "PROBE_INTERVAL", "PROBE_TIMEOUT", "PROBE_CONCURRENCY", "PROBE_INSECURE_SKIP_VERIFY",
}

func clearSettingsEnv(t *testing.T) {
	for _, k := range settingsEnvKeys {
		t.Setenv(k, "")
	}
}

func TestLoadSettingsDefaults(t *testing.T) {
	clearSettingsEnv(t)
	s := LoadSettings()
	if s.Host != "" || s.Port != "8080" || s.DBPath != "data/opsportal.db" || s.JWTSecret != "" || s.AdminPassword != "" {
		t.Errorf("unexpected defaults: %+v", s)
	}
	if len(s.CORSAllowOrigins) != 0 || s.MetricsToken != "" {
		t.Errorf("CORS and metrics token should be empty by default: %+v", s)
	}
	want := ProbeSettings{Enabled: true, Interval: 60 * time.Second, Timeout: 5 * time.Second, Concurrency: 10}
	if s.Probe != want {
		t.Errorf("probe defaults = %+v, want %+v", s.Probe, want)
	}
}

func TestLoadSettingsOverrides(t *testing.T) {
	clearSettingsEnv(t)
	t.Setenv("HOST", "127.0.0.1")
	t.Setenv("PORT", "9000")
	t.Setenv("DB_PATH", "/var/lib/opsportal/portal.db")
	t.Setenv("CORS_ALLOW_ORIGINS", " https://a.example.com, ,https://b.example.com ")
	t.Setenv("PROBE_ENABLED", "false")
	t.Setenv("PROBE_INTERVAL", "2m")
	t.Setenv("PROBE_TIMEOUT", "8")
	t.Setenv("PROBE_CONCURRENCY", "20")
	t.Setenv("PROBE_INSECURE_SKIP_VERIFY", "true")

	s := LoadSettings()
	if s.Host != "127.0.0.1" || s.Port != "9000" || s.DBPath != "/var/lib/opsportal/portal.db" {
		t.Errorf("unexpected port/db: %+v", s)
	}
	if !reflect.DeepEqual(s.CORSAllowOrigins, []string{"https://a.example.com", "https://b.example.com"}) {
		t.Errorf("CORS origins = %v", s.CORSAllowOrigins)
	}
	want := ProbeSettings{Enabled: false, Interval: 2 * time.Minute, Timeout: 8 * time.Second, Concurrency: 20, InsecureSkipVerify: true}
	if s.Probe != want {
		t.Errorf("probe = %+v, want %+v", s.Probe, want)
	}
}

func TestLoadSettingsClampsInvalidValues(t *testing.T) {
	clearSettingsEnv(t)
	t.Setenv("PROBE_ENABLED", "maybe")
	t.Setenv("PROBE_INTERVAL", "5s")
	t.Setenv("PROBE_TIMEOUT", "1m")
	t.Setenv("PROBE_CONCURRENCY", "0")

	s := LoadSettings()
	if !s.Probe.Enabled {
		t.Error("invalid PROBE_ENABLED should fall back to default true")
	}
	if s.Probe.Interval != 10*time.Second {
		t.Errorf("interval should be clamped to 10s, got %s", s.Probe.Interval)
	}
	if s.Probe.Timeout != s.Probe.Interval {
		t.Errorf("timeout should not exceed interval, got %s", s.Probe.Timeout)
	}
	if s.Probe.Concurrency != 10 {
		t.Errorf("invalid concurrency should fall back to 10, got %d", s.Probe.Concurrency)
	}
}
