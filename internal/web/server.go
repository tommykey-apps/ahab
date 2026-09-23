package web

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"runtime"
	"strings"

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
	s.mux.HandleFunc("POST /api/containers/{id}/{action}", s.action)
	s.mux.HandleFunc("GET /api/containers/{id}/logs", s.logs)
	s.mux.HandleFunc("GET /debug/goroutines", s.debugGoroutine)
	s.mux.Handle("GET /static/", http.FileServerFS(assets))
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	if err := tmpl.ExecuteTemplate(w, "index.html", nil); err != nil {
		log.Printf("render %v", err)
	}
}

func (s *Server) action(w http.ResponseWriter, r *http.Request) {
	action := r.PathValue("action")
	switch action {
	case "start", "stop", "restart":
	default:
		http.Error(w, "unknown action", http.StatusBadRequest)
		return
	}
	if err := s.docker.Action(r.Context(), r.PathValue("id"), action); err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) logs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	rc := http.NewResponseController(w)

	lines := make(chan string)
	go func() { _ = s.docker.Logs(r.Context(), r.PathValue("id"), lines) }()
	for {
		select {
		case l := <-lines:
			l = strings.ReplaceAll(strings.TrimRight(l, "\n"), "\n", "\ndata: ")
			fmt.Fprintf(w, "data: %s\n\n", l)
			if err := rc.Flush(); err != nil {
				return
			}
		case <-r.Context().Done():
			return
		}
	}
}

func (s *Server) debugGoroutine(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, runtime.NumGoroutine())
}
