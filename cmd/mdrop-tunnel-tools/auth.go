package main

import (
	"fmt"
	"os"
)

func AuthCommand(privToken string, args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: mdrop-tunnel auth <token>")
		os.Exit(1)
	}

	token := args[0]
	if privToken != token {
		fmt.Println("Token unmatched. Not authenticated.")
		os.Exit(1)
	}

	fmt.Println("Token matched!")
}
