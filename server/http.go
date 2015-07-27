package server

import (
	"fmt"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strings"

	"github.com/sigmonsays/go-apachelog"

	"github.com/sigmonsays/voyager/asset"
	"github.com/sigmonsays/voyager/cache"
	"github.com/sigmonsays/voyager/config"
	"github.com/sigmonsays/voyager/filetype"
	"github.com/sigmonsays/voyager/handler"
	"github.com/sigmonsays/voyager/voy"
)

type Server struct {
	Addr string
	Conf *config.ApplicationConfig

	Cache *cache.FileCache

	srv *http.Server
}

func NewServer(addr string) *Server {
	mux := http.NewServeMux()
	handler := apachelog.NewHandler(mux, os.Stderr)

	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	s := &Server{
		Addr: addr,
		srv:  srv,
	}

	mux.Handle("/", s)

	mux.HandleFunc("/favicon.ico", asset.Handler)
	mux.HandleFunc("/s/", asset.Handler)
	mux.HandleFunc("/image/", s.ImageHandler)
	mux.HandleFunc("/c/", s.CacheHandler)
	return s
}

func (s *Server) Start() error {
	log.Infof("starting server")
	return s.srv.ListenAndServe()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	tmp := strings.Split(path, "/")
	username := ""
	if strings.HasPrefix(tmp[1], "~") {
		username = tmp[1][1:]
	}
	if username == "" {
		return
	}
	log.Tracef("user=%s: path %s", username, path)
	user_ent, err := user.Lookup(username)
	if err != nil {
		log.Warnf("user lookup %s: %s", username, err)
		return
	}

	homedir := user_ent.HomeDir

	voy := voy.DefaultConfig()
	voyfile := filepath.Join(homedir, ".voyager")

	err = voy.LoadYaml(voyfile)
	if err != nil {
		WriteError(w, r, "load voyfile %s: %s", voyfile, err)
		return
	}

	relpath, err := filepath.Rel("/~"+username+"/", path)
	if err != nil {
		log.Warnf("relpath %s", err)
	}

	if voy.Allowed(relpath) == false {
		w.WriteHeader(403)
		WriteError(w, r, "bye bye")
		return
	}

	log.Infof("user=%s path=%s", username, relpath)

	// dispatch handler to appropriate handler based on content
	hndlr := &handler.Handler{
		Username:  username,
		Homedir:   homedir,
		Path:      relpath,
		UrlPrefix: "/~" + username,
	}

	localpath := filepath.Join(homedir, relpath)

	fh, err := os.Open(localpath)
	if err != nil {
		w.WriteHeader(404)
		WriteError(w, r, "%s", err)
		return
	}
	defer fh.Close()
	st, err := fh.Stat()
	if err != nil {
		w.WriteHeader(404)
		WriteError(w, r, "%s", err)
		return
	}

	if st.IsDir() == false {
		// serve the object directly
		handler.NewListHandler(hndlr).ServeHTTP(w, r)
		return
	}

	directories := make([]string, 0)
	filenames := make([]string, 0)

	files, err := fh.Readdir(-1)
	if err != nil {
		WriteError(w, r, "readdir: %s", err)
		return
	}
	for _, file := range files {
		if file.IsDir() {
			directories = append(directories, file.Name())
		} else {
			filenames = append(filenames, file.Name())
		}
	}
	sort.Strings(directories)
	sort.Strings(filenames)

	hndlr.Layout = filetype.GuessLayout(localpath, filenames)
	hndlr.Filenames = filenames
	hndlr.Directories = directories

	log.Infof("path %s %d files, %d dirs, layout %s", localpath, len(filenames), len(directories), hndlr.Layout)

	if hndlr.Layout == filetype.PictureFile {
		handler.NewPictureHandler(hndlr).ServeHTTP(w, r)
	} else {
		handler.NewListHandler(hndlr).ServeHTTP(w, r)
	}

}

func WriteError(w http.ResponseWriter, r *http.Request, s string, args ...interface{}) {
	log.Warnf(s, args...)
	fmt.Fprintf(w, s, args...)
}
