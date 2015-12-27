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
	Title       string
	LocalPath   string
	Files       []*File
	Directories []*File
}

type File struct {
	Url  string
	Name string
	Size int64
}

func (f *File) Basename() string {
	return filepath.Base(f.Name)
}
