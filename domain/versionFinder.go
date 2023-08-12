/*
Copyright © 2023 Piotr Kozak <piotrkrzysztofkozak@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package domain

import (
	"errors"
	"sort"
	"strconv"
	"strings"
)

func FindVersion(version string, versions []string) (Version, int, error) {
	parsedVersions, err := parseVersions(versions)
	if err != nil {
		return Version{}, -1, err
	}
	sort.Slice(parsedVersions, func(i, j int) bool {
		return CompareDesc(parsedVersions[i], parsedVersions[j])
	})
	parsedVersion, err := parseVersion(version)
	if err != nil {
		return Version{}, -1, err
	}

	matchingVersion := Version{}
	for _, v := range parsedVersions {
		if (v.major == parsedVersion.major || parsedVersion.major == -1) &&
			(v.minor == parsedVersion.minor || parsedVersion.minor == -1) &&
			(v.patch == parsedVersion.patch || parsedVersion.patch == -1) {
			matchingVersion = v
			break
		}
	}
	for i, v := range versions {
		if v == matchingVersion.Value {
			return matchingVersion, i, nil
		}

	}
	return Version{}, -1, errors.New("No version found for: " + version)

}

func parseVersions(versions []string) ([]Version, error) {
	parsedVersions := make([]Version, 0)
	for _, version := range versions {
		v, err := parseVersion(version)
		if err != nil {
			return nil, err
		}
		parsedVersions = append(parsedVersions, v)
	}
	return parsedVersions, nil
}

func parseVersion(version string) (Version, error) {
	var versionWithoutV string
	if strings.HasPrefix(version, "v") {
		versionWithoutV = strings.TrimPrefix(version, "v")
	} else {
		versionWithoutV = version
	}

	splitVersion := strings.Split(versionWithoutV, ".")
	var major, minor, patch, build = -1, -1, -1, -1
	var errVersion = Version{Value: version, major: major, minor: minor, patch: patch, build: build}
	var err error
	if len(splitVersion) > 0 && splitVersion[0] != "" {
		major, err = strconv.Atoi(splitVersion[0])
	}
	if err != nil {
		return errVersion, err
	}
	if len(splitVersion) > 1 && splitVersion[1] != "" {
		minor, err = strconv.Atoi(splitVersion[1])
	}
	if err != nil {
		return errVersion, err
	}
	if len(splitVersion) > 2 && splitVersion[2] != "" {
		patch, err = strconv.Atoi(splitVersion[2])
	}
	if err != nil {
		return errVersion, err
	}
	if len(splitVersion) > 3 && splitVersion[3] != "" {
		build, err = strconv.Atoi(splitVersion[3])
	}
	if err != nil {
		return errVersion, err
	}

	return Version{Value: version, major: major, minor: minor, patch: patch, build: build}, nil
}
