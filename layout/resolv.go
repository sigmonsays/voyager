package layout

import (
	"strings"

	"github.com/sigmonsays/voyager/filetype"
	"github.com/sigmonsays/voyager/types"
	"github.com/sigmonsays/voyager/voy"
)

type LayoutResolver interface {
	Resolve(voy *voy.VoyFile, localpath string, files []*types.File) (filetype.FileType, error)
}

func NewLayoutResolver() *layoutResolver {
	return &layoutResolver{}
}

type layoutResolver struct {
}

func (l *layoutResolver) Resolve(voy *voy.VoyFile, localpath string, files []*types.File) (filetype.FileType, error) {
	var customLayout string
	ltmp := strings.Split(localpath, "/")
Layout:
	for i := len(ltmp); i > 1; i-- {
		p := strings.Join(ltmp[:i], "/")
		log.Tracef("check custom layout %s", p)
		l, found := voy.Layouts[p]
		if found {
			log.Tracef("found custom layout %s for %s", l, p)
			customLayout = l
			break Layout
		}
	}

	var layout filetype.FileType

	if customLayout == "" {
		layout = filetype.GuessLayout(localpath, files)
	} else {
		layout = filetype.TypeFromString(customLayout)
		log.Debugf("using custom layout %s (%q) for %s", layout, customLayout, localpath)
	}

	return layout, nil
}
