package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	"github.com/go-resty/resty"
	"io/ioutil"
	"os"
	re "regexp"
	"strings"
)

const VERSION = "0.0.1"

const DOCOPT = `
etcdnew.

Replaces references on discovery.etcd.io on files

Usage: etcdnew [--url URL] FILES...

Options:
  --url URL  Instead of generating a new use, use URL for one file

`

var discoveryMatcher = re.MustCompile("https://discovery.etcd.io/[a-z0-9]{32}")

func main() {
	args, _ := docopt.Parse(DOCOPT, nil, true, VERSION, true, true)

	fileList := args["FILES"].([]string)
	urlToUse, hasUrlToUse := args["--url"].(string)

	if hasUrlToUse && urlToUse != "" && 1 != len(fileList) {
		panic(fmt.Errorf("--url only accepts one FILE"))
	}

	for _, file := range fileList {
		err := processFile(urlToUse, file)

		if nil != err {
			panic(err)
		}
	}
}

func processFile(urlToUse string, file string) error {
	if "" == urlToUse {
		resp, err := resty.
		R().
			Get("https://discovery.etcd.io/new")

		if nil != err {
			return err
		}

		if resp.StatusCode() != 200 {
			err = fmt.Errorf("Invalid resultcode returned by discovery.etcd.io/new: %d", resp.StatusCode())

			return err
		}

		urlToUse = string(resp.Body())
	}

	if !discoveryMatcher.MatchString(urlToUse) {
		err := fmt.Errorf("Invalid pattern for urlToUse: %s", urlToUse)

		return err
	}

	lines, err := readAllLines(file)

	if nil != err {
		return err
	}

	noOfModifiedLines := 0

	for i, line := range lines {
		originalLine := line

		modifiedLine := discoveryMatcher.ReplaceAllLiteralString(line, urlToUse)

		if originalLine != modifiedLine {
			lines[i] = modifiedLine

			noOfModifiedLines++
		}
	}

	if 0 != noOfModifiedLines {
		err = writeAllLines(file, lines)

		if nil != err {
			return err
		}
	}

	return nil
}

func readAllLines(file string) (lines []string, err error) {
	f, err := os.OpenFile(file, os.O_RDONLY, 0)

	if nil != err {
		return
	}

	defer f.Close()

	byteContent, err := ioutil.ReadAll(f)

	if nil != err {
		return
	}

	strContent := string(byteContent)

	lines = strings.Split(strContent, "\n")

	for i, x := range lines {
		lines[i] = strings.TrimRight(x, "\n")
	}

	return
}

func writeAllLines(file string, lines []string) (err error) {
	stringContent := strings.Join(lines, "\n")

	f, err := os.OpenFile(file, os.O_TRUNC|os.O_WRONLY, 0666)

	if nil != err {
		return
	}

	defer f.Close()

	f.WriteString(stringContent)

	return
}
