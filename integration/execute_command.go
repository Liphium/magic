package integration

import (
	"bufio"
	"fmt"
	"os/exec"
)

func ExecCmdWithFunc(funcPrint func(string), shouldReturn bool, name string, args ...string) error {
	cmd := exec.Command(name, args...)

	// Set up the output pipe
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

	if shouldReturn {
		// Start the command
		if err := cmd.Start(); err != nil {
			return err
		}
	} else {
		fmt.Println("hi")
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	fmt.Println("finished running")

	return nil
}
