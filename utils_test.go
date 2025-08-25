package main

import "testing"

func TestClamp(t *testing.T) {
	tests := []struct {
		s      string
		start  int
		max    int
		expect string
	}{
		{"hello", 0, 0, "hello"},  // no max clamp
		{"hello", 0, 3, "hel"},    // clamp to 3
		{"hello", 0, 10, "hello"}, // clamp larger than len
		{"", 0, 5, ""},            // empty string
	}
	for _, tt := range tests {
		got := clamp(tt.s, tt.start, tt.max)
		if got != tt.expect {
			t.Errorf("clamp(%q,%d,%d) = %q, want %q", tt.s, tt.start, tt.max, got, tt.expect)
		}
	}
}

func TestNz(t *testing.T) {
	if got := nz(""); got != "" {
		t.Errorf("nz(\"\") = %q, want \"\"", got)
	}
	if got := nz("abc"); got != "abc" {
		t.Errorf("nz(\"abc\") = %q, want \"abc\"", got)
	}
}

func TestNzf(t *testing.T) {
	tests := []struct {
		in     float64
		expect float64
	}{
		{0, 1},
		{-3, 1},
		{2.5, 2.5},
	}
	for _, tt := range tests {
		got := nzf(tt.in)
		if got != tt.expect {
			t.Errorf("nzf(%v) = %v, want %v", tt.in, got, tt.expect)
		}
	}
}

func TestMax0(t *testing.T) {
	tests := []struct {
		in     int
		expect int
	}{
		{-5, 0},
		{0, 0},
		{10, 10},
	}
	for _, tt := range tests {
		got := max0(tt.in)
		if got != tt.expect {
			t.Errorf("max0(%d) = %d, want %d", tt.in, got, tt.expect)
		}
	}
}

func TestBoolToInt(t *testing.T) {
	if got := boolToInt(true); got != 1 {
		t.Errorf("boolToInt(true) = %d, want 1", got)
	}
	if got := boolToInt(false); got != 0 {
		t.Errorf("boolToInt(false) = %d, want 0", got)
	}
}
