package blocks

import (
	"gmisail.me/gbar/config"
	"gmisail.me/gbar/modules"
	"gmisail.me/gbar/style"
	"github.com/valyala/fasttemplate"

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
 *	Substitute the values from a given module into a block's template 
 *	and return the result
 */
func RenderTemplate(template string, data map[string] interface{}) string {
	t := fasttemplate.New(template, "<", ">")
	return t.ExecuteString(data)
}

func RenderBlock(id string, block config.ConfigBlock, module modules.Module) string {
	var content string = block.Template

	if module != nil {
		content = RenderTemplate(block.Template, module.Run())
	}

	/* this block is a button, so wrap it in a button with the correct click event ID */
	if len(block.OnClick) != 0 {
		return style.Button(content, id)
	}

	return content
}

/*
 *	Handles blocks that are generated once and then never update, i.e. number of cores
 */
func RenderStatic(id string, block config.ConfigBlock, modules map[string] modules.Module, renderer chan Block) {
	module := modules[block.Module]

	var content string = RenderBlock(id, block, module) //RenderTemplate(block.Template, module.Run())
	renderer <- Block{ Name: id, Content: content }
}

/*
 *	Update the bar based on the STDOUT from the given command.
 */
func UpdateOnData(id string, block config.ConfigBlock, modules map[string] modules.Module, renderer chan Block) {
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
		fmt.Fprintf(os.Stderr, "Could not start block '%s' with invalid interval %s.\n", id, block.Interval)
		return
	}

	module := modules[block.Module]

	var content string

	for {
		content = RenderBlock(id, block, module)
		renderer <- Block{ Name: id, Content: content }
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func RunBlock(id string, block config.ConfigBlock, modules map[string] modules.Module, renderer chan Block) {
	/*
	 *	There are two types of blocks: ones that run on an interval and ones that wait for data. Since
	 *	these are handled differently, the logic is split into two different functions.
	 */
	if block.Interval == "ondata" {
		UpdateOnData(id, block, modules, renderer)
	} else if len(block.Interval) > 0 {
		UpdateOnInterval(id, block, modules, renderer)
	} else {
		RenderStatic(id, block, modules, renderer)
	}
}
