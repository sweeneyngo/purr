package main

// Return s clamped to max length if max > 0, else return s.
func clamp(s string, start, max int) string {
	if s == "" {
		return ""
	}
	if max > 0 && len(s) > max {
		return s[:max]
	}
	return s
}

// Return s if it's non-empty, else return "".
func nz(s string) string {
	if s == "" {
		return ""
	}
	return s
}

// Return f if it's > 0, else return 1.
func nzf(f float64) float64 {
	if f <= 0 {
		return 1
	}
	return f
}

// Return the value of i if it's >= 0, else return 0.
func max0(i int) int {
	if i < 0 {
		return 0
	}
	return i
}

// Convert bool to int (1 for true, 0 for false).
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
