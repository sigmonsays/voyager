package main

import (
	"github.com/sigmonsays/voyager/config"
	"github.com/sigmonsays/voyager/server"
	"github.com/sigmonsays/voyager/util/devrestarter"

	reload_git "github.com/sigmonsays/git-watch/reload/git"
)

func main() {
	devrestarter.Init()

	gw := reload_git.NewGitWatch(".", "master")
	gw.Interval = 30

	err := gw.Start()
	if err != nil {
		log.Errorf("starting git watch %s", err)
		return
	}

	cfg := config.GetDefaultConfig()
	srv := server.NewServer(cfg.Http.BindAddr)

	log.Infof("%s", cfg.StartupBanner)
	err = srv.Start()
	if err != nil {
		log.Errorf("starting server %s", err)
		return
	}

}
