package helps

import (
	"net/http"
	"reflect"
	"testing"
)

func TestClaudeRetryHeadersAllowlist(t *testing.T) {
	want := http.Header{
		"Anthropic-Ratelimit-Unified-Reset":      {"1900000000"},
		"Anthropic-Ratelimit-Requests-Remaining": {"0"},
		"Retry-After":                            {"30"},
	}
	src := want.Clone()
	for _, name := range []string{"Cookie", "Set-Cookie", "Authorization", "Proxy-Authorization", "X-Api-Key", "Retry-After-Ms", "X-Unrelated"} {
		src.Set(name, "synthetic")
	}
	got := ClaudeRetryHeaders(src)
	if !reflect.DeepEqual(got, want) {
		t.Fatal("retry headers do not match the allowlist")
	}
	got["Retry-After"][0] = "changed"
	if src.Get("Retry-After") != "30" {
		t.Fatal("retry header values alias the upstream header")
	}
}
