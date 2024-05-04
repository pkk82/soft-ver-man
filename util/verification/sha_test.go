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

package verification

import (
	"github.com/pkk82/soft-ver-man/util/test"
	"path/filepath"
	"testing"
)

func TestVerifySha256(t *testing.T) {
	type args struct {
		createFile   bool
		fileName     string
		fileContent  []string
		expectedHash string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "no file",
			args: args{
				createFile:   false,
				fileName:     "myFile.txt",
				expectedHash: "hash",
			},
			wantErr: true,
		},
		{
			name: "empty file - correct hash",
			args: args{
				createFile:   true,
				fileName:     "myFile.txt",
				fileContent:  []string{},
				expectedHash: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			},
			wantErr: false,
		},

		{
			name: "empty file - incorrect hash",
			args: args{
				createFile:   true,
				fileName:     "myFile.txt",
				fileContent:  []string{},
				expectedHash: "73cb3858a687a8494ca3323053016282f3dad39d42cf62ca4e79dda2aac7d9ac",
			},
			wantErr: true,
		},
		{
			name: "correct hash",
			args: args{
				createFile:   true,
				fileName:     "myFile.txt",
				fileContent:  []string{"line1", "line2"},
				expectedHash: "683376e290829b482c2655745caffa7a1dccfa10afaa62dac2b42dd6c68d0f83",
			},
			wantErr: false,
		},

		{
			name: "incorrect hash",
			args: args{
				createFile:   true,
				fileName:     "myFile.txt",
				fileContent:  []string{"line1", "line2", ""},
				expectedHash: "683376e290829b482c2655745caffa7a1dccfa10afaa62dac2b42dd6c68d0f83",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDir := test.CreateTestDir(t)
			if tt.args.createFile {
				test.CreateFile(testDir, tt.args.fileName, tt.args.fileContent, t)
			}

			if err := VerifySha256(filepath.Join(testDir, tt.args.fileName), tt.args.expectedHash); (err != nil) != tt.wantErr {
				t.Errorf("VerifySha256() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
