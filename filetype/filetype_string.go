// Code generated by "stringer -type=FileType"; DO NOT EDIT

package filetype

import "fmt"

const _FileType_name = "UnknownFilePictureFileVideoFileAudioFileListFile"

var _FileType_index = [...]uint8{0, 11, 22, 31, 40, 48}

func (i FileType) String() string {
	if i < 0 || i >= FileType(len(_FileType_index)-1) {
		return fmt.Sprintf("FileType(%d)", i)
	}
	return _FileType_name[_FileType_index[i]:_FileType_index[i+1]]
}
