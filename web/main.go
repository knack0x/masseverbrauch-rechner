package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
)

//go:embed assets
var assets embed.FS

var webVersion string

func init() {
	for _, path := range []string{"VERSION", "../VERSION"} {
		if data, err := os.ReadFile(path); err == nil {
			webVersion = strings.TrimSpace(string(data))
			return
		}
	}
	webVersion = "dev"
}

func apiBaseURL() string {
	u := os.Getenv("API_URL")
	if u == "" {
		u = "http://localhost:8080/api/calculate"
	}
	return strings.TrimSuffix(u, "/api/calculate")
}

func Serve() {
	if err := initTemplates(); err != nil {
		log.Fatalf("[fatal] Failed to parse templates: %v", err)
	}

	assetsFS, err := fs.Sub(assets, "assets")
	if err != nil {
		log.Fatalf("[fatal] Failed to get assets sub FS: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(assetsFS))))
	mux.HandleFunc("/manifest.json", func(w http.ResponseWriter, r *http.Request) {
		data, err := assets.ReadFile("assets/manifest.json")
		if err != nil {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/calculate", calculateHandler)
	mux.HandleFunc("/api/version", versionHandler)

	handler := loggingMiddleware(cacheMiddleware(mux))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("[info] Server starting on :%s (version: %s)", port, webVersion)
	log.Fatalf("[fatal] %v", http.ListenAndServe(":"+port, handler))
}

func main() {
	Serve()
}
