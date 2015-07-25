package main

import (
	"time"

	"github.com/sigmonsays/voyager/config"
	"github.com/sigmonsays/voyager/server"
)

func main() {
	cfg := config.GetDefaultConfig()
	srv := server.NewServer(cfg.Http.BindAddr)
	err := srv.Start()
	if err != nil {
		log.Errorf("starting server %s", err)
		return
	}

	for {
		time.Sleep(time.Second)
	}
}
