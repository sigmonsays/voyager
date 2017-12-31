package voy

import (
	"fmt"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/sigmonsays/voyager/types"
)

type VoyLoader interface {
	Load(req *types.ListPathRequest) (*VoyFile, error)
	ResolvePath(voy *VoyFile, req *types.ListPathRequest) (*types.PathRequest, error)
}

func NewVoyLoader() *fileLoader {
	l := &fileLoader{}
	return l
}

type fileLoader struct {
}

func (l *fileLoader) Load(req *types.ListPathRequest) (*VoyFile, error) {
	user_ent, err := user.Lookup(req.User)
	if err != nil {
		return nil, fmt.Errorf("user lookup %s: %s", req.User, err)
	}
	homedir := user_ent.HomeDir

	cfg := DefaultConfig()
	voyfile := filepath.Join(homedir, ".voyager.cfg")

	err = cfg.LoadYaml(voyfile)
	if err != nil {
		return nil, fmt.Errorf("load voyfile %s: %s", voyfile, err)
	}

	return cfg, nil
}

func (l *fileLoader) ResolvePath(voy *VoyFile, req *types.ListPathRequest) (*types.PathRequest, error) {

	user_ent, err := user.Lookup(req.User)
	if err != nil {
		return nil, fmt.Errorf("user lookup %s: %s", req.User, err)
	}
	homedir := user_ent.HomeDir

	tmp := strings.Split(req.Path, "/")
	if len(tmp) < 3 {
		// TODO: Do we want to support any kind of top level index?
		return nil, fmt.Errorf("incomplete path")
	}

	var localpath string
	var rootpath string
	var relpath string
	var urlprefix string

	topdir := tmp[1]
	alias, is_alias := voy.Alias[topdir]

	log.Tracef("topdir:%s alias:%s is_alias:%v", topdir, alias, is_alias)

	if is_alias {
		localpath = filepath.Join(alias, strings.Join(tmp[2:], "/"))
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
			return nil, fmt.Errorf("nothing to see here. bye bye.")
		}
		urlprefix = "/~" + req.User
		log.Debugf("%s is regular path (localpath:%s relpath:%s urlprefix:%s)", req.Path, localpath, relpath, urlprefix)
	}

	log.Infof("request user:%s rootpath:%s path:%s localpath:%s urlprefix:%s",
		req.User, rootpath, relpath, localpath, urlprefix)

	preq := &types.PathRequest{
		LocalPath: localpath,
		RootPath:  rootpath,
		RelPath:   relpath,
		UrlPrefix: urlprefix,
	}

	return preq, nil

}
