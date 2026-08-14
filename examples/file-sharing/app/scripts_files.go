package app

import (
	"fmt"
	"log"
	"os"

	"github.com/Liphium/magic/v4/mconfig"
	"github.com/Liphium/magic/v4/mrunner"
	"resty.dev/v3"
)

// This file defines scripts, tools you can run when your app is active or use during tests.
//
// Learn more at https://liphium.dev/magic/documentation/magic-scripts/

type UploadData struct {
	UploadPath string `prompt:"File to upload"`
	Target     string `prompt:"Path to upload to"`
}

// This is a function, registered as a script, that can upload a file from somewhere on your machine.
func UploadFile(data UploadData) error {
	content, err := os.ReadFile(data.UploadPath)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	client := resty.New()
	defer client.Close()

	res, err := client.R().
		SetBody(content).
		Post("http://" + os.Getenv("LISTEN") + "/" + data.Target)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	log.Println("Output:", res.Status())
	return nil
}

type DownloadData struct {
	Path   string `prompt:"Path to download"`
	SaveTo string `prompt:"Path to save the file at"`
}

// This is a function, registered as a script, that can download a certain path to somewhere on your machine.
func DownloadFile(data DownloadData) error {
	client := resty.New()
	defer client.Close()

	res, err := client.R().
		Get("http://" + os.Getenv("LISTEN") + "/" + data.Path)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	if !res.IsStatusSuccess() {
		return fmt.Errorf("server returned error status: %s", res.Status())
	}

	if err := os.WriteFile(data.SaveTo, res.Bytes(), 0644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	log.Println("Output:", res.Status())
	return nil
}

// This is a function, registered as a script, that clears all of the files out of SeaweedFS.
func ClearFiles(runner *mrunner.Runner) error {
	return runner.RunInstruction(mconfig.InstructionClearFiles)
}
