package main

import ( 
	"gmisail.me/gbar/blocks"
	"gmisail.me/gbar/modules"
	"gmisail.me/gbar/config"

	"fmt" 
	"strings"
	"bufio"
	"os"
	"os/exec"
	"io"
)


type BlockLocation struct {
	Alignment int 
	Order int
}

/*
func CurrentWorkspace(renderChannel chan []Module) {
	workspaceIds, _ := exec.Command("bspc", "query", "-D").Output()
	workspacesList := strings.Split(string(workspaceIds), "\n")

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
*/

// given a list of modules, renders their text with alignment
func RenderBlocks(config *blocks.Blocks, sep string, stdin io.WriteCloser) {
	var buffer strings.Builder

	buffer.WriteString("%{l}")
	for i, block := range config.Left {
		buffer.WriteString(block.Content)
	
		if i < len(config.Left) - 1 {
			buffer.WriteString(sep)
		}
	}

	buffer.WriteString(" %{c}")
	for i, block := range config.Center {
		buffer.WriteString(block.Content)
		
		if i < len(config.Center) - 1 {
			buffer.WriteString(sep)
		}
	}

	buffer.WriteString(" %{r}")
	for i, block := range config.Right {
		buffer.WriteString(block.Content)
		
		if i < len(config.Right) - 1 {
			buffer.WriteString(sep)
		}
	}

	// writes the data to lemonbar's STDIN
	io.WriteString(stdin, buffer.String())
}

/* Renders the status bar once it receives updated modules */
func RenderStatus(renderer chan blocks.Block, blocks *blocks.Blocks, sep string, stdin io.WriteCloser) {
	locations := make(map[string] BlockLocation)

	/*
     *	Lookup table for block locations, avoids iterating through every block
	 *	on update. Instead, we can use a little more memory to cache where they're stored
	 *	for easy lookup.
	 */
	for i, block := range blocks.Left {
		locations[block.Name] = BlockLocation{ Alignment: 0, Order: i }
	}

	for i, block := range blocks.Center {
		locations[block.Name] = BlockLocation{ Alignment: 1, Order: i }
	}

	for i, block := range blocks.Right {
		locations[block.Name] = BlockLocation{ Alignment: 2, Order: i }
	}

	for {
		updatedBlock := <-renderer

		// Update the value of a given block (or blocks)
		location := locations[updatedBlock.Name]

		switch location.Alignment {
		case 0:
			blocks.Left[location.Order].Content = updatedBlock.Content
		case 1:
			blocks.Center[location.Order].Content = updatedBlock.Content
		case 2:
			blocks.Right[location.Order].Content = updatedBlock.Content
		}

		RenderBlocks(blocks, sep, stdin)
	}
}

func CreateBlocksFromConfig(config config.Configuration) *blocks.Blocks{
	blockConfig := &blocks.Blocks{}
	blockConfig.Left = make([]blocks.Block, len(config.Template.Left))
	blockConfig.Center = make([]blocks.Block, len(config.Template.Center))
	blockConfig.Right = make([]blocks.Block, len(config.Template.Right))

	for i, block := range config.Template.Left {
		blockConfig.Left[i] = blocks.Block{ Name: block, Content: "" }
	}

	for i, block := range config.Template.Center {
		blockConfig.Center[i] = blocks.Block{ Name: block, Content: "" }
	}
	
	for i, block := range config.Template.Right {
		blockConfig.Right[i] = blocks.Block{ Name: block, Content: "" }
	}

	return blockConfig
}

// creates a lemonbar instance and connects the stdin pipe to gbar's renderer
func StartBar(renderer chan blocks.Block, config config.Configuration) {
	buttons := make(map[string] *exec.Cmd)

	/* load events from a configuration file */
	for _, button := range config.Buttons {
		commandArgs := strings.Split(button.Command, " ")
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

	configuration := CreateBlocksFromConfig(config)	

	go RenderStatus(renderer, configuration, config.Settings.Separator, barStdin)
	
	// Wait for any button events
	barScanner := bufio.NewScanner(barStdout)
	err = bar.Start()

	if err != nil {
		panic(err)
	}

	// Listens for any events being sent by lemonbar and then processes them accordingly
	for barScanner.Scan() {
		if command, exists := buttons[barScanner.Text()]; exists {
			err = command.Start()
		}
	}
}

func main() {
	renderer := make(chan blocks.Block)
	config := config.LoadConfig("config.json")

	modulesConfig := make(map[string] modules.Module)
	modulesConfig["cpu"] = modules.CPU{}
	modulesConfig["ram"] = modules.RAM{}
	modulesConfig["time"] = modules.Time{}

/*	configuration := []Module {
		Module{ ID: 3, Name: "Workspace", Align: Right, Content: "          " },
		Module{ ID: 4, Name: "Power", Align: Right, Content: Button(Color("-", "#EE4B2B", "   "), "power-menu") },
	}*/

	go StartBar(renderer, config)

	blocks.CreateBlocks(config.Blocks, modulesConfig, renderer)

	select { }
}
