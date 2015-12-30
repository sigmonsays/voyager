package handler

import (
	"net/http"
	"text/template"

	"golang.org/x/net/context"

	"github.com/sigmonsays/voyager/asset"
	"github.com/sigmonsays/voyager/proto/vapi"
	"github.com/sigmonsays/voyager/types"
	"github.com/sigmonsays/voyager/voy"
)

func NewDashboardHandler() *DashboardHandler {
	d := &DashboardHandler{}
	return d
}

type DashboardHandler struct {
	Username string
	VoyFile  voy.VoyLoader
	Ctx      context.Context
}

type DashboardData struct {
	Username   string
	UrlPrefix  string
	ServerData map[string]*ServerData
}
type ServerData struct {
	Server     string
	ServerAddr string
	Voy        *voy.VoyFile
}

func (me *DashboardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Tracef("request %s", r.URL)
	log.Tracef("%#v", me)

	data := &DashboardData{
		Username:   me.Username,
		UrlPrefix:  "/~" + me.Username,
		ServerData: make(map[string]*ServerData, 0),
	}

	// local host server data
	req := &types.ListPathRequest{
		User: me.Username,
	}
	v, err := me.VoyFile.Load(req)
	if err != nil {
		return
	}
	data.ServerData["localhost"] = &ServerData{
		Server:     "localhost",
		ServerAddr: "localhost",
		Voy:        v,
	}

	dopts := vapi.DefaultDialOptions()
	for server, server_addr := range v.Servers {
		hostport := vapi.HostDefaultPort(server_addr, vapi.DefaultPortString)
		conn, err := vapi.Connect(hostport, dopts)
		if err != nil {
			log.Warnf("unable to connect to %s: %s", hostport, err)
			continue
		}

		creq := &vapi.ConfigRequest{
			User: me.Username,
		}

		cres, err := conn.Client.GetConfig(me.Ctx, creq)
		if err != nil {
			log.Warnf("unable to get config to %s: %s", hostport, err)
			continue
		}

		vf := &voy.VoyFile{
			Allow:   cres.Allow,
			Alias:   cres.Alias,
			Servers: cres.Servers,
		}
		data.ServerData[server] = &ServerData{
			Server:     server,
			ServerAddr: server_addr,
			Voy:        vf,
		}

	}

	tmplData, err := asset.Asset("dashboard.html")
	if err != nil {
		WriteError(w, r, "template: %s", err)
		return
	}

	tmpl := template.Must(template.New("dashboard.html").Parse(string(tmplData)))

	err = tmpl.Execute(w, data)
	if err != nil {
		WriteError(w, r, "template render: %s", err)
		return
	}

}
