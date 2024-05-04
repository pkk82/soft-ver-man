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
	"fmt"
	"os"
	"path"
	"strings"
)

type EnvVariable struct {
	Name           string
	PrefixVariable *EnvVariable
	SuffixValue    string
}

func (v EnvVariable) ToExport() string {
	return fmt.Sprintf("export %s", v.toString())
}

func (v EnvVariable) ToEnv() string {
	return fmt.Sprintf("env %s", v.toString())
}

func (v EnvVariable) toString() string {
	if v.PrefixVariable == nil {
		return fmt.Sprintf("%v=\"%v\"", v.Name, v.SuffixValue)
	}
	if v.SuffixValue == "" {
		return fmt.Sprintf("%v=\"$%v\"", v.Name, v.PrefixVariable.Name)
	}

	return fmt.Sprintf("%v=\"$%v%v%v\"", v.Name, v.PrefixVariable.Name, string(os.PathSeparator), v.SuffixValue)
}

type EnvVariables struct {
	Variables              []EnvVariable
	MainVariable           *EnvVariable
	ExecutableRelativePath string
}

func (envVariables EnvVariables) ToExport() []string {

	lines := make([]string, len(envVariables.Variables)+1)

	for i, envVariable := range envVariables.Variables {
		lines[i] = envVariable.ToExport()
	}

	mainVariable := *envVariables.MainVariable

	var executableRelativePath string
	if envVariables.ExecutableRelativePath == "" ||
		strings.HasPrefix(envVariables.ExecutableRelativePath, string(os.PathSeparator)) {
		executableRelativePath = envVariables.ExecutableRelativePath
	} else {
		executableRelativePath = path.Join(string(os.PathSeparator), envVariables.ExecutableRelativePath)
	}

	lines[len(lines)-1] = fmt.Sprintf("export PATH=\"$%v%v%v$PATH\"",
		mainVariable.Name,
		executableRelativePath,
		string(os.PathListSeparator))

	return lines
}
