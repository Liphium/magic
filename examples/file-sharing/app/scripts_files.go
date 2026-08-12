package app

import (
	"fmt"
	"log"
	"os"

	"resty.dev/v3"
)

func UploadFile(data struct {
	UploadPath string `prompt:"File to upload"`
	Target     string `prompt:"Path to upload to"`
}) error {
	content, err := os.ReadFile(data.UploadPath)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	client := resty.New()
	defer client.Close()

	res, err := client.R().
		SetBody(content).
		Post("http://" + os.Getenv("LISTEN") + "/" + data.UploadPath)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	log.Println("Output:", res.Status())
	return nil
}

func DownloadFile(data struct {
	Path   string `prompt:"Path to download"`
	SaveTo string `prompt:"Path to save the file at"`
}) error {
	client := resty.New()
	defer client.Close()

	res, err := client.R().
		Get("http://" + os.Getenv("LISTEN") + "/" + data.Path)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}

	if err := os.WriteFile(data.SaveTo, res.Bytes(), 0644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	log.Println("Output:", res.Status())
	return nil
}
