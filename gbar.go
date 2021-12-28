package main

import ( 
	"blocks"

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

type BlockLocation struct {
	Alignment int 
	Order int
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
func RenderBlocks(config *[]blocks.Blocks, stdin io.WriteCloser) {
	var buffer strings.Builder

	// for each module, update its content
	for i := 0; i < len(config); i++ {
		switch modules[i].Align {
			case Left:
				buffer.WriteString("%{l}")
			case Center:
				buffer.WriteString("%{c}")
			case Right:
				buffer.WriteString("%{r}")
		}

		for _, block := range(config[i]) {
			buffer.WriteString(modules[i].Content)

			// TODO: add separators between blocks
		}
	}

	// end of line, updates the bar with the new data
	buffer.WriteString("\n")

	// writes the data to lemonbar's STDIN
	io.WriteString(stdin, buffer.String())
}

/* Renders the status bar once it receives updated modules */
func RenderStatus(renderer chan []Module, config []blocks.Blocks, stdin io.WriteCloser) {
	locations := map[string] BlockLocation

	/*
     *	Lookup table for block locations, avoids iterating through every block
	 *	on update. Instead, we can use a little more memory to cache where they're stored
	 *	for easy lookup.
	 */
	for i, group := range config {
		for j, block := range group {
			locations[block.Name] = BlockLocation{ Alignment: i, Order: j }
		}
	}

	for {
		updatedBlocks := <-renderer

		// Update the value of a given block (or blocks)
		for _, block := range updatedBlocks {
			location := locations[block.Name]
			config[location.Alignment][location.Order].Content = block.Content
		}

		RenderBlocks(&config, stdin)
	}
}

// creates a lemonbar instance and connects the stdin pipe to gbar's renderer
func StartBar(renderer chan []Module, configuration []Module, config Configuration) {
	buttons := make(map[string] []*exec.Cmd)

	/* load events from a configuration file */
	for _, button := range config.Buttons {
		commandArgs := strings.Split(button.OnClick, " ")
		buttons[button.Name] = exec.Command(commandArgs[0], commandArgs[1:]...)
	}

	barExec := strings.Split(config.Settings.Lemonbar, " ")

	if len(config.Settings.Font) > 0 {
		barExec = append(barExec, "-f", config.Settings.Font)
	}

	bar := exec.Command(barExec[0], barExec[1:]...)

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
	
	// Wait for any button events
	barScanner := bufio.NewScanner(barStdout)
	err = bar.Start()

	if err != nil {
		panic(err)
	}

	// Listens for any events being sent by lemonbar and then processes them accordingly
	for barScanner.Scan() {
		if command, exists := events[barScanner.Text()]; exists {
			err = command.Start()
		}
	}
}

func main() {
	renderer := make(chan []Module)
	config := LoadConfig("config.json")

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

	blocks := &blocks.Blocks{}

	/* switch these with config blocks? */
	blocks.left = []blocks.Block{ 
		blocks.Block{ Name: "CPU", Content: "0.0%" } 
	}

	blocks.center = []blocks.Block{
		blocks.Block{ Name: "CPU", Content: "0.0%" } 
	}
	
	blocks.right = []blocks.Block{
		blocks.Block{ Name: "CPU", Content: "0.0%" } 
	}

	go StartBar(renderer, configuration, config)

	CreateBlocks(config.Blocks, renderer)

	select { }

	return
}
