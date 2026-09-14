package handlers

import (
	"log"
	"net/http"
)

func serverError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf(
		"request failed: method=%s path=%s error=%v",
		r.Method,
		r.URL.Path,
		err,
	)

	http.Error(
		w,
		"internal server error",
		http.StatusInternalServerError,
	)
}
