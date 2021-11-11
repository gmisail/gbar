package main

import ( 
	"fmt" 
	"time"
	"strings"
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

func SetBackground(color string) {
	fmt.Printf("%%{B%s}", color)
}

func Button(text string, command string) {
	fmt.Printf("%%{A:%s:}%s%%{A}", command, text)
}

func Statistics(renderChannel chan []Module) {
	systemStats := Module{ ID: 0, Name: "CPU_RAM", Align: Left, Content: "Default" }
	date := Module{ ID: 1, Name: "Date", Align: Center, Content: "Default" }

	for {
		times, _ := cpu.Percent(0, false)
		cpuUsage := times[0]

		memory, _ := mem.VirtualMemory()

		systemStats.Content = fmt.Sprintf("[CPU: %.2f%%] [RAM: %.2f%%]", cpuUsage, memory.UsedPercent)

		currentTime := time.Now()
		date.Content = fmt.Sprintf(currentTime.Format("[Monday, January 2] [3:04:05 PM]"))

		copyModules := make([]Module, 2)
		copyModules[0] = systemStats
		copyModules[1] = date
		renderChannel <- copyModules

		time.Sleep(time.Second)
	}
}

func AlignRight() {
	fmt.Print("%{r}")

	workspace, err := exec.Command("bspc", "query", "-D", "-d", "focused", "--names").Output()

	if err != nil {
		fmt.Println("Can't find current workspace.")
	}

	/* remove the newline from the workspace ID */
	workspaceId := strings.TrimSuffix(string(workspace), "\n")

	fmt.Printf("[%s] ", workspaceId)

	Button("[App Launcher] ", "rofi -show drun")
}

func Process() {
//	AlignLeft()
//	AlignCenter()
//	AlignRight()

	// flush the buffers
//	fmt.Println()
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

func RenderStatus(renderChannel chan []Module, config []Module) {
	// date => { name: "date", content: "Feb 1 2021" }
	var modules []Module

	for i := 0; i < len(config); i++ {
		modules = append(modules, config[i])
	}

	RenderModules(modules)

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
		Module{ ID: 0, Name: "CPU", Align: Left, Content: "cpu stuff" },
		Module{ ID: 1, Name: "Date", Align: Center, Content: "date stuff" },
		Module{ ID: 2, Name: "Workspace", Align: Right, Content: "workspace" },
	}

	go RenderStatus(renderChannel, configuration)
	go Statistics(renderChannel)

	select { }

	return
}
