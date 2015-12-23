package crawl

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	TypeDir = iota
	TypeFile
)

type DirectoryIndex struct {
	Path  string
	f     *os.File
	s     *bufio.Reader
	files map[string]*Entry
}

func NewDirectoryIndex(path string) (*DirectoryIndex, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	s := bufio.NewReader(f)
	d := &DirectoryIndex{
		Path:  path,
		f:     f,
		s:     s,
		files: make(map[string]*Entry, 0),
	}
	return d, nil
}
func (d *DirectoryIndex) Close() error {
	return d.f.Close()
}

type Entry struct {
	Path  string
	Type  int
	CTime int64
}

func (e *Entry) String() string {
	return fmt.Sprintf("%s %d %d", e.Path, e.Type, e.CTime)
}

func ParseEntry(line string) (*Entry, error) {
	idx := strings.LastIndex(line, " ")
	if idx == -1 {
		return nil, fmt.Errorf("unexpected format")
	}

	tmp := strings.Fields(line)
	tmplen := len(tmp)

	typ := int64(0)
	typ, err := strconv.ParseInt(tmp[tmplen-1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("typ error: %s", err)
	}

	ctime, err := strconv.ParseInt(line[idx+1:len(line)-1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("ctime error: %s", err)
	}

	e := &Entry{
		Path:  line[0:idx],
		Type:  int(typ),
		CTime: ctime,
	}
	return e, nil
}
func (d *DirectoryIndex) Stat(fname string) (*Entry, error) {
	entry, ok := d.files[fname]
	if ok == false {
		return nil, fmt.Errorf("%s not found", fname)
	}
	return entry, nil
}

func (d *DirectoryIndex) RemovePrefix(path string) error {
	tmpname := d.Path + ".work"
	newindex, err := os.Create(tmpname)
	if err != nil {
		return err
	}
	defer newindex.Close()

	infile, err := os.Open(d.Path)
	if err != nil {
		return err
	}
	defer infile.Close()

	s := bufio.NewReader(infile)
	for {
		bline, err := s.ReadBytes('\n')
		if err != nil {
			break
		}

		line := string(bline)
		if strings.HasPrefix(line, path) {
			continue
		}
		newindex.Write(bline)
	}
	return nil
}

func (d *DirectoryIndex) Next() (*Entry, error) {

	line, err := d.s.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	entry, err := ParseEntry(string(line))
	if err == nil {
		d.files[entry.Path] = entry
	}
	return entry, err

}
