package state

import (
	"testing"
	"time"

	"github.com/tommykey-apps/ahab/internal/docker"
)

func TestStore_統計はコンテナIDで結合される(t *testing.T) {
	s := New()
	s.Replace([]docker.Container{{ID: "a", Name: "web"}, {ID: "b", Name: "db"}})
	s.SetStat(docker.Stat{ID: "b", CPU: 12.5, Memory: 1024})
	got := s.Get()
	if len(got) != 2 || got[0].CPU != 0 || got[1].CPU != 12.5 || got[1].Memory != 1024 {
		t.Fatalf("got %+v", got)
	}
}

func TestStore_購読者に更新が届く(t *testing.T) {
	s := New()
	ch := s.Subscribe()
	defer s.Unsubscribe(ch)
	s.Replace([]docker.Container{{ID: "a", Name: "web"}})
	select {
	case v := <-ch:
		if len(v) != 1 || v[0].Name != "web" {
			t.Fatalf("got %+v", v)
		}
	case <-time.After(time.Second):
		t.Fatal("no broadcast")
	}
}

func TestStore_遅い購読者がいても止まらない(t *testing.T) {
	s := New()
	ch := s.Subscribe()
	defer s.Unsubscribe(ch)
	for i := 0; i < 10; i++ {
		s.Replace([]docker.Container{{ID: "a"}})
	}
	if got := s.Get(); len(got) != 1 {
		t.Fatalf("got %+v", got)
	}
}
