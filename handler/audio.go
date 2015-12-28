package handler

import (
	"net/http"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/sigmonsays/voyager/asset"
	"github.com/sigmonsays/voyager/filetype"
	"github.com/sigmonsays/voyager/types"
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
	Files       []*types.File
	Directories []*types.File
	Breadcrumb  *Breadcrumb
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
		Title:       "Audio",
		Files:       make([]*types.File, 0),
		Directories: make([]*types.File, 0),
		Path:        h.Path,
		LocalPath:   h.LocalPath(),
		Breadcrumb:  NewBreadcrumb(),
	}

	tmp := strings.Split(filepath.Join(h.UrlPrefix, h.Path), "/")
	var crumbUrl string
	for i := 1; i < len(tmp); i++ {
		crumbUrl = strings.Join(tmp[0:i+1], "/") + "/"
		log.Debugf("breadcrumb #%d - %s - %s", i, tmp[i], crumbUrl)
		data.Breadcrumb.Add(crumbUrl, tmp[i])
	}

	for _, f := range h.Files {
		if f.IsDir {
			data.Directories = append(data.Directories, f)
			continue
		}
		if filetype.FileMatch(f, filetype.AudioFile) {
			data.Files = append(data.Files, f)
		}

	}

	log.Tracef("handler path:%s localpath:%s files:%d dirs:%d",
		data.Path, data.LocalPath, len(data.Files), len(data.Directories))

	err = tmpl.Execute(w, data)
	if err != nil {
		WriteError(w, r, "template render: %s", err)
		return
	}

}
