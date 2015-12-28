package filetype

import (
	"github.com/sigmonsays/voyager/types"
)

func FileMatch(file *types.File, types ...FileType) bool {
	keep := make(map[FileType]bool, 0)
	for _, t := range types {
		keep[t] = true
	}
	ftype, _ := Determine(file.Name)
	_, match := keep[ftype]

	return match

}

func Filter(files []*types.File, ftypes ...FileType) []*types.File {
	result := make([]*types.File, 0)
	keep := make(map[FileType]bool, 0)
	for _, t := range ftypes {
		keep[t] = true
	}
	for _, f := range files {
		ftype, _ := Determine(f.Name)
		if _, ok := keep[ftype]; ok {
			result = append(result, f)
		}
	}
	return result
}
