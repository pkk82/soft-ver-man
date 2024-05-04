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

	tests := []struct {
		name    string
		arg     EnvVariables
		want    EnvVariables
		wantErr bool
	}{
		{
			name:    "empty",
			arg:     EnvVariables{Variables: make([]EnvVariable, 0)},
			want:    EnvVariables{Variables: make([]EnvVariable, 0)},
			wantErr: false,
		},
		{
			name:    "one resolved",
			arg:     EnvVariables{Variables: []EnvVariable{{Name: "JAVA_HOME", SuffixValue: "/usr/lib/jvm/java-11-openjdk-amd64"}}},
			want:    EnvVariables{Variables: []EnvVariable{{Name: "JAVA_HOME", SuffixValue: "/usr/lib/jvm/java-11-openjdk-amd64"}}},
			wantErr: false,
		},
		{
			name: "not resolved",
			arg: EnvVariables{Variables: []EnvVariable{
				usrLibDirEnvVariable,
				javaDirEnvVariable,
				java11HomeEnvVariable,
				javaHomeEnvVariable,
			}},
			want:    EnvVariables{Variables: []EnvVariable{usrLibDirEnvVariable, javaDirEnvVariableResolved, java11HomeEnvVariableResolved, javaHomeEnvVariableResolved}},
			wantErr: false,
		},
		{
			name: "inconsistent",
			arg: EnvVariables{Variables: []EnvVariable{
				usrLibDirEnvVariable,
				java11HomeEnvVariable,
				javaHomeEnvVariable,
			}},
			want:    EnvVariables{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			envVariables := tt.arg
			got, err := envVariables.Resolve()
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
