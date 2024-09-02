package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const KEY_FILE_LOCATION string = "/etc/mdrop-tunnel-key"

func main() {
	// Parse args
	flag.Parse()
	args := flag.Args()

	// Check key file
	if _, err := os.Stat(KEY_FILE_LOCATION); os.IsNotExist(err) {
		fmt.Println("Key/token file not found at /etc/mdrop-tunnel-key")
		os.Exit(10)
	}
	byteKeyFile, err := os.ReadFile(KEY_FILE_LOCATION)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	keyFile := strings.TrimSpace(string(byteKeyFile))

	// Check is command contains key or not
	if keyFile != "" {
		if !implContains(args, keyFile) {
			fmt.Println("Invalid key")
			os.Exit(10)
		}
	}

	if len(args) < 1 {
		fmt.Println("Usage: mdrop-tunnel <subcommand> [key]")
		os.Exit(1)
	}

	cmd := args[0]
	switch cmd {
	case "ping":
		PingCommand()
	case "get-port":
		GetPortCommand()
	case "connect":
		ConnectCommand()
	default:
		fmt.Println("Subcommand: ping, get-port")
	}
}

func ConnectCommand() {
	for true {
		fmt.Println("Pong.....")
		time.Sleep(10 * time.Second)
	}
}

func GetPortCommand() {
	portFound := 0
	for i := 5000; i <= 59999; i++ {
		timeout := time.Second
		conn, err := net.DialTimeout("tcp", net.JoinHostPort("127.0.0.1", strconv.Itoa(i)), timeout)
		if err != nil {
			fmt.Println("Error [" + strconv.Itoa(i) + "]: " + err.Error())
		}
		if conn == nil {
			portFound = i
			break
		}
		fmt.Println("Get-Port-Closed [" + strconv.Itoa(i) + "]: Continue.")
	}

	if portFound != 0 {
		fmt.Println("Get-Port-Found " + strconv.Itoa(portFound))
	} else {
		fmt.Println("Get-Port-Closed [" + strconv.Itoa(portFound) + "]: Not found.")
		os.Exit(2)
	}
}

func PingCommand() {
	fmt.Println("Pong!")
}

func implContains(sl []string, name string) (found bool) {
	found = false
	for _, value := range sl {
		if value == name {
			found = true
		}
	}
	return found
}
