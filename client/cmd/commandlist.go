package main

import (
	"fmt"
)

type CommandListStructure struct {
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
	for k, v := range c.list {
		fmt.Println(" ", k, ":", v.desc)
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
