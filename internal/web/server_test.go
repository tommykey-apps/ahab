package web

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tommykey-apps/ahab/internal/docker"
	"github.com/tommykey-apps/ahab/internal/state"
)

func fakeDocker(t *testing.T) http.Handler {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1.44/images/json", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[{"Id":"sha256:img1","RepoTags":["nginx:latest"],"Size":10,"Created":1700000000}]`))
	})
	mux.HandleFunc("DELETE /v1.44/images/{id...}", func(w http.ResponseWriter, r *http.Request) {
		switch r.PathValue("id") {
		case "sha256:img1":
			http.Error(w, `{"message":"image is being used by running container abc"}`, http.StatusConflict)
		case "repo/name:tag":
			w.WriteHeader(http.StatusOK)
		default:
			http.Error(w, `{"message":"No such image: `+r.PathValue("id")+`"}`, http.StatusNotFound)
		}
	})
	mux.HandleFunc("GET /v1.44/volumes", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("filters") != "" {
			w.Write([]byte(`{"Volumes":[]}`))
			return
		}
		w.Write([]byte(`{"Volumes":[{"Name":"pgdata","Driver":"local","CreatedAt":"2026-09-24T00:00:00Z"}]}`))
	})
	mux.HandleFunc("DELETE /v1.44/volumes/{name}", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"message":"volume is in use"}`, http.StatusConflict)
	})
	mux.HandleFunc("DELETE /v1.44/containers/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.PathValue("id") == "missing":
			http.Error(w, `{"message":"No such container: missing"}`, http.StatusNotFound)
		case r.PathValue("id") == "running" && r.URL.Query().Get("force") != "1":
			http.Error(w, `{"message":"cannot remove a running container"}`, http.StatusConflict)
		default:
			w.WriteHeader(http.StatusNoContent)
		}
	})
	mux.HandleFunc("POST /v1.44/containers/{id}/{action}", func(w http.ResponseWriter, r *http.Request) {
		if r.PathValue("id") == "missing" {
			http.Error(w, `{"message":"No such container"}`, http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
	return mux
}

func newTestServer(t *testing.T) (*httptest.Server, *state.Store) {
	t.Helper()
	upstream := httptest.NewServer(fakeDocker(t))
	t.Cleanup(upstream.Close)
	dc := docker.NewWithTransport(&http.Transport{
		DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
			return (&net.Dialer{}).DialContext(ctx, "tcp", upstream.Listener.Addr().String())
		},
	})
	store := state.New()
	srv := httptest.NewServer(NewServer(dc, store))
	t.Cleanup(srv.Close)
	return srv, store
}

func do(t *testing.T, method, url string) (*http.Response, string) {
	t.Helper()
	req, _ := http.NewRequest(method, url, nil)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	b, _ := io.ReadAll(res.Body)
	return res, string(b)
}

func TestIndex_HTMLを返す(t *testing.T) {
	srv, _ := newTestServer(t)
	res, body := do(t, "GET", srv.URL+"/")
	if res.StatusCode != 200 || !strings.Contains(body, `<html lang="ja">`) {
		t.Fatalf("%d %q", res.StatusCode, body[:60])
	}
}

func TestAction_許可外と存在しないコンテナ(t *testing.T) {
	srv, _ := newTestServer(t)
	if res, _ := do(t, "POST", srv.URL+"/api/containers/abc/kill"); res.StatusCode != http.StatusBadRequest {
		t.Errorf("kill: %d", res.StatusCode)
	}
	if res, _ := do(t, "POST", srv.URL+"/api/containers/abc/start"); res.StatusCode != http.StatusNoContent {
		t.Errorf("start: %d", res.StatusCode)
	}
	if res, _ := do(t, "POST", srv.URL+"/api/containers/missing/start"); res.StatusCode != http.StatusBadGateway {
		t.Errorf("missing: %d", res.StatusCode)
	}
}

func TestRemoveContainer_Dockerの404と409をそのまま返しforceで通る(t *testing.T) {
	srv, _ := newTestServer(t)
	if res, body := do(t, "DELETE", srv.URL+"/api/containers/running"); res.StatusCode != http.StatusConflict || !strings.Contains(body, "running container") {
		t.Errorf("409: %d %q", res.StatusCode, body)
	}
	if res, _ := do(t, "DELETE", srv.URL+"/api/containers/running?force=1"); res.StatusCode != http.StatusNoContent {
		t.Errorf("force: %d", res.StatusCode)
	}
	if res, body := do(t, "DELETE", srv.URL+"/api/containers/missing"); res.StatusCode != http.StatusNotFound || !strings.Contains(body, "No such container: missing") {
		t.Errorf("404: %d %q", res.StatusCode, body)
	}
	if res, _ := do(t, "DELETE", srv.URL+"/api/containers/abc"); res.StatusCode != http.StatusNoContent {
		t.Errorf("ok: %d", res.StatusCode)
	}
}

func TestImages_使用中のコンテナ名を付ける(t *testing.T) {
	srv, store := newTestServer(t)
	store.Replace([]docker.Container{{ID: "c1", Name: "web", ImageID: "sha256:img1"}, {ID: "c2", Name: "other", ImageID: "sha256:zzz"}})
	res, body := do(t, "GET", srv.URL+"/api/images")
	if res.StatusCode != 200 || !strings.HasPrefix(res.Header.Get("Content-Type"), "application/json") {
		t.Fatalf("%d %s", res.StatusCode, res.Header.Get("Content-Type"))
	}
	var got []struct {
		ID     string
		Tags   []string
		UsedBy []string
	}
	if err := json.Unmarshal([]byte(body), &got); err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || len(got[0].UsedBy) != 1 || got[0].UsedBy[0] != "web" {
		t.Errorf("got %+v", got)
	}
}

func TestRemoveImage_Dockerの404と409をそのまま返す(t *testing.T) {
	srv, _ := newTestServer(t)
	if res, body := do(t, "DELETE", srv.URL+"/api/images/sha256:img1"); res.StatusCode != http.StatusConflict || !strings.Contains(body, "being used") {
		t.Errorf("409: %d %q", res.StatusCode, body)
	}
	if res, body := do(t, "DELETE", srv.URL+"/api/images/nope"); res.StatusCode != http.StatusNotFound || !strings.Contains(body, "No such image: nope") {
		t.Errorf("404: %d %q", res.StatusCode, body)
	}
	if res, _ := do(t, "DELETE", srv.URL+"/api/images/repo/name:tag"); res.StatusCode != http.StatusNoContent {
		t.Errorf("slash in id: %d", res.StatusCode)
	}
}

func TestVolumes_一覧と使用中の削除(t *testing.T) {
	srv, store := newTestServer(t)
	store.Replace([]docker.Container{{ID: "c1", Name: "db", Volumes: []string{"pgdata"}}, {ID: "c2", Name: "web", Volumes: []string{"other"}}})
	res, body := do(t, "GET", srv.URL+"/api/volumes")
	if res.StatusCode != 200 || !strings.Contains(body, `"Name":"pgdata"`) || !strings.Contains(body, `"InUse":true`) || !strings.Contains(body, `"UsedBy":["db"]`) {
		t.Errorf("list: %d %s", res.StatusCode, body)
	}
	if res, _ := do(t, "DELETE", srv.URL+"/api/volumes/pgdata"); res.StatusCode != http.StatusConflict {
		t.Errorf("delete: %d", res.StatusCode)
	}
}

func TestGraph_storeの内容をMermaidで返す(t *testing.T) {
	srv, store := newTestServer(t)
	store.Replace([]docker.Container{{ID: "c1", Name: "web", Networks: []string{"front"}}})
	res, body := do(t, "GET", srv.URL+"/api/graph")
	if res.StatusCode != 200 || !strings.HasPrefix(res.Header.Get("Content-Type"), "text/plain") {
		t.Fatalf("%d %s", res.StatusCode, res.Header.Get("Content-Type"))
	}
	if !strings.Contains(body, "c_web --> n_front") {
		t.Errorf("body %q", body)
	}
}

func TestStatic_書体と部品が同梱されている(t *testing.T) {
	srv, _ := newTestServer(t)
	for _, p := range []string{"/static/app.js", "/static/dads/tokens.css", "/static/fonts/noto-sans-jp-japanese-400-normal.woff2", "/static/mermaid.min.js"} {
		if res, _ := do(t, "GET", srv.URL+p); res.StatusCode != 200 {
			t.Errorf("%s: %d", p, res.StatusCode)
		}
	}
}
