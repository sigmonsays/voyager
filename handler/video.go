package handler

import (
	"net/http"
	"strings"
	"text/template"

	"github.com/sigmonsays/voyager/asset"
	"github.com/sigmonsays/voyager/filetype"
	"github.com/sigmonsays/voyager/types"
)

// provides listing Videos and auto thumbnailing
type VideoHandler struct {
	*Handler
}

func NewVideoHandler(handler *Handler) *VideoHandler {
	return &VideoHandler{
		Handler: handler,
	}
}

func (h *VideoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debugf("path:%s localpath:%s", h.Path, h.LocalPath)

	tmplData, err := asset.Asset("video.html")
	if err != nil {
		WriteError(w, r, "template: %s", err)
		return
	}

	tmpl := template.Must(template.New("Videos.html").Parse(string(tmplData)))

	data := &Gallery{
		Title:        "Videos",
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
		if filetype.FileMatch(file, filetype.VideoFile) {
			data.Files = append(data.Files, file)
		}
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		WriteError(w, r, "template render: %s", err)
		return
	}

}
