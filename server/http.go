package server

import (
	"fmt"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/sigmonsays/go-apachelog"

	"github.com/sigmonsays/voyager/asset"
	"github.com/sigmonsays/voyager/cache"
	"github.com/sigmonsays/voyager/config"
	"github.com/sigmonsays/voyager/filetype"
	"github.com/sigmonsays/voyager/handler"
	"github.com/sigmonsays/voyager/types"
	"github.com/sigmonsays/voyager/voy"
)

type Server struct {
	Addr string
	Conf *config.ApplicationConfig

	Cache      *cache.FileCache
	Factory    *handler.HandlerFactory
	PathLoader handler.PathLoader

	srv *http.Server
}

func NewServer(addr string) *Server {
	mux := http.NewServeMux()
	hndlr := apachelog.NewHandler(mux, os.Stderr)
	srv := &http.Server{
		Addr:    addr,
		Handler: hndlr,
	}

	s := &Server{
		Addr: addr,
		srv:  srv,
	}

	mux.Handle("/", s)

	mux.HandleFunc("/image/", s.ImageHandler)
	mux.HandleFunc("/c/", s.CacheHandler)

	static := http.FileServer(
		&assetfs.AssetFS{
			Asset:     asset.Asset,
			AssetDir:  asset.AssetDir,
			AssetInfo: asset.AssetInfo,
			Prefix:    ""},
	)
	mux.Handle("/s/", http.StripPrefix("/s", static))
	mux.Handle("/favicon.ico", static)

	return s
}

func (s *Server) Start() error {
	log.Infof("starting server")
	return s.srv.ListenAndServe()
}

// parse a url path and turn it into a list path request
func (s *Server) parsePath(path string) (*types.ListPathRequest, error) {
	tmp := strings.Split(path, "/")
	username := ""
	if strings.HasPrefix(tmp[1], "~") {
		username = tmp[1][1:]
	}
	if username == "" {
		return nil, fmt.Errorf("empty username")
	}

	res := &types.ListPathRequest{
		User: username,
		Path: "/" + strings.Join(tmp[2:], "/"),
	}

	log.Tracef("user=%s: path %s return %+v", username, path, res)

	return res, nil

}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	tmp := strings.Split(path, "/")

	req, err := s.parsePath(path)
	if err != nil {
		WriteError(w, r, "parse path %s: %s", path, err)
		return
	}

	user_ent, err := user.Lookup(req.User)
	if err != nil {
		WriteError(w, r, "user lookup %s: %s", req.User, err)
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

	if len(tmp) < 3 {
		// TODO: Do we want to support any kind of top level index?
		WriteError(w, r, "incomplete path")
		return
	}

	var localpath string
	var rootpath string
	var relpath string
	var urlprefix string

	topdir := tmp[2]
	alias, is_alias := voy.Alias[topdir]

	if is_alias {
		localpath = filepath.Join(alias, strings.Join(tmp[3:], "/"))
		rootpath = alias
		relpath, err = filepath.Rel(rootpath, localpath)
		if err != nil {
			log.Warnf("relpath %s", err)
		}
		urlprefix = "/~" + filepath.Join(req.User, topdir)

		log.Debugf("%s is an alias for %s: new path %s (relpath:%s urlprefix:%s)", topdir, alias, localpath, relpath, urlprefix)
	} else {
		rootpath = homedir
		localpath = filepath.Join(homedir, strings.Join(tmp[2:], "/"))
		relpath, err = filepath.Rel(rootpath, localpath)
		if err != nil {
			log.Warnf("relpath rootpath:%s localpath:%s : %s", rootpath, localpath, err)
		}
		if voy.Allowed(relpath) == false {
			log.Warnf("rootpath:%s localpath:%s relpath:%s not allowed", rootpath, localpath, relpath)
			w.WriteHeader(403)
			WriteError(w, r, "nothing to see here. bye bye.")
			return
		}
		urlprefix = "/~" + req.User
	}

	log.Infof("request user:%s rootpath:%s path:%s localpath:%s urlprefix:%s",
		req.User, rootpath, relpath, localpath, urlprefix)

	// dispatch handler to appropriate handler based on content
	hndlr := &handler.Handler{
		Username:  req.User,
		RootPath:  rootpath,
		Path:      relpath,
		UrlPrefix: urlprefix,
	}

	// call the path loader
	files, err := s.PathLoader.GetFiles(localpath)
	if err != nil {
		WriteError(w, r, "path loader GetFiles: %s", err)
		return
	}

	// resolve layout
	var customLayout string
	ltmp := strings.Split(localpath, "/")
Layout:
	for i := len(ltmp); i > 1; i-- {
		p := strings.Join(ltmp[:i], "/")
		log.Tracef("check custom layout %s", p)
		l, found := voy.Layouts[p]
		if found {
			log.Tracef("found custom layout %s for %s", l, p)
			customLayout = l
			break Layout
		}
	}

	if customLayout == "" {
		hndlr.Layout = filetype.GuessLayout(localpath, files)
	} else {
		hndlr.Layout = filetype.TypeFromString(customLayout)
		log.Debugf("using custom layout %s (%q) for %s", hndlr.Layout, customLayout, localpath)
	}

	hndlr.Files = files

	log.Debugf("dispatch %s user:%s rootpath:%s path:%s localpath:%s urlprefix:%s files:%d",
		hndlr.Layout, req.User, rootpath, relpath, localpath, urlprefix, len(files))

	reqHandler := s.Factory.MakeHandler(hndlr)
	if reqHandler == nil {
		handler.NewListHandler(hndlr).ServeHTTP(w, r)
	} else {
		reqHandler.ServeHTTP(w, r)
	}

}

func WriteError(w http.ResponseWriter, r *http.Request, s string, args ...interface{}) {
	log.Warnf(s, args...)
	fmt.Fprintf(w, s, args...)
}
