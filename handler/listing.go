package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/sigmonsays/voyager/asset"
	"github.com/sigmonsays/voyager/types"
)

// provides listing Lists and auto thumbnailing
type ListHandler struct {
	*Handler
}

func NewListHandler(handler *Handler) *ListHandler {
	return &ListHandler{
		Handler: handler,
	}
}

func (h *ListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debugf("path:%s localpath:%s", h.Path, h.LocalPath)

	q := r.URL.Query()
	format := q.Get("format")

	if format == "json" {
		h.jsonList(w, r)
		return
	}

	h.templateList(w, r)
}

func (h *ListHandler) jsonList(w http.ResponseWriter, r *http.Request) {

	log.Debugf("jsonList path:%s localpath:%s", h.Path, h.LocalPath)
	w.Header().Set("Content-Type", "application/json")

	data := h.buildGallery()

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		WriteError(w, r, "encode: %s", err)
		return
	}
}

func (h *ListHandler) templateList(w http.ResponseWriter, r *http.Request) {
	tmplData, err := asset.Asset("list.html")
	if err != nil {
		WriteError(w, r, "template: %s", err)
		return
	}

	st, err := os.Stat(h.LocalPath)
	if err == nil && st.IsDir() == false {
		dirname := filepath.Dir(h.LocalPath)
		prefix := filepath.Dir(r.URL.Path)
		log.Debugf("Serving static file path %s: dirname:%s prefix:%s localpath:%s", r.URL.Path, dirname, prefix, h.LocalPath)
		http.StripPrefix(prefix, http.FileServer(http.Dir(dirname))).ServeHTTP(w, r)
		return
	}

	tmpl := template.Must(template.New("list.html").Parse(string(tmplData)))

	log.Tracef("handler %s", h.Handler.Path)

	data := h.buildGallery()

	err = tmpl.Execute(w, data)
	if err != nil {
		WriteError(w, r, "template render: %s", err)
		return
	}

}

func (h *ListHandler) buildGallery() *Gallery {
	data := &Gallery{
		Title:        "Lists",
		Files:        make([]*types.File, 0),
		Path:         h.Path,
		LocalPath:    h.LocalPath,
		UrlPrefix:    h.UrlPrefix,
		RelPath:      h.RelPath,
		RemoteServer: h.RemoteServer,
		Breadcrumb:   NewBreadcrumb(),
	}

	tmp := strings.Split(h.Path, "/")
	for i := 0; i < len(tmp); i++ {
		data.Breadcrumb.Add(h.Url(strings.Join(tmp[0:i+1], "/")), tmp[i])
	}

	for _, file := range h.Files {
		if file.IsDir {
			data.Directories = append(data.Directories, file)
		} else {
			data.Files = append(data.Files, file)
		}
	}
	return data
}
