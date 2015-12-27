package main

import (
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
	"golang.org/x/net/context"

	"github.com/sigmonsays/voyager/config"
	"github.com/sigmonsays/voyager/util"

	gologging "github.com/sigmonsays/go-logging"
)

const AppName = "voy"
const AppVersion = "0.0.1"

type Application struct {
	*cli.App
	Ctx context.Context
	Cfg *config.ApplicationConfig
}

func main() {

	cli_app := cli.NewApp()
	cli_app.Name = AppName
	cli_app.Version = AppVersion

	app := &Application{
		App: cli_app,
	}

	cfg := config.GetDefaultConfig()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "level, l",
			Value: "WARN",
			Usage: "change log level",
		},
	}

	app.Before = func(c *cli.Context) error {
		var err error

		cfgfile := filepath.Join(os.Getenv("HOME"), ".voyager")
		if util.FileExists(cfgfile) {
			err = cfg.LoadYaml(cfgfile)
			if err != nil {
				log.Errorf("load config %s: %s", cfgfile, err)
				return err
			}
		}
		app.Cfg = cfg
		gologging.SetLogLevel(c.String("level"))

		app.Ctx = context.Background()

		return nil
	}

	ping := cli.Command{
		Name:        "ping",
		Usage:       "ping rpc",
		Description: PingDescription,
		Flags:       PingFlags,
		Action:      app.Ping,
	}

	app.Commands = []cli.Command{
		ping,
	}
	app.Run(os.Args)

}
