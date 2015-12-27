package main

import (
	"fmt"

	"github.com/codegangsta/cli"

	"golang.org/x/net/context"
	// "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/sigmonsays/voyager/proto/vapi"
)

var PingFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "message, m",
		Value: "ping",
		Usage: "the ping message",
	},
	cli.StringFlag{
		Name:  "host, H",
		Usage: "the host to ping",
	},
}

var PingDescription = `
ping the rpc service
`

func (ac *Application) Ping(c *cli.Context) {
	message := c.String("message")
	host := vapi.HostDefaultPort(c.String("host"), vapi.DefaultPortString)

	fmt.Printf("ping: host %s: message: %s\n", host, message)

	md := metadata.Pairs("secret", ac.Cfg.Rpc.Secret)
	ctx := context.Background()
	ctx = metadata.NewContext(ctx, md)

	dopts := vapi.DefaultDialOptions()

	con, err := vapi.Connect(host, dopts)
	if err != nil {
		log.Warnf("Connect %s: %s", host, err)
	}

	ping := &vapi.PingRequest{
		Message: message,
	}

	res, err := con.Client.Ping(ctx, ping)
	if err != nil {
		log.Warnf("Ping %s: %s", host, err)
	}

	fmt.Printf("%#v\n", res)

}
