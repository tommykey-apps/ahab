package web

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *Server) sse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")

	ch := s.store.Subscribe()
	defer s.store.Unsubscribe(ch)

	rc := http.NewResponseController(w)
	send := func(v any) error {
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, "data: %s\n\n", b); err != nil {
			return err
		}
		return rc.Flush()
	}

	if err := send(s.store.Get()); err != nil {
		return
	}

	for {
		select {
		case cs := <-ch:
			if err := send(cs); err!=nil {
				return
			}
		case <-r.Context().Done():
			return
		}
	}
}
