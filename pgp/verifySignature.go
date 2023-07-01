package pgp

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/pkk82/soft-ver-man/console"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
	"os"
)

func VerifySignature(filePath, signatureFilePath string, publicKeyFilePaths []string) {
	errs := make(chan error)
	resultCh := make(chan bool)

	go func() {
		for i := 0; i < len(publicKeyFilePaths); i++ {
			err := <-errs
			if err == nil {
				close(errs)
				resultCh <- true
				return
			}
		}
		resultCh <- false
	}()

	go func() {
		for _, publicKeyFilePath := range publicKeyFilePaths {
			go verifySignatureAsync(filePath, signatureFilePath, publicKeyFilePath, errs)
		}
	}()

	result := <-resultCh

	if result {
		console.Info(filePath + " is correct file")
	} else {
		console.Fatal(errors.New(filePath + " is corrupted file"))
	}

}

func verifySignatureAsync(filePath, signatureFilePath, publicKeyFilePath string, ch chan<- error) {
	err := verifySignature(filePath, signatureFilePath, publicKeyFilePath)
	ch <- err
}

func verifySignature(filePath, signatureFilePath, publicKeyFilePath string) error {

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	signatureFile, err := os.Open(signatureFilePath)
	if err != nil {
		return err
	}

	defer func() {
		if err := signatureFile.Close(); err != nil {
			console.Error(err)
		}
	}()

	pack, err := packet.Read(signatureFile)
	if err != nil {
		return err
	}

	signature, ok := pack.(*packet.Signature)
	if !ok {
		return errors.New(signatureFilePath + " is not a valid signature file")
	}

	publicKeyContent, err := os.ReadFile(publicKeyFilePath)
	if err != nil {
		return err
	}

	block, err := armor.Decode(bytes.NewReader(publicKeyContent))
	if err != nil {
		return fmt.Errorf("error decoding public key: %s", err)
	}
	if block.Type != "PGP PUBLIC KEY BLOCK" {
		return errors.New(publicKeyFilePath + " not an armored public key")
	}

	pack, err = packet.Read(block.Body)
	if err != nil {
		return fmt.Errorf("error reading public key: %s", err)
	}

	publicKey, ok := pack.(*packet.PublicKey)
	if !ok {
		return errors.New("invalid public key")
	}

	println(publicKey)

	hash := signature.Hash.New()

	_, err = hash.Write(fileContent)
	if err != nil {
		return err
	}

	err = publicKey.VerifySignature(hash, signature)
	if err != nil {
		return err
	}

	return nil
}
