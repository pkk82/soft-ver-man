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

package maven_test

import (
	"fmt"
	"github.com/pkk82/soft-ver-man/domain"
	"github.com/pkk82/soft-ver-man/software/maven"
	"testing"
)

func Test_calculateDownloadUrl(t *testing.T) {
	plugin := domain.GetPlugin(maven.Name)

	tests := []struct {
		version      domain.Version
		os           string
		arch         string
		expectedPath string
		expectedType domain.Type
	}{
		{
			version:      domain.Ver("3.9.9", t),
			os:           "linux",
			arch:         "amd64",
			expectedPath: "https://archive.apache.org/dist/maven/maven-3/3.9.9/binaries/apache-maven-3.9.9-bin.zip",
			expectedType: domain.ZIP,
		}, {
			version:      domain.Ver("2.2.1", t),
			os:           "linux",
			arch:         "arm64",
			expectedPath: "https://archive.apache.org/dist/maven/maven-2/2.2.1/binaries/apache-maven-2.2.1-bin.zip",
			expectedType: domain.ZIP,
		},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("Checking download path for version %s, os %s and arch %s", tt.version.Value, tt.os, tt.arch)
		t.Run(name, func(t *testing.T) {
			actualPath, actualType := plugin.CalculateDownloadUrl(tt.version, tt.os, tt.arch)
			if actualPath != tt.expectedPath {
				t.Errorf("calculateDownloadPath().path = %v, want %v", actualPath, tt.expectedPath)
			}
			if actualType != tt.expectedType {
				t.Errorf("calculateDownloadPath().type = %v, want %v", actualType, tt.expectedType)
			}
		})
	}
}
