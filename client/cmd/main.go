package main

import (
	"fmt"
	"os"
)

func main() {
	cmdList := NewCommandList()
	cmdList.list["auth"] = CommandListStructure{
		desc: "Get file from target server",
		start: func () {
			fmt.Println("Not implemented")
		},
	}

	args := os.Args
	if len(args) <= 1 {
		fmt.Println("No command provided.")
		cmdList.PrintUsageCommand()
		os.Exit(1)
	}

	subCommand := &args[1]
	if !cmdList.IsCommandAvailable(*subCommand) {
		fmt.Print(fmt.Sprintf("No command available for \"%s\".\n\n", *subCommand))
		cmdList.PrintUsageCommand()
		os.Exit(1)
	}

	cmdList.list[*subCommand].start()
}
