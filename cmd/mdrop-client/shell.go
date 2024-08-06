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

var sshErrGlobal chan error = make(chan error)
var sshPidGlobal chan int = make(chan int)

func KillShell(pid int) error {
	cmd := exec.Command("kill", "-9", strconv.Itoa(pid))
	err := cmd.Start()
	if err != nil {
		return err
	}
	cmd.Wait()
	return nil
}

func StartShellTunnel(s internal.SSHParameter, isTunnel bool) {
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
		sshErrGlobal <- err
	}
	cmd.Start()

	go func() {
		s := bufio.NewScanner(stdout)
		s.Split(bufio.ScanLines)
		for s.Scan() {
			m := s.Text()
			if strings.Contains(m, "Pong!") {
				sshErrGlobal <- nil
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
				sshErrGlobal <- errors.New("Duplicate remote on bridge server")
			}
			if strings.Contains(m, "EXCEPTION") {
				sshErrGlobal <- errors.New("Error from Server: " + m)
			}
			if strings.Contains(m, "exec") && strings.Contains(m, "not found") {
				sshErrGlobal <- errors.New("Proxy app not found. Did you install it?")
			}
			fmt.Println(m)
		}
	}()

	sshPidGlobal <- cmd.Process.Pid
	cmd.Wait()
}

