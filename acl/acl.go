package acl

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

type handler struct {
	http.Handler
	nets []*net.IPNet
}

var privateNetworks = []string{
	"127.0.0.0/8",
	"10.0.0.0/8",
	"172.16.0.0/12",
	"192.168.0.0/16",
}

func PrivateNetworks() []string {
	return privateNetworks
}

func parseNetworks(networks []string) ([]*net.IPNet, error) {
	nets := make([]*net.IPNet, 0)
	for _, network := range networks {
		_, net, err := net.ParseCIDR(network)
		if err != nil {
			return nil, err
		}
		nets = append(nets, net)
	}
	return nets, nil
}
func NewHandler(h http.Handler, nets []*net.IPNet) http.Handler {
	aclhandler := &handler{
		Handler: h,
		nets:    nets,
	}
	return aclhandler
}

func NewHandlerWithNetworks(h http.Handler, networks []string) (http.Handler, error) {
	nets, err := parseNetworks(networks)
	if err != nil {
		return nil, fmt.Errorf("parseNetworks: %s", err)
	}
	return NewHandler(h, nets), nil
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	tmp := strings.Split(r.RemoteAddr, ":")
	ip := net.ParseIP(tmp[0])
	allowed := false
	for _, net := range h.nets {
		if net.Contains(ip) {
			allowed = true
			break
		}
	}
	// fmt.Printf("%s access=%v\n", tmp[0], allowed)
	if allowed == false {
		rw.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(rw, "Forbidden")
		return
	}
	h.Handler.ServeHTTP(rw, r)
}
