package filetype

//go:generate stringer -type=FileType

import (
	"fmt"
	"path/filepath"
	"strings"
)

type FileType int

const (
	UnknownFile FileType = iota
	PictureFile
	VideoFile
	AudioFile
)

var Picture = map[string]bool{
	"jpg":  true,
	"gif":  true,
	"png":  true,
	"tiff": true,
}

var Video = map[string]bool{
	"asf":  true,
	"avi":  true,
	"flv":  true,
	"m4v":  true,
	"mkv":  true,
	"mov":  true,
	"mp4":  true,
	"mpg":  true,
	"vob":  true,
	"webm": true,
	"wmv":  true,
}

var Audio = map[string]bool{
	"mp3": true,
	"wav": true,
	"m3u": true,
}

// extension mapping
var FileTypes map[string]FileType

var FileTypeNames map[string]FileType
var FileTypeIds map[FileType]string

func TypeToString(ftype FileType) string {
	name, ok := FileTypeIds[ftype]
	if ok == false {
		log.Warnf("unknown type %s", name)
		return "unknown"
	}
	return name
}
func TypeFromString(name string) FileType {
	ftype, ok := FileTypeNames[strings.ToLower(name)]
	if ok == false {
		log.Warnf("unknown type %s", name)
		return UnknownFile
	}
	return ftype
}

func merge(src map[string]bool, dst map[string]FileType, filetype FileType) {
	for ext, enabled := range src {
		if enabled {
			dst[ext] = filetype
		}
	}
}

func init() {
	FileTypes = make(map[string]FileType, 0)
	FileTypeIds = make(map[FileType]string, 0)
	merge(Picture, FileTypes, PictureFile)
	merge(Video, FileTypes, VideoFile)
	merge(Audio, FileTypes, AudioFile)

	FileTypeNames = map[string]FileType{
		"unknown":  UnknownFile,
		"pictures": PictureFile,
		"video":    VideoFile,
		"audio":    AudioFile,
	}

	for name, ftype := range FileTypeNames {
		FileTypeIds[ftype] = name
	}
}

func Determine(path string) (FileType, error) {

	ext := filepath.Ext(path)

	if ext == "" {
		return UnknownFile, fmt.Errorf("no extension")
	}

	ext = strings.ToLower(ext[1:])

	filetype, ok := FileTypes[ext]

	if ok == false {
		return UnknownFile, fmt.Errorf("unknown extension %s", ext)
	}
	return filetype, nil
}
