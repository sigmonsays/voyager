package handler

import (
	"net/http"
	"strings"
	"text/template"

	"github.com/sigmonsays/voyager/asset"
	"github.com/sigmonsays/voyager/filetype"
)

type AudioHandler struct {
	*Handler
}

func NewAudioHandler(handler *Handler) *AudioHandler {
	return &AudioHandler{
		Handler: handler,
	}
}

type Playlist struct {
	Path        string
	LocalPath   string
	Title       string
	Files       []*File
	Directories []*File
	Breadcrum   []*Breadcrum
}

func (h *AudioHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debugf("path:%s localpath:%s", h.Path, h.LocalPath())

	tmplData, err := asset.Asset("audio.html")
	if err != nil {
		WriteError(w, r, "template: %s", err)
		return
	}

	tmpl := template.Must(template.New("audio.html").Parse(string(tmplData)))

	data := &Playlist{
		Title:     "Audio",
		Files:     make([]*File, 0),
		Path:      h.Path,
		LocalPath: h.LocalPath(),
		Breadcrum: make([]*Breadcrum, 0),
	}

	log.Tracef("handler %#v data %+v", h.Handler, data)

	tmp := strings.Split(h.Path, "/")
	for i := 0; i < len(tmp); i++ {
		b := &Breadcrum{
			Url:  h.Url(strings.Join(tmp[0:i+1], "/")),
			Name: tmp[i],
		}
		data.Breadcrum = append(data.Breadcrum, b)
	}

	for _, dirname := range h.Directories {
		f := &File{
			Url:  h.Url(h.Path, dirname),
			Name: dirname,
		}
		data.Directories = append(data.Directories, f)
	}

	for _, filename := range filetype.Filter(h.Filenames, filetype.AudioFile) {
		f := &File{
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
