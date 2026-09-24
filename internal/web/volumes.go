package web

import (
	"net/http"

	"github.com/tommykey-apps/ahab/internal/docker"
)

type volumeView struct {
	docker.Volume
	UsedBy []string
}

func (s *Server) volumes(w http.ResponseWriter, r *http.Request) {
	vols, err := s.docker.Volumes(r.Context())
	if err != nil {
		writeDockerError(w, err)
		return
	}
	usedBy := map[string][]string{}
	for _, c := range s.store.Get() {
		for _, v := range c.Volumes {
			usedBy[v] = append(usedBy[v], c.Name)
		}
	}
	views := make([]volumeView, 0, len(vols))
	for _, v := range vols {
		views = append(views, volumeView{Volume: v, UsedBy: usedBy[v.Name]})
	}
	writeJSON(w, views)
}

func (s *Server) removeVolume(w http.ResponseWriter, r *http.Request) {
	force := r.URL.Query().Get("force") == "1"
	if err := s.docker.RemoveVolume(r.Context(), r.PathValue("name"), force); err != nil {
		writeDockerError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
