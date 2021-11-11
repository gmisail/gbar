# gbar
`gbar` is a command-line utility written in Go that feeds `lemonbar`, a status bar for Linux. It uses `goroutines` so that different modules can generate their 
content concurrently. This also means that modules that have different timers (e.g. one module can update once a second and another can update once a minute) can operate at the same time and report back changes once they're complete.

```
go build gbar.go
./gbar | lemonbar
```
Since `./gbar` will not terminate, the status bar will run until it is manually terminated *or* the system turns off.
