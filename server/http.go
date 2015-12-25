package server

import (
	"fmt"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strings"

	"github.com/elazarl/go-bindata-assetfs"
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

	Cache   *cache.FileCache
	Factory *handler.HandlerFactory

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
		Addr:    addr,
		srv:     srv,
		Factory: handler.NewHandlerFactory(),
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
		urlprefix = "/~" + filepath.Join(username, topdir)

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
		urlprefix = "/~" + username
	}

	log.Infof("request user:%s rootpath:%s path:%s localpath:%s urlprefix:%s",
		username, rootpath, relpath, localpath, urlprefix)

	// dispatch handler to appropriate handler based on content
	hndlr := &handler.Handler{
		Username:  username,
		RootPath:  rootpath,
		Path:      relpath,
		UrlPrefix: urlprefix,
	}

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
		log.Debugf("dispatch ListHandler user:%s rootpath:%s path:%s localpath:%s urlprefix:%s",
			username, rootpath, relpath, localpath, urlprefix)

		http.StripPrefix(urlprefix, http.FileServer(http.Dir(rootpath))).ServeHTTP(w, r)

		// objectHandler := handler.NewListHandler(hndlr)
		// objectHandler.ServeHTTP(w, r)

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
		// should this be an option? skip hidden files..
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}
		if file.IsDir() {
			directories = append(directories, file.Name())
		} else {
			filenames = append(filenames, file.Name())
		}
	}
	sort.Strings(directories)
	sort.Strings(filenames)

	var customLayout string
	ltmp := strings.Split(localpath, "/")
Layout:
	for i := len(ltmp); i > 1; i-- {
		p := strings.Join(ltmp[:i], "/")
		log.Tracef("check custom layout %s", p)
		l, found := voy.Layouts[p]
		if found {
			customLayout = l
			break Layout
		}
	}

	if customLayout == "" {
		hndlr.Layout = filetype.GuessLayout(localpath, filenames)
	} else {
		hndlr.Layout = filetype.TypeFromString(customLayout)
		log.Debugf("using custom layout %s for %s", hndlr.Layout, localpath)
	}

	hndlr.Filenames = filenames
	hndlr.Directories = directories

	log.Debugf("dispatch %s user:%s rootpath:%s path:%s localpath:%s urlprefix:%s files:%d dirs:%d",
		hndlr.Layout, username, rootpath, relpath, localpath, urlprefix, len(filenames), len(directories))

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
