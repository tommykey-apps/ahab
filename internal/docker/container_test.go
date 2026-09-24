package docker

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestContainers_公開ポートの重複を除き名前付きボリュームだけ拾う(t *testing.T) {
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/"+apiVersion+"/containers/json" || r.URL.Query().Get("all") != "1" {
			t.Errorf("unexpected request %s", r.URL)
		}
		w.Write([]byte(`[
		 {"Id":"b","Names":["/web"],"Image":"nginx","ImageID":"sha256:x","State":"running",
		  "Ports":[{"PrivatePort":80,"PublicPort":8080},{"PrivatePort":80,"PublicPort":8080},{"PrivatePort":443}],
		  "NetworkSettings":{"Networks":{"front":{},"back":{}}},
		  "Mounts":[{"Type":"volume","Name":"zdata"},{"Type":"bind","Name":""},
		            {"Type":"volume","Name":"adata"},
		            {"Type":"volume","Name":"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"}]},
		 {"Id":"a","Names":["/db"],"Image":"postgres","ImageID":"sha256:y","State":"exited"}
		]`))
	}))
	got, err := c.Containers(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	want := []Container{
		{ID: "a", Name: "db", Image: "postgres", ImageID: "sha256:y", State: "exited"},
		{ID: "b", Name: "web", Image: "nginx", ImageID: "sha256:x", State: "running",
			Ports: []Port{{Public: 8080, Private: 80}}, Networks: []string{"back", "front"}, Volumes: []string{"adata", "zdata"}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v\nwant %+v", got, want)
	}
}

func TestIsAnonymous(t *testing.T) {
	cases := map[string]bool{
		"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef": true,
		"0123456789ABCDEF0123456789abcdef0123456789abcdef0123456789abcdef": false,
		"pgdata": false,
		"":       false,
	}
	for name, want := range cases {
		if got := isAnonymous(name); got != want {
			t.Errorf("isAnonymous(%q) = %v, want %v", name, got, want)
		}
	}
}

func TestRemoveContainer_vは常に付けforceは指定時だけ(t *testing.T) {
	var gotPath string
	var gotQuery url.Values
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotQuery = r.URL.Path, r.URL.Query()
		if r.Method != http.MethodDelete {
			t.Errorf("method %s", r.Method)
		}
		switch r.URL.Path {
		case "/" + apiVersion + "/containers/running":
			http.Error(w, `{"message":"cannot remove a running container"}`, http.StatusConflict)
		case "/" + apiVersion + "/containers/none":
			http.Error(w, `{"message":"No such container"}`, http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusNoContent)
		}
	}))
	ctx := context.Background()
	if err := c.RemoveContainer(ctx, "ok", false); err != nil {
		t.Fatal(err)
	}
	if gotPath != "/"+apiVersion+"/containers/ok" || gotQuery.Get("v") != "1" || gotQuery.Has("force") {
		t.Errorf("request %s?%s", gotPath, gotQuery.Encode())
	}
	if err := c.RemoveContainer(ctx, "ok", true); err != nil || gotQuery.Get("force") != "1" || gotQuery.Get("v") != "1" {
		t.Errorf("force: %s %v", gotQuery.Encode(), err)
	}
	var apiErr *APIError
	if err := c.RemoveContainer(ctx, "running", false); !errors.As(err, &apiErr) || apiErr.Status != http.StatusConflict || apiErr.Message != "cannot remove a running container" {
		t.Errorf("409: %v", err)
	}
	if err := c.RemoveContainer(ctx, "none", false); !errors.As(err, &apiErr) || apiErr.Status != http.StatusNotFound {
		t.Errorf("404: %v", err)
	}
}
