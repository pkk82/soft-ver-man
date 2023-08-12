/*
Copyright © 2023 Piotr Kozak <piotrkrzysztofkozak@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package domain_test

import (
	"fmt"
	ver "github.com/pkk82/soft-ver-man/domain"
	"testing"
)

func TestFindVersionForNode(t *testing.T) {
	versions := []string{
		"v19.13.19",
		"v21.0.0",
		"v20.0.0",
		"v20.0.1",
		"v20.0.2",
		"v20.0.3",
		"v20.0.4",
		"v20.0.5",
		"v20.0.6",
		"v20.0.7",
		"v20.0.8",
		"v20.0.9",
		"v20.0.10",
		"v20.1.0",
		"v20.2.0",
		"v20.3.0",
		"v20.3.1",
		"v20.4.0",
		"v20.5.0",
		"v20.6.0",
		"v20.7.0",
		"v20.8.0",
		"v20.9.0",
		"v20.10.0",
		"v20.10.1",
	}
	terms := map[string]string{
		"":        "v21.0.0",
		"v":       "v21.0.0",
		"v20":     "v20.10.1",
		"20":      "v20.10.1",
		"v20.":    "v20.10.1",
		"20.":     "v20.10.1",
		"v20.3":   "v20.3.1",
		"20.3":    "v20.3.1",
		"v20.3.1": "v20.3.1",
		"20.3.1":  "v20.3.1",
		"v20.2":   "v20.2.0",
		"20.2":    "v20.2.0",
		"v20.1":   "v20.1.0",
		"20.1":    "v20.1.0",
		"v20.0":   "v20.0.10",
		"20.0":    "v20.0.10",
		"v20.0.1": "v20.0.1",
		"20.0.1":  "v20.0.1",
	}
	for version, expected := range terms {
		actual, _, err := ver.FindVersion(version, versions)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		if actual.Value != expected {
			t.Errorf("Expected: %v,. but got: %v", expected, actual)
		}
	}
}
func TestFindVersionForJava(t *testing.T) {
	versions := []string{
		"20.0.2",
		"19.0.2",
		"18.0.2.1",
	}
	terms := map[string]string{
		"":         "20.0.2",
		"v":        "20.0.2",
		"v20":      "20.0.2",
		"20":       "20.0.2",
		"18":       "18.0.2.1",
		"18.":      "18.0.2.1",
		"18.0.2.1": "18.0.2.1",
	}
	for version, expected := range terms {
		actual, _, err := ver.FindVersion(version, versions)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		if actual.Value != expected {
			t.Errorf("Expected: %v,. but got: %v", expected, actual)
		}
	}
}

func TestFindNoVersion(t *testing.T) {
	versions := []string{
		"v20.0.9",
	}
	terms := []string{
		"v19",
		"v21",
		"v20.1",
		"v20.2",
		"v20.0.1",
		"v20.0.8",
		"v20.0.10",
	}
	for _, version := range terms {
		actual, _, err := ver.FindVersion(version, versions)
		if err == nil {
			t.Errorf("Expected error, but got response: %v", actual)
		}
		expectedError := fmt.Sprintf("No version found for: %v", version)
		if err.Error() != expectedError {
			t.Errorf("Expected: %v,. but got: %v", expectedError, err.Error())
		}
	}
}

func TestFindIndex(t *testing.T) {
	versions := []string{
		"v20.0.9",
		"v20.1.9",
		"v20.2.9",
	}
	terms := []string{
		"v20.0",
		"v20.1",
		"v20.2",
	}
	for index, version := range terms {
		_, actualIndex, err := ver.FindVersion(version, versions)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		if index != actualIndex {
			t.Errorf("Expected: %v, but got: %v", index, actualIndex)
		}
	}
}
