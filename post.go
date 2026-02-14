package main

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"syscall"
)

var (
	FILE_SPEC_JSON = "spec.json"
)

// Makes a post request to the `Endpoint` provided in cmd
// A spec.json file is first opened anytime a post request
// If the file does not exist, the request is aborted
//
// The file should contain the body to send to the server,
// it could be empty or not
//
// Redirects are followed and the default client is used
func (cmd Command) MakePostRequest() (ServerResponse, error) {
	// so we need to read spec.json
	jsonContent, err := readSpecFile()
	if err != nil {
		return ServerResponse{}, err
	}

	res, err := http.Post(cmd.Endpoint, cmd.ContentType, bytes.NewReader(jsonContent))
	if err != nil && errors.Is(err, syscall.ECONNREFUSED) {
		log.Fatalf("  Server is not active, please make sure it's running")
	} else if err != nil {
		return ServerResponse{}, err
	}

	defer res.Body.Close()
	content, err := io.ReadAll(res.Body)
	if err != nil {
		return ServerResponse{}, err
	}

	return ServerResponse{StatusCode: res.StatusCode, Body: string(content)}, nil
}

// Makes a post request to the `Endpoint` provided in cmd
// A spec.json file is first opened anytime a post request
// If the file does not exist, the request is aborted
//
// The file should contain the body to send to the server,
// it could be empty or not
// 
// A custom client is made to not follow 
// redirects except configured in the `cmd` receiver
//
func (cmd Command) MakePostRequestNoRedirect() (ServerResponse, error) {
	client := &http.Client{}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	jsonContent, err := readSpecFile()
	if err != nil {
		return ServerResponse{}, err
	}
	request, err := http.NewRequest(cmd.Method, cmd.Endpoint, bytes.NewReader(jsonContent))
	if err != nil {
		return ServerResponse{}, err
	}

	res, err := client.Do(request)
	if err != nil && errors.Is(err, syscall.ECONNREFUSED) {
		log.Fatalf("  Server is not active, please make sure it's running")
	} else if err != nil {
		return ServerResponse{}, err
	}

	defer res.Body.Close()
	content, err := io.ReadAll(res.Body)
	if err != nil {
		return ServerResponse{}, err
	}

	return ServerResponse{StatusCode: res.StatusCode, Body: string(content)}, nil
}

func readSpecFile() ([]byte, error) {
	content, err := os.ReadFile(FILE_SPEC_JSON)
	if err != nil {
		return []byte{}, nil
	}
	return content, nil
}
