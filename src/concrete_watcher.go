package src

import "path/filepath"

type concreteWatcher struct {
	paths, files, commands []string
}

func (c *concreteWatcher) SetCommands(cmds []string) { c.commands = cmds }
func (c *concreteWatcher) SetFiles(f []string) error {
	f, err := toAbs(f)
	if err != nil {
		return err
	}
	c.files = f
	return nil
}

func (c *concreteWatcher) SetPaths(p []string) error {
	p, err := toAbs(p)
	if err != nil {
		return err
	}
	c.paths = p
	return nil
}

func (c *concreteWatcher) Watch() error {
	if err := Watch(c.commands, c.files, c.paths); err != nil {
		return err
	}
	return nil
}

func toAbs(filePaths []string) ([]string, error) {
	var absPaths []string
	for _, f := range filePaths {
		abs, err := filepath.Abs(f)
		if err != nil {
			return nil, err
		}
		absPaths = append(absPaths, abs)
	}

	return absPaths, nil
}
