package cache

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type FileCache struct {
	Path string
}

func key(path string) string {
	k := base64.URLEncoding.EncodeToString([]byte(path))

	return k[0:2] + "/" + k[3:]
}

func NewFileCache(path string) *FileCache {
	c := &FileCache{
		Path: path,
	}
	return c
}

func (c *FileCache) filepath(path string) string {
	k := key(path)
	return filepath.Join(c.Path, k)
}

func (c *FileCache) Get(path string) (io.ReadCloser, error) {
	k := c.filepath(path)
	f, err := os.Open(k)
	if err != nil {
		return nil, err
	}
	log.Tracef("get path=%s (%s)", path, k)
	return f, nil
}

func (c *FileCache) Set(path string) (io.WriteCloser, error) {
	k := c.filepath(path)
	dir := filepath.Dir(k)

	if st, err := os.Stat(dir); err == nil {
		if st.IsDir() == false {
			return nil, fmt.Errorf("not a directory: %s", dir)
		}
	} else {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return nil, err
		}
	}

	f, err := os.Create(k)
	if err != nil {
		return nil, err
	}

	log.Tracef("set path=%s (%s)", path, k)
	return f, nil
}
