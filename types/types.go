package types

import (
	"path/filepath"
)

type ListPathRequest struct {
	User string
	Path string
}
type ListPathResponse struct {
	*ListPathRequest
	Title     string
	LocalPath string
	Files     []*File
}

type File struct {
	IsDir bool
	Name  string
	Size  int64
	MTime int64
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
