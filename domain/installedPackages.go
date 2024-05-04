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
	"errors"
	"sort"
)

type InstalledPackages struct {
	Plugin Plugin
	Items  []InstalledPackage
}

func (installedPackages *InstalledPackages) IsInstalled(version Version) bool {
	for _, item := range installedPackages.Items {
		if item.Version == version {
			return true
		}
	}
	return false
}

func (installedPackages *InstalledPackages) Add(installedPackage InstalledPackage) {
	installedPackages.Items = append(installedPackages.Items, installedPackage)
}

func (installedPackages *InstalledPackages) RemoveByVersion(version Version) *InstalledPackage {
	newItems := make([]InstalledPackage, 0)
	var itemToRemove *InstalledPackage = nil
	for _, item := range installedPackages.Items {
		if item.Version != version {
			newItems = append(newItems, item)
		} else {
			copyItem := item
			itemToRemove = &copyItem
		}
	}
	installedPackages.Items = newItems
	return itemToRemove
}

func (installedPackages *InstalledPackages) FoundMain() (*InstalledPackage, error) {

	var items = make([]InstalledPackage, len(installedPackages.Items))
	copy(items, installedPackages.Items)

	if len(items) == 0 {
		return nil, errors.New("no packages found")
	}

	// find latest main
	sort.Slice(items, func(i, j int) bool {
		return items[i].InstalledOn > items[j].InstalledOn
	})

	for _, p := range items {
		if p.Main {
			return &p, nil
		}
	}

	// find latest version
	sort.Slice(items, func(i, j int) bool {
		return CompareDesc(items[i].Version, items[j].Version)
	})

	return &items[0], nil

}