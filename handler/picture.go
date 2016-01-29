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

var tmpl *template.Template

func init() {
	tmplData, err := asset.Asset("picture.html")
	if err != nil {
		log.Errorf("template: %s", err)
	}
	tmpl = template.Must(template.New("pictures").Parse(string(tmplData)))
}

func (h *PictureHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debugf("path:%s localpath:%s", h.Path, h.LocalPath)

	data := &Gallery{
		Title:        "Pictures",
		Files:        make([]*types.File, 0),
		Path:         h.Path,
		LocalPath:    h.LocalPath,
		UrlPrefix:    h.UrlPrefix,
		RelPath:      h.RelPath,
		RemoteServer: h.RemoteServer,
		Breadcrumb:   NewBreadcrumb(),
	}

	log.Tracef("handler %s", h.Handler.Path)

	tmp := strings.Split(h.Path, "/")
	for i := 0; i < len(tmp); i++ {
		data.Breadcrumb.Add(h.Url(strings.Join(tmp[0:i+1], "/")), tmp[i])
	}

	for _, file := range h.Files {
		if file.IsDir {
			data.Directories = append(data.Directories, file)
			continue
		}
		if filetype.FileMatch(file, filetype.PictureFile) {
			data.Files = append(data.Files, file)
		}
	}

	err := tmpl.Execute(w, data)
	if err != nil {
		WriteError(w, r, "template render: %s", err)
		return
	}

}
