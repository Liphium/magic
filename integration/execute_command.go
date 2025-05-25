package integration

import (
	"bufio"
	"os/exec"
)

func ExecCmdWithFunc(funcPrint func(string), name string, args ...string) error {
	cmd, err := execHelper(funcPrint, name, args...)
	if err != nil{
		return err
	}
	return cmd.Run()
}

func ExecCmdWithFuncStart(funcPrint func(string), funcStart func(), name string, args ...string) error {
	cmd, err := execHelper(funcPrint, name, args...)
	if err != nil{
		return err
	}
	if err = cmd.Start(); err != nil{
		return err
	}
	funcStart()
	return cmd.Wait()
}

func execHelper(funcPrint func(string), name string, args ...string) (*exec.Cmd, error){
	cmd := exec.Command(name, args...)

	// Read the normal logs from the app
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil,err
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
		return nil, err
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
	return cmd, nil
}
