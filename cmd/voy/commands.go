package main

import "github.com/codegangsta/cli"

func GetCommands(app *Application) []cli.Command {
	ping := cli.Command{
		Name:        "ping",
		Usage:       "ping rpc",
		Description: PingDescription,
		Flags:       PingFlags,
		Action:      app.Ping,
	}
	list := cli.Command{
		Name:        "list",
		Usage:       "list rpc",
		Description: ListDescription,
		Flags:       ListFlags,
		Action:      app.List,
	}
	ret := []cli.Command{
		ping,
		list,
	}
	return ret
}
