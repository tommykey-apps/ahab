package graph

import (
	"strings"
	"testing"

	"github.com/tommykey-apps/ahab/internal/docker"
)

func TestRender_コンテナが無ければ見出し行だけ(t *testing.T) {
	if got := Render(nil); got != "graph LR\n" {
		t.Fatalf("got %q", got)
	}
}

func TestRender_コンテナ同士は直接結ばずネットワーク経由(t *testing.T) {
	got := Render([]docker.Container{
		{Name: "web", Networks: []string{"front"}},
		{Name: "db", Networks: []string{"front"}},
	})
	for _, want := range []string{
		`  n_front(("front"))`,
		`  c_web["web"]`,
		`  c_web --> n_front`,
		`  c_db --> n_front`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in:\n%s", want, got)
		}
	}
	if strings.Contains(got, "c_web --> c_db") || strings.Contains(got, "c_db --> c_web") {
		t.Errorf("containers must not be linked directly:\n%s", got)
	}
}

func TestRender_公開ポートは_host_から結ぶ(t *testing.T) {
	got := Render([]docker.Container{{Name: "web", Ports: []docker.Port{{Public: 8080, Private: 80}}}})
	for _, want := range []string{`  host{{"host"}}`, `  host -->|":8080→80"| c_web`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in:\n%s", want, got)
		}
	}
	if got := Render([]docker.Container{{Name: "web"}}); strings.Contains(got, "host") {
		t.Errorf("host node must not appear without published ports:\n%s", got)
	}
}

func TestRender_名前付きボリュームは点線で結ぶ(t *testing.T) {
	got := Render([]docker.Container{{Name: "db", Volumes: []string{"pgdata"}}})
	for _, want := range []string{`  v_pgdata[("pgdata")]`, `  c_db -.-> v_pgdata`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in:\n%s", want, got)
		}
	}
}

func TestRender_入力の順序によらず同じ出力(t *testing.T) {
	a := Render([]docker.Container{
		{Name: "web", Networks: []string{"front", "back"}, Volumes: []string{"b", "a"}},
		{Name: "db", Networks: []string{"back"}},
	})
	b := Render([]docker.Container{
		{Name: "db", Networks: []string{"back"}},
		{Name: "web", Networks: []string{"back", "front"}, Volumes: []string{"a", "b"}},
	})
	if a != b {
		t.Errorf("output differs:\n%s\n---\n%s", a, b)
	}
}

func TestNodeID_記号は_アンダースコアに置き換える(t *testing.T) {
	if got := nodeID("c", "my-app.v2/x"); got != "c_my_app_v2_x" {
		t.Fatalf("got %q", got)
	}
}
