package types

import (
	"fmt"
	"path/filepath"
)

// used to resolve the path from a voy file
type PathRequest struct {
	LocalPath string
	RootPath  string
	RelPath   string
	UrlPrefix string
}

func (p *PathRequest) String() string {
	return fmt.Sprintf("localpath:%s rootpath:%s relpath:%s urlprefix:%s",
		p.LocalPath, p.RootPath, p.RelPath, p.UrlPrefix)
}

type ListPathRequest struct {
	// the server name (alias if used)
	Server string
	// the resolved server name (real IP or dns name)
	ServerName string

	User string
	Path string
}

func (p *ListPathRequest) String() string {
	return fmt.Sprintf("server:%s user:%s path:%s", p.Server, p.User, p.Path)
}

type ListPathResponse struct {
	*ListPathRequest
	Title     string
	LocalPath string
	Files     []*File
}

type Files []*File

type File struct {
	IsDir bool
	Name  string
	Size  int64
	MTime int64
}

type SortByName struct {
	Files
}

func (me SortByName) Len() int {
	return len(me.Files)
}
func (me SortByName) Less(i, j int) bool {
	return me.Files[i].Name < me.Files[j].Name
}
func (me SortByName) Swap(i, j int) {
	me.Files[i], me.Files[j] = me.Files[j], me.Files[i]
}

func (f *File) String() string {
	return f.Name
}

// convenience for templating
type FileData struct {
	*File
	Url string
}

func (f *File) Basename() string {
	return filepath.Base(f.Name)
}
