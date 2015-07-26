package asset

import (
	"fmt"
	"net/http"
	"strings"
)

//go:generate becky -wrap Blob style.css favicon.ico

type blob struct {
	asset
}

func Blob(a asset) blob {
	return blob{a}
}

var assets = map[string]blob{
	"favicon.ico": favicon,
	"style.css":   style,
}

func Get(path string) (blob, error) {
	a, ok := assets[path]
	if ok == false {
		return a, fmt.Errorf("no such asset: %s", path)
	}

	return a, nil

}
func notFound(w http.ResponseWriter, r *http.Request) {
	http.NotFoundHandler().ServeHTTP(w, r)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s\n", r.URL.Path)
	offset := 1
	if strings.HasPrefix(r.URL.Path, "/s/") {
		offset = 3
	}
	if len(r.URL.Path) < offset {
		return
	}
	path := r.URL.Path[offset:]
	b, err := Get(path)
	if err != nil {
		notFound(w, r)
		return
	}

	b.ServeHTTP(w, r)

}
