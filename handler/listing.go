package handler

import (
	"net/http"
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
	strip_prefix := h.Handler.UrlPrefix
	log.Debugf("serve path:%s (stripPrefix:%s rootPath:%s)",
		r.URL.Path, strip_prefix, h.RootPath)
	handler := http.StripPrefix(strip_prefix, http.FileServer(http.Dir(h.RootPath)))
	handler.ServeHTTP(w, r)
}
