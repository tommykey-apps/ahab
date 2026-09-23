package web

import (
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/tommykey-apps/ahab/internal/docker"
	"github.com/tommykey-apps/ahab/internal/graph"
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
	s.mux.HandleFunc("GET /api/images", s.images)
	s.mux.HandleFunc("DELETE /api/images/{id...}", s.removeImage)
	s.mux.HandleFunc("GET /api/volumes", s.volumes)
	s.mux.HandleFunc("DELETE /api/volumes/{name}", s.removeVolume)
	s.mux.HandleFunc("GET /api/graph", s.graph)
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

	since := ""
	var lastT time.Time
	if last := r.Header.Get("Last-Event-ID"); last != "" {
		if t, err := time.Parse(time.RFC3339Nano, last); err == nil {
			since = fmt.Sprintf("%d.%09d", t.Unix(), t.Nanosecond())
			lastT = t
		}
	}

	lines := make(chan string)
	errc := make(chan error, 1)
	go func() { errc <- s.docker.Logs(r.Context(), r.PathValue("id"), since, lines) }()
	for {
		select {
		case l := <-lines:
			ts, body, _ := strings.Cut(l, " ")
			if t, err := time.Parse(time.RFC3339Nano, ts); err == nil && !t.After(lastT) {
				continue
			}
			body = strings.ReplaceAll(strings.TrimRight(body, "\n"), "\n", "\ndata: ")
			fmt.Fprintf(w, "id: %s\ndata: %s\n\n", ts, body)
			if err := rc.Flush(); err != nil {
				return
			}
		case err := <-errc:
			if !errors.Is(err, io.EOF) && r.Context().Err() == nil {
				fmt.Fprintf(w, "event: logerror\ndata: %s\n\n", err)
				rc.Flush()
			}
			return
		case <-r.Context().Done():
			return
		}
	}
}

func (s *Server) graph(w http.ResponseWriter, r *http.Request) {
	views := s.store.Get()
	cs := make([]docker.Container, 0, len(views))
	for _, v := range views {
		cs = append(cs, v.Container)
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(w, graph.Render(cs))
}

func (s *Server) debugGoroutine(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, runtime.NumGoroutine())
}
