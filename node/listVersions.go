/*
Copyright © 2023 Piotr Kozak <piotrkrzysztofkozak@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package node

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"
)

type Version struct {
	Id    string   `json:"version"`
	Date  string   `json:"date"`
	Files []string `json:"files"`
}

func ListVersions(jsonFileUrl string) {
	resp, err := http.Get(jsonFileUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)
	var versions []Version
	err = json.NewDecoder(resp.Body).Decode(&versions)
	if err != nil {
		log.Fatal(err)
	}
	supportedFile := supportedFile()
	for _, version := range versions {
		files := version.Files
		if includes(files, supportedFile) {
			println(version.Id)
		}
	}
}

func includes(c []string, term string) bool {
	for _, elem := range c {
		if elem == term {
			return true
		}
	}
	return false
}

func supportedFile() string {
	os := os()
	var supportedFile string
	switch os {
	case "win":
		supportedFile = fmt.Sprintf("%s-%s-zip", os, arch())
	case "osx":
		supportedFile = fmt.Sprintf("%s-%s-tar", os, arch())
	default:
		supportedFile = fmt.Sprintf("%s-%s", os, arch())
	}
	return supportedFile
}

func os() string {
	goOS := runtime.GOOS
	var os string
	switch goOS {
	case "darwin":
		os = "osx"
	case "windows":
		os = "win"
	default:
		os = goOS
	}
	return os
}

func arch() string {
	goArch := runtime.GOARCH
	var arch string
	switch goArch {
	case "amd64":
		arch = "x64"
	case "386":
		arch = "x86"
	default:
		arch = goArch
	}
	return arch
}
