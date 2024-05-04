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
	"github.com/pkk82/soft-ver-man/config"
	"github.com/pkk82/soft-ver-man/domain"
	"github.com/pkk82/soft-ver-man/util/test"
	"path"
	"testing"
)

type expectedContentProvider func(dir string) []string

func Test_initVariablesInSvmRc(t *testing.T) {
	type args struct {
		installedPackages domain.InstalledPackages
	}
	tests := []struct {
		name            string
		args            args
		expectedContent expectedContentProvider
	}{
		{
			name: "node installation",
			args: args{
				installedPackages: domain.InstalledPackages{
					Plugin: domain.Plugin{Name: "node"},
					Items: []domain.InstalledPackage{
						{
							Version:     domain.Ver("v20.1.3", t),
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
				installedPackages: domain.InstalledPackages{
					Plugin: domain.Plugin{Name: "go", VersionGranularity: domain.VersionGranularityMajor, ExecutableRelativePath: "/bin"},
					Items: []domain.InstalledPackage{

						{
							Version:     domain.Ver("1.20.1", t),
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
			err := initVariables(test.TestDirs{Home: dir}, tt.args.installedPackages)
			if err != nil {
				t.Errorf("Failed to init variables: %s", err)
			}
			test.AssertFileContent(path.Join(dir, config.HomeConfigDir), config.RcFile, tt.expectedContent(dir), t)
		})
	}
}
func Test_initVariablesInSpecificRc(t *testing.T) {
	type args struct {
		installedPackages domain.InstalledPackages
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
				installedPackages: domain.InstalledPackages{
					Plugin: domain.Plugin{Name: "node", VersionGranularity: domain.VersionGranularityMajor, ExecutableRelativePath: "bin"},
					Items: []domain.InstalledPackage{
						{
							Version:     domain.Ver("v20.1.3", t),
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
				installedPackages: domain.InstalledPackages{
					Plugin: domain.Plugin{Name: "node", VersionGranularity: domain.VersionGranularityMajor, ExecutableRelativePath: "bin"},
					Items: []domain.InstalledPackage{
						{
							Version:     domain.Ver("v20.1.3", t),
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
				installedPackages: domain.InstalledPackages{
					Plugin: domain.Plugin{Name: "node", VersionGranularity: domain.VersionGranularityMajor, ExecutableRelativePath: "bin"},
					Items: []domain.InstalledPackage{

						{
							Version:     domain.Ver("19.1.3", t),
							Path:        "/home/user/pf/node/node-v19.1.3-linux-x64",
							Main:        true,
							InstalledOn: 1689017267000,
						},
						{
							Version:     domain.Ver("19.1.4", t),
							Path:        "/home/user/pf/node/node-v19.1.4-linux-x64",
							Main:        false,
							InstalledOn: 168901728000,
						},

						{
							Version:     domain.Ver("v20.1.3", t),
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
				installedPackages: domain.InstalledPackages{
					Plugin: domain.Plugin{Name: "node", VersionGranularity: domain.VersionGranularityMajor, ExecutableRelativePath: "bin"},
					Items: []domain.InstalledPackage{

						{
							Version:     domain.Ver("19.1.3", t),
							Path:        "/home/user/pf/node/node-v19.1.3-linux-x64",
							Main:        false,
							InstalledOn: 168901728000,
						},
						{
							Version:     domain.Ver("19.1.4", t),
							Path:        "/home/user/pf/node/node-v19.1.4-linux-x64",
							Main:        false,
							InstalledOn: 168901727000,
						},

						{
							Version:     domain.Ver("v20.1.3", t),
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
			name: "latest main installation - major granularity",
			args: args{
				installedPackages: domain.InstalledPackages{
					Plugin: domain.Plugin{Name: "node", VersionGranularity: domain.VersionGranularityMajor, ExecutableRelativePath: "bin"},
					Items: []domain.InstalledPackage{
						{
							Version:     domain.Ver("19.1.3", t),
							Path:        "/home/user/pf/node/node-v19.1.3-linux-x64",
							Main:        true,
							InstalledOn: 1689017268000,
						}, {
							Version:     domain.Ver("v20.1.3", t),
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
		}, {
			name: "latest main installation - minor granularity",
			args: args{
				installedPackages: domain.InstalledPackages{
					Plugin: domain.Plugin{Name: "kotlin", VersionGranularity: domain.VersionGranularityMinor, ExecutableRelativePath: "bin"},
					Items: []domain.InstalledPackage{
						{
							Version:     domain.Ver("v1.8.22", t),
							Path:        "/home/user/pf/kotlin/kotlin-compiler-1.8.22",
							Main:        true,
							InstalledOn: 1689017268000,
						}, {
							Version:     domain.Ver("v1.9.0", t),
							Path:        "/home/user/pf/kotlin/kotlin-compiler-1.9.0",
							Main:        true,
							InstalledOn: 1689017267000,
						},
					},
				},
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
				installedPackages: domain.InstalledPackages{
					Plugin: domain.Plugin{Name: "soft-ver-man", VersionGranularity: domain.VersionGranularityMinor, ExecutableRelativePath: ""},
					Items: []domain.InstalledPackage{
						{
							Version:     domain.Ver("v0.5.0", t),
							Path:        "/home/user/pf/soft-ver-man/soft-ver-man-0.5.0",
							Main:        true,
							InstalledOn: 1689017268000,
						},
					},
				},
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
			err := initVariables(test.TestDirs{Home: dir}, tt.args.installedPackages)
			if err != nil {
				t.Errorf("Failed to init variables: %s", err)
			}
			test.AssertFileContent(path.Join(dir, config.HomeConfigDir), tt.expectedSpecificFileName, tt.expectedSpecificContent, t)
		})
	}
}
