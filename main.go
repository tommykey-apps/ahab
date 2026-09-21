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

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	flag.Parse()
	
	ctx := context.Background()
	dc := docker.New()
	store := state.New()

	if _, err := dc.Containers(context.Background()); err != nil {
		log.Fatalf("dockerに接続できません: %v", err)
	}

	refresh := func() {
		if cs, err := dc.Containers(ctx); err == nil {
			store.Replace(cs)
		}
	}
	refresh()

	evs := make(chan docker.Event)

	go func() {
		for {
			if err := dc.Events(ctx, evs); err != nil {
				log.Printf("events: %v(再接続します)", err)
				time.Sleep(time.Second)
			}
		}
	}()

	go func(){
		for range evs {
			refresh()
		}
	}()

	log.Printf("listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, web.NewServer(dc)))
}
