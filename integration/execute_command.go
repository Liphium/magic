package integration

import (
	"bufio"
	"os/exec"
)

func ExecCmdWithFunc(funcPrint func(string), shouldReturn bool, name string, args ...string) error {
	cmd := exec.Command(name, args...)

	// Read the normal logs from the app
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			funcPrint(scanner.Text())
		}
	}()

	// Read the errors output from the app
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			if name != "go" {
				funcPrint("ERROR: " + scanner.Text())
			} else {
				funcPrint(scanner.Text())
			}
		}
	}()

	if shouldReturn {
		// Start the command
		if err := cmd.Start(); err != nil {
			return err
		}
	} else {
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}
