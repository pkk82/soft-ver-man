package shell

import (
	"github.com/pkk82/soft-ver-man/config"
	"github.com/pkk82/soft-ver-man/history"
	"github.com/pkk82/soft-ver-man/test"
	"path"
	"testing"
)

type expectedContentProvider func(dir string) []string

func Test_initVariablesInSvmRc(t *testing.T) {
	type args struct {
		history history.PackageHistory
	}
	tests := []struct {
		name            string
		args            args
		expectedContent expectedContentProvider
	}{
		{
			name: "node installation",
			args: args{
				history: history.PackageHistory{
					Name: "node",
					Items: []history.PackageHistoryItem{
						{
							Version:     "v20.1.3",
							Path:        "/home/user/pf/node/node-v20.1.3-linux-x64",
							Main:        true,
							InstalledOn: 1689017267000,
						},
					},
				},
			},
			expectedContent: func(dir string) []string {
				return []string{"export SVM_SOFT_DIR=\"" + path.Join(dir, "pf\""), "[[ -s \"$HOME/.soft-ver-man/.noderc\" ]] && source \"$HOME/.soft-ver-man/.noderc\"\n"}
			},
		}, {
			name: "go installation",
			args: args{
				history: history.PackageHistory{
					Name: "go",
					Items: []history.PackageHistoryItem{

						{
							Version:     "1.20.1",
							Path:        "/home/user/pf/go/go-v1.20.1-linux-x64",
							Main:        true,
							InstalledOn: 1689017267000,
						},
					},
				},
			},
			expectedContent: func(dir string) []string {
				return []string{"export SVM_SOFT_DIR=\"" + path.Join(dir, "pf\""), "[[ -s \"$HOME/.soft-ver-man/.gorc\" ]] && source \"$HOME/.soft-ver-man/.gorc\"\n"}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := test.CreateTestDir(t)
			err := initVariables(test.TestDirs{Home: dir}, tt.args.history)
			if err != nil {
				t.Errorf("Failed to init variables: %s", err)
			}
			test.AssertFileContent(path.Join(dir, config.HomeConfigDir), config.RcFile, tt.expectedContent(dir), t)
		})
	}
}
func Test_initVariablesInNodeRc(t *testing.T) {
	type args struct {
		history history.PackageHistory
	}
	tests := []struct {
		name                     string
		args                     args
		expectedSpecificFileName string
		expectedSpecificContent  []string
	}{
		{
			name: "first installation as main",
			args: args{
				history: history.PackageHistory{
					Name: "node",
					Items: []history.PackageHistoryItem{
						{
							Version:     "v20.1.3",
							Path:        "/home/user/pf/node/node-v20.1.3-linux-x64",
							Main:        true,
							InstalledOn: 1689017267000,
						},
					},
				},
			},
			expectedSpecificFileName: ".noderc",
			expectedSpecificContent: []string{
				"export SVM_SOFT_NODE_DIR=\"$SVM_SOFT_DIR/node\"",
				"export NODE_20_HOME=\"$SVM_SOFT_NODE_DIR/node-v20.1.3-linux-x64\"",
				"export NODE_HOME=\"$SVM_SOFT_NODE_DIR/node-v20.1.3-linux-x64\"",
				"export PATH=\"$NODE_HOME/bin:$PATH\"",
			},
		}, {
			name: "first installation - no main",
			args: args{
				history: history.PackageHistory{
					Name: "node",
					Items: []history.PackageHistoryItem{
						{
							Version:     "v20.1.3",
							Path:        "/home/user/pf/node/node-v20.1.3-linux-x64",
							Main:        false,
							InstalledOn: 1689017267000,
						},
					},
				},
			},
			expectedSpecificFileName: ".noderc",
			expectedSpecificContent: []string{
				"export SVM_SOFT_NODE_DIR=\"$SVM_SOFT_DIR/node\"",
				"export NODE_20_HOME=\"$SVM_SOFT_NODE_DIR/node-v20.1.3-linux-x64\"",
				"export NODE_HOME=\"$SVM_SOFT_NODE_DIR/node-v20.1.3-linux-x64\"",
				"export PATH=\"$NODE_HOME/bin:$PATH\"",
			},
		}, {
			name: "another main installation",
			args: args{
				history: history.PackageHistory{
					Name: "node",
					Items: []history.PackageHistoryItem{

						{
							Version:     "19.1.3",
							Path:        "/home/user/pf/node/node-v19.1.3-linux-x64",
							Main:        true,
							InstalledOn: 1689017267000,
						},
						{
							Version:     "19.1.4",
							Path:        "/home/user/pf/node/node-v19.1.4-linux-x64",
							Main:        false,
							InstalledOn: 168901728000,
						},

						{
							Version:     "v20.1.3",
							Path:        "/home/user/pf/node/node-v20.1.3-linux-x64",
							Main:        true,
							InstalledOn: 1689017269000,
						},
					},
				},
			},
			expectedSpecificFileName: ".noderc",
			expectedSpecificContent: []string{
				"export SVM_SOFT_NODE_DIR=\"$SVM_SOFT_DIR/node\"",
				"export NODE_19_HOME=\"$SVM_SOFT_NODE_DIR/node-v19.1.4-linux-x64\"",
				"export NODE_20_HOME=\"$SVM_SOFT_NODE_DIR/node-v20.1.3-linux-x64\"",
				"export NODE_HOME=\"$SVM_SOFT_NODE_DIR/node-v20.1.3-linux-x64\"",
				"export PATH=\"$NODE_HOME/bin:$PATH\"",
			},
		}, {
			name: "another installation",
			args: args{
				history: history.PackageHistory{
					Name: "node",
					Items: []history.PackageHistoryItem{

						{
							Version:     "19.1.3",
							Path:        "/home/user/pf/node/node-v19.1.3-linux-x64",
							Main:        false,
							InstalledOn: 168901728000,
						},
						{
							Version:     "19.1.4",
							Path:        "/home/user/pf/node/node-v19.1.4-linux-x64",
							Main:        false,
							InstalledOn: 168901727000,
						},

						{
							Version:     "v20.1.3",
							Path:        "/home/user/pf/node/node-v20.1.3-linux-x64",
							Main:        false,
							InstalledOn: 168901726000,
						},
					},
				},
			},
			expectedSpecificFileName: ".noderc",
			expectedSpecificContent: []string{
				"export SVM_SOFT_NODE_DIR=\"$SVM_SOFT_DIR/node\"",
				"export NODE_19_HOME=\"$SVM_SOFT_NODE_DIR/node-v19.1.4-linux-x64\"",
				"export NODE_20_HOME=\"$SVM_SOFT_NODE_DIR/node-v20.1.3-linux-x64\"",
				"export NODE_HOME=\"$SVM_SOFT_NODE_DIR/node-v20.1.3-linux-x64\"",
				"export PATH=\"$NODE_HOME/bin:$PATH\"",
			},
		}, {
			name: "latest main installation",
			args: args{
				history: history.PackageHistory{
					Name: "node",
					Items: []history.PackageHistoryItem{
						{
							Version:     "19.1.3",
							Path:        "/home/user/pf/node/node-v19.1.3-linux-x64",
							Main:        true,
							InstalledOn: 1689017268000,
						}, {
							Version:     "v20.1.3",
							Path:        "/home/user/pf/node/node-v20.1.3-linux-x64",
							Main:        true,
							InstalledOn: 1689017267000,
						},
					},
				},
			},
			expectedSpecificFileName: ".noderc",
			expectedSpecificContent: []string{
				"export SVM_SOFT_NODE_DIR=\"$SVM_SOFT_DIR/node\"",
				"export NODE_19_HOME=\"$SVM_SOFT_NODE_DIR/node-v19.1.3-linux-x64\"",
				"export NODE_20_HOME=\"$SVM_SOFT_NODE_DIR/node-v20.1.3-linux-x64\"",
				"export NODE_HOME=\"$SVM_SOFT_NODE_DIR/node-v19.1.3-linux-x64\"",
				"export PATH=\"$NODE_HOME/bin:$PATH\"",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := test.CreateTestDir(t)
			err := initVariables(test.TestDirs{Home: dir}, tt.args.history)
			if err != nil {
				t.Errorf("Failed to init variables: %s", err)
			}
			test.AssertFileContent(path.Join(dir, config.HomeConfigDir), tt.expectedSpecificFileName, tt.expectedSpecificContent, t)
		})
	}
}
