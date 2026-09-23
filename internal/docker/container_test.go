package docker

import (
	"context"
	"net/http"
	"reflect"
	"testing"
)

func TestContainers_名前付きボリュームだけ拾いソートする(t *testing.T) {
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/"+apiVersion+"/containers/json" || r.URL.Query().Get("all") != "1" {
			t.Errorf("unexpected request %s", r.URL)
		}
		w.Write([]byte(`[
		 {"Id":"b","Names":["/web"],"Image":"nginx","ImageID":"sha256:x","State":"running",
		  "Ports":[{"PrivatePort":80,"PublicPort":8080},{"PrivatePort":443}],
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
