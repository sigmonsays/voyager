package handler

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/sigmonsays/voyager/filetype"
)

type Handler struct {
	Layout      filetype.FileType
	Username    string
	Homedir     string
	Path        string
	Directories []string
	Filenames   []string
}

func (h *Handler) LocalPath() string {
	offset := len(h.Username) + 2
	return filepath.Join(h.Homedir, h.Path[offset:])
}

func WriteError(w http.ResponseWriter, r *http.Request, s string, args ...interface{}) {
	log.Warnf(s, args...)
	fmt.Fprintf(w, s, args...)
}
