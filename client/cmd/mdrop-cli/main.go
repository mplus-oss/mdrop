package main

import (
	"fmt"
	"os"

	serviceGet "github.com/mplus-oss/mdrop/internal/service_get"
)

func main() {
	cmdList := NewCommandList()
	cmdList.list["get"] = CommandListStructure{
		name: "get",
		desc: "Get file from target server",
		start: serviceGet.StartGetService,
	}
	cmdList.list["transfer"] = CommandListStructure{
		name: "transfer",
		desc: "Transfer file from target client",
		start: func() {
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

type CommandListStructure struct {
	name string
	desc string
	start func()
}

type CommandList struct {
	list map[string]CommandListStructure
}

func NewCommandList() (c CommandList) {
	c = CommandList{}
	c.list = make(map[string]CommandListStructure)
	return c
}

func (c CommandList) PrintUsageCommand() {
	fmt.Print("Usage:\n  mdrop <command> <args> [options]\n\n")
	fmt.Println("Command available:")
	for _, v := range c.list {
		fmt.Println(" ", v.name, ":", v.desc)
	}
	fmt.Println()
}

func (c CommandList) IsCommandAvailable(cmd string) (result bool) {
	result = false
	if _, ok := c.list[cmd]; ok {
		result = true
	}

	return result
}
