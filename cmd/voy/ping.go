package main

import (
	"fmt"

	"github.com/urfave/cli"

	"golang.org/x/net/context"
	// "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/sigmonsays/voyager/proto/vapi"
)

var PingFlags = []cli.Flag{
	&cli.StringFlag{
		Name:  "message, m",
		Value: "ping",
		Usage: "the ping message",
	},
	&cli.StringFlag{
		Name:  "host, H",
		Usage: "the host to ping",
	},
}

var PingDescription = `
ping the rpc service
`

func (ac *Application) Ping(c *cli.Context) error {
	message := c.String("message")
	host := vapi.HostDefaultPort(c.String("host"), vapi.DefaultPortString)

	fmt.Printf("ping: host %s: message: %s\n", host, message)

	log.Tracef("request-secret:%s", ac.Cfg.Rpc.Secret)

	md := metadata.Pairs("request-secret", ac.Cfg.Rpc.Secret)
	ctx := context.Background()
	ctx = metadata.NewOutgoingContext(ctx, md)

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

	return nil
}
