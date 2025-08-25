package main

import (
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	t.Helper()
	db := mustOpenDB(":memory:")
	mustMigrate(db)
	return db
}

func TestMustOpenDB(t *testing.T) {
	db := mustOpenDB(":memory:")
	if db == nil {
		t.Fatal("expected db, got nil")
	}
	if err := db.Ping(); err != nil {
		t.Fatalf("db.Ping() failed: %v", err)
	}
}

func TestMustMigrate_CreatesSchema(t *testing.T) {
	db := setupTestDB(t)

	// Check table exists
	rows, err := db.Query(`SELECT name FROM sqlite_master WHERE type='table' AND name='events'`)
	if err != nil {
		t.Fatalf("query failed: %v", err)
	}
	defer rows.Close()

	if !rows.Next() {
		t.Fatal("expected events table to exist")
	}
}

func TestInsertEvent_RoundTrip(t *testing.T) {
	db := setupTestDB(t)

	ev := Event{
		TS:      time.Now(),
		Path:    "/test",
		Title:   "Hello World",
		RefHost: "example.com",
		RefPath: "/ref",
		Browser: "Chrome",
		ScreenW: 1920,
		ScreenH: 1080,
		DPR:     2.0,
		Lang:    "en-US",
		JS:      true,
	}

	if err := insertEvent(db, ev); err != nil {
		t.Fatalf("insertEvent failed: %v", err)
	}

	// Fetch it back
	var got Event
	err := db.Get(&got, `SELECT ts, path, title, ref_host, ref_path, browser, screen_w, screen_h, dpr, lang, js FROM events LIMIT 1`)
	if err != nil {
		t.Fatalf("failed to query event: %v", err)
	}

	// Validate fields
	if got.Path != ev.Path {
		t.Errorf("Path = %q, want %q", got.Path, ev.Path)
	}
	if got.Title != ev.Title {
		t.Errorf("Title = %q, want %q", got.Title, ev.Title)
	}
	if got.RefHost != ev.RefHost {
		t.Errorf("RefHost = %q, want %q", got.RefHost, ev.RefHost)
	}
	if got.Browser != ev.Browser {
		t.Errorf("Browser = %q, want %q", got.Browser, ev.Browser)
	}
	if got.ScreenW != ev.ScreenW {
		t.Errorf("ScreenW = %d, want %d", got.ScreenW, ev.ScreenW)
	}
	if got.ScreenH != ev.ScreenH {
		t.Errorf("ScreenH = %d, want %d", got.ScreenH, ev.ScreenH)
	}
	if got.DPR != ev.DPR {
		t.Errorf("DPR = %f, want %f", got.DPR, ev.DPR)
	}
	if got.Lang != ev.Lang {
		t.Errorf("Lang = %q, want %q", got.Lang, ev.Lang)
	}
	if got.JS != ev.JS {
		t.Errorf("JS = %v, want %v", got.JS, ev.JS)
	}
}
