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
