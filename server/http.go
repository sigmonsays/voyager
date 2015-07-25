package server

import (
	"fmt"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/sigmonsays/go-apachelog"
	"github.com/sigmonsays/voyager/config"
	"github.com/sigmonsays/voyager/voy"
)

type Server struct {
	Addr string

	Conf *config.ApplicationConfig

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

	strip_prefix := fmt.Sprintf("/~%s/", username)
	homedir := user_ent.HomeDir

	voy := voy.DefaultConfig()
	voyfile := filepath.Join(homedir, ".voyager")

	err = voy.LoadYaml(voyfile)
	if err != nil {
		log.Warnf("load voyfile %s: %s", voyfile, err)
		return
	}

	relpath := filepath.Join(tmp[2:]...)

	if voy.Allowed(relpath) == false {
		w.WriteHeader(403)
		fmt.Fprintf(w, "path not allowed")
		return
	}

	log.Infof("user=%s path=%s", username, relpath)

	handler := http.StripPrefix(strip_prefix, http.FileServer(http.Dir(homedir)))
	handler.ServeHTTP(w, r)
}
