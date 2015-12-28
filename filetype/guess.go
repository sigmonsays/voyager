package filetype

import (
	"path/filepath"

	"github.com/sigmonsays/voyager/types"
)

func GuessLayout(localpath string, files []*types.File) FileType {
	totals := make(map[FileType]int)
	for _, f := range files {
		if f.IsDir {
			continue
		}
		ftype, err := Determine(filepath.Join(localpath, f.Name))
		if err != nil {
			log.Warnf("determine filetype %s: %s", f.Name, err)
		}
		if _, ok := totals[ftype]; ok == false {
			totals[ftype] = 0
		}
		totals[ftype]++
	}
	log.Tracef("localpath:%s totals:%+v", localpath, totals)
	var layout FileType
	var numfiles int
	for ftype, cnt := range totals {
		if ftype != UnknownFile && cnt > numfiles {
			layout = ftype
			numfiles = cnt
		}
	}
	log.Debugf("guessed layout for %s as %s", localpath, layout)
	return layout
}
