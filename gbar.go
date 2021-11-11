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
	Left Alignment = 0
	Center
	Right
)

type Module struct {
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

func AlignLeft() {
	fmt.Print("%{l}")

	times, _ := cpu.Percent(0, false)
	cpuUsage := times[0]

	memory, _ := mem.VirtualMemory()

	fmt.Printf(" [CPU: %.2f%%] [RAM: %.2f%%]", cpuUsage, memory.UsedPercent)
}

func AlignCenter() {
	fmt.Print("%{c}")
	currentTime := time.Now()
	fmt.Print(currentTime.Format("[Monday, January 2] [3:04:05 PM]"))
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
	AlignLeft()
	AlignCenter()
	AlignRight()

	// flush the buffers
	fmt.Println()
}

func RenderStatus(renderChannel chan Module[]) {
	// date => { name: "date", content: "Feb 1 2021" }
	modules := map[string]Module

	for {
		updatedModules := <- modules:

		// update modified modules
		// print the status bar
	}
}

func main() {
	renderChannel := make(chan Module[])

	go RenderStatus(renderChannel)

	for {
		SetBackground("#162424")
		Process()
		time.Sleep(time.Second)
	}
}
