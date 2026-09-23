package docker

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestClient(t *testing.T, h http.Handler) *Client {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	return &Client{http: &http.Client{Transport: &http.Transport{
		DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
			return (&net.Dialer{}).DialContext(ctx, "tcp", srv.Listener.Addr().String())
		},
	}}}
}

func TestSocketPath_DOCKER_HOSTがunixなら優先する(t *testing.T) {
	t.Setenv("DOCKER_HOST", "unix:///tmp/x.sock")
	if got := socketPath(); got != "/tmp/x.sock" {
		t.Fatalf("got %q", got)
	}
	t.Setenv("DOCKER_HOST", "tcp://localhost:2375")
	if got := socketPath(); got != "/var/run/docker.sock" {
		t.Fatalf("got %q", got)
	}
}

func TestGet_200以外はエラー(t *testing.T) {
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"message":"boom"}`, http.StatusInternalServerError)
	}))
	if _, err := c.get(context.Background(), "/info"); err == nil {
		t.Fatal("expected error")
	}
}
