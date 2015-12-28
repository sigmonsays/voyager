package http_api

import (
	"net/http"
)

func (s *Server) CacheHandler(w http.ResponseWriter, r *http.Request) {
	handler := http.StripPrefix("/c/", http.FileServer(http.Dir(s.Conf.CacheDir)))
	handler.ServeHTTP(w, r)
}
