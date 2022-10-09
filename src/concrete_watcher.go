package src

import "path/filepath"

type concreteWatcher struct {
	paths, files, commands []string
}

func (cw *concreteWatcher) SetCommands(c []string) { cw.commands = c }
func (cw *concreteWatcher) SetFiles(f []string) error {
	f, err := toAbs(f)
	if err != nil {
		return err
	}
	cw.files = f
	return nil
}

func (cw *concreteWatcher) SetPaths(p []string) error {
	p, err := toAbs(p)
	if err != nil {
		return err
	}
	cw.paths = p
	return nil
}

func (cw *concreteWatcher) Watch() error {
	if err := Watch(cw.commands, cw.files, cw.paths); err != nil {
		return err
	}
	return nil
}

func toAbs(filePaths []string) ([]string, error) {
	var absSlice []string
	for _, f := range filePaths {
		abs, err := filepath.Abs(f)
		if err != nil {
			return nil, err
		}
		absSlice = append(absSlice, abs)
	}

	return absSlice, nil
}
