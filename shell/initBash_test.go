package shell

import (
	"github.com/pkk82/soft-ver-man/test"
	"testing"
)

const bashrc = ".bashrc"

func TestInitBash(t *testing.T) {

	tests := []struct {
		name            string
		bashRcProvider  interface{}
		expectedContent []string
	}{
		{
			name:            "no bashrc",
			bashRcProvider:  func(dir string) {},
			expectedContent: []string{"### soft-ver-man", `[[ -s "$HOME/.soft-ver-man/.svmrc" ]] && source "$HOME/.soft-ver-man/.svmrc"`, ""},
		},
		{
			name:            "empty bashrc",
			bashRcProvider:  func(dir string) { test.CreateEmptyFile(dir, bashrc, t) },
			expectedContent: []string{"### soft-ver-man", `[[ -s "$HOME/.soft-ver-man/.svmrc" ]] && source "$HOME/.soft-ver-man/.svmrc"`, ""},
		},
		{
			name:            "existing bashrc",
			bashRcProvider:  func(dir string) { test.CreateFile(dir, bashrc, []string{"### soft-ver-man", ""}, t) },
			expectedContent: []string{"### soft-ver-man", "", "### soft-ver-man", `[[ -s "$HOME/.soft-ver-man/.svmrc" ]] && source "$HOME/.soft-ver-man/.svmrc"`, ""},
		},
		{
			name:            "existing bashrc without new line",
			bashRcProvider:  func(dir string) { test.CreateFile(dir, bashrc, []string{"### soft-ver-man"}, t) },
			expectedContent: []string{"### soft-ver-man", "", "### soft-ver-man", `[[ -s "$HOME/.soft-ver-man/.svmrc" ]] && source "$HOME/.soft-ver-man/.svmrc"`, ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := test.CreateTestDir(t)
			tt.bashRcProvider.(func(string))(dir)
			err := initBash(test.TestDirs{Home: dir})
			if err != nil {
				t.Errorf("Failed to init bash: %s", err)
			}
			test.AssertFileContent(dir, bashrc, tt.expectedContent, t)
		})
	}
}
