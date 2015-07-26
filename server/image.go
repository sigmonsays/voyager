package server

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"image"
	"image/gif"
	"image/jpeg"
	"image/png"

	"github.com/nfnt/resize"
)

type DecodeFunc func(f io.Reader) (image.Image, error)

var DecodeTypes = map[string]DecodeFunc{
	"jpg":  jpeg.Decode,
	"jpeg": jpeg.Decode,
	"gif":  gif.Decode,
	"png":  png.Decode,
}

func (s *Server) ImageHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	width_s := q.Get("width")
	if width_s == "" {
		width_s = "128"
	}
	var width uint

	if val, err := strconv.Atoi(width_s); err == nil {
		width = uint(val)
	} else {
		w.WriteHeader(400)
		WriteError(w, r, "width error: %s: %s", width_s, err)
		return
	}

	localpath := q.Get("path")
	if localpath == "" {
		w.WriteHeader(400)
		WriteError(w, r, "path required. pass path=")
		return
	}

	ext := filepath.Ext(localpath)
	if ext != "" {
		ext = strings.ToLower(ext[1:])
	}

	decoder, ok := DecodeTypes[ext]

	if ok == false {
		w.WriteHeader(400)
		WriteError(w, r, "Unable to open file format ext=%s path=%s", ext, localpath)
		return
	}

	file, err := os.Open(localpath)
	if err != nil {
		w.WriteHeader(400)
		WriteError(w, r, "Unable to open file ext=%s path=%s: %s", ext, localpath, err)
		return
	}
	defer file.Close()

	img, err := decoder(file)
	if err != nil {
		w.WriteHeader(400)
		WriteError(w, r, "Unable to decode file ext=%s path=%s: %s", ext, localpath, err)
		return
	}

	m := resize.Resize(width, 0, img, resize.Lanczos3)

	opts := &jpeg.Options{
		Quality: 50,
	}
	w.Header().Add("Content-Type", "image/jpeg")
	err = jpeg.Encode(w, m, opts)
	if err != nil {
		w.WriteHeader(400)
		WriteError(w, r, "Unable to encode file ext=%s path=%s: %s", ext, localpath, err)
		return
	}

}
