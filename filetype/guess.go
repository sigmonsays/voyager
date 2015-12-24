package filetype

import (
	"path/filepath"
)

func GuessLayout(localpath string, filenames []string) FileType {
	found := make(map[FileType]int)
	for _, filename := range filenames {
		ftype, err := Determine(filepath.Join(localpath, filename))
		if err != nil {
			log.Warnf("determine filetype %s: %s", filename, err)
		}
		if _, ok := found[ftype]; ok == false {
			found[ftype] = 0
		}
		found[ftype]++
	}
	var layout FileType
	var numfiles int
	for ftype, cnt := range found {
		if ftype != UnknownFile && cnt > numfiles {
			layout = ftype
			cnt = numfiles
		}
	}
	log.Debugf("guessed layout for %s as %s", localpath, layout)
	return layout
}
