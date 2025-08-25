package main

import (
	_ "embed"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed pixel.gif
var pixelGIF []byte

func main() {
	db := mustOpenDB("analytics.db")
	defer db.Close()
	mustMigrate(db)

	mux := http.NewServeMux()
	mux.HandleFunc("/pixel.gif", func(w http.ResponseWriter, r *http.Request) {
		handlePixel(w, r, db)
	})
	mux.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		handleStats(w, r, db)
	})

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", withCommonHeaders(mux)))
}

func withCommonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Disallow caching of the pixel response
		w.Header().Set("Cache-Control", "no-store, max-age=0")
		// CORS is usually not required for <img>, but safe to allow
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}
