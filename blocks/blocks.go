package blocks

import "strconv"

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

func CreateCommand(command string) *Cmd {
	// TODO: parse command into exec.Cmd
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

func UpdateOnInterval(block ConfigBlock, renderer chan []Block) {
	interval, err := strconv.Atoi(block.Interval)

	hasCommand := len(block.Command) > 0
	command := hasCommand ? CreateCommand(block.Command) : nil
		
	// TODO: get module type and assign it to an internal module. If one does 
	// not exist then throw an error.

	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not start block '%s' with invalid interval %s.", block.Name, block.Interval)
		return
	}

	for {
		content := ""
		
		if hasCommand {
			// TODO: run the command that is given and apply the result to the template
			// command.Run()
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