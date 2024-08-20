package internal

import (
	"errors"
	"fmt"
	"os"
)

func PrintErrorWithExit(instance string, err error, exitCode int) {
	debug_env := os.Getenv("MDROP_DEBUG")

	message := fmt.Sprintf("Error: %v", err.Error())
	if debug_env == "1" {
		message = fmt.Sprintf("\nFatal Error on %v:\n\t%v", instance, err.Error())
	}

	fmt.Println(message)
	os.Exit(exitCode)
}

func CustomizeError(code string, err error) error {
	return errors.New(
		fmt.Sprintf("%v: %v", code, err.Error()),
	)
}
