package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"syscall"

	"log/slog"
	"net/http"
)

/*
 So we want to build a tool that helps for
 running api-calls, specifically our custom servers
 For now this is the scope:
 	1. We read a json file that shows the schema
	2. The file should also contain the data schema
	3. All those are for normal post requests

	For Get requests,
	1. Provide the simple endpoint
	2. If we need params we simply include flags for that
*/

type Command struct {
	Endpoint       string
	ContentType    string
	Method         string
	FollowRedirect bool
}

var (
	contentTypeJson           = "application/json"
	contentTypeTextHTML       = "text/html"
	contentTypeTextCSS        = "text/css"
	contentTypeTextJAVASCRIPT = "text/css"
)

// This automatically uses bat to read
// show content on the terminal to avoid
// dumping huge respones on the terminal
//
// By default, this uses bat, looking into
// /opt/homebrew/bin/bat
func GetArgs() Command {
	method := flag.Bool("post", false, "Make POST request to endpoint")
	endpoint := flag.String("ep", "http://localhost:3000", "Endpoint to make request to, default is localhost:3000")
	contentType := flag.String("ct", contentTypeJson, "Content Type format to send data")
	followRedirect := flag.Bool("redirect", false, "Follow subsequent redirects")

	flag.Parse()

	cmd := Command{}
	if *method {
		cmd.Method = http.MethodPost
	} else {
		cmd.Method = http.MethodGet
	}

	cmd.Endpoint = *endpoint
	cmd.ContentType = *contentType
	cmd.FollowRedirect = *followRedirect
	return cmd
}

type ServerResponse struct {
	StatusCode int
	Body       string
	// PagerType   string
	ContentType string
}

var inProgress = errors.New("Server still under construction")

func makeRequest(cmd Command) (ServerResponse, error) {
	// if cmd.Method != http.MethodGet {
	// 	slog.Error("Method still under construction", "for", cmd.Method)
	// 	return ServerResponse{}, inProgress
	// }

	switch cmd.Method {
	case http.MethodPost:
		var res ServerResponse
		var err error
		if cmd.FollowRedirect {
			fmt.Println("warning! following redirects from server")
			res, err = cmd.MakePostRequest()
		} else {
			res, err = cmd.MakePostRequestNoRedirect()
		}
		if err != nil {
			return ServerResponse{}, err
		}
		return res, nil

	}

	res, err := http.Get(cmd.Endpoint)
	if err != nil && errors.Is(err, syscall.ECONNREFUSED) {
		log.Fatalf("  Server is not active, please make sure it's running")
	} else if err != nil {
		return ServerResponse{}, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return ServerResponse{}, err
	}

	serverRes := ServerResponse{}
	serverRes.StatusCode = res.StatusCode
	serverRes.Body = string(body)
	// serverRes.PagerType = contentTypeTextHTML

	// contentType := res.Header.Get("content-type")
	// switch true {
	// case strings.Contains(contentType, contentTypeJson):
	// 	serverRes.PagerType = contentTypeJson
	// }
	serverRes.ContentType = res.Header.Get("content-type")

	return serverRes, nil
}

func main() {
	command := GetArgs()
	slog.Info("Making request, ", "method", command.Method, "endpoint", command.Endpoint)
	fmt.Println()

	res, err := makeRequest(command)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(" status-code: %d\n", res.StatusCode)
	if err := Pager(res); err != nil {
		log.Fatal(err)
	}

}
