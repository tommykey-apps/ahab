package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/tommykey-apps/ahab/internal/docker"
	"github.com/tommykey-apps/ahab/internal/web"
)

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	flag.Parse()

	dc := docker.New()
	if _, err := dc.Containers(context.Background()); err != nil {
		log.Fatalf("dockerに接続できません: %v", err)
	}

	log.Printf("listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, web.NewServer(dc)))
}
