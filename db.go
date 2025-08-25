package main

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

type Event struct {
	TS      time.Time `db:"ts" json:"ts"`
	Path    string    `db:"path" json:"path"`
	Title   string    `db:"title" json:"title"`
	RefHost string    `db:"ref_host" json:"ref_host"`
	RefPath string    `db:"ref_path" json:"ref_path"`
	Browser string    `db:"browser" json:"browser"`
	ScreenW int       `db:"screen_w" json:"screen_w"`
	ScreenH int       `db:"screen_h" json:"screen_h"`
	DPR     float64   `db:"dpr" json:"dpr"`
	Lang    string    `db:"lang" json:"lang"`
	JS      bool      `db:"js" json:"js"`
}

// Open (or create) the SQLite database at the given path.
// Set WAL journal mode and enable foreign keys.
func mustOpenDB(path string) *sqlx.DB {
	db, err := sqlx.Open("sqlite3", path+"?_journal_mode=WAL&_fk=1")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// Create or migrate the database schema (events table and indexes).
// Schema:
// - ts: 	timestamp of event
// - path: 	page path (from JS or derived from referer)
// - title:	page title (from JS)
// - ref_host:	referer host (if any)
// - ref_path:	referer path (if any)
// - browser:	browser family (Edge, Chrome, Firefox, Safari, Other)
// - screen_w:	screen width in pixels (from JS)
// - screen_h:	screen height in pixels (from JS)
// - dpr:		device pixel ratio (from JS)
// - lang:		language (from JS)
// - js:		whether JS was enabled (from JS)
func mustMigrate(db *sqlx.DB) {
	schema := `
CREATE TABLE IF NOT EXISTS events (
	ts        DATETIME NOT NULL,
	path      TEXT     NOT NULL,
	title     TEXT     NOT NULL,
	ref_host  TEXT     NOT NULL,
	ref_path  TEXT     NOT NULL,
	browser   TEXT     NOT NULL,
	screen_w  INTEGER  NOT NULL,
	screen_h  INTEGER  NOT NULL,
	dpr       REAL     NOT NULL,
	lang      TEXT     NOT NULL,
	js        INTEGER  NOT NULL
);
CREATE INDEX IF NOT EXISTS ix_events_ts ON events(ts);
CREATE INDEX IF NOT EXISTS ix_events_path ON events(path);
CREATE INDEX IF NOT EXISTS ix_events_refhost ON events(ref_host);
`
	if _, err := db.Exec(schema); err != nil {
		log.Fatal(err)
	}
}

// Insert a new event into the database.
func insertEvent(db *sqlx.DB, e Event) error {
	_, err := db.Exec(`
INSERT INTO events (ts, path, title, ref_host, ref_path, browser, screen_w, screen_h, dpr, lang, js)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		e.TS, nz(e.Path), nz(e.Title), nz(e.RefHost), nz(e.RefPath), nz(e.Browser),
		max0(e.ScreenW), max0(e.ScreenH), nzf(e.DPR), nz(e.Lang), boolToInt(e.JS),
	)
	return err
}
