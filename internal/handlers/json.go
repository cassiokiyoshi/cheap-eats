package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func readJSON(
	w http.ResponseWriter,
	r *http.Request,
	destination any,
) error {
	const maxBodyBytes = 64 * 1024

	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(destination); err != nil {
		return err
	}

	var extra any
	err := decoder.Decode(&extra)

	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err
	}

	return fmt.Errorf("body must contain exactly one JSON value")
}
