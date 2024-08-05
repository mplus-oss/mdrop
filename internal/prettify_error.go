package internal

import (
	"fmt"
	"os"
)

func PrintErrorWithExit(instance string, err error, exitCode int) {
    fmt.Println("\nFatal Error on "+instance+":\n\t"+err.Error())
    os.Exit(exitCode)
}

func PrintError(instance string, err error) {
    fmt.Println("\nFatal Error on "+instance+":\n\t"+err.Error())
}
