package handler

import (
	"net/http"
	"strings"
)

// provides most basic file listing when no other handler has been detected
type ListHandler struct {
	*Handler
}

func NewListHandler(handler *Handler) *ListHandler {
	return &ListHandler{
		Handler: handler,
	}
}
func (h *ListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	tmp := strings.Split(r.URL.Path, "/")
	strip_prefix := tmp[1]

	log.Debugf("serve path:%s (stripPrefix:%s relpath:%s urlPrefix:%s localPath:%s rootPath:%s)",
		r.URL.Path, strip_prefix, h.RelPath, h.UrlPrefix, h.LocalPath, h.RootPath)

	handler := http.StripPrefix(strip_prefix, http.FileServer(http.Dir(h.LocalPath)))

	handler.ServeHTTP(w, r)
}
