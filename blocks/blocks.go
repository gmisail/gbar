package blocks

import (
	"gmisail.me/gbar/config"
	"gmisail.me/gbar/modules"

	"os"
	"os/exec"
	"fmt"
	"strings"
	"strconv"
	"time"
)

type Blocks struct {
	Left []Block
	Center []Block
	Right []Block
}

type Block struct {
	Name string
	Content string
}

/*
 *	Setup blocks by creating a goroutine for each
 */
func CreateBlocks(blocks map[string] config.ConfigBlock, modules map[string] modules.Module, renderer chan Block) {
	for id, block := range blocks {
		go RunBlock(id, block, modules, renderer)
	}	
}

/*
 *	Create a command from a given string
 */
func CreateCommand(command string) *exec.Cmd {
	components := strings.Fields(command)
	executable := components[0]
	args := components[1:]

	return exec.Command(executable, args...)
}

/*
 *	Update the bar based on the STDOUT from the given command.
 */
func UpdateOnCommandData(block config.ConfigBlock, renderer chan Block) {
//	command := CreateCommand(block.Command)

	// TODO: get the stdout pipe from the command and 
	// watch for any output. When output is received, 
	// apply it to the given template.
}

/*
 *	Update a given block after a given increment of time. For
 *	instance, "1" will update every second while "5" will update
 *  every 5 seconds
 */
func UpdateOnInterval(id string, block config.ConfigBlock, modules map[string] modules.Module, renderer chan Block) {
	interval, err := strconv.Atoi(block.Interval)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not start block '%s' with invalid interval %s.", block.Name, block.Interval)
		return
	}

	hasCommand := len(block.Command) > 0
	//	command := hasCommand ? CreateCommand(block.Command) : nil
	module := modules[block.Module]

	/* Over a given interval, either call a function or run a module and assign the bar's content */
	content := ""

	for {
		if hasCommand {
//			content := command.Output()
		} else {
			if module != nil {
				content = module.Run()	
			}
		}

		renderer <- Block{ Name: id, Content: content }
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func RunBlock(id string, block config.ConfigBlock, modules map[string] modules.Module, renderer chan Block) {
	/*
	 *	If this block waits for data and has a command, then spawn a special process
	 *	which will update the block only when the command updates
	 */
	if block.Interval == "ondata" && len(block.Command) > 0 {
		UpdateOnCommandData(block, renderer)
	} else {
		UpdateOnInterval(id, block, modules, renderer)
	}
}
