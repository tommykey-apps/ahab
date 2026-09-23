package docker

import (
	"context"
	"net/http"
	"testing"
)

func TestVolumes_danglingに無いものを使用中にする(t *testing.T) {
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("filters") == `{"dangling":["true"]}` {
			w.Write([]byte(`{"Volumes":[{"Name":"free","Driver":"local","Mountpoint":"/v/free","CreatedAt":"2026-09-24T00:00:00+09:00"}],"Warnings":null}`))
			return
		}
		w.Write([]byte(`{"Volumes":[
		 {"Name":"used","Driver":"local","Mountpoint":"/v/used","CreatedAt":"2026-09-24T00:00:00+09:00"},
		 {"Name":"free","Driver":"local","Mountpoint":"/v/free","CreatedAt":"bad"}
		],"Warnings":null}`))
	}))
	got, err := c.Volumes(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0].Name != "free" || got[1].Name != "used" {
		t.Fatalf("order: %+v", got)
	}
	if got[0].InUse || !got[1].InUse {
		t.Errorf("InUse: free=%v used=%v", got[0].InUse, got[1].InUse)
	}
	if !got[0].Created.IsZero() || got[1].Created.IsZero() {
		t.Errorf("created: %v %v", got[0].Created, got[1].Created)
	}
}

func TestRemoveVolume_名前をパスに入れる(t *testing.T) {
	var got string
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.URL.RequestURI()
		w.WriteHeader(http.StatusNoContent)
	}))
	if err := c.RemoveVolume(context.Background(), "pg data", false); err != nil {
		t.Fatal(err)
	}
	if got != "/"+apiVersion+"/volumes/pg%20data?" {
		t.Errorf("got %q", got)
	}
}
