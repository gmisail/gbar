package main

import ( 
	"fmt" 
	"time"
	"strings"
	"bufio"
	"os"
	"os/exec"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
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
	fmt.Printf("%%{A:%s:}%s%%{A}", command, text)
}

func Statistics(renderChannel chan []Module) {
	systemStats := Module{ ID: 0, Name: "CPU_RAM", Align: Left, Content: "Default" }
	date := Module{ ID: 1, Name: "Date", Align: Center, Content: "Default" }

	for {
		times, _ := cpu.Percent(0, false)
		cpuUsage := times[0]

		memory, _ := mem.VirtualMemory()

		systemStats.Content = ""
		systemStats.Content += Color("#D6E3F8", "#000000", "   ")
		systemStats.Content += SetColor("#141414", "#ffffff") + fmt.Sprintf(" %.2f%% ", cpuUsage) + SetColor("-", "-")

		systemStats.Content += " "

		systemStats.Content += SetColor("#CC3F0C", "#000000") + "  " + SetColor("-", "-")
		systemStats.Content += SetColor("#141414", "#ffffff") + fmt.Sprintf(" %.2f%% ", memory.UsedPercent) + SetColor("-", "-")

		currentTime := time.Now()
		
		date.Content = ""
		date.Content += SetColor("#F3EFF5", "#000000") + "  " + SetColor("-", "-")
		date.Content += SetColor("#141414", "#ffffff") + fmt.Sprintf(currentTime.Format(" Monday, January 2 ")) + SetColor("-", "-")
		date.Content += " "

		date.Content += SetColor("#D0C0D8", "#000000") + "  " + SetColor("-", "-")
		date.Content += SetColor("#141414", "#ffffff") + fmt.Sprintf(currentTime.Format(" 3:04:05 PM ")) + SetColor("-", "-")

		copyModules := make([]Module, 2)
		copyModules[0] = systemStats
		copyModules[1] = date
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

	workspace := Module{ ID: 2, Align: Right, Content: "" }
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
func RenderModules(modules []Module) {
	var alignment Alignment = None

	for i := 0; i < len(modules); i++ {
		if modules[i].Align != alignment {
			alignment = modules[i].Align
			switch modules[i].Align {
			case Left:
				fmt.Print("%{l}")
			case Center:
				fmt.Print("%{c}")
			case Right:
				fmt.Print("%{r}")
			}
		}
		fmt.Print(modules[i].Content)
	}

	fmt.Println()
}

/* Renders the status bar once it receives updated modules */
func RenderStatus(renderChannel chan []Module, config []Module) {
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

		RenderModules(modules)
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
		Module{ ID: 1, Name: "Date", Align: Center, Content: "" },
		Module{ ID: 2, Name: "Workspace", Align: Right, Content: "" },
	}

	go RenderStatus(renderChannel, configuration)

	/* generate each module's status concurrently */
	go Statistics(renderChannel)
	go CurrentWorkspace(renderChannel)

	select { }

	return
}
