package env_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/cautem/quarshell-core/env"
)

func TestFromHostForGuestNoGlobalAPIKeys(t *testing.T) {
	t.Setenv("TERM", "xterm")
	t.Setenv("CURSOR_API_KEY", "crsr_real")
	t.Setenv("ANTHROPIC_API_KEY", "sk-real")
	got := env.FromHostForGuest()
	for _, e := range got {
		k, _, _ := strings.Cut(e, "=")
		if k == "CURSOR_API_KEY" || k == "ANTHROPIC_API_KEY" {
			t.Fatalf("must not auto-inject %s without provider/env_allow: %v", k, got)
		}
	}
	got2 := env.FromHostForGuest("ANTHROPIC_API_KEY")
	var key string
	for _, e := range got2 {
		k, v, _ := strings.Cut(e, "=")
		if k == "ANTHROPIC_API_KEY" {
			key = v
		}
	}
	if key != env.PlaceholderPrefix+"ANTHROPIC_API_KEY" {
		t.Fatalf("key=%q", key)
	}
	// Cursor must still not appear unless explicitly requested
	for _, e := range got2 {
		if strings.HasPrefix(e, "CURSOR_API_KEY=") {
			t.Fatal("CURSOR_API_KEY leaked via unrelated extraKeys")
		}
	}
}

func TestFromHostForGuestPlaceholderWithoutHostEnv(t *testing.T) {
	t.Setenv("TERM", "xterm")
	_ = os.Unsetenv("GITHUB_TOKEN")
	got := env.FromHostForGuest("GITHUB_TOKEN")
	var key string
	for _, e := range got {
		k, v, _ := strings.Cut(e, "=")
		if k == "GITHUB_TOKEN" {
			key = v
		}
	}
	if key != env.PlaceholderPrefix+"GITHUB_TOKEN" {
		t.Fatalf("expected placeholder without host env, got %q in %v", key, got)
	}
}

func TestWarnCredentialEnv(t *testing.T) {
	var buf bytes.Buffer
	env.WarnCredentialEnv(&buf, map[string]string{
		"GITHUB_TOKEN": "x",
		"TERM":         "xterm",
	}, map[string][]env.ProfileHint{
		"GITHUB_TOKEN": env.HintKey("GITHUB_TOKEN"),
	}, false)
	out := buf.String()
	if !strings.Contains(out, "GITHUB_TOKEN") || !strings.Contains(out, "provider create") {
		t.Fatalf("warn output: %s", out)
	}
	if strings.Contains(out, "TERM") {
		t.Fatal("TERM should not warn")
	}
}
