package pathfinderhttp

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/cassiokiyoshi/cheap-eats/internal/pathfinder"
)

const maxRequestBytes = 2 << 20

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	request, err := decodeSearchRequest(r)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := request.Validate(); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		writeJSONError(
			w,
			http.StatusInternalServerError,
			"streaming is not supported",
		)
		return
	}

	w.Header().Set("Content-Type", "application/x-ndjson")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	encoder := json.NewEncoder(w)
	startedAt := time.Now()
	visitedCount := 0

	path, found, err := pathfinder.Run(
		r.Context(),
		request,
		func(position pathfinder.Position) error {
			visitedCount++

			event := pathfinder.StreamEvent{
				Type:     "visited",
				Position: &position,
			}

			if err := encoder.Encode(event); err != nil {
				return err
			}

			flusher.Flush()
			return nil
		},
	)

	if err != nil {
		if errors.Is(err, r.Context().Err()) {
			return
		}

		_ = encoder.Encode(pathfinder.StreamEvent{
			Type:    "error",
			Message: err.Error(),
		})
		flusher.Flush()
		return
	}

	if found {
		if err := encoder.Encode(pathfinder.StreamEvent{
			Type:      "path",
			Positions: path,
		}); err != nil {
			return
		}

		flusher.Flush()
	}

	duration := float64(time.Since(startedAt).Microseconds()) / 1000
	pathLength := 0

	if found {
		pathLength = len(path) - 1
	}

	if err := encoder.Encode(pathfinder.StreamEvent{
		Type:       "complete",
		Visited:    visitedCount,
		PathLength: pathLength,
		DurationMS: duration,
		Found:      boolPointer(found),
	}); err != nil {
		return
	}

	flusher.Flush()
}

func decodeSearchRequest(r *http.Request) (pathfinder.SearchRequest, error) {
	var request pathfinder.SearchRequest

	decoder := json.NewDecoder(
		io.LimitReader(r.Body, maxRequestBytes),
	)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		return request, err
	}

	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return request, errors.New("request body must contain one JSON object")
	}

	return request, nil
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(pathfinder.StreamEvent{
		Type:    "error",
		Message: message,
	})
}

func boolPointer(value bool) *bool {
	return &value
}
