package main

import (
	"github.com/cassiokiyoshi/cheap-eats/internal/pathfinderhttp"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func mountPathfinder(router chi.Router, directory string) {
	router.Get("/pathfinder", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/pathfinder/", http.StatusPermanentRedirect)
	})
	router.Post("/pathfinder/api/search", pathfinderhttp.SearchHandler)
	router.Handle("/pathfinder/*", http.StripPrefix("/pathfinder", http.FileServer(http.Dir(directory))))
}
