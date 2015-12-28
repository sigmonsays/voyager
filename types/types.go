package types

import (
	"path/filepath"
)

type ListPathRequest struct {
	Server string
	User   string
	Path   string
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
