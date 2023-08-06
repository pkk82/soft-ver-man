package verification

import (
	"crypto/sha256"
	"fmt"
	"github.com/pkk82/soft-ver-man/console"
	"io"
	"os"
)

func VerifySha256(filePath, expectedHash string) error {

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			console.Error(err)
		}
	}(file)

	h := sha256.New()
	if _, err := io.Copy(h, io.Reader(file)); err != nil {
		return err
	}

	fileHash := fmt.Sprintf("%x", h.Sum(nil))

	if fileHash == expectedHash {
		return nil
	} else {
		return fmt.Errorf(filePath + " is corrupted file (expected hash: " + expectedHash + ", actual hash: " + fileHash + ")")
	}
}
