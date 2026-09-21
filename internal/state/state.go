package state

import "github.com/tommykey-apps/ahab/internal/docker"

type list = []docker.Container

type Store struct {
	replace   chan list
	queries   chan chan list
	subscribe chan chan list
	unsub     chan chan list
}

func New() *Store {
	s := &Store{
		replace:   make(chan list),
		queries:   make(chan chan list),
		subscribe: make(chan chan list),
		unsub:     make(chan chan list),
	}
	go s.loop()
	return s
}

func (s *Store) loop() {
	var containers list
	subs := map[chan list]bool{}

	for {
		select {
		case cs := <-s.replace:
			containers = cs
			for ch := range subs {
				select {
				case ch <- containers:
				default:
				}
			}
		case reply := <-s.queries:
			reply <- containers
		case ch := <-s.subscribe:
			subs[ch] = true
		case ch := <-s.unsub:
			delete(subs, ch)
			close(ch)
		}
	}
}

func (s *Store) Replace(cs list) { s.replace <- cs }

func (s *Store) Get() list {
	reply := make(chan list)
	s.queries <- reply
	return <-reply
}

func (s *Store) Subscribe() chan list {
	ch := make(chan list, 1)
	s.subscribe <- ch
	return ch
}

func (s *Store) Unsubscribe(ch chan list) { s.unsub <- ch }
