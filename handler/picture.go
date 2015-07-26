package handler

import (
	"net/http"
)

// provides listing pictures and auto thumbnailing
type PictureHandler struct {
	*Handler
}

func NewPictureHandler(handler *Handler) *PictureHandler {
	return &PictureHandler{
		Handler: handler,
	}
}
func (h *PictureHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}
