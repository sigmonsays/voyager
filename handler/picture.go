package handler

import (
	"net/http"
	"path/filepath"

	"github.com/sigmonsays/voyager/asset"
	"github.com/sigmonsays/voyager/filetype"
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
	Files       []*File
	Directories []*File
}

type File struct {
	Url  string
	Name string
}

func (f *File) Basename() string {
	return filepath.Base(f.Name)
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
	for _, dirname := range h.Directories {
		f := &File{
			Url:  filepath.Join(h.Path, dirname),
			Name: dirname,
		}
		data.Directories = append(data.Directories, f)
	}

	for _, filename := range filetype.Filter(h.Filenames, filetype.PictureFile) {
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
