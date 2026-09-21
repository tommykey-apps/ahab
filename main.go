package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/tommykey-apps/ahab/internal/docker"
	"github.com/tommykey-apps/ahab/internal/state"
	"github.com/tommykey-apps/ahab/internal/web"
)

type statsRunner struct {
	dc      *docker.Client
	out     chan docker.Stat
	cancels map[string]context.CancelFunc
}

func (r *statsRunner) sync(ctx context.Context, cs []docker.Container) {
	running := map[string]bool{}
	for _, c := range cs {
		if c.State == "running" {
			running[c.ID] = true
		}
	}
	for id, cancel := range r.cancels {
		if !running[id] {
			cancel()
			delete(r.cancels, id)
		}
	}

	for id := range running {
		if _, ok := r.cancels[id]; ok {
			continue
		}
		cctx, cancel := context.WithCancel(ctx)
		r.cancels[id] = cancel
		go func(id string) { _ = r.dc.Stats(cctx, id, r.out) }(id)
	}
}

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	flag.Parse()

	ctx := context.Background()
	dc := docker.New()
	store := state.New()
	statCh := make(chan docker.Stat)
	runner := &statsRunner{dc: dc, out: statCh, cancels: map[string]context.CancelFunc{}}
	evs := make(chan docker.Event)

	if _, err := dc.Containers(context.Background()); err != nil {
		log.Fatalf("dockerに接続できません: %v", err)
	}

	refresh := func() {
		cs, err := dc.Containers(ctx)
		if err != nil {
			log.Printf("containers: %v", err)
			return
		}
		store.Replace(cs)
		runner.sync(ctx, cs)
	}

	go func() {
		for {
			if err := dc.Events(ctx, evs); err != nil {
				log.Printf("events: %v(再接続します)", err)
				time.Sleep(time.Second)
			}
		}
	}()

	go func() {
		refresh()
		for range evs {
			refresh()
		}
	}()

	go func() {
		for st := range statCh {
			store.SetStat(st)
		}
	}()

	log.Printf("listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, web.NewServer(dc, store)))
}
