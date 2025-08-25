package main

import "testing"

func TestSplitRef(t *testing.T) {
	tests := []struct {
		ref  string
		host string
		path string
	}{
		{"https://example.com/foo/bar", "example.com", "/foo/bar"},
		{"http://sub.domain.com/", "sub.domain.com", "/"},
		{"not a url", "", "not a url"},
		{"", "", ""},
		{"ftp://example.org/resource", "example.org", "/resource"},
	}
	for _, tt := range tests {
		host, path := splitRef(tt.ref)
		if host != tt.host || path != tt.path {
			t.Errorf("splitRef(%q) = (%q,%q), want (%q,%q)", tt.ref, host, path, tt.host, tt.path)
		}
	}
}

func TestParseBrowserFamily(t *testing.T) {
	tests := []struct {
		ua     string
		expect string
	}{
		{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) Edg/116.0", "Edge"},
		{"Mozilla/5.0 (X11; Linux x86_64) Chrome/115.0.0.0 Safari/537.36", "Chrome"},
		{"Mozilla/5.0 (iPhone; CPU iPhone OS 15_5) CriOS/103.0.5060.63", "Chrome"},
		{"Mozilla/5.0 (Windows NT 10.0; rv:102.0) Gecko/20100101 Firefox/102.0", "Firefox"},
		{"Mozilla/5.0 (iPhone; CPU iPhone OS 15_5) FxiOS/102.0", "Firefox"},
		{"Mozilla/5.0 (Macintosh; Intel Mac OS X) Safari/605.1.15", "Safari"},
		{"SomethingRandomBrowser/1.0", "Other"},
	}
	for _, tt := range tests {
		got := parseBrowserFamily(tt.ua)
		if got != tt.expect {
			t.Errorf("parseBrowserFamily(%q) = %q, want %q", tt.ua, got, tt.expect)
		}
	}
}
