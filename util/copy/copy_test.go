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

package copy

import (
	"github.com/pkk82/soft-ver-man/domain"
	"github.com/pkk82/soft-ver-man/util/test"
	"path"
	"path/filepath"
	"reflect"
	"testing"
)

func TestCopy(t *testing.T) {
	type args struct {
		fetchedPackage domain.FetchedPackage
	}
	tests := []struct {
		name         string
		args         args
		wantPath     string
		wantFileName string
		wantContent  string
	}{
		{
			name: "copy file",
			args: args{
				fetchedPackage: domain.FetchedPackage{
					Version:  domain.Version{Value: "v0.5.0"},
					FilePath: filepath.Join("testdata", "soft-ver-man-v0.5.0"),
					Type:     domain.RAW,
				},
			},
			wantPath:     "soft-ver-man-v0.5.0-dir",
			wantFileName: "svm",
			wantContent:  "soft-ver-man-v0.5.0-content",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDir := test.CreateTestDir(t)
			got, err := Copy(tt.args.fetchedPackage, path.Join(testDir, "soft-ver-man-v0.5.0-dir"), "svm")
			if err != nil {
				t.Errorf("Extract() error = %v", err)
				return
			}
			expected := CopiedPackage{
				Version:    tt.args.fetchedPackage.Version,
				PathToFile: path.Join(testDir, tt.wantPath),
				FileName:   tt.wantFileName,
			}
			if !reflect.DeepEqual(got, expected) {
				t.Errorf("Copy() got = %v, want %v", got, expected)
			}
			test.AssertFileContent(expected.PathToFile, expected.FileName, []string{tt.wantContent}, t)
			test.AssertFileMode(expected.PathToFile, expected.FileName, 0755, t)
		})
	}
}
