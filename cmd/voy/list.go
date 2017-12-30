package main

import (
	"fmt"

	"github.com/codegangsta/cli"

	"golang.org/x/net/context"
	// "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/sigmonsays/voyager/proto/vapi"
)

var ListFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "host, H",
		Usage: "the host to list",
	},
}

var ListDescription = `
list contents
`

func (ac *Application) List(c *cli.Context) {
	message := c.String("host")
	host := vapi.HostDefaultPort(c.String("host"), vapi.DefaultPortString)

	md := metadata.Pairs("request-secret", ac.Cfg.Rpc.Secret)
	ctx := context.Background()
	ctx = metadata.NewOutgoingContext(ctx, md)

	dopts := vapi.DefaultDialOptions()

	con, err := vapi.Connect(host, dopts)
	if err != nil {
		log.Warnf("Connect %s: %s", host, err)
	}

	// todo
	ping := &vapi.PingRequest{
		Message: message,
	}

	res, err := con.Client.Ping(ctx, ping)
	if err != nil {
		log.Warnf("Ping %s: %s", host, err)
	}

	fmt.Printf("%#v\n", res)

}
