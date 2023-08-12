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
			expectedContent: map[string][]string{bashrc: {"### soft-ver-man", `[[ -s "$HOME/.soft-ver-man/.svmrc" ]] && source "$HOME/.soft-ver-man/.svmrc"`, ""}},
		}, {
			name:            "empty .zshrc",
			files:           map[string][]string{zshrc: {}},
			expectedContent: map[string][]string{zshrc: {"### soft-ver-man", `[[ -s "$HOME/.soft-ver-man/.svmrc" ]] && source "$HOME/.soft-ver-man/.svmrc"`, ""}},
		}, {
			name:  "empty .bashrc and .zshrc",
			files: map[string][]string{zshrc: {}, bashrc: {}},
			expectedContent: map[string][]string{
				zshrc:  {"### soft-ver-man", `[[ -s "$HOME/.soft-ver-man/.svmrc" ]] && source "$HOME/.soft-ver-man/.svmrc"`, ""},
				bashrc: {"### soft-ver-man", `[[ -s "$HOME/.soft-ver-man/.svmrc" ]] && source "$HOME/.soft-ver-man/.svmrc"`, ""}},
		},
		{
			name:            "existing bashrc",
			files:           map[string][]string{bashrc: {"### soft-ver-man", ""}},
			expectedContent: map[string][]string{bashrc: {"### soft-ver-man", "", "### soft-ver-man", `[[ -s "$HOME/.soft-ver-man/.svmrc" ]] && source "$HOME/.soft-ver-man/.svmrc"`, ""}},
		},
		{
			name:            "existing bashrc without new line",
			files:           map[string][]string{bashrc: {"### soft-ver-man"}},
			expectedContent: map[string][]string{bashrc: {"### soft-ver-man", "", "### soft-ver-man", `[[ -s "$HOME/.soft-ver-man/.svmrc" ]] && source "$HOME/.soft-ver-man/.svmrc"`, ""}},
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
