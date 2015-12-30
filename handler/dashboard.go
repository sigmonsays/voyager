package handler

import (
	"net/http"
	"text/template"

	"github.com/sigmonsays/voyager/asset"
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
}

type DashboardData struct {
	Username  string
	UrlPrefix string
	Voy       *voy.VoyFile
}

func (me *DashboardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Tracef("request %s", r.URL)
	log.Tracef("%#v", me)

	req := &types.ListPathRequest{
		User: me.Username,
	}
	v, err := me.VoyFile.Load(req)
	if err != nil {
		return
	}

	tmplData, err := asset.Asset("dashboard.html")
	if err != nil {
		WriteError(w, r, "template: %s", err)
		return
	}

	data := &DashboardData{
		Username:  me.Username,
		UrlPrefix: "/~" + req.User,
		Voy:       v,
	}

	tmpl := template.Must(template.New("dashboard.html").Parse(string(tmplData)))

	err = tmpl.Execute(w, data)
	if err != nil {
		WriteError(w, r, "template render: %s", err)
		return
	}

}
