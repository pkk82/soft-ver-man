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

func (envVariables EnvVariables) Resolve(extraEnvVariables []EnvVariable) (EnvVariables, error) {

	envVariablesByName := make(map[string]EnvVariable)
	for _, envVariable := range envVariables.Variables {
		envVariablesByName[envVariable.Name] = envVariable
	}
	for _, extraEnvVariables := range extraEnvVariables {
		envVariablesByName[extraEnvVariables.Name] = extraEnvVariables
	}

	resolvedVariables := make([]EnvVariable, len(envVariables.Variables))
	for index := range resolvedVariables {
		resolveVariable, err := resolve(envVariables.Variables[index], envVariablesByName)
		if err != nil {
			return EnvVariables{}, err
		}
		resolvedVariables[index] = resolveVariable
	}

	return EnvVariables{
		Variables:              resolvedVariables,
		MainVariable:           nil,
		ExecutableRelativePath: envVariables.ExecutableRelativePath,
	}, nil

}

func resolve(variable EnvVariable, envVariablesByName map[string]EnvVariable) (EnvVariable, error) {

	suffixValue := variable.SuffixValue
	iVariable := variable
	for iVariable.PrefixVariable != nil {

		mapVariable, ok := envVariablesByName[iVariable.PrefixVariable.Name]
		if !ok {
			return EnvVariable{}, fmt.Errorf("variable %v not found", iVariable.PrefixVariable.Name)
		}
		iVariable = mapVariable

		if suffixValue == "" {
			suffixValue = iVariable.SuffixValue
		} else {
			suffixValue = iVariable.SuffixValue + string(os.PathSeparator) + suffixValue
		}

	}

	return EnvVariable{
		Name:        variable.Name,
		SuffixValue: suffixValue,
	}, nil
}
