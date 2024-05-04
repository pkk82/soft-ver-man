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

package download

import (
	"fmt"
	"github.com/pkk82/soft-ver-man/util/console"
	"github.com/schollz/progressbar/v3"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func FetchFile(url, downloadDir, fileName string) string {
	return fetchFile(url, downloadDir, fileName, true)
}

func FetchFileSilently(url, downloadDir, fileName string) string {
	return fetchFile(url, downloadDir, fileName, false)
}

func fetchFile(url, downloadDir, fileName string, useProgressBar bool) string {
	response, err := http.Get(url)
	if err != nil {
		console.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(response.Body)

	err = os.MkdirAll(downloadDir, os.ModePerm)
	if err != nil {
		console.Fatal(err)
	}

	file, err := os.Create(filepath.Join(downloadDir, fileName))
	if err != nil {
		console.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)

	if useProgressBar {
		bar := progressbar.DefaultBytes(response.ContentLength, "Downloading "+fileName)
		_, err = io.Copy(io.MultiWriter(file, bar), response.Body)
	} else {
		_, err = io.Copy(file, response.Body)
	}

	if err != nil {
		console.Fatal(err)
	}

	return filepath.Join(downloadDir, fileName)
}
