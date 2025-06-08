package integration

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
)

// Build and then run a go program.
func BuildThenRun(funcPrint func(string), funcStart func(*exec.Cmd), directory string, args ...string) error {

	// Get the old working directory
	workDir, err := os.Getwd()
	if err != nil {
		return err
	}

	// Change directory to the file
	if err := os.Chdir(directory); err != nil {
		return err
	}

	// Build the program
	if err := ExecCmdWithFuncStart(funcPrint, func(c *exec.Cmd) {}, "go", "build", "-o", "program.exe"); err != nil {
		return err
	}

	// Change back to the original working directory
	if err := os.Chdir(workDir); err != nil {
		return err
	}

	// Execute and return the process
	if err := ExecCmdWithFuncStart(funcPrint, funcStart, filepath.Join(directory, "program.exe"), args...); err != nil {
		return err
	}

	return nil
}

func ExecCmdWithFunc(funcPrint func(string), name string, args ...string) error {
	cmd, err := execHelper(funcPrint, name, args...)
	if err != nil {
		return err
	}
	return cmd.Run()
}

func ExecCmdWithFuncStart(funcPrint func(string), funcStart func(*exec.Cmd), name string, args ...string) error {
	cmd, err := execHelper(funcPrint, name, args...)
	if err != nil {
		return err
	}
	if err = cmd.Start(); err != nil {
		return err
	}
	funcStart(cmd)
	return cmd.Wait()
}

func execHelper(funcPrint func(string), name string, args ...string) (*exec.Cmd, error) {
	cmd := exec.Command(name, args...)

	// Read the normal logs from the app
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
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
			funcPrint(scanner.Text())
		}
	}()
	return cmd, nil
}
