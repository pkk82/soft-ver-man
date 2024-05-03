package shell

import (
	"github.com/pkk82/soft-ver-man/config"
	"github.com/pkk82/soft-ver-man/domain"
	"github.com/pkk82/soft-ver-man/util/test"
	"path"
	"testing"
)

type expectedContentProvider func(dir string) []string

func Test_initVariablesInSvmRc(t *testing.T) {
	type args struct {
		history domain.PackageHistory
	}
	tests := []struct {
		name            string
		args            args
		expectedContent expectedContentProvider
	}{
		{
			name: "node installation",
			args: args{
				history: domain.PackageHistory{
					Name: "node",
					Items: []domain.PackageHistoryItem{
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
				history: domain.PackageHistory{
					Name: "go",
					Items: []domain.PackageHistoryItem{

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
			err := initVariables(test.TestDirs{Home: dir}, tt.args.history, "/bin", domain.VersionGranularityMajor)
			if err != nil {
				t.Errorf("Failed to init variables: %s", err)
			}
			test.AssertFileContent(path.Join(dir, config.HomeConfigDir), config.RcFile, tt.expectedContent(dir), t)
		})
	}
}
func Test_initVariablesInSpecificRc(t *testing.T) {
	type args struct {
		history           domain.PackageHistory
		granularity       domain.VersionGranularity
		executableDirName string
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
				history: domain.PackageHistory{
					Name: "node",
					Items: []domain.PackageHistoryItem{
						{
							Version:     "v20.1.3",
							Path:        "/home/user/pf/node/node-v20.1.3-linux-x64",
							Main:        true,
							InstalledOn: 1689017267000,
						},
					},
				},
				granularity:       domain.VersionGranularityMajor,
				executableDirName: "bin",
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
				history: domain.PackageHistory{
					Name: "node",
					Items: []domain.PackageHistoryItem{
						{
							Version:     "v20.1.3",
							Path:        "/home/user/pf/node/node-v20.1.3-linux-x64",
							Main:        false,
							InstalledOn: 1689017267000,
						},
					},
				},
				granularity:       domain.VersionGranularityMajor,
				executableDirName: "bin",
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
				history: domain.PackageHistory{
					Name: "node",
					Items: []domain.PackageHistoryItem{

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
				granularity:       domain.VersionGranularityMajor,
				executableDirName: "bin",
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
				history: domain.PackageHistory{
					Name: "node",
					Items: []domain.PackageHistoryItem{

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
				granularity:       domain.VersionGranularityMajor,
				executableDirName: "bin",
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
			name: "latest main installation - major granularity",
			args: args{
				history: domain.PackageHistory{
					Name: "node",
					Items: []domain.PackageHistoryItem{
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
				granularity:       domain.VersionGranularityMajor,
				executableDirName: "bin",
			},
			expectedSpecificFileName: ".noderc",
			expectedSpecificContent: []string{
				"export SVM_SOFT_NODE_DIR=\"$SVM_SOFT_DIR/node\"",
				"export NODE_19_HOME=\"$SVM_SOFT_NODE_DIR/node-v19.1.3-linux-x64\"",
				"export NODE_20_HOME=\"$SVM_SOFT_NODE_DIR/node-v20.1.3-linux-x64\"",
				"export NODE_HOME=\"$SVM_SOFT_NODE_DIR/node-v19.1.3-linux-x64\"",
				"export PATH=\"$NODE_HOME/bin:$PATH\"",
			},
		}, {
			name: "latest main installation - minor granularity",
			args: args{
				history: domain.PackageHistory{
					Name: "kotlin",
					Items: []domain.PackageHistoryItem{
						{
							Version:     "v1.8.22",
							Path:        "/home/user/pf/kotlin/kotlin-compiler-1.8.22",
							Main:        true,
							InstalledOn: 1689017268000,
						}, {
							Version:     "v1.9.0",
							Path:        "/home/user/pf/kotlin/kotlin-compiler-1.9.0",
							Main:        true,
							InstalledOn: 1689017267000,
						},
					},
				},
				granularity:       domain.VersionGranularityMinor,
				executableDirName: "bin",
			},
			expectedSpecificFileName: ".kotlinrc",
			expectedSpecificContent: []string{
				"export SVM_SOFT_KOTLIN_DIR=\"$SVM_SOFT_DIR/kotlin\"",
				"export KOTLIN_1_8_HOME=\"$SVM_SOFT_KOTLIN_DIR/kotlin-compiler-1.8.22\"",
				"export KOTLIN_1_9_HOME=\"$SVM_SOFT_KOTLIN_DIR/kotlin-compiler-1.9.0\"",
				"export KOTLIN_HOME=\"$SVM_SOFT_KOTLIN_DIR/kotlin-compiler-1.8.22\"",
				"export PATH=\"$KOTLIN_HOME/bin:$PATH\"",
			},
		},
		{
			name: "soft-ver-man",
			args: args{
				history: domain.PackageHistory{
					Name: "soft-ver-man",
					Items: []domain.PackageHistoryItem{
						{
							Version:     "v0.5.0",
							Path:        "/home/user/pf/soft-ver-man/soft-ver-man-0.5.0",
							Main:        true,
							InstalledOn: 1689017268000,
						},
					},
				},
				executableDirName: "",
				granularity:       domain.VersionGranularityMinor,
			},
			expectedSpecificFileName: ".svmrc",
			expectedSpecificContent: []string{
				"export SVM_SOFT_SVM_DIR=\"$SVM_SOFT_DIR/soft-ver-man\"",
				"export SVM_0_5_HOME=\"$SVM_SOFT_SVM_DIR/soft-ver-man-0.5.0\"",
				"export SVM_HOME=\"$SVM_SOFT_SVM_DIR/soft-ver-man-0.5.0\"",
				"export PATH=\"$SVM_HOME:$PATH\"",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := test.CreateTestDir(t)
			err := initVariables(test.TestDirs{Home: dir}, tt.args.history, tt.args.executableDirName, tt.args.granularity)
			if err != nil {
				t.Errorf("Failed to init variables: %s", err)
			}
			test.AssertFileContent(path.Join(dir, config.HomeConfigDir), tt.expectedSpecificFileName, tt.expectedSpecificContent, t)
		})
	}
}
