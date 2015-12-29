package handler

import (
	"net/http"
	"strings"
	"text/template"

	"github.com/sigmonsays/voyager/asset"
	"github.com/sigmonsays/voyager/filetype"
	"github.com/sigmonsays/voyager/types"
)

// provides listing pictures and auto thumbnailing
type PictureHandler struct {
	*Handler
}

func NewPictureHandler(handler *Handler) *PictureHandler {
	return &PictureHandler{
		Handler: handler,
	}
}

type Gallery struct {
	Path         string
	LocalPath    string
	UrlPrefix    string
	RelPath      string
	RemoteServer string
	Title        string
	Files        []*types.File
	Directories  []*types.File
	Breadcrumb   *Breadcrumb
}

func (h *PictureHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debugf("path:%s localpath:%s", h.Path, h.LocalPath())

	tmplData, err := asset.Asset("picture.html")
	if err != nil {
		WriteError(w, r, "template: %s", err)
		return
	}

	tmpl := template.Must(template.New("pictures.html").Parse(string(tmplData)))

	data := &Gallery{
		Title:        "Pictures",
		Files:        make([]*types.File, 0),
		Path:         h.Path,
		LocalPath:    h.LocalPath(),
		UrlPrefix:    h.UrlPrefix,
		RelPath:      h.RelPath,
		RemoteServer: h.RemoteServer,
		Breadcrumb:   NewBreadcrumb(),
	}

	log.Tracef("handler %#v data %+v", h.Handler, data)

	tmp := strings.Split(h.Path, "/")
	for i := 0; i < len(tmp); i++ {
		data.Breadcrumb.Add(h.Url(strings.Join(tmp[0:i+1], "/")), tmp[i])
	}

	for _, file := range h.Files {
		if file.IsDir == false {
			continue
		}
		data.Directories = append(data.Directories, file)
	}

	for _, file := range filetype.Filter(h.Files, filetype.PictureFile) {
		if file.IsDir {
			continue
		}
		data.Files = append(data.Files, file)
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		WriteError(w, r, "template render: %s", err)
		return
	}

}
