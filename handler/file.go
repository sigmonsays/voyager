package handler

import (
	"path/filepath"
)

type File struct {
	Url  string
	Name string
}

func (f *File) Basename() string {
	return filepath.Base(f.Name)
}
