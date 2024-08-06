package main

import (
	"bufio"
	"errors"
	"fmt"
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

	args := ""
	if isTunnel {
		args = s.GenerateRemoteSSHArgs()
	} else {
		args = s.GenerateConnectSSHArgs()
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
		}
	}()

	go func() {
		s := bufio.NewScanner(stderr)
		s.Split(bufio.ScanLines)
		for s.Scan() {
			m := s.Text()
			// False flag #1
			if strings.Contains(m, "chdir") || strings.Contains(m, "Permanently added") {
				continue
			}
			if strings.Contains(m, "remote port forwarding failed") {
				errChan <- errors.New("Duplicate remote on bridge server")
			}
			if strings.Contains(m, "EXCEPTION") {
				errChan <- errors.New("Error from Server: " + m)
			}
			if strings.Contains(m, "exec") && strings.Contains(m, "not found") {
				errChan <- errors.New("Proxy app not found. Did you install it?")
			}
			fmt.Println(m)
		}
	}()

	output = <-outputChan
	err = <-errChan
	cmd.Wait()

	return output, err
}

