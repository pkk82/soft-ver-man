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
package domain_test

import (
	"fmt"
	ver "github.com/pkk82/soft-ver-man/domain"
	"testing"
)

func TestCorrectVersions(t *testing.T) {
	versions := []string{"v10.11.12", "10.11.12", "v10.11", "10.11", "v10", "10",
		"v10.11.", "10.11.", "v10.", "10.", "1.0", "1.1", "0.1", "0.2",
	}
	for _, v := range versions {
		err := ver.ValidateVersion(v)
		if err != nil {
			t.Errorf("Validate should pass, but got exception: %v", err)
		}
	}
}

func TestIncorrectVersions(t *testing.T) {
	versions := []string{"a.b.c", "1.2.a", "x1.2.3"}
	for _, v := range versions {
		err := ver.ValidateVersion(v)
		expectedErr := fmt.Sprintf("%s is not a valid version", v)
		if err.Error() != expectedErr {
			t.Errorf("Got unexpected error: %v, expected: %v", err, expectedErr)
		}
	}
}
