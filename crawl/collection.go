package crawl

//go:generate stringer -type=Change

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"
)

type Change int

const (
	ChangeUnknown Change = iota
	ChangeNone
	ChangeUpdated
	ChangeRemoved
)

// collection represents a single file system path
// that we're going to crawl and maintain state of
func NewCollection(path string) (*Collection, error) {
	c := &Collection{
		Path:      path,
		DbPath:    filepath.Join(path, ".voyager/db"),
		IndexFile: filepath.Join(path, ".voyager/index"),
	}

	err := c.init()
	if err != nil {
		return nil, err
	}

	return c, nil
}

type Collection struct {
	Path      string
	DbPath    string
	IndexFile string
}

func (c *Collection) init() error {

	// create database path
	st, err := os.Stat(c.DbPath)
	if err != nil {
		err = os.MkdirAll(c.DbPath, 0755)
		if err != nil {
			return err
		}
		st, err = os.Stat(c.DbPath)
	}
	if err == nil && st.IsDir() == false {
		return fmt.Errorf("%s is a file", c.DbPath)
	}
	return nil

}

func (c *Collection) CreateIndex() error {
	tmpname := c.IndexFile + ".tmp"
	f, err := os.Create(tmpname)
	if err != nil {
		return err
	}
	defer f.Close()

	walkfn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() == false {
			return nil
		}

		e := &Entry{
			Path:  path,
			Type:  TypeDir,
			CTime: StatCtime(info),
		}

		fmt.Fprintf(f, "%s\n", e)
		return nil
	}
	err = filepath.Walk(c.Path, walkfn)
	if err != nil {
		return err
	}

	err = os.Rename(tmpname, c.IndexFile)
	return err
}

type Result struct {
	*Entry
	Err    error
	Change Change
}

func StatCtime(st os.FileInfo) int64 {
	if s, ok := st.Sys().(*syscall.Stat_t); ok {
		return s.Ctim.Sec
	}
	return -1
}

func (c *Collection) Update() error {
	newIndexName := c.IndexFile + ".new"
	if _, err := os.Stat(newIndexName); err == nil {
		return fmt.Errorf("update already running")
	}

	firstRun := false
	if _, err := os.Stat(c.IndexFile); err != nil {
		err = c.CreateIndex()
		if err != nil {
			return err
		}
		firstRun = true
	}

	// open directory index file
	dirlist, err := NewDirectoryIndex(c.IndexFile)
	if err != nil {
		return err
	}
	defer dirlist.Close()

	// start a bunch of workers who will check each directory and return a path if the ctimes differ
	concurrency := 10
	control := make(chan int, 0)
	results := make(chan *Result, 0)
	work := make(chan *Entry, concurrency)
	worker := func(work chan *Entry) {
		var entry *Entry
		var ok bool
	Loop:
		for {
			select {
			case entry, ok = <-work:
				if ok == false {
					break Loop
				}
				st, err := os.Stat(entry.Path)
				if err != nil {
					result := &Result{
						Entry: entry,
					}
					if os.IsNotExist(err) {
						result.Change = ChangeRemoved
					} else {
						result.Err = err
						log.Warnf("stat %s: %s", entry.Path, err)
					}
					results <- result
					continue
				}
				ctime := StatCtime(st)

				if ctime == entry.CTime {
					results <- &Result{
						Entry:  entry,
						Change: ChangeNone,
					}
				} else {
					entry.CTime = ctime
					results <- &Result{
						Entry:  entry,
						Change: ChangeUpdated,
					}
				}
			}
		}
	}
	for i := 0; i < concurrency; i++ {
		go worker(work)
	}

	sourcefn := func() {
		processed := 0
		for {
			entry, err := dirlist.Next()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Errorf("%s", err)
				break
			}
			if entry.Path == c.DbPath {
				continue
			}
			processed++
			work <- entry
		}
		control <- processed

	}
	go sourcefn()

	// make a new index file
	newindex, err := os.Create(newIndexName)
	if err != nil {
		return err
	}

	changed := make(map[string]*Result, 0)

	process := func(result *Result) error {
		if result.Entry != nil {
			if result.Change != ChangeRemoved {
				fmt.Fprintf(newindex, "%s\n", result.Entry.String())
			}
		}

		if firstRun {
			fmt.Printf("NEWDIR %+v\n", result.Entry)
			return nil
		}
		if result.Err != nil {
			fmt.Printf("ERROR %+v\n", result.Err)
			return nil
		}
		if result.Change == ChangeNone {
			return nil
		}

		fmt.Printf("CHANGE(%s) %+v\n", result.Change, result.Entry)

		if result.Entry.Type == TypeDir {

			f, err := os.Open(result.Entry.Path)
			if err != nil {
				log.Warnf("Open %s: %s", result.Entry.Path, err)
				return nil
			}
			directories, err := f.Readdir(-1)
			f.Close()

			for _, f := range directories {
				ctime := StatCtime(f)
				if ctime >= result.Entry.CTime {
					path := filepath.Join(result.Entry.Path, f.Name())
					fmt.Printf("CHANGED %s\n", path)
					e := &Entry{
						Path:  filepath.Join(result.Entry.Path, f.Name()),
						CTime: ctime,
					}
					changed[path] = result
					if result.Change != ChangeRemoved {
						fmt.Fprintf(newindex, "%s\n", e)
					}
				}
			}
		}

		return nil
	}

	var result *Result
	received := 0
	finished := 0
	source_finished := false
Loop:
	for {
		select {
		case result = <-results:
			received++
			process(result)
			if source_finished && finished == received {
				break Loop
			}

		case finished = <-control:
			source_finished = true
			log.Infof("source finished, %d processed, %d finished, %d pending", received, finished, finished-received)
			close(work)

			if finished == received {
				break Loop
			}
		}
	}

	newindex.Close()
	err = os.Rename(newIndexName, c.IndexFile)

	// go through what has changed
	fmt.Printf("inspecting..\n")
	for _, change := range changed {
		fmt.Printf("%+v\n", change)
		if change.Change == ChangeRemoved {
			err = dirlist.RemovePrefix(change.Entry.Path)
			if err != nil {
				log.Errorf("RemovePrefix %s: %s", change.Entry.Path, err)
			}
		} else {
			log.Warnf("unsupported %s", change)
		}

		_, err := dirlist.Stat(change.Entry.Path)
		if err == nil {
			// already in the index, just do a partial update

			fmt.Printf("partial update %s\n", change)

		} else {
			fmt.Printf("full update %s\n", change)
		}

	}
	return err
}
