package util

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/bytedance/sonic"
)

var WorkingDirectory = ""
var BackendURL = ""

func InitWorkingDirectory() {
	if os.Getenv("SC_WORKDIR") == "" {
		log.Fatalln("Please set the working directory using the SC_WORKDIR environment variable.")
	}

	// Change to the working directory
	os.Chdir(os.Getenv("SC_WORKDIR"))
	WorkingDirectory = os.Getenv("SC_WORKDIR")
}

// Send a post request to the backend
func PostRequestBackend(path string, body map[string]interface{}) (map[string]interface{}, error) {

	byteBody, err := sonic.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader := strings.NewReader(string(byteBody))

	res, err := http.Post(BackendURL+path, "application/json", bodyReader)
	if err != nil {
		return nil, err
	}

	// Decrypt the request body
	defer res.Body.Close()
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, res.Body)
	if err != nil {
		return nil, err
	}

	// Parse decrypted body into JSON
	var data map[string]interface{}
	err = sonic.Unmarshal(buf.Bytes(), &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Send a post request to the backend
func PostRequestBackendToken(path string, token string, body map[string]interface{}) (map[string]interface{}, error) {

	// Encode the
	byteBody, err := sonic.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader := strings.NewReader(string(byteBody))

	// Create a new request
	req, err := http.NewRequest(http.MethodPost, BackendURL+path, bodyReader)
	if err != nil {
		return nil, err
	}

	// Set the required headers
	req.Header["Content-Type"] = []string{"application/json"}
	req.Header["SC-Token"] = []string{token}

	// Execute the request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	// Copy into a buffer
	defer res.Body.Close()
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, res.Body)
	if err != nil {
		return nil, err
	}

	// Parse decoded body into JSON
	var data map[string]interface{}
	err = sonic.Unmarshal(buf.Bytes(), &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Delete all the files in the specified directory.
func DeleteAllFiles(dir string) error {

	// List all the files in the directory
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	// Delete all of them
	for _, entry := range entries {
		err = os.RemoveAll(filepath.Join(dir, entry.Name()))
		if err != nil {
			return err
		}
	}
	return nil
}
