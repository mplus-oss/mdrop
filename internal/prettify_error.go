package internal

import (
	"errors"
	"fmt"
	"os"
)

func PrintErrorWithExit(instance string, err error, exitCode int) {
	fmt.Println("\nFatal Error on " + instance + ":\n\t" + err.Error())
	os.Exit(exitCode)
}

func CustomizeError(code string, err error) error {
	return errors.New(
		fmt.Sprintf("%v: %v", code, err.Error()),
	)
}
