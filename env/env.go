// Package env maps host environment variables into sandbox guests and
// resolves credential placeholders during egress rewrite.
package env

import (
	"os"
	"strings"
)

// PlaceholderPrefix is the guest-visible credential marker written into sandbox env.
const PlaceholderPrefix = "cautem:resolve:env:"

// LegacyCredentialPlaceholderPrefix preserves the alternate guest credential marker.
const LegacyCredentialPlaceholderPrefix = "openshell:resolve:env:"

// OpenShellPlaceholderPrefix is a deprecated alias retained for source compatibility.
const OpenShellPlaceholderPrefix = LegacyCredentialPlaceholderPrefix

// rewritePrefixes are accepted when the proxy rewrites requests (canonical first).
var rewritePrefixes = []string{
	PlaceholderPrefix,
	LegacyCredentialPlaceholderPrefix,
}

// DefaultPassthrough is forwarded into the sandbox with real host values (never secrets).
// Credential keys are NOT auto-injected from a global allowlist —
// only via policy credentials.env_allow / provider GuestEnvKeys (see FromHostForGuest extraKeys).
var DefaultPassthrough = []string{
	"TERM",
	"LANG",
	"LC_ALL",
	"LC_CTYPE",
	"COLORTERM",
	"DISPLAY",
}

// KnownCredentialKeys are well-known provider env names used for --env warnings
// and LooksLikeCredential. They are NOT auto-injected into the guest.
var KnownCredentialKeys = []string{
	"ANTHROPIC_API_KEY",
	"CLAUDE_API_KEY",
	"OPENAI_API_KEY",
	"OPENROUTER_API_KEY",
	"GOOGLE_API_KEY",
	"GEMINI_API_KEY",
	"GROQ_API_KEY",
	"CURSOR_API_KEY",
	"GITHUB_TOKEN",
	"GH_TOKEN",
	"NVIDIA_API_KEY",
	"NGC_API_KEY",
}

// passthroughKeys keep real host values in the guest (not secrets).
var passthroughKeys = map[string]struct{}{
	"TERM": {}, "LANG": {}, "LC_ALL": {}, "LC_CTYPE": {}, "COLORTERM": {}, "DISPLAY": {},
}

// IsPassthrough reports keys that should never be placeholder-rewritten.
func IsPassthrough(key string) bool {
	_, ok := passthroughKeys[key]
	return ok
}

// Filter keeps only allowlisted keys from environ (KEY=VAL entries).
func Filter(environ []string, keys ...string) []string {
	if len(keys) == 0 {
		return nil
	}
	allowed := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		k = strings.TrimSpace(k)
		if k != "" {
			allowed[k] = struct{}{}
		}
	}
	out := make([]string, 0, len(allowed))
	for _, entry := range environ {
		key, _, ok := strings.Cut(entry, "=")
		if !ok || key == "" {
			continue
		}
		if _, ok := allowed[key]; ok {
			out = append(out, entry)
		}
	}
	return out
}

// MergeAllowlists returns base plus extra keys (deduped, stable order).
func MergeAllowlists(base []string, extra ...string) []string {
	seen := make(map[string]struct{}, len(base)+len(extra))
	out := make([]string, 0, len(base)+len(extra))
	for _, k := range append(append([]string{}, base...), extra...) {
		k = strings.TrimSpace(k)
		if k == "" {
			continue
		}
		if _, ok := seen[k]; ok {
			continue
		}
		seen[k] = struct{}{}
		out = append(out, k)
	}
	return out
}

// FromHost builds allowlisted KEY=VAL from the current process environment (real values).
// Prefer FromHostForGuest for sandbox create when using credential placeholders.
func FromHost(extraKeys ...string) []string {
	keys := MergeAllowlists(DefaultPassthrough, extraKeys...)
	return Filter(os.Environ(), keys...)
}

// FromHostForGuest injects placeholders for credential keys and real values for passthrough keys.
//
// Passthrough keys are emitted only when present on the host.
// Credential extraKeys always get cautem:resolve:env:KEY placeholders — the real
// secret lives on the gateway/proxy (host env may be unset after provider refresh).
//
// OpenShell semantics: pass provider/policy keys as extraKeys. Do not rely on a
// global dump of API keys — Cursor (inject_env:false) must not see placeholders.
func FromHostForGuest(extraKeys ...string) []string {
	keys := MergeAllowlists(DefaultPassthrough, extraKeys...)
	out := make([]string, 0, len(keys))
	for _, k := range keys {
		if IsPassthrough(k) {
			v, ok := os.LookupEnv(k)
			if !ok {
				continue
			}
			out = append(out, k+"="+v)
			continue
		}
		out = append(out, k+"="+PlaceholderPrefix+k)
	}
	return out
}

// SecretsFromHost returns real KEY=VAL for credential keys present on the host (proxy sidecar).
// Only explicit extraKeys (plus non-passthrough) — not a global API-key allowlist.
func SecretsFromHost(extraKeys ...string) []string {
	keys := MergeAllowlists(nil, extraKeys...)
	var secretKeys []string
	for _, k := range keys {
		if IsPassthrough(k) {
			continue
		}
		secretKeys = append(secretKeys, k)
	}
	return Filter(os.Environ(), secretKeys...)
}

// PlaceholderFor returns cautem:resolve:env:KEY.
func PlaceholderFor(key string) string {
	return PlaceholderPrefix + key
}

// ContainsPlaceholder reports whether s embeds a credential marker.
func ContainsPlaceholder(s string) bool {
	for _, p := range rewritePrefixes {
		if strings.Contains(s, p) {
			return true
		}
	}
	return false
}

// ParsePlaceholder extracts KEY from a full placeholder token.
func ParsePlaceholder(token string) (key string, ok bool) {
	token = strings.TrimSpace(token)
	for _, p := range rewritePrefixes {
		if !strings.HasPrefix(token, p) {
			continue
		}
		key = strings.TrimPrefix(token, p)
		if key == "" || strings.ContainsAny(key, ":/ \t") {
			return "", false
		}
		return key, true
	}
	return "", false
}

// IsPlaceholder reports guest/host marker values (not real secrets).
func IsPlaceholder(v string) bool {
	v = strings.TrimSpace(v)
	for _, p := range rewritePrefixes {
		if strings.HasPrefix(v, p) {
			return true
		}
	}
	return false
}

// IndexPlaceholder returns the earliest marker start index and prefix length in s.
// If none, idx is -1.
func IndexPlaceholder(s string) (idx, prefixLen int) {
	best := -1
	plen := 0
	for _, p := range rewritePrefixes {
		i := strings.Index(s, p)
		if i < 0 {
			continue
		}
		if best < 0 || i < best {
			best = i
			plen = len(p)
		}
	}
	return best, plen
}
