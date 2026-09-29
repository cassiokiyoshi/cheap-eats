package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestPathfinderMount(t *testing.T) {
	directory := t.TempDir()
	if err := os.WriteFile(filepath.Join(directory, "index.html"), []byte("Pathfinder interface"), 0600); err != nil {
		t.Fatal(err)
	}
	router := chi.NewRouter()
	mountPathfinder(router, directory)
	router.Get("/api/health", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("cheap eats healthy")) })
	for _, tc := range []struct {
		url    string
		status int
		body   string
	}{
		{"/pathfinder", 308, ""},
		{"/pathfinder/", 200, "Pathfinder interface"},
		{"/pathfinder/missing.js", 404, ""},
		{"/api/health", 200, "cheap eats healthy"},
	} {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, httptest.NewRequest("GET", tc.url, nil))
		if w.Code != tc.status || !strings.Contains(w.Body.String(), tc.body) {
			t.Fatalf("%s: %d %s", tc.url, w.Code, w.Body.String())
		}
	}
	for _, algorithm := range []string{"astar", "dijkstra", "bfs"} {
		w := httptest.NewRecorder()
		body := `{"algorithm":"` + algorithm + `","rows":5,"cols":5,"start":{"row":0,"col":0},"goal":{"row":4,"col":4},"walls":[]}`
		router.ServeHTTP(w, httptest.NewRequest("POST", "/pathfinder/api/search", strings.NewReader(body)))
		if w.Code != 200 || w.Header().Get("Content-Type") != "application/x-ndjson" {
			t.Fatalf("stream failed: %d %s", w.Code, w.Body.String())
		}
		lines := strings.Split(strings.TrimSpace(w.Body.String()), "\n")
		var result struct {
			Type       string
			Found      bool
			PathLength int
		}
		if err := json.Unmarshal([]byte(lines[len(lines)-1]), &result); err != nil {
			t.Fatal(err)
		}
		if result.Type != "complete" || !result.Found || result.PathLength != 8 {
			t.Fatalf("%s: %+v", algorithm, result)
		}
	}
}

func TestPathfinderCORS(t *testing.T) {
	handler := pathfinderCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) }))
	for _, tc := range []struct {
		path, origin string
		allowed      bool
	}{
		{"/pathfinder/api/search", "https://cassiokiyoshi.github.io", true},
		{"/pathfinder/api/search", "https://example.com", false},
		{"/api/dishes", "https://cassiokiyoshi.github.io", false},
	} {
		req := httptest.NewRequest("OPTIONS", tc.path, nil)
		req.Header.Set("Origin", tc.origin)
		req.Header.Set("Access-Control-Request-Method", "POST")
		req.Header.Set("Access-Control-Request-Headers", "content-type")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		allowed := w.Header().Get("Access-Control-Allow-Origin") == tc.origin
		if allowed != tc.allowed {
			t.Fatalf("%s %s: allowed=%v", tc.path, tc.origin, allowed)
		}
	}
}
