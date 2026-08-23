/*
 * Copyright © 2026 Piotr Kozak <piotrkrzysztofkozak@gmail.com>
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
package java

import (
	"encoding/xml"
	"fmt"
	"github.com/pkk82/soft-ver-man/config"
	"github.com/pkk82/soft-ver-man/domain"
	"github.com/pkk82/soft-ver-man/util/file"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const xmlDeclaration = `<?xml version="1.0" encoding="UTF-8"?>`

const defaultToolchainsOpeningTag = `<toolchains xmlns="http://maven.apache.org/TOOLCHAINS/1.1.0"
            xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
            xsi:schemaLocation="http://maven.apache.org/TOOLCHAINS/1.1.0 https://maven.apache.org/xsd/toolchains-1.1.0.xsd">`

var toolchainsOpeningTagPattern = regexp.MustCompile(`(?s)<toolchains(\s[^>]*)?>`)

type toolchainsDoc struct {
	XMLName xml.Name         `xml:"toolchains"`
	Items   []toolchainEntry `xml:"toolchain"`
}

type toolchainEntry struct {
	Type     string            `xml:"type"`
	Provides toolchainProvides `xml:"provides"`
	InnerXML string            `xml:",innerxml"`
}

type toolchainProvides struct {
	Version string `xml:"version"`
}

func registerInToolchains(installedPackage domain.InstalledPackage) error {
	toolchainsPath, err := toolchainsFilePath()
	if err != nil {
		return err
	}
	return upsertToolchainEntry(toolchainsPath, installedPackage.Version.Major(), installedPackage.Path)
}

func unregisterFromToolchains(version domain.Version) error {
	toolchainsPath, err := toolchainsFilePath()
	if err != nil {
		return err
	}
	exists, err := file.FileExists(toolchainsPath)
	if err != nil || !exists {
		return err
	}
	installedPackages, err := config.LoadInstalledPackages(Name)
	if err != nil {
		return err
	}
	replacement := findLatestWithMajor(installedPackages.Items, version.Major())
	if replacement != nil {
		return upsertToolchainEntry(toolchainsPath, version.Major(), replacement.Path)
	}
	return removeToolchainEntry(toolchainsPath, version.Major())
}

func toolchainsFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".m2", "toolchains.xml"), nil
}

func findLatestWithMajor(items []domain.InstalledPackage, major int) *domain.InstalledPackage {
	var latest *domain.InstalledPackage
	for i, item := range items {
		if item.Version.Major() != major {
			continue
		}
		if latest == nil || domain.CompareDesc(item.Version, latest.Version) {
			latest = &items[i]
		}
	}
	return latest
}

func upsertToolchainEntry(toolchainsPath string, major int, jdkHome string) error {
	doc, content, err := readToolchains(toolchainsPath)
	if err != nil {
		return err
	}
	newEntry := renderJdkToolchain(major, jdkHome)
	entries := make([]string, 0, len(doc.Items)+1)
	replaced := false
	for _, item := range doc.Items {
		if isJdkEntry(item, major) {
			if !replaced {
				entries = append(entries, newEntry)
				replaced = true
			}
			continue
		}
		entries = append(entries, renderVerbatim(item))
	}
	if !replaced {
		entries = append(entries, newEntry)
	}
	return writeToolchains(toolchainsPath, content, entries)
}

func removeToolchainEntry(toolchainsPath string, major int) error {
	exists, err := file.FileExists(toolchainsPath)
	if err != nil || !exists {
		return err
	}
	doc, content, err := readToolchains(toolchainsPath)
	if err != nil {
		return err
	}
	entries := make([]string, 0, len(doc.Items))
	for _, item := range doc.Items {
		if !isJdkEntry(item, major) {
			entries = append(entries, renderVerbatim(item))
		}
	}
	return writeToolchains(toolchainsPath, content, entries)
}

func readToolchains(toolchainsPath string) (toolchainsDoc, string, error) {
	exists, err := file.FileExists(toolchainsPath)
	if err != nil || !exists {
		return toolchainsDoc{}, "", err
	}
	content, err := file.ReadFile(toolchainsPath)
	if err != nil {
		return toolchainsDoc{}, "", err
	}
	if strings.TrimSpace(content) == "" {
		return toolchainsDoc{}, "", nil
	}
	var doc toolchainsDoc
	err = xml.Unmarshal([]byte(content), &doc)
	if err != nil {
		return toolchainsDoc{}, "", err
	}
	return doc, content, nil
}

func writeToolchains(toolchainsPath, originalContent string, entries []string) error {
	lines := []string{xmlDeclaration, toolchainsOpeningTag(originalContent)}
	lines = append(lines, entries...)
	lines = append(lines, "</toolchains>", "")
	return file.OverrideFileWithContent(toolchainsPath, lines)
}

func toolchainsOpeningTag(originalContent string) string {
	match := toolchainsOpeningTagPattern.FindString(originalContent)
	if match == "" {
		return defaultToolchainsOpeningTag
	}
	return match
}

func isJdkEntry(item toolchainEntry, major int) bool {
	return strings.TrimSpace(item.Type) == "jdk" &&
		strings.TrimSpace(item.Provides.Version) == strconv.Itoa(major)
}

func renderVerbatim(item toolchainEntry) string {
	return "  <toolchain>" + item.InnerXML + "</toolchain>"
}

func renderJdkToolchain(major int, jdkHome string) string {
	var escapedJdkHome strings.Builder
	err := xml.EscapeText(&escapedJdkHome, []byte(jdkHome))
	if err != nil {
		escapedJdkHome.Reset()
		escapedJdkHome.WriteString(jdkHome)
	}
	return strings.Join([]string{
		"  <toolchain>",
		"    <type>jdk</type>",
		"    <provides>",
		fmt.Sprintf("      <version>%d</version>", major),
		"      <vendor>zulu</vendor>",
		"    </provides>",
		"    <configuration>",
		fmt.Sprintf("      <jdkHome>%s</jdkHome>", escapedJdkHome.String()),
		"    </configuration>",
		"  </toolchain>",
	}, "\n")
}
