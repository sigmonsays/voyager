package handler

import (
	"net/http"

	"github.com/sigmonsays/voyager/filetype"
)

func NewHandlerFactory() *HandlerFactory {
	f := &HandlerFactory{}
	return f
}

type HandlerFactory struct {
}

type HandlerFunc func(handler *Handler) http.Handler

var Handlers map[filetype.FileType]HandlerFunc

func init() {
	Handlers = make(map[filetype.FileType]HandlerFunc, 0)

	// Handlers[UnknownFile ] = ...
	// Handlers[filetype.PictureFile] = func(h *Handler) ContentHandler { return NewPictureHandler(h) }

	Handlers[filetype.PictureFile] = func(h *Handler) http.Handler { return NewPictureHandler(h) }
	// VideoFile
	Handlers[filetype.AudioFile] = func(h *Handler) http.Handler { return NewAudioHandler(h) }
	Handlers[filetype.VideoFile] = func(h *Handler) http.Handler { return NewVideoHandler(h) }

	// browsable listing
	Handlers[filetype.ListFile] = func(h *Handler) http.Handler { return NewListHandler(h) }
}

func (h *HandlerFactory) MakeHandler(handler *Handler) http.Handler {

	newHandler, ok := Handlers[handler.Layout]
	if ok == false {
		log.Debugf("no handler for filetype %s", handler.Layout)
		return nil
	}

	return newHandler(handler)

}
