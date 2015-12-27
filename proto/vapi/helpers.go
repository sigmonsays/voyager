package vapi

import (
	"fmt"
	"net"
	"time"

	grpc "google.golang.org/grpc"
)

var DefaultPort = 8191

var DefaultPortString = fmt.Sprintf("%d", DefaultPort)

// put port at the end of the host if its not present
func HostDefaultPort(hostport, default_port string) string {
	if default_port == "" {
		default_port = DefaultPortString
	}

	host, port, err := net.SplitHostPort(hostport)
	if err != nil {
		host = hostport
	}
	if port == "" {
		port = default_port
	}
	rpc_host := fmt.Sprintf("%s:%s", host, port)

	return rpc_host

}

func DefaultDialOptions() []grpc.DialOption {
	rpc_timeout := 2
	dopts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithTimeout(time.Duration(rpc_timeout) * time.Second),
	}
	return dopts
}

type Connection struct {
	ClientConn *grpc.ClientConn
	Client     VApiClient
}

func Connect(hostport string, dopts []grpc.DialOption) (*Connection, error) {
	_, _, err := net.SplitHostPort(hostport)
	if err != nil {
		panic(err)
	}

	if dopts == nil {
		rpc_timeout := 2
		dopts = []grpc.DialOption{
			grpc.WithInsecure(),
			grpc.WithTimeout(time.Duration(rpc_timeout) * time.Second),
		}
	}
	conn, err := grpc.Dial(hostport, dopts...)
	if err != nil {
		log.Warnf("rpc dial: %s: %s", hostport, err)
	}

	client := NewVApiClient(conn)

	c := &Connection{
		ClientConn: conn,
		Client:     client,
	}

	return c, err
}
