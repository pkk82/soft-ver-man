package node

import (
	"bufio"
	"fmt"
	"github.com/pkk82/soft-ver-man/util/console"
	"github.com/pkk82/soft-ver-man/util/verification"
	"os"
	"path/filepath"
	"regexp"
)

func verifySha(filePath, signatureFilePath string) {

	expectedHash := readHashes(signatureFilePath)[filepath.Base(filePath)]
	err := verification.VerifySha256(filePath, expectedHash)

	if err == nil {
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
