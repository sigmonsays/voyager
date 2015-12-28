package handler

import (
	"os"
	"sort"
	"strings"

	"github.com/sigmonsays/voyager/types"
)

type PathLoader interface {
	GetFiles(path string) ([]*types.File, error)
}

func NewFilesystemPathLoader() *FilesystemPathLoader {
	p := &FilesystemPathLoader{}
	return p
}

type FilesystemPathLoader struct {
}

func (me *FilesystemPathLoader) GetFiles(localpath string) ([]*types.File, error) {

	fl := make([]*types.File, 0)

	log.Tracef("getfiles %s", localpath)
	fh, err := os.Open(localpath)
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	st, err := fh.Stat()
	if err != nil {
		return nil, err
	}

	// TODO: is this right?
	if st.IsDir() == false {
		return nil, nil
	}

	/*
		if st.IsDir() == false {

			// serve the object directly
			log.Debugf("dispatch ListHandler user:%s rootpath:%s path:%s localpath:%s urlprefix:%s",
				username, rootpath, relpath, localpath, urlprefix)

			http.StripPrefix(urlprefix, http.FileServer(http.Dir(rootpath))).ServeHTTP(w, r)

			// objectHandler := handler.NewListHandler(hndlr)
			// objectHandler.ServeHTTP(w, r)

			return
		}
	*/

	files, err := fh.Readdir(-1)
	if err != nil {
		return nil, err
	}

	filecnt := 0
	dircnt := 0
	for _, file := range files {

		if file.IsDir() {
			dircnt++
		} else {
			filecnt++
		}

		// should this be an option? skip hidden files..
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}

		f := &types.File{
			IsDir: file.IsDir(),
			Name:  file.Name(),
			Size:  file.Size(),
			MTime: file.ModTime().Unix(),
		}
		fl = append(fl, f)
	}

	sorter := types.SortByName{
		Files: fl,
	}
	sort.Sort(sorter)

	log.Debugf("loaded files from %s files:%d dirs:%d", localpath, filecnt, dircnt)

	return fl, nil
}
