package main

import (
	"net/url"
	"strings"
)

// Splits a URL reference into its host and path components.
// If the reference is not a valid URL, it returns the entire ref as the path and an empty host.
func splitRef(ref string) (host, path string) {
	u, err := url.Parse(ref)
	if err != nil || u == nil {
		return "", ""
	}
	return u.Host, u.Path
}

// Parses the User-Agent string to identify the browser family.
// Returns one of "Edge", "Chrome", "Firefox", "Safari", or "Other".
func parseBrowserFamily(ua string) string {
	s := strings.ToLower(ua)
	switch {
	case strings.Contains(s, "edg/"):
		return "Edge"
	case strings.Contains(s, "chrome/") || strings.Contains(s, "crios/"):
		return "Chrome"
	case strings.Contains(s, "firefox/") || strings.Contains(s, "fxios/"):
		return "Firefox"
	case strings.Contains(s, "safari/"):
		return "Safari"
	default:
		return "Other"
	}
}
