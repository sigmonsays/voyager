package handler

import (
	"github.com/sigmonsays/voyager/types"
)

// supports video and pictures
type Gallery struct {
	Path         string
	LocalPath    string
	UrlPrefix    string
	RelPath      string
	RemoteServer string
	Title        string
	Files        []*types.File
	Directories  []*types.File
	Breadcrumb   *Breadcrumb
}
