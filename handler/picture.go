package handler

import (
	"net/http"
	"path/filepath"

	"github.com/sigmonsays/voyager/asset"
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
	Path      string
	LocalPath string
	Title     string
	Files     []*File
}

type File struct {
	Url  string
	Name string
}

func (h *PictureHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tmpl, err := asset.GetTemplate("picture.html")
	if err != nil {
		WriteError(w, r, "template: %s", err)
		return
	}
	log.Tracef("handler %#v", h.Handler)
	data := &Gallery{
		Title:     "Pictures",
		Files:     make([]*File, 0),
		Path:      h.Path,
		LocalPath: h.LocalPath(),
	}
	for _, filename := range h.Filenames {
		f := &File{
			Url:  filepath.Join(h.Path, filename),
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
