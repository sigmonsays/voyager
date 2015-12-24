package handler

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/sigmonsays/voyager/filetype"
)

type Handler struct {
	Layout   filetype.FileType
	Username string

	// ~/username portion
	UrlPrefix string

	// local path
	RootPath string

	Path        string
	Directories []string
	Filenames   []string
}

func (h *Handler) LocalPath() string {
	return filepath.Join(h.RootPath, h.Path)
}
func (h *Handler) Url(paths ...string) string {
	return filepath.Join(h.UrlPrefix, filepath.Join(paths...))
}

func WriteError(w http.ResponseWriter, r *http.Request, s string, args ...interface{}) {
	log.Warnf(s, args...)
	fmt.Fprintf(w, s, args...)
}
