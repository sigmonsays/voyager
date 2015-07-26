package main

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sigmonsays/voyager/config"
	"github.com/sigmonsays/voyager/server"
	"github.com/sigmonsays/voyager/util"
	"github.com/sigmonsays/voyager/util/devrestarter"

	reload_git "github.com/sigmonsays/git-watch/reload/git"
)

func main() {
	var err error
	cfg := config.GetDefaultConfig()

	cfgfile := filepath.Join(os.Getenv("HOME"), ".voyager")
	if util.FileExists(cfgfile) {
		err = cfg.LoadYaml(cfgfile)
		if err != nil {
			log.Errorf("load config %s: %s", cfgfile, err)
			return
		}
	}

	if cfg.AutoRestart {
		devrestarter.Init()
	}

	if cfg.AutoUpgrade {
		gw := reload_git.NewGitWatch(".", "master")
		gw.Interval = 30
		gw.OnChange = func(dir, branch, lhash, rhash string) error {
			cmd := exec.Command("go", "install", "-v", "./...")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				log.Warnf("go install error: %s", err)
			}
			return err
		}

		err := gw.Start()
		if err != nil {
			log.Errorf("starting git watch %s", err)
			return
		}
	}

	srv := server.NewServer(cfg.Http.BindAddr)

	log.Infof("%s", cfg.StartupBanner)
	err = srv.Start()
	if err != nil {
		log.Errorf("starting server %s", err)
		return
	}

}
