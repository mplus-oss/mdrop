package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	internal "github.com/mplus-oss/mdrop/client/internal"
	"golang.org/x/term"
)

func SendCommand(args []string) {
	flag := flag.NewFlagSet("mdrop send", flag.ExitOnError)
	var (
		help	  = flag.Bool("help", false, "Print this message")
		localPort = flag.Int("localPort", 6000, "Specified reciever port on local.")
	)
	flag.Parse(args)

	file := flag.Arg(0)
	if *help || file == "" {
		fmt.Println("Command: mdrop send [options] <file>")
		flag.Usage()	
		os.Exit(1)
	}	

	var c internal.ConfigFile
	err := c.ParseConfig(&c)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Token Prompt
	fmt.Print("Enter Token: ")
	token, err := readToken()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Authenticating
	fmt.Println("Authenticating...")
	joinData, err := sendToken(c, string(token))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Connecting to tunnel...")
	go StartShellTunnel(false, c, *localPort, joinData.Port)
	err = <- sshErrGlobal
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Try to send data to reciever
	fmt.Println("Authenticated. Sending file...")
	uri := fmt.Sprintf(
		"http://localhost:%v/receive?token=%v",
		*localPort,
		token,
	)
	val := map[string]io.Reader{
		"file": mustOpen(file),
	}
	err = Upload(uri, val)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Success!")

	pid := <- sshPidGlobal
	err = KillShell(pid)	
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func readToken() (string, error) {
	stdin := int(syscall.Stdin)
	oldState, err := term.GetState(stdin)
	if err != nil {
		return "", err
	}
	defer term.Restore(stdin, oldState)

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	go func() {
		for range sigch {
			term.Restore(stdin, oldState)
			os.Exit(1)
		}
	}()

	token, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}

	return string(token), nil
}
