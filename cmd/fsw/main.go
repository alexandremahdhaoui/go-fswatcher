package main

import (
	"fmt"
	"gitlab.com/alexandre.mahdhaoui/go-fswatcher/src"
	"os"
	"path/filepath"
)

var usage = `
fswatcher is a library providing a simple file watcher that can
execute specified commands on change.

https://github.com/alexandre.mahdhaoui/go-fswatcher

Options:

	-f 	--file 		[filepath]
	-h 	--help
	-p 	--path		[path]
	-x 	--execute	[command]
`

func main() {
	if len(os.Args) == 0 {
		help()
	}
	checkHelpFlag()

	c, f, p := parseFlags()

	w, err := src.NewWatcher()
	if err != nil {
		exit("error while creating a new watcher: %s", err)
	}

	w.SetCommands(c)
	w.SetFiles(f)
	w.SetPaths(p)

	if err = w.Watch(); err != nil {
		exit("%s", err)
	}
}

func appendFlag(position int, optArray []string) []string {
	if len(os.Args) <= position {
		exit("option was specified but expected 1 argument")
	}
	return append(optArray, os.Args[position+1])
}

func checkHelpFlag() {
	for _, flag := range os.Args[1:] {
		switch flag {
		case "-h", "--help":
			help()
		}
	}
}

func exit(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, filepath.Base(os.Args[0])+": "+format+"\n", a...)
	help()
}

func help() {
	fmt.Printf("%s [options]\n", os.Args[0])
	fmt.Print(usage)
	os.Exit(1)
}

func parseFlags() ([]string, []string, []string) {
	var (
		cmds  []string
		files []string
		paths []string
	)
	for i, flag := range os.Args[1:] {
		switch flag {
		case "-f", "--file":
			files = appendFlag(i, files)
		case "-p", "--path":
			paths = appendFlag(i, paths)
		case "-x", "--execute":
			cmds = appendFlag(i, cmds)
		}
	}
	validateFlags(files, paths, cmds)
	return cmds, files, paths
}

func validateFlags(files, paths, cmds []string) {
	if len(cmds) == 0 {
		exit("Please specify at least one executable command, using the `-x`,`--execute` flag.")
	}
	if len(files) == 0 || len(paths) == 0 {
		exit("Please specify at least one `path` or `file` to watch.")
	}
}
