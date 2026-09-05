/*
 * Copyright © 2026 Piotr Kozak <piotrkrzysztofkozak@gmail.com>
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

package gradle

import (
	"github.com/pkk82/soft-ver-man/domain"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_expectedSha256_fromApiProperties(t *testing.T) {
	asset := domain.Asset{
		Version:         "9.7.1",
		Url:             "https://services.gradle.org/distributions/gradle-9.7.1-bin.zip",
		Type:            domain.ZIP,
		ExtraProperties: map[string]string{Sha256Property: "acd53f1edaf02f1a8ff99879f8a34b302661a057d9b063ae9e35b552f804d20a"},
	}
	actual, err := expectedSha256(asset)
	if err != nil {
		t.Fatalf("expectedSha256() returned error: %v", err)
	}
	expected := "acd53f1edaf02f1a8ff99879f8a34b302661a057d9b063ae9e35b552f804d20a"
	if actual != expected {
		t.Errorf("expectedSha256() = %v, want %v", actual, expected)
	}
}

func Test_expectedSha256_fromChecksumFile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/distributions/gradle-9.7.1-bin.zip.sha256" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		_, _ = w.Write([]byte("acd53f1edaf02f1a8ff99879f8a34b302661a057d9b063ae9e35b552f804d20a\n"))
	}))
	defer server.Close()

	asset := domain.Asset{
		Version: "9.7.1",
		Url:     server.URL + "/distributions/gradle-9.7.1-bin.zip",
		Type:    domain.ZIP,
	}
	actual, err := expectedSha256(asset)
	if err != nil {
		t.Fatalf("expectedSha256() returned error: %v", err)
	}
	expected := "acd53f1edaf02f1a8ff99879f8a34b302661a057d9b063ae9e35b552f804d20a"
	if actual != expected {
		t.Errorf("expectedSha256() = %v, want %v", actual, expected)
	}
}

func Test_expectedSha256_checksumFileMissing(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	asset := domain.Asset{
		Version: "9.7.1",
		Url:     server.URL + "/distributions/gradle-9.7.1-bin.zip",
		Type:    domain.ZIP,
	}
	_, err := expectedSha256(asset)
	if err == nil {
		t.Errorf("expectedSha256() should fail when the checksum file is missing")
	}
}
