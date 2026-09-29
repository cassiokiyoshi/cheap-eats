package main

import (
	"github.com/cassiokiyoshi/cheap-eats/internal/pathfinderhttp"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"net/http"
)

func mountPathfinder(router chi.Router, directory string) {
	router.Get("/pathfinder", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/pathfinder/", http.StatusPermanentRedirect)
	})
	router.Post("/pathfinder/api/search", pathfinderhttp.SearchHandler)
	router.Handle("/pathfinder/*", http.StripPrefix("/pathfinder", http.FileServer(http.Dir(directory))))
}

// Only the stateless search endpoint accepts cross-origin Pages requests.
func pathfinderCORS(next http.Handler) http.Handler {
	search := cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://cassiokiyoshi.github.io"},
		AllowedMethods: []string{http.MethodPost, http.MethodOptions},
		AllowedHeaders: []string{"Content-Type"},
		MaxAge:         300,
	})(next)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/pathfinder/api/search" {
			search.ServeHTTP(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
