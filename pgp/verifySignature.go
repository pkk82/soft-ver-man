package pgp

import (
	"bytes"
	"context"
	"errors"
	"github.com/pkk82/soft-ver-man/console"
	"golang.org/x/crypto/openpgp"
	"io"
	"os"
)

func VerifySignature(filePath, signatureFilePath string, publicKeyFilePaths []string) {
	errs := make(chan error)
	resultCh := make(chan bool)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for i := 0; i < len(publicKeyFilePaths); i++ {
			err := <-errs
			if err == nil {
				resultCh <- true
				return
			}
		}
		resultCh <- false
	}()

	go func() {
		for _, publicKeyFilePath := range publicKeyFilePaths {
			go func(pkFilePath string) {
				select {
				case <-ctx.Done():
					return
				default:
				}
				err := verifySignature(filePath, signatureFilePath, pkFilePath)
				if err == nil {
					errs <- nil
					cancel()
				} else {
					errs <- err
				}
			}(publicKeyFilePath)
		}
	}()

	result := <-resultCh

	if result {
		console.Info(filePath + " is correct file")
	} else {
		console.Fatal(errors.New(filePath + " is corrupted file"))
	}

}

func verifySignature(filePath, signatureFilePath, publicKeyFilePath string) error {

	publicKeyContent, err := os.ReadFile(publicKeyFilePath)
	if err != nil {
		return err
	}

	entityList, err := openpgp.ReadArmoredKeyRing(bytes.NewReader(publicKeyContent))
	if err != nil {
		return err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	signatureFile, err := os.Open(signatureFilePath)
	if err != nil {
		return err
	}

	_, err = openpgp.CheckDetachedSignature(entityList, io.Reader(file), io.Reader(signatureFile))
	return err
}
