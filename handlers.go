package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
)

// Handle requests to /pixel.gif
// Get parameters from query string and headers, store event in DB, return 1x1 transparent GIF.
func handlePixel(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	q := r.URL.Query()

	// We prefer params from JS; for no-JS we rely on the Referer header.
	path := clamp(q.Get("p"), 0, 2048)
	title := clamp(q.Get("t"), 0, 256)
	jsFlag := q.Get("js") == "1"

	// Screen info (optional)
	wStr, hStr := q.Get("w"), q.Get("h")
	dprStr := q.Get("dpr")
	lang := clamp(q.Get("lang"), 0, 32)

	sw, _ := strconv.Atoi(wStr)
	sh, _ := strconv.Atoi(hStr)
	dpr, _ := strconv.ParseFloat(dprStr, 64)

	// Referrer from query (JS) OR header (no-JS)
	ref := q.Get("ref")
	if ref == "" {
		ref = r.Referer()
	}
	refHost, refPath := splitRef(ref)

	// If no JS path given, try to derive from referer URL
	if path == "" {
		if u, err := url.Parse(ref); err == nil && u != nil {
			path = u.Path
		}
	}

	// Minimal browser family detection from UA (do NOT store full UA for privacy)
	ua := r.Header.Get("User-Agent")
	browser := parseBrowserFamily(ua)

	evt := Event{
		TS:      time.Now().UTC(),
		Path:    path,
		Title:   title,
		RefHost: refHost,
		RefPath: refPath,
		Browser: browser,
		ScreenW: sw,
		ScreenH: sh,
		DPR:     dpr,
		Lang:    lang,
		JS:      jsFlag,
	}

	if err := insertEvent(db, evt); err != nil {
		// Never fail the pixel; just log
		log.Printf("insert error: %v", err)
	}

	// Return the pixel
	w.Header().Set("Content-Type", "image/gif")
	w.Header().Set("Content-Length", strconv.Itoa(len(pixelGIF)))
	_, _ = w.Write(pixelGIF)
}

// Handle requests to /stats
// Return JSON with aggregated stats: top pages, referrers, browsers, screen sizes.
func handleStats(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	type Row struct {
		Key   string `db:"k" json:"key"`
		Count int    `db:"c" json:"count"`
	}

	var pages []Row
	_ = db.Select(&pages, `SELECT path AS k, COUNT(*) AS c FROM events GROUP BY path ORDER BY c DESC LIMIT 50`)

	var refs []Row
	_ = db.Select(&refs, `SELECT ref_host AS k, COUNT(*) AS c FROM events WHERE ref_host!='' GROUP BY ref_host ORDER BY c DESC LIMIT 50`)

	var browsers []Row
	_ = db.Select(&browsers, `SELECT browser AS k, COUNT(*) AS c FROM events GROUP BY browser ORDER BY c DESC`)

	var sizes []Row
	_ = db.Select(&sizes, `SELECT printf('%dx%d@%g', screen_w, screen_h, dpr) AS k, COUNT(*) AS c FROM events GROUP BY k ORDER BY c DESC LIMIT 50`)

	resp := map[string]any{
		"pages":     pages,
		"referrers": refs,
		"browsers":  browsers,
		"screens":   sizes,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
