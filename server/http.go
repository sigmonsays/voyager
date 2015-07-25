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
	"github.com/sigmonsays/voyager/filetype"
	"github.com/sigmonsays/voyager/handler"
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

	homedir := user_ent.HomeDir

	voy := voy.DefaultConfig()
	voyfile := filepath.Join(homedir, ".voyager")

	err = voy.LoadYaml(voyfile)
	if err != nil {
		WriteError(w, r, "load voyfile %s: %s", voyfile, err)
		return
	}

	relpath := filepath.Join(tmp[2:]...)

	if voy.Allowed(relpath) == false {
		w.WriteHeader(403)
		WriteError(w, r, "path not allowed")
		return
	}

	log.Infof("user=%s path=%s", username, relpath)

	// dispatch handler to appropriate handler based on content

	localpath := filepath.Join(homedir, relpath)

	st, err := os.Stat(localpath)
	if err != nil {
		w.WriteHeader(404)
		WriteError(w, r, "%s", err)
		return
	}

	var filenames []string
	if st.IsDir() {
		f, err := os.Open(localpath)
		if err != nil {
			WriteError(w, r, "open: %s", err)
			return
		}
		defer f.Close()

		filenames, err = f.Readdirnames(-1)
		if err != nil {
			WriteError(w, r, "readdir: %s", err)
			return
		}
	}

	found := make(map[filetype.FileType]int)
	for _, filename := range filenames {
		ftype, err := filetype.Determine(filepath.Join(localpath, filename))
		if err != nil {
			log.Warnf("determine filetype %s: %s", filename, err)
		}
		if _, ok := found[ftype]; ok == false {
			found[ftype] = 0
		}
		found[ftype]++
	}
	var layout filetype.FileType
	var numfiles int
	for ftype, cnt := range found {
		if cnt > numfiles {
			layout = ftype
			cnt = numfiles
		}
	}

	log.Infof("layout type %s", layout)

	handler.NewListHandler(username, homedir, relpath, filenames).ServeHTTP(w, r)

}

func WriteError(w http.ResponseWriter, r *http.Request, s string, args ...interface{}) {
	log.Warnf(s, args...)
	fmt.Fprintf(w, s, args...)
}