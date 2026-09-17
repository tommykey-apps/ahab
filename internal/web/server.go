package web

import (
	"embed"
	"html/template"
	"log"
	"net/http"

	"github.com/tommykey-apps/ahab/internal/docker"
)

var assets embed.FS

var tmpl = template.Must(templateParseFS(assets, "templates/*.html"))

type Server struct {
	docker *docker.Client
	mux *http.ServeMux
}

func NewServer(dc *docker.Client) *Server {
	s := &Server{docker: dc, mux: http.NewServMux()}
	s.mux.HandleFunc("GET /", s.index)
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	cs, err := s.docker.Containers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	if err := tmpl.ExecuteTemplate(w, "index.html", cs); err != nil {
		log.Printf("render %v", err)	
	}
}
