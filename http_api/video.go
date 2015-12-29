package http_api

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/sigmonsays/voyager/filetype"
)

type VideoServer struct {
	*Server
}

func (s *Server) VideoHandler(w http.ResponseWriter, r *http.Request) {
	tmp := strings.Split(r.URL.Path, "/")

	vs := &VideoServer{s}

	switch tmp[1] {
	case "thumb":
		vs.thumbnail(w, r)
		return
	case "transcode":
		vs.transcode(w, r)
		return
	default:
		w.WriteHeader(500)
	}
}

type request struct {
	width, height int
	localpath     string
}

func (s *VideoServer) parseRequest(w http.ResponseWriter, r *http.Request) (*request, error) {
	q := r.URL.Query()
	width := q.Get("width")
	if width == "" {
		width = "128"
	}
	height := q.Get("height")
	if height == "" {
		height = "128"
	}
	req := &request{}

	if val, err := strconv.Atoi(width); err == nil {
		req.width = int(val)
	} else {
		return nil, fmt.Errorf("width error: %s: %s", width, err)
	}

	if val, err := strconv.Atoi(height); err == nil {
		req.height = int(val)
	} else {
		return nil, fmt.Errorf("width error: %s: %s", width, err)
	}

	req.localpath = q.Get("path")
	if req.localpath == "" {
		return nil, fmt.Errorf("empty path argument")
	}

	ext := filepath.Ext(req.localpath)

	supported, ok := filetype.Video[ext]
	if ok == false || supported == false {
		return nil, fmt.Errorf("format not supported")
	}

	return req, nil
}

func (s *VideoServer) thumbnail(w http.ResponseWriter, r *http.Request) {
}

func (s *VideoServer) transcode(w http.ResponseWriter, r *http.Request) {

	req, err := s.parseRequest(w, r)
	if err != nil {
		s.WriteError(w, r, "request %s", r.URL.Path, err)
		return
	}
	log.Tracef("request %+v", req)

	// 	ffmpeg -i source.mp4 -c:v libx264 -ar 22050 -crf 28 destinationfile.flv

}
