package util

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type Headers = map[string]string

// Send a post request to any URL with headers attached
func Post[T any](url string, body interface{}, headers Headers) (T, error) {

	// Declared here so it can be returned as nil before it's actually used
	var data T

	// Encode body to JSON
	byteBody, err := json.Marshal(body)
	if err != nil {
		return data, err
	}

	// Set headers
	reqHeaders := http.Header{}
	reqHeaders.Set("Content-Type", "application/json")
	for key, value := range headers {
		reqHeaders.Set(key, value)
	}

	// Send the request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(byteBody))
	if err != nil {
		return data, err
	}
	req.Header = reqHeaders

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return data, err
	}

	// Use the extracted function to read and parse the response body
	return readResponseBody[T](res)
}

// Send a get request to any URL with headers attached
func Get[T any](url string, headers Headers) (T, error) {

	// Declared here so it can be returned as nil before it's actually used
	var data T

	// Set headers
	reqHeaders := http.Header{}
	for key, value := range headers {
		reqHeaders.Set(key, value)
	}

	// Send the request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return data, err
	}
	req.Header = reqHeaders

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return data, err
	}

	// Use the extracted function to read and parse the response body
	return readResponseBody[T](res)
}

// readResponseBody reads the HTTP response body and unmarshals it into the provided type
func readResponseBody[T any](res *http.Response) (T, error) {
	var data T

	// Grab all bytes from the buffer
	defer res.Body.Close()
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, res.Body)
	if err != nil {
		return data, err
	}

	// Parse body into JSON
	err = json.Unmarshal(buf.Bytes(), &data)
	if err != nil {
		return data, err
	}
	return data, nil
}
