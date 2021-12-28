package blocks

import "strconv"

type Blocks struct {
	left []Block
	center []Block
	right []Block
}

type Block struct {
	Name string
	Content string
}


func CreateBlocks(blocks []ConfigBlock) *[]Block {
	components := make([]Block, len(blocks))

	for i, block := range blocks {
		// 
	}
}

/*
 *	Create a command from a given string
 */
func CreateCommand(command string) *Cmd {
	components := strings.Fields(command)
	executable := components[0]
	args := components[1:]

	return exec.Command(executable, args...)
}

/*
 *	Update the bar based on the STDOUT from the given command.
 */
func UpdateOnCommandData(block ConfigBlock, renderer chan []Block) {
	command := CreateCommand(block.Command)

	// TODO: get the stdout pipe from the command and 
	// watch for any output. When output is received, 
	// apply it to the given template.
}

/*
 *	Update a given block after a given increment of time. For
 *	instance, "1" will update every second while "5" will update
 *  every 5 seconds
 */
func UpdateOnInterval(block ConfigBlock, renderer chan []Block) {
	interval, err := strconv.Atoi(block.Interval)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not start block '%s' with invalid interval %s.", block.Name, block.Interval)
		return
	}

	hasCommand := len(block.Command) > 0
	command := hasCommand ? CreateCommand(block.Command) : nil
		
	// TODO: get module type and assign it to an internal module. If one does 
	// not exist then throw an error.

	/* Over a given interval, either call a function or run a module and assign the bar's content */
	content := ""

	for {
		if hasCommand {
			content := command.Output()
		}

		renderer <- Block{ Name: block.Name, Content: content }
		time.Sleep(interval * time.Second)
	}
}

func Run(block ConfigBlock, renderer chan []Block) {
	/*
	 *	If this block waits for data and has a command, then spawn a special process
	 *	which will update the block only when the command updates
	 */
	if block.Interval == "ondata" && len(block.Command) > 0 {
		UpdateOnCommandData(block)
	} else {
		UpdateOnInterval(block)
	}
}
