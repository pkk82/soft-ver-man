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
package software

import (
	"errors"
	"github.com/pkk82/soft-ver-man/config"
	"github.com/pkk82/soft-ver-man/domain"
	"os"
)

func Uninstall(plugin domain.Plugin, inputVersion string) error {
	installedPackages, err := config.LoadInstalledPackages(plugin.Name)
	if err != nil {
		return err
	}

	versions := make([]string, len(installedPackages.Items))
	for i, item := range installedPackages.Items {
		versions[i] = item.Version.Value
	}

	version, _, err := domain.FindVersion(inputVersion, versions)
	if err != nil {
		return errors.New("Version " + inputVersion + " is not installed")
	}

	removedItem := installedPackages.RemoveByVersion(version)
	if removedItem == nil {
		return errors.New("Version " + inputVersion + " is not installed")
	}

	err = os.RemoveAll(removedItem.Path)
	if err != nil {
		return err
	}

	err = config.StoreInstalledPackages(installedPackages)
	if err != nil {
		return err
	}

	err = plugin.PostUninstall(version)
	if err != nil {
		return err
	}

	return nil

}
