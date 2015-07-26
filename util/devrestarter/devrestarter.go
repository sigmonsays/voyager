// Package devrestarter automatically restarts your server when it's updated.
// This is suitable for use during development. When you recompile your server,
// it will be restarted.
package devrestarter

import (
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"github.com/go-fsnotify/fsnotify"
)

var debug func(v ...interface{})

func noop(v ...interface{}) {}

func init() {
	if os.Getenv("RELOADER_DEBUG") == "" {
		debug = noop
	} else {
		debug = log.Debug
	}
}

func watch() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	argv0, err := exec.LookPath(os.Args[0])
	if err != nil {
		return err
	}

	debug("reloader:", argv0)
	if err := watcher.Add(argv0); err != nil {
		return err
	}

	var doit <-chan time.Time
	for {
		select {
		case <-doit:
			log.Info("Restarting", argv0)
			if err := syscall.Exec(argv0, os.Args, os.Environ()); err != nil {
				return err
			}
		case event := <-watcher.Events:
			debug("watcher.Event:", event)
			doit = time.After(500 * time.Millisecond)
			if err := watcher.Add(event.Name); err != nil {
				return err
			}
		case err := <-watcher.Errors:
			debug("watcher.Error:", err)
		}
	}
}

// make Init idempotent.
var once sync.Once

// Init initializes the background goroutine that will restart the binary when
// it changes.
func Init() {
	once.Do(func() {
		go func() {
			if err := watch(); err != nil {
				panic(err)
			}
		}()
	})
}
