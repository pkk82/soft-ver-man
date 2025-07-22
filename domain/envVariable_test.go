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

package domain

import (
	"reflect"
	"testing"
)

func Test_EnvVariables_Resolve(t *testing.T) {
	usrLibDirEnvVariable := EnvVariable{Name: "USR_LIB_DIR", SuffixValue: "/usr/lib"}
	javaDirEnvVariable := EnvVariable{Name: "JAVA_DIR", PrefixVariable: &usrLibDirEnvVariable, SuffixValue: "jvm"}
	javaDirEnvVariableResolved := EnvVariable{Name: "JAVA_DIR", SuffixValue: "/usr/lib/jvm"}
	java11HomeEnvVariable := EnvVariable{Name: "JAVA_11", PrefixVariable: &javaDirEnvVariable, SuffixValue: "java-11-openjdk-amd64"}
	java11HomeEnvVariableResolved := EnvVariable{Name: "JAVA_11", SuffixValue: "/usr/lib/jvm/java-11-openjdk-amd64"}
	javaHomeEnvVariable := EnvVariable{Name: "JAVA", PrefixVariable: &java11HomeEnvVariable}
	javaHomeEnvVariableResolved := EnvVariable{Name: "JAVA", SuffixValue: "/usr/lib/jvm/java-11-openjdk-amd64"}

	type args struct {
		envVariables   EnvVariables
		extraVariables []EnvVariable
	}

	tests := []struct {
		name    string
		args    args
		want    EnvVariables
		wantErr bool
	}{
		{
			name: "empty",
			args: args{
				envVariables:   EnvVariables{Variables: make([]EnvVariable, 0)},
				extraVariables: make([]EnvVariable, 0),
			},
			want:    EnvVariables{Variables: make([]EnvVariable, 0)},
			wantErr: false,
		},
		{
			name: "one resolved",
			args: args{
				envVariables:   EnvVariables{Variables: []EnvVariable{{Name: "JAVA_HOME", SuffixValue: "/usr/lib/jvm/java-11-openjdk-amd64"}}},
				extraVariables: make([]EnvVariable, 0),
			},
			want:    EnvVariables{Variables: []EnvVariable{{Name: "JAVA_HOME", SuffixValue: "/usr/lib/jvm/java-11-openjdk-amd64"}}},
			wantErr: false,
		},
		{
			name: "not resolved",
			args: args{
				envVariables: EnvVariables{Variables: []EnvVariable{
					usrLibDirEnvVariable,
					javaDirEnvVariable,
					java11HomeEnvVariable,
					javaHomeEnvVariable,
				}},
				extraVariables: make([]EnvVariable, 0),
			},
			want: EnvVariables{Variables: []EnvVariable{
				usrLibDirEnvVariable,
				javaDirEnvVariableResolved,
				java11HomeEnvVariableResolved,
				javaHomeEnvVariableResolved,
			}},
			wantErr: false,
		},
		{
			name: "inconsistent",
			args: args{
				envVariables: EnvVariables{Variables: []EnvVariable{
					usrLibDirEnvVariable,
					java11HomeEnvVariable,
					javaHomeEnvVariable,
				}},
				extraVariables: make([]EnvVariable, 0),
			},
			want:    EnvVariables{},
			wantErr: true,
		}, {
			name: "extra resolved",
			args: args{
				envVariables: EnvVariables{Variables: []EnvVariable{
					usrLibDirEnvVariable,
					java11HomeEnvVariable,
					javaHomeEnvVariable,
				}},
				extraVariables: []EnvVariable{javaDirEnvVariable},
			},
			want: EnvVariables{Variables: []EnvVariable{
				usrLibDirEnvVariable,
				java11HomeEnvVariableResolved,
				javaHomeEnvVariableResolved},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			envVariables := tt.args.envVariables
			got, err := envVariables.Resolve(tt.args.extraVariables)
			if (err != nil) != tt.wantErr {
				t.Errorf("EnvVariables.Resolve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EnvVaariables.Resolve() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_EvnVariables_extractToHere(t *testing.T) {
	dirVariable := EnvVariable{Name: "SVM_SOFT_GO_DIR", SuffixValue: "/go"}
	go115 := EnvVariable{Name: "GO_1_15_ROOT", PrefixVariable: &dirVariable, SuffixValue: "go1.15.15.linux-amd64"}
	go114 := EnvVariable{Name: "GO_1_14_ROOT", PrefixVariable: &dirVariable, SuffixValue: "go1.14.15.linux-amd64"}
	goMainTo115 := EnvVariable{Name: "GOROOT", PrefixVariable: &go115}

	type fields struct {
		Variables              []EnvVariable
		MainVariable           *EnvVariable
		ExecutableRelativePath string
	}
	type args struct {
		suffixValue string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "empty",
			fields: fields{
				Variables: make([]EnvVariable, 0),
			},
			args: args{
				suffixValue: "test",
			},
			want:    []string{},
			wantErr: true,
		},
		{
			name: "not empty - not found",
			fields: fields{
				Variables: []EnvVariable{
					dirVariable, goMainTo115, go114, go115,
				},
				MainVariable:           &goMainTo115,
				ExecutableRelativePath: "bin",
			},
			args: args{
				suffixValue: "go1.14.14.linux-amd64",
			},
			want:    []string{},
			wantErr: true,
		}, {
			name: "not empty - found",
			fields: fields{
				Variables: []EnvVariable{
					dirVariable, goMainTo115, go114, go115,
				},
				MainVariable:           &goMainTo115,
				ExecutableRelativePath: "bin",
			},
			args: args{
				suffixValue: "go1.14.15.linux-amd64",
			},
			want: []string{
				"export GOROOT=\"$GO_1_14_ROOT\"",
				"export PATH=\"$GOROOT/bin:$PATH\"",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			envVariables := EnvVariables{
				Variables:              tt.fields.Variables,
				MainVariable:           tt.fields.MainVariable,
				ExecutableRelativePath: tt.fields.ExecutableRelativePath,
			}
			got, err := envVariables.ExtractToHere(tt.args.suffixValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("EnvVariables.ExtractToHere() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.ToExport(), tt.want) {
				t.Errorf("EnvVariables.ExtractToHere().ToExport() = %v, want %v", got.ToExport(), tt.want)
			}
		})
	}
}
