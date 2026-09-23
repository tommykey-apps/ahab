package web

import "net/http"

func (s *Server) volumes(w http.ResponseWriter, r *http.Request) {
	vols, err := s.docker.Volumes(r.Context())
	if err != nil {
		writeDockerError(w, err)
		return
	}
	writeJSON(w, vols)
}

func (s *Server) removeVolume(w http.ResponseWriter, r *http.Request) {
	force := r.URL.Query().Get("force") == "1"
	if err := s.docker.RemoveVolume(r.Context(), r.PathValue("name"), force); err != nil {
		writeDockerError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
