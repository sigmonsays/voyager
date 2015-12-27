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
	Path        string
	LocalPath   string
	Title       string
	Files       []*types.File
	Directories []*types.File
	Breadcrumb  *Breadcrumb
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
		Title:      "Pictures",
		Files:      make([]*types.File, 0),
		Path:       h.Path,
		LocalPath:  h.LocalPath(),
		Breadcrumb: NewBreadcrumb(),
	}

	log.Tracef("handler %#v data %+v", h.Handler, data)

	tmp := strings.Split(h.Path, "/")
	for i := 0; i < len(tmp); i++ {
		data.Breadcrumb.Add(h.Url(strings.Join(tmp[0:i+1], "/")), tmp[i])
	}

	for _, dirname := range h.Directories {
		f := &types.File{
			Url:  h.Url(h.Path, dirname),
			Name: dirname,
		}
		data.Directories = append(data.Directories, f)
	}

	for _, filename := range filetype.Filter(h.Filenames, filetype.PictureFile) {
		f := &types.File{
			// Url:  h.Url(h.Path, filename),
			Url:  h.Url(h.Path, filename),
			Name: filename,
		}
		data.Files = append(data.Files, f)
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		WriteError(w, r, "template render: %s", err)
		return
	}

}
