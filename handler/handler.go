package handler

import (
	"github.com/sigmonsays/voyager/filetype"
)

type Handler struct {
	Layout    filetype.FileType
	Username  string
	Homedir   string
	Path      string
	Filenames []string
}
