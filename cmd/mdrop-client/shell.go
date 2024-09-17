package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/mplus-oss/mdrop/internal"
)

func KillShell(pid int) error {
	cmd := exec.Command("kill", "-9", strconv.Itoa(pid))
	err := cmd.Start()
	if err != nil {
		return err
	}
	cmd.Wait()
	return nil
}

func StartShellTunnel(s internal.SSHParameter, isTunnel bool) (output string, err error) {
	errChan := make(chan error, 0)
	outputChan := make(chan string, 0)
	debugEnv := os.Getenv("MDROP_SSH_DEBUG")

	args := ""
	if isTunnel {
		args = s.GenerateRemoteSSHArgs()
	} else {
		args = s.GenerateConnectSSHArgs()
	}

	if debugEnv == "1" {
		fmt.Println(args)
	}
	cmd := exec.Command("sh", "-c", args)
	stdout, err := cmd.StdoutPipe()
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return output, err
	}
	cmd.Start()

	go func() {
		s := bufio.NewScanner(stdout)
		s.Split(bufio.ScanLines)
		for s.Scan() {
			m := s.Text()
			if debugEnv == "1" {
				fmt.Println(m)
			}

			// NOTE: There's some reasone we're keeping this code. But we'll see it later.
			// Case: If connected
			if strings.Contains(m, "Pong!") {
				outputChan <- m
				errChan <- nil
			}
			// Case: If port found
			if strings.Contains(m, "Get-Port-Found") {
				outputChan <- m
				errChan <- nil
			}
			// Case: If port full
			if strings.Contains(m, "Get-Port-Closed") && strings.Contains(m, "Not found") {
				outputChan <- "Not found"
				errChan <- errors.New("Port full on server")
			}
			// Case: Invalid key
			if strings.Contains(m, "Invalid key") {
				outputChan <- "Invalid key"
				errChan <- errors.New("Invalid key on execution")
			}
		}
	}()

	go func() {
		s := bufio.NewScanner(stderr)
		s.Split(bufio.ScanLines)
		for s.Scan() {
			m := s.Text()
			if debugEnv == "1" {
				fmt.Println(m)
			}

			// False flag #1
			if strings.Contains(m, "chdir") || strings.Contains(m, "Permanently added") {
				continue
			}

			// TODO: Remove it for testing. If stable, remove this.
			//if strings.Contains(m, "remote port forwarding failed") {
			//	errChan <- errors.New("Duplicate remote on bridge server")
			//}
			//if strings.Contains(m, "EXCEPTION") {
			//	errChan <- errors.New("Error from Server: " + m)
			//}
			//if strings.Contains(m, "exec") && strings.Contains(m, "not found") {
			//	errChan <- errors.New("Proxy app not found. Did you install it?")
			//}
			fmt.Println(m)
		}
	}()

	err = cmd.Wait()

	// Invoke this if there's an error
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			outputChan <- "Exit with error code " + strconv.Itoa(exitErr.ExitCode())
			errChan <- exitErr
		}
	}
	output = <-outputChan
	err = <-errChan

	return output, err
}
