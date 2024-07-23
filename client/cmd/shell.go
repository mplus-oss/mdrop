package main

import (
	"bufio"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/mplus-oss/mdrop/client/internal"
)

var sshErrGlobal chan error = make(chan error)
var sshPidGlobal chan int = make(chan int)

func StartShellTunnel(isRemote bool, c internal.ConfigFile, localPort int, remotePort int) {
	args := strings.Split(internal.GenerateSSHArgs(isRemote, c, localPort, remotePort), " ")
	cmd := exec.Command("ssh", args...)

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
			if strings.Contains(m, "Connected!") {
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
			fmt.Println(m)
		}
	}()

	sshPidGlobal <- cmd.Process.Pid
	cmd.Wait()
}

func KillShell(pid int) error {
	cmd := exec.Command("kill", "-9", strconv.Itoa(pid))
	err := cmd.Start()
	if err != nil {
		return err
	}
	cmd.Wait()
	return nil
}