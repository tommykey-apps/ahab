package web

import (
	"embed"
	"html/template"
	"log"
	"net/http"

	"github.com/tommykey-apps/ahab/internal/docker"
	"github.com/tommykey-apps/ahab/internal/state"
)

//go:embed templates static
var assets embed.FS

var tmpl = template.Must(template.ParseFS(assets, "templates/*.html"))

type Server struct {
	docker *docker.Client
	store  *state.Store
	mux    *http.ServeMux
}

func NewServer(dc *docker.Client, store *state.Store) *Server {
	s := &Server{docker: dc, store: store, mux: http.NewServeMux()}
	s.mux.HandleFunc("GET /", s.index)
	s.mux.HandleFunc("GET /events", s.sse)
	s.mux.Handle("GET /static/", http.FileServerFS(assets))
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
