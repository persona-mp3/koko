package main

import (
	"os"
	"os/exec"
	"strings"
)

/*
 The following pagers are installed via homebrew
 1. jq
 2. bat
 And these are system pagers
 1. less
 2. cat but is not used
*/

var (
	jqPagerPath   = "/opt/homebrew/bin/jq"
	batPagerPath  = "/opt/homebrew/bin/bat"
	lessPagerPath = "/usr/bin/less"
)

// Sends output from res.Body to stdout depending on the
// content-type and pager assigned to it
func Pager(res ServerResponse) error {
	// cmd := exec.Command("/opt/homebrew/bin/bat")

	pathToPager := canonPager(res)

	cmd := exec.Command(pathToPager)
	cmd.Stdout = os.Stdout
	cmd.Stdin = strings.NewReader(res.Body)
	return cmd.Run()
}

// Returns the pager path to use depending on the
// content-type of application.
func canonPager(res ServerResponse) string {
	switch true {
	case strings.Contains(res.ContentType, contentTypeJson):
		return jqPagerPath
	case strings.Contains(res.ContentType, contentTypeTextHTML):
		return batPagerPath

	default:
		return lessPagerPath
	}
}
