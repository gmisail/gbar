package main

import ( 
	"fmt" 
	"time"
	"strings"
	"bufio"
	"os"
	"os/exec"
	"io"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/host"
)

type Alignment int

const (
	Left Alignment = iota
	Center
	Right
	None
)

type Module struct {
	ID int
	Name string
	Align Alignment
	Content string
}

func SetColor(bg string, fg string) string {
	return fmt.Sprintf("%%{B%s}%%{F%s}", bg, fg)
}

func Color(bg string, fg string, text string) string {
	return SetColor(bg, fg) + text + SetColor("-", "-")
}

func Underline(text string) string {
	return fmt.Sprintf("%%{+u}%s%%{-u}", text)
}

func Overline(text string) string {
	return fmt.Sprintf("%%{+o}%s%%{-o}", text)
}

func Button(text string, command string) string {
	return fmt.Sprintf("%%{A:%s:}%s%%{A}", command, text)
}

func Statistics(renderChannel chan []Module) {
	systemStats := Module{ ID: 0, Name: "CPU_RAM", Align: Left, Content: "Default" }
	date := Module{ ID: 2, Name: "Date", Align: Center, Content: "Default" }
	tempMod := Module{ ID: 1, Name: "CPU_Temp", Align: Left, Content: "Default" }

	for {
		times, _ := cpu.Percent(0, false)
		cpuUsage := times[0]
		temperatures, _ := host.SensorsTemperatures()
		memory, _ := mem.VirtualMemory()

		systemStats.Content = ""
		systemStats.Content += Color("-", "#D6E3F8", "  ")
		systemStats.Content += fmt.Sprintf("%.2f%% ", cpuUsage)

		systemStats.Content += " "
		systemStats.Content += Color("-", "#CC3F0C", "  ")
		systemStats.Content += fmt.Sprintf("%.2f%% ", memory.UsedPercent)

		currentTime := time.Now()

		date.Content = Color("-", "#F3EFF5", "  ")
		date.Content += fmt.Sprintf(currentTime.Format("Monday, January 2 "))
		date.Content += " "

		date.Content += Color("-", "#D5A021", "  ")
		date.Content += fmt.Sprintf(currentTime.Format("3:04:05 PM ")) 

		cpuTemperature := 0.0
		for i := 0; i < len(temperatures); i++ {
			if temperatures[i].SensorKey == "k10temp_tdie" {
				cpuTemperature = temperatures[i].Temperature
				break
			}
		}

		tempMod.Content = Color("-", "#FCFC62", "   ")
		tempMod.Content += fmt.Sprintf("%.2f C", cpuTemperature)

		copyModules := make([]Module, 3)
		copyModules[0] = systemStats
		copyModules[1] = tempMod
		copyModules[2] = date
		renderChannel <- copyModules

		time.Sleep(time.Second)
	}
}

/* watches for a BSPWM workspace changes (i.e. switching) */
func CurrentWorkspace(renderChannel chan []Module) {
	workspaceIds, _ := exec.Command("bspc", "query", "-D").Output()
	workspacesList := strings.Split(string(workspaceIds), "\n")

	/* converts each desktop ID to an integer (0-9) */
	workspaces := make(map[string]int)
	for i, id := range workspacesList {
		workspaces[id] = i
	}

	workspace := Module{ ID: 3, Align: Right, Content: "" }
	modules := make([]Module, 1)

	subscription := exec.Command("bspc", "subscribe", "desktop")
	subscriptionPipe, err := subscription.StdoutPipe()

	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed establishing pipe to bspc", err)
		return
	}

	// subscribe to bspwm for changes
	subscriptionScanner := bufio.NewScanner(subscriptionPipe)
	subscription.Start()

	// block for changes. When a change arrives, update it and push to the bar
	for subscriptionScanner.Scan() {
		// split up the command into arguments
		workspaceEvent := strings.Split(subscriptionScanner.Text(), " ")

		// look up which desktop the new ID belongs to
		currentWorkspace := workspaces[workspaceEvent[2]]

		workspace.Content = ""
		for i := 0; i < len(workspaces) - 1; i++ {
			if i == currentWorkspace {
				workspace.Content += " "
			} else {
				workspace.Content += " "
			}
		}

		modules[0] = workspace
		renderChannel <- modules
    }
}

// given a list of modules, renders their text with alignment
func RenderModules(modules []Module, stdin io.WriteCloser) {
	var alignment Alignment = None
	var buffer strings.Builder

	// for each module, update its content
	for i := 0; i < len(modules); i++ {
		if modules[i].Align != alignment {
			alignment = modules[i].Align
			switch modules[i].Align {
			case Left:
				buffer.WriteString("%{l}")
			case Center:
				buffer.WriteString("%{c}")
			case Right:
				buffer.WriteString("%{r}")
			}
		}

		buffer.WriteString(modules[i].Content)
	}

	// end of line, updates the bar with the new data
	buffer.WriteString("\n")

	// writes the data to lemonbar's STDIN
	io.WriteString(stdin, buffer.String())
}

/* Renders the status bar once it receives updated modules */
func RenderStatus(renderChannel chan []Module, config []Module, stdin io.WriteCloser) {
	// date => { name: "date", content: "Feb 1 2021" }
	var modules []Module

	for i := 0; i < len(config); i++ {
		modules = append(modules, config[i])
	}

	for {
		updatedModules := <-renderChannel

		for _, module := range updatedModules {
			modules[module.ID] = module
		}

		RenderModules(modules, stdin)
	}
}

// creates a lemonbar instance and connects the stdin pipe to gbar's renderer
func StartBar(renderer chan []Module, configuration []Module) {
	events := make(map[string] []string)
	events["power-menu"] = []string {
		"rofi", "-show", "p", "-modi", "p:rofi-power-menu" }

	bar := exec.Command(
		"/home/gmisail/Documents/Development/gbar/lemonbar",
		"-U", "#0A0A0A",
		"-u", "4",
		"-B", "#0A0A0A",
		"-g", "x24",
		"-f", "Iosevka Nerd Font", "-p")

	barStdout, err := bar.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed establishing stdout pipe to lemonbar", err)
		return
	}

	barStdin, err := bar.StdinPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed establishing stdin pipe to lemonbar", err)
		barStdout.Close()

		return
	}

	go RenderStatus(renderer, configuration, barStdin)
	
	// wait for any button events
	barScanner := bufio.NewScanner(barStdout)
	err = bar.Start()

	if err != nil {
		panic(err)
	}

	// listens for any events being sent by lemonbar and then processes them accordingly
	for barScanner.Scan() {
		if command, exists := events[barScanner.Text()]; exists {
			eventCommand := exec.Command(command[0], command[1:]...)
			err = eventCommand.Start()
		}
	}
}

func main() {
	renderChannel := make(chan []Module)

	/*
		In the RenderStatus goroutine, it will replace any of the existing blocks 
		with ones that it receives. So, by populating it beforehand we can ensure 
		that the ordering is correct
	*/
	configuration := []Module {
		Module{ ID: 0, Name: "CPU", Align: Left, Content: "" },
		Module{ ID: 1, Name: "CPU_Temp", Align: Left, Content: "" },
		Module{ ID: 2, Name: "Date", Align: Center, Content: "" },
		Module{ ID: 3, Name: "Workspace", Align: Right, Content: "          " },
		Module{ ID: 4, Name: "Power", Align: Right, Content: Button(Color("-", "#EE4B2B", "   "), "power-menu") },
	}

	go StartBar(renderChannel, configuration)

	/* generate each module's status concurrently */
	go Statistics(renderChannel)
	go CurrentWorkspace(renderChannel)

	select { }

	return
}
