# gbar
`gbar` is a command-line utility written in Go that feeds `lemonbar`, a status bar for Linux. It uses `goroutines` so that different modules can generate their 
content concurrently. This means that modules that have different timers (e.g. one module can update once a second and another can update once a minute) can operate at the same time and report back changes once they're complete.

Currently, all modules that you want to use must be added to the `gbar.go` file and then recompiled with `go build`. Once the executable is built, you can organize the blocks in any configuration without re-compiling by using a `config.json` file. 

The design of your status bar consists of three parts: the template, blocks, and modules. The template defines where blocks are positioned. Blocks contain information about the chunks of your status bar, such as update intervals and modules. Modules define what is actually rendered.

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
Below is an example configuration file.
```
{
	"settings": {
		"lemonbar": "lemonbar -U #0A0A0A -u 4 -B #0A0A0A -g x24 -p",
		"font": "Iosevka Nerd Font",
		"separator": "%{B-}%{F#1f1f1f} | %{B-}%{F-}"
	},

	"template": {
		"left": ["cpu", "ram"],
		"center": ["time"],
		"right": ["workspaces", "power"]
	},

	"blocks": {
		"cpu": {
			"module": "cpu",
			"interval": "2"
		},

		"ram": {
			"module": "ram",
			"interval": "5"
		},

		"time": {
			"module": "time",
			"interval": "1"
		}
	},

	"buttons": [
		{
			"name": "power",
			"onclick": "rofi ...",
			"label": "Power"
		}
	]
}
```
