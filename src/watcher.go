package src

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const RefractoryPeriod = 500 * time.Millisecond

type Watcher interface {
	SetCommands(c []string)
	SetFiles(f []string) error
	SetPaths(p []string) error
	Watch() error
}

func NewWatcher() (Watcher, error) {
	return &concreteWatcher{}, nil
}

func Watch(commands, files, paths []string) error {
	if err := ValidateFields(commands, files, paths); err != nil {
		return err
	}

	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer w.Close()

	if err = RegisterFiles(files, w); err != nil {
		return err
	}
	if err = RegisterPaths(paths, w); err != nil {
		return err
	}
	if err = WatchLoop(commands, files, paths, w); err != nil {
		return err
	}
	return nil
}

func Execute(commands []string) {
	for _, c := range commands {
		split := strings.Split(c, " ")
		name := split[0]
		args := split[1:]

		cmd := exec.Command(name, args...)
		cmd.Stdout = os.Stdout
		err := cmd.Run()

		if err != nil {
			_, _ = fmt.Fprint(os.Stderr, err)
		}
	}
}

func RegisterFiles(files []string, w *fsnotify.Watcher) error {
	for _, f := range files {
		st, err := os.Lstat(f)
		if err != nil {
			return err
		}

		if st.IsDir() {
			return fmt.Errorf("expected `file` received `directory`")
		}

		// We will be watching the directory where the specified files are located, discarding events belonging to other
		// files.
		if err = w.Add(filepath.Dir(f)); err != nil {
			return err
		}
	}
	return nil
}

func RegisterPaths(paths []string, w *fsnotify.Watcher) error {
	for _, p := range paths {
		if err := w.Add(p); err != nil {
			return err
		}
	}
	return nil
}

func ValidateFields(commands, files, paths []string) error {
	if len(commands) == 0 {
		return fmt.Errorf("should specified at least one `command` to execute on change")
	}
	if len(files) == 0 && len(paths) == 0 {
		return fmt.Errorf("should specify at least one `path` or `file` to watch")
	}
	return nil
}

// WatchLoop is a controller that triggers specific cmds on fs events
func WatchLoop(commands, files, paths []string, w *fsnotify.Watcher) error {
	var (
		timer         = time.AfterFunc(math.MaxInt64, func() { Execute(commands) })
		validateEvent = func(e fsnotify.Event) bool {
			for _, f := range files {
				if e.Name == f {
					return true
				}
			}
			for _, p := range paths {
				if strings.HasPrefix(e.Name, p) {
					return true
				}
			}
			return false
		}
	)

	for {
		select {
		case err, ok := <-w.Errors:
			// Read `w.Errors` & if not `ok`, it means the channel was closed & we can return err
			if !ok {
				return err
			}
		case e, ok := <-w.Events:
			if !ok {
				return nil
			}
			// discard event if it's not a Create or Write event
			if !(e.Op == fsnotify.Create) && !(e.Op == fsnotify.Write) {
				continue
			}
			// discard event if event is not validated.
			if !validateEvent(e) {
				continue
			}

			timer.Reset(RefractoryPeriod)
		}
	}
}
