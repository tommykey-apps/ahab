package docker

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"
)

func TestImages_タグ無しは末尾でnoneタグは除く(t *testing.T) {
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[
		 {"Id":"sha256:3","RepoTags":null,"Size":30,"Created":1700000000},
		 {"Id":"sha256:2","RepoTags":["<none>:<none>"],"Size":20,"Created":1700000000},
		 {"Id":"sha256:1","RepoTags":["nginx:latest","nginx:1.27"],"Size":10,"Created":1700000000},
		 {"Id":"sha256:0","RepoTags":["alpine:3"],"Size":5,"Created":1700000000}
		]`))
	}))
	got, err := c.Images(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 4 || got[0].ID != "sha256:0" || got[1].ID != "sha256:1" || got[2].ID != "sha256:2" || got[3].ID != "sha256:3" {
		t.Fatalf("order: %+v", got)
	}
	if len(got[1].Tags) != 2 || got[1].Tags[0] != "nginx:1.27" {
		t.Errorf("tags not sorted: %v", got[1].Tags)
	}
	if got[2].Tags != nil {
		t.Errorf("<none> tag must be dropped: %v", got[2].Tags)
	}
	if !got[0].Created.Equal(time.Unix(1700000000, 0)) {
		t.Errorf("created: %v", got[0].Created)
	}
}

func TestRemoveImage_forceとステータス(t *testing.T) {
	var gotPath, gotQuery string
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotQuery = r.URL.Path, r.URL.RawQuery
		switch r.URL.Path {
		case "/" + apiVersion + "/images/used":
			http.Error(w, `{"message":"image is being used"}`, http.StatusConflict)
		case "/" + apiVersion + "/images/none":
			http.Error(w, `{"message":"No such image"}`, http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusOK)
		}
	}))
	ctx := context.Background()
	if err := c.RemoveImage(ctx, "ok", true); err != nil {
		t.Fatal(err)
	}
	if gotPath != "/"+apiVersion+"/images/ok" || gotQuery != "force=1" {
		t.Errorf("request %s?%s", gotPath, gotQuery)
	}
	if err := c.RemoveImage(ctx, "ok", false); err != nil || gotQuery != "" {
		t.Errorf("force must be omitted: %q %v", gotQuery, err)
	}
	var apiErr *APIError
	if err := c.RemoveImage(ctx, "used", false); !errors.As(err, &apiErr) || apiErr.Status != http.StatusConflict || apiErr.Message != "image is being used" {
		t.Errorf("409: %v", err)
	}
	if err := c.RemoveImage(ctx, "none", false); !errors.As(err, &apiErr) || apiErr.Status != http.StatusNotFound {
		t.Errorf("404: %v", err)
	}
}

func TestRemove_本文がJSONでなければステータス文を使う(t *testing.T) {
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "plain text", http.StatusBadGateway)
	}))
	var apiErr *APIError
	err := c.remove(context.Background(), "/x")
	if !errors.As(err, &apiErr) || apiErr.Status != http.StatusBadGateway || apiErr.Message != "502 Bad Gateway" {
		t.Errorf("got %v", err)
	}
}
