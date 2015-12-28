package http_api

import (
	"fmt"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/sigmonsays/go-apachelog"
	"golang.org/x/net/context"

	"github.com/sigmonsays/voyager/asset"
	"github.com/sigmonsays/voyager/cache"
	"github.com/sigmonsays/voyager/config"
	"github.com/sigmonsays/voyager/handler"
	"github.com/sigmonsays/voyager/layout"
	"github.com/sigmonsays/voyager/proto/vapi"
	"github.com/sigmonsays/voyager/types"
	"github.com/sigmonsays/voyager/voy"
)

type Server struct {
	Addr string
	Conf *config.ApplicationConfig

	Ctx        context.Context
	Cache      *cache.FileCache
	Factory    *handler.HandlerFactory
	PathLoader handler.PathLoader
	Layout     layout.LayoutResolver
	VoyFile    voy.VoyLoader

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
func (s *Server) parseRequest(path string) (*types.ListPathRequest, error) {
	tmp := strings.Split(path, "/")

	identity := ""
	if strings.HasPrefix(tmp[1], "~") {
		identity = tmp[1][1:]
	}
	if identity == "" {
		return nil, fmt.Errorf("empty identity")
	}

	var username string
	var server string

	idx := strings.Index(identity, "@")
	if idx >= 0 {
		username = identity[:idx]
		server = identity[idx+1:]
		server = vapi.HostDefaultPort(server, vapi.DefaultPortString)
	} else {
		username = identity
	}

	res := &types.ListPathRequest{
		Server: server,
		User:   username,
		Path:   "/" + strings.Join(tmp[2:], "/"),
	}

	log.Tracef("user=%s: path %s return %+v", username, path, res)

	return res, nil

}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	req, err := s.parseRequest(path)
	if err != nil {
		WriteError(w, r, "parse path %s: %s", path, err)
		return
	}

	if req.Server == "" {
		s.LocalRequest(w, r, req)
	} else {
		s.RemoteRequest(w, r, req)
	}

}

func (s *Server) RemoteRequest(w http.ResponseWriter, r *http.Request, req *types.ListPathRequest) {

	// dispatch handler to appropriate handler based on content

	log.Tracef("remote request server:%s", req.Server)

	dopts := vapi.DefaultDialOptions()

	c, err := vapi.Connect(req.Server, dopts)
	if err != nil {
		WriteError(w, r, "remote:%s connect %s: %s", req.Server, req.Path, err)
		return
	}

	list_req := &vapi.ListRequest{
		User: req.User,
		Path: req.Path,
	}

	res, err := c.Client.ListFiles(s.Ctx, list_req)
	if err != nil {
		WriteError(w, r, "remote:%s list files %s: %s", req.Server, req.Path, err)
		return
	}
	log.Debugf("response %+v", res)
}

func (s *Server) LocalRequest(w http.ResponseWriter, r *http.Request, req *types.ListPathRequest) {

	voy, err := s.VoyFile.Load(req)
	if err != nil {
		WriteError(w, r, "load voyfile %s: %s", req.User, err)
		return
	}

	log.Tracef("voyfile allow:%+v alias:%+v servers:%+v", voy.Allow, voy.Alias, voy.Servers)

	// we still need to lookup the user to get the home directory
	// because paths that do not start with a / are relative to the home directory
	user_ent, err := user.Lookup(req.User)
	if err != nil {
		WriteError(w, r, "user lookup %s: %s", req.User, err)
		return
	}
	homedir := user_ent.HomeDir

	tmp := strings.Split(req.Path, "/")
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
		localpath = filepath.Join(homedir, strings.Join(tmp[1:], "/"))
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
		log.Debugf("%s is regular path (localpath:%s relpath:%s urlprefix:%s)", req.Path, localpath, relpath, urlprefix)
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

	// determine the layout
	layout, err := s.Layout.Resolve(voy, localpath, files)
	if err != nil {
		WriteError(w, r, "resolve layout %s: %s", localpath, err)
		return
	}

	hndlr.Layout = layout
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
