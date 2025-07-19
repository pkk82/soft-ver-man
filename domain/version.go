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

type Version struct {
	Value       string
	major       int
	minor       int
	patchString string
	patch       int
	buildString string
	build       int
}

func NewVersion(version string) (Version, error) {
	v, err := parseVersion(version)
	if err != nil {
		return Version{}, err
	}
	return v, nil
}

func (version Version) Major() int {
	return version.major
}

func (version Version) Minor() int {
	return version.minor
}

func CompareDesc(v1, v2 Version) bool {
	return v1.major > v2.major ||
		(v1.major == v2.major && v1.minor > v2.minor) ||
		(v1.major == v2.major && v1.minor == v2.minor && compareEmptyStringDesc(v1.patchString, v2.patchString)) ||
		(v1.major == v2.major && v1.minor == v2.minor && v1.patchString == v2.patchString && v1.patch > v2.patch) ||
		(v1.major == v2.major && v1.minor == v2.minor && v1.patchString == v2.patchString && v1.patch == v2.patch && compareEmptyStringDesc(v1.buildString, v2.buildString)) ||
		(v1.major == v2.major && v1.minor == v2.minor && v1.patchString == v2.patchString && v1.patch == v2.patch && v1.buildString == v2.buildString && v1.build > v2.build)
}

func compareEmptyStringDesc(v1, v2 string) bool {
	if v1 == "" && v2 != "" {
		return true
	}
	if v1 != "" && v2 == "" {
		return false
	}
	if v1 > v2 {
		return true
	}
	return false
}

func CompareAsc(v1, v2 Version) bool {
	return CompareDesc(v2, v1)
}
