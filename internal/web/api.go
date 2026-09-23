package web

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/tommykey-apps/ahab/internal/docker"
)

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("encode: %v", err)
	}
}

func writeDockerError(w http.ResponseWriter, err error) {
	var apiErr *docker.APIError
	if errors.As(err, &apiErr) {
		switch apiErr.Status {
		case http.StatusNotFound, http.StatusConflict:
			http.Error(w, apiErr.Message, apiErr.Status)
			return
		}
	}
	http.Error(w, err.Error(), http.StatusBadGateway)
}
