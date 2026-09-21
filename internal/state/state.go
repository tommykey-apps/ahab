package state

import "github.com/tommykey-apps/ahab/internal/docker"

type View struct {
	docker.Container
	CPU    float64
	Memory uint64
}

type Store struct {
	replace   chan []docker.Container
	stat      chan docker.Stat
	queries   chan chan []View
	subscribe chan chan []View
	unsub     chan chan []View
}

func New() *Store {
	s := &Store{
		replace:   make(chan []docker.Container),
		stat:      make(chan docker.Stat),
		queries:   make(chan chan []View),
		subscribe: make(chan chan []View),
		unsub:     make(chan chan []View),
	}
	go s.loop()
	return s
}

func (s *Store) loop() {
	var containers []docker.Container
	stats := map[string]docker.Stat{}
	subs := map[chan []View]bool{}

	broadcast := func() {
		v := view(containers, stats)
		for ch := range subs {
			select {
			case ch <- v:
			default:
			}
		}
	}

	for {
		select {
		case cs := <-s.replace:
			containers = cs
			broadcast()
		case st := <-s.stat:
			stats[st.ID] = st
			broadcast()
		case reply := <-s.queries:
			reply <- view(containers, stats)
		case ch := <-s.subscribe:
			subs[ch] = true
		case ch := <-s.unsub:
			delete(subs, ch)
			close(ch)
		}
	}
}

func (s *Store) Replace(cs []docker.Container) { s.replace <- cs }

func (s *Store) Get() []View {
	reply := make(chan []View)
	s.queries <- reply
	return <-reply
}

func (s *Store) Subscribe() chan []View {
	ch := make(chan []View, 1)
	s.subscribe <- ch
	return ch
}

func (s *Store) Unsubscribe(ch chan []View) { s.unsub <- ch }

func view(cs []docker.Container, stats map[string]docker.Stat) []View {
	out := make([]View, 0, len(cs))
	for _, c := range cs {
		st := stats[c.ID]
		out = append(out, View{Container: c, CPU: st.CPU, Memory: st.Memory})
	}
	return out
}

func (s *Store) SetStat(st docker.Stat) { s.stat <- st }
