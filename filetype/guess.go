package filetype

import (
	"path/filepath"
)

func GuessLayout(localpath string, filenames []string) FileType {
	totals := make(map[FileType]int)
	for _, filename := range filenames {
		ftype, err := Determine(filepath.Join(localpath, filename))
		if err != nil {
			log.Warnf("determine filetype %s: %s", filename, err)
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
