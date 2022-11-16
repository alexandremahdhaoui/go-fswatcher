# go-fswatcher

fsw is a library providing a simple file watcher that can
execute specified commands on change.

https://github.com/alexandre.mahdhaoui/go-fswatcher

# Install

```shell
go install gitlab.com/alexandre.mahdhaoui/go-fswatcher/cmd/fsw@latest
```

Usage: 	
```shell
fsw [options]
```

Options:

	-f	--file 		[file]		file to watch.
	-h	--help				print the helper.
	-p	--path		[path]		path to a directory to watch.
	-x	--execute	[command]	executable command.


# Example

```shell
fsw -p . -x "echo some changes triggered this message"
```
