# gbar
`gbar` is a command-line utility written in Go that feeds `lemonbar`, a status bar for Linux. 

To start `gbar` with `lemonbar`, simply pipe the output from one process to the other:
```
./gbar | lemonbar
```
Since `./gbar` will not terminate, the status bar will run until it is manually terminated *or* the system turns off.
