package defaults

import "testing"

func TestPublishedImage(t *testing.T) {
	cases := map[string]string{
		ImageLocal:               ImageBaseRef,
		"cautem-sandbox:base":    ImageBaseRef,
		ImageGUI:                 ImageGUIRef,
		ImageGPU:                 ImageGPURef,
		ImageGUIGPU:              ImageGUIGPURef,
		ImageCursor:              ImageCursorRef,
		" cautem-Sandbox:Claude": ImageClaudeRef,
		ImageCodex:               ImageCodexRef,
	}
	for in, want := range cases {
		got, ok := PublishedImage(in)
		if !ok || got != want {
			t.Fatalf("PublishedImage(%q) = %q, %v; want %q", in, got, ok, want)
		}
	}
	for _, in := range []string{"", "debian:bookworm", "cautem-sandbox:mine", "ghcr.io/cautem/cautem/sandboxes/base:latest"} {
		if got, ok := PublishedImage(in); ok {
			t.Fatalf("PublishedImage(%q) = %q, want no mapping", in, got)
		}
	}
}
