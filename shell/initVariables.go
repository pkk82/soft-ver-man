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
	"github.com/pkk82/soft-ver-man/util/file"
	"path"
	"strings"
)

func initVariables(finder DirFinder, installedPackages domain.InstalledPackages) error {
	plugin := installedPackages.Plugin
	name := plugin.Name

	homeDir, err := finder.HomeDir()
	if err != nil {
		return err
	}

	err = initSvmRcFile(name, homeDir, finder)
	if err != nil {
		return err
	}

	err = initSpecificRcRile(installedPackages, homeDir, plugin)
	if err != nil {
		return err
	}

	return nil

}

func initSvmRcFile(name, homeDir string, finder DirFinder) error {

	softDir, err := finder.SoftDir()
	if err != nil {
		return err
	}

	exportLine := PrepareSvmSoftDirEnvVariable(softDir).ToExport()
	err = assertFileWithContent(path.Join(homeDir, config.HomeConfigDir, config.RcFile), exportLine,
		[]string{exportLine})
	if err != nil {
		return err
	}
	initLine := bashToLoad(rcName(name))
	err = assertFileWithContent(path.Join(homeDir, config.HomeConfigDir, config.RcFile), initLine,
		[]string{initLine})
	if err != nil {
		return err
	}
	return nil
}

func initSpecificRcRile(installedPackages domain.InstalledPackages, homeDir string, plugin domain.Plugin) error {

	variables, err := installedPackages.PrepareEnvVariables(plugin)
	if err != nil {
		return err
	}

	lines := variables.ToExport()

	err = file.OverrideFileWithContent(path.Join(homeDir, config.HomeConfigDir, rcName(plugin.Name)), lines)
	if err != nil {
		return err
	}
	return nil
}

func transformName(name string) string {
	if strings.Contains(name, "-") {
		parts := strings.Split(name, "-")
		short := ""
		for _, part := range parts {
			short += string(part[0])
		}
		return short
	} else {
		return name
	}
}

func rcName(name string) string {
	transformedName := transformName(name)
	return makeRcName(transformedName)
}
