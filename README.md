# gbar
`gbar` is a command-line utility written in Go that feeds `lemonbar`, a status bar for Linux. It uses `goroutines` so that different modules can generate their 
content concurrently. This means that modules that have different timers (e.g. one module can update once a second and another can update once a minute) can operate at the same time and report back changes once they're complete.

## Installation
```
go build gbar.go
./gbar
```

Note that `gbar` is completely self-contained, i.e. it will spawn the `lemonbar` process for you. 

## Future

- [x] Modules
- [x] Goroutine-based render loop
- [x] Embedded `lemonbar` process
- [ ] JSON configuration file (module organization, general settings, etc.)
- [ ] Custom modules, i.e. shell scripts

## Limitations
Currently, executable paths are hardcoded to wherever they are on my system. This is due to how it currently does not support configuration files; to change where the paths point to, you will need to clone this repo and change the paths accordingly.
