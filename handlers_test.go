package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// helper to create a temporary DB for tests
func setupTestTempDB(t *testing.T) *sqlx.DB {
	t.Helper()

	tmpfile := t.TempDir() + "/test.db"
	db := mustOpenDB(tmpfile)
	mustMigrate(db)
	return db
}

func TestHandlePixel_InsertsEvent(t *testing.T) {
	db := setupTestTempDB(t)

	q := url.Values{}
	q.Set("p", "/home")
	q.Set("t", "Homepage")
	q.Set("js", "1")
	q.Set("w", "1920")
	q.Set("h", "1080")
	q.Set("dpr", "2")
	q.Set("lang", "en-US")
	q.Set("ref", "https://example.com/ref")

	req := httptest.NewRequest("GET", "/pixel.gif?"+q.Encode(), nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/115.0")

	rr := httptest.NewRecorder()
	handlePixel(rr, req, db)

	res := rr.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status %d", res.StatusCode)
	}
	if ct := res.Header.Get("Content-Type"); ct != "image/gif" {
		t.Errorf("expected image/gif, got %s", ct)
	}

	// check DB for inserted event
	var count int
	err := db.Get(&count, `SELECT COUNT(*) FROM events WHERE path=?`, "/home")
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Errorf("expected 1 event inserted, got %d", count)
	}
}

func TestHandlePixel_NoJSRefererFallback(t *testing.T) {
	db := setupTestTempDB(t)

	// No path/title in query -> fallback to referer header
	req := httptest.NewRequest("GET", "/pixel.gif", nil)
	req.Header.Set("Referer", "https://foo.com/test/path")
	req.Header.Set("User-Agent", "Firefox/123")

	rr := httptest.NewRecorder()
	handlePixel(rr, req, db)

	var evt Event
	err := db.Get(&evt, `SELECT * FROM events LIMIT 1`)
	if err != nil {
		t.Fatal(err)
	}

	if evt.Path != "/test/path" {
		t.Errorf("expected path from referer '/test/path', got %q", evt.Path)
	}
	if evt.Browser != "Firefox" {
		t.Errorf("expected Firefox, got %q", evt.Browser)
	}
}

func TestHandleStats_ReturnsJSON(t *testing.T) {
	db := setupTestTempDB(t)

	// Insert some events manually
	insertEvent(db, Event{
		TS:      time.Now(),
		Path:    "/home",
		Title:   "Homepage",
		RefHost: "google.com",
		RefPath: "/search",
		Browser: "Chrome",
		ScreenW: 1920,
		ScreenH: 1080,
		DPR:     2,
		Lang:    "en",
		JS:      true,
	})

	req := httptest.NewRequest("GET", "/stats", nil)
	rr := httptest.NewRecorder()
	handleStats(rr, req, db)

	res := rr.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status %d", res.StatusCode)
	}
	if ct := res.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected application/json, got %s", ct)
	}

	var body map[string]any
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	// basic sanity check: at least one page entry
	pages := body["pages"].([]any)
	if len(pages) == 0 {
		t.Errorf("expected at least 1 page stat, got 0")
	}
}
