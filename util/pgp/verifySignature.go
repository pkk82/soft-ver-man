/*
 * Copyright © 2024 Piotr Kozak <piotrkrzysztofkozak@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 */

package pgp

import (
	"bytes"
	"context"
	"errors"
	"github.com/pkk82/soft-ver-man/util/console"
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
