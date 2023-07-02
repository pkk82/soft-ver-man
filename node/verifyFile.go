package node

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"github.com/pkk82/soft-ver-man/console"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

func verifySha(filePath, signatureFilePath string) {

	file, err := os.Open(filePath)
	if err != nil {
		console.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			console.Error(err)
		}
	}(file)

	h := sha256.New()
	if _, err := io.Copy(h, io.Reader(file)); err != nil {
		console.Fatal(err)
	}

	fileHash := fmt.Sprintf("%x", h.Sum(nil))
	expectedHash := readHashes(signatureFilePath)[filepath.Base(file.Name())]

	if fileHash == expectedHash {
		console.Info(filePath + " is correct file")
	} else {
		console.Fatal(fmt.Errorf(filePath + " is corrupted file"))
	}
}

func readHashes(signatureFilePath string) map[string]string {
	file, err := os.Open(signatureFilePath)
	if err != nil {
		console.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			console.Error(err)
		}
	}(file)

	regex, err := regexp.Compile("\\s+")
	if err != nil {
		console.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	result := make(map[string]string)
	for scanner.Scan() {
		line := scanner.Text()
		split := regex.Split(line, -1)
		result[split[1]] = split[0]
	}
	return result
}
