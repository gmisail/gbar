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

		systemStats.Content = fmt.Sprintf("%%{B#D6E3F8}%%{F#000000} CPU: %.2f%% %%{B-}%%{F-}%%{B#FEF5EF}%%{F#000000} RAM: %.2f%% %%{B-}%%{F-}", cpuUsage, memory.UsedPercent)

		currentTime := time.Now()
		date.Content = fmt.Sprintf(currentTime.Format("[Monday, January 2] [3:04:05 PM]"))

		copyModules := make([]Module, 2)
		copyModules[0] = systemStats
		copyModules[1] = date
		renderChannel <- copyModules

		time.Sleep(time.Second)
	}
}

func CurrentWorkspace(renderChannel chan []Module) {
	occupiedWorkspace := ""
	//emptyWorkspace := "_"

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

	subscriptionScanner := bufio.NewScanner(subscriptionPipe)
	subscription.Start()

	// block for changes. When a change arrives, update it and push to the bar
	for subscriptionScanner.Scan() {
		// split up the command into arguments
		workspaceEvent := strings.Split(subscriptionScanner.Text(), " ")

		// look up which desktop the new ID belongs to
		currentWorkspace := workspaces[workspaceEvent[2]]

		workspace.Content = ""
//		for i := 0; i < len(workspaces); i++ {
			workspace.Content = fmt.Sprintf("%s %d ", occupiedWorkspace, currentWorkspace)
//		}

		modules[0] = workspace
		renderChannel <- modules
    }
}

/*
Example of a module:

func Test(renderChannel chan []Module) {
	i := 0
	for {
		testModules := make([]Module, 1)
		testModules[0] = Module{ ID: 3, Name: "Test", Align: Right, Content: fmt.Sprintf("Current: %d", i) }
		renderChannel <- testModules
		i = i + 1
		time.Sleep(time.Second / 2)
	}
}
*/

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
