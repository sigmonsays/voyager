package filetype

//go:generate stringer -type=FileType

import (
	"fmt"
	"path/filepath"
	"strings"
)

var Picture = map[string]bool{
	"jpg":  true,
	"gif":  true,
	"png":  true,
	"tiff": true,
}

var Video = map[string]bool{
	"avi": true,
	"mp4": true,
	"m4v": true,
}

var Audio = map[string]bool{
	"mp3": true,
	"wav": true,
	"m3u": true,
}

type FileType int

const (
	UnknownFile FileType = iota
	PictureFile
	VideoFile
	AudioFile
)

var FileTypes map[string]FileType

func merge(src map[string]bool, dst map[string]FileType, filetype FileType) {
	for ext, enabled := range src {
		if enabled {
			dst[ext] = filetype
		}
	}
}

func init() {
	FileTypes = make(map[string]FileType, 0)
	merge(Picture, FileTypes, PictureFile)
	merge(Video, FileTypes, VideoFile)
	merge(Audio, FileTypes, AudioFile)
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
