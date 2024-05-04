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

func verifySha(filePath, signatureFilePath string) error {

	expectedHash := readHashes(signatureFilePath)[filepath.Base(filePath)]
	err := verification.VerifySha256(filePath, expectedHash)

	if err == nil {
		console.Info(filePath + " is correct file")
	} else {
		console.Error(fmt.Errorf(filePath + " is corrupted file"))
	}
	return err
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
