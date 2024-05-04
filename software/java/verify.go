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

package java

import (
	"encoding/json"
	"github.com/pkk82/soft-ver-man/util/console"
	"io"
	"net/http"
)

type ExtendedPackage struct {
	Id     string `json:"package_uuid"`
	Sha256 string `json:"sha256_hash"`
}

func getExtendedPackage(packageId string) (ExtendedPackage, error) {
	url := PackagesAPIURL + "/" + packageId
	resp, err := http.Get(url)
	if err != nil {
		return ExtendedPackage{}, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			console.Fatal(err)
		}
	}(resp.Body)
	var extendedPackage ExtendedPackage
	err = json.NewDecoder(resp.Body).Decode(&extendedPackage)
	if err != nil {
		return ExtendedPackage{}, err
	}
	return extendedPackage, nil
}
