package util

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/bytedance/sonic"
)

const BackendURL = "http://localhost:8081"

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
