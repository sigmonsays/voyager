package handler

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/sigmonsays/voyager/filetype"
	"github.com/sigmonsays/voyager/types"
)

type ContentHandler interface {
	ListPath(req types.ListPathRequest) (*types.ListPathResponse, error)
}

type Handler struct {
	Layout   filetype.FileType
	Username string

	// ~/username portion
	UrlPrefix string

	// local path
	RootPath string

	Path  string
	Files []*types.File
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
