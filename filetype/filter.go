package filetype

func Filter(names []string, types ...FileType) []string {
	result := make([]string, 0)
	keep := make(map[FileType]bool, 0)
	for _, t := range types {
		keep[t] = true
	}
	for _, name := range names {
		ftype, _ := Determine(name)
		if _, ok := keep[ftype]; ok {
			result = append(result, name)
		}
	}
	return result
}
