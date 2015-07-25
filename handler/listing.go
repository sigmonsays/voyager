package handler

import (
	"fmt"
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
	strip_prefix := fmt.Sprintf("/~%s/", h.Username)
	handler := http.StripPrefix(strip_prefix, http.FileServer(http.Dir(h.Homedir)))
	handler.ServeHTTP(w, r)
}
