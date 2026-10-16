package env_test

import (
	"testing"

	"github.com/cautem/quarshell-core/env"
)

func TestLooksLikeCredential(t *testing.T) {
	yes := []string{
		"GITHUB_TOKEN",
		"DB_TOKEN",
		"MY_ACCESS_KEY",
		"FOO_API_KEY",
		"SERVICE_PASSWORD",
		"CURSOR_API_KEY",
		"ANTHROPIC_API_KEY",
		"client_secret",
	}
	no := []string{
		"TOKENIZERS_PARALLELISM",
		"PASSWORDLESS_LOGIN",
		"TERM",
		"LANG",
		"PATH",
		"",
		"FEATURE_FLAG",
	}
	for _, k := range yes {
		if !env.LooksLikeCredential(k) {
			t.Fatalf("%q should look like credential", k)
		}
	}
	for _, k := range no {
		if env.LooksLikeCredential(k) {
			t.Fatalf("%q should NOT look like credential", k)
		}
	}
}
