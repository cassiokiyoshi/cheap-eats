package main

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestWebRoutes(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte("app shell"), 0600); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		path   string
		status int
	}{
		{"/", 200},
		{"/dishes/1", 200},
		{"/api/missing", 404},
		{"/missing.js", 404},
	} {
		t.Run(tc.path, func(t *testing.T) {
			response := httptest.NewRecorder()
			webHandler(dir).ServeHTTP(response, httptest.NewRequest("GET", tc.path, nil))
			if response.Code != tc.status {
				t.Fatalf("got %d, want %d", response.Code, tc.status)
			}
		})
	}
}
