package web

import (
	"net/http"

	"github.com/tommykey-apps/ahab/internal/docker"
)

type imageView struct {
	docker.Image
	UsedBy []string
}

func (s *Server) images(w http.ResponseWriter, r *http.Request) {
	imgs, err := s.docker.Images(r.Context())
	if err != nil {
		writeDockerError(w, err)
		return
	}
	usedBy := map[string][]string{}
	for _, c := range s.store.Get() {
		usedBy[c.ImageID] = append(usedBy[c.ImageID], c.Name)
	}
	views := make([]imageView, 0, len(imgs))
	for _, img := range imgs {
		views = append(views, imageView{Image: img, UsedBy: usedBy[img.ID]})
	}
	writeJSON(w, views)
}

func (s *Server) removeImage(w http.ResponseWriter, r *http.Request) {
	force := r.URL.Query().Get("force") == "1"
	if err := s.docker.RemoveImage(r.Context(), r.PathValue("id"), force); err != nil {
		writeDockerError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
