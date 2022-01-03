# gbar
`gbar` is a command-line utility written in Go that feeds `lemonbar`, a status bar for Linux. It uses `goroutines` so that different modules can generate their 
content concurrently. This means that modules that have different timers (e.g. one module can update once a second and another can update once a minute) can operate at the same time and report back changes once they're complete.

Currently, all modules that you want to use must be added to the `gbar.go` file and then recompiled with `go build`. Once the executable is built, you can organize the blocks in any configuration without re-compiling by using a `config.json` file. 

## Installation
```
go build
./gbar
```

Note that `gbar` is completely self-contained, i.e. it will spawn the `lemonbar` process for you. 

## Features

- [x] Modules
- [x] Goroutine-based render loop
- [x] Embedded `lemonbar` process
- [x] JSON configuration file (module organization, general settings, etc.)
- [x] Custom modules, i.e. shell scripts
	- Custom modules are written in Go, however you can run shell scripts by running them from within Go 
- [ ] Module scripting
	- Custom modules no longer require recompilation of `gbar`, but instead use Lua / Tengo / etc...

## Configuration
For an example of a configuration file, refer to the `config.json` file.
