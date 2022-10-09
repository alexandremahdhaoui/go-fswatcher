package src

type concreteWatcher struct {
	paths, files, commands []string
}

func (cw *concreteWatcher) SetCommands(c []string) { cw.commands = c }
func (cw *concreteWatcher) SetFiles(f []string)    { cw.files = f }
func (cw *concreteWatcher) SetPaths(p []string)    { cw.paths = p }

func (cw *concreteWatcher) Watch() error {
	if err := Watch(cw.commands, cw.files, cw.paths); err != nil {
		return err
	}
	return nil
}