package handler

import (
	"fmt"
	"net/http"
)

// provides most basic file listing when no other handler has been detected
type ListHandler struct {
	username  string
	homedir   string
	path      string
	filenames []string
}

func NewListHandler(username, homedir, path string, filenames []string) *ListHandler {
	return &ListHandler{
		username:  username,
		homedir:   homedir,
		path:      path,
		filenames: filenames,
	}
}
func (h *ListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	strip_prefix := fmt.Sprintf("/~%s/", h.username)
	handler := http.StripPrefix(strip_prefix, http.FileServer(http.Dir(h.homedir)))
	handler.ServeHTTP(w, r)
}
