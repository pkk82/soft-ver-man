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

package shell

import (
	"github.com/pkk82/soft-ver-man/util/test"
	"testing"
)

const bashrc = ".bashrc"
const zshrc = ".zshrc"

func TestInitBash(t *testing.T) {

	tests := []struct {
		name            string
		files           map[string][]string
		expectedError   string
		expectedContent map[string][]string
	}{
		{
			name:          "no rc files",
			expectedError: "at least of of the following files must exist: .bashrc, .zshrc",
		},
		{
			name:            "empty .bashrc",
			files:           map[string][]string{bashrc: {}},
			expectedContent: map[string][]string{bashrc: {"### soft-ver-man", `[[ -s "$HOME/.soft-ver-man/.svmmainrc" ]] && source "$HOME/.soft-ver-man/.svmmainrc"`, ""}},
		}, {
			name:            "empty .zshrc",
			files:           map[string][]string{zshrc: {}},
			expectedContent: map[string][]string{zshrc: {"### soft-ver-man", `[[ -s "$HOME/.soft-ver-man/.svmmainrc" ]] && source "$HOME/.soft-ver-man/.svmmainrc"`, ""}},
		}, {
			name:  "empty .bashrc and .zshrc",
			files: map[string][]string{zshrc: {}, bashrc: {}},
			expectedContent: map[string][]string{
				zshrc:  {"### soft-ver-man", `[[ -s "$HOME/.soft-ver-man/.svmmainrc" ]] && source "$HOME/.soft-ver-man/.svmmainrc"`, ""},
				bashrc: {"### soft-ver-man", `[[ -s "$HOME/.soft-ver-man/.svmmainrc" ]] && source "$HOME/.soft-ver-man/.svmmainrc"`, ""}},
		},
		{
			name:            "existing bashrc",
			files:           map[string][]string{bashrc: {"### soft-ver-man", ""}},
			expectedContent: map[string][]string{bashrc: {"### soft-ver-man", "", "### soft-ver-man", `[[ -s "$HOME/.soft-ver-man/.svmmainrc" ]] && source "$HOME/.soft-ver-man/.svmmainrc"`, ""}},
		},
		{
			name:            "existing bashrc without new line",
			files:           map[string][]string{bashrc: {"### soft-ver-man"}},
			expectedContent: map[string][]string{bashrc: {"### soft-ver-man", "", "### soft-ver-man", `[[ -s "$HOME/.soft-ver-man/.svmmainrc" ]] && source "$HOME/.soft-ver-man/.svmmainrc"`, ""}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := test.CreateTestDir(t)
			if tt.files != nil {
				for file, content := range tt.files {
					test.CreateFile(dir, file, content, t)
				}
			}
			err := initShell(test.TestDirs{Home: dir})
			if tt.expectedError != "" {
				if err.Error() != tt.expectedError {
					t.Errorf("Expected error: %s, got: %s", tt.expectedError, err.Error())
				}
			} else {
				for file, content := range tt.expectedContent {
					test.AssertFileContent(dir, file, content, t)
				}
			}
		})
	}
}
