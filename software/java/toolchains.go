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

const toolchainsIndent = "  "

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
		entries = append(entries, renderToolchain(item))
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
			entries = append(entries, renderToolchain(item))
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
	for _, entry := range entries {
		lines = append(lines, formatToolchainEntry(entry))
	}
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

func renderToolchain(item toolchainEntry) string {
	return "<toolchain>" + item.InnerXML + "</toolchain>"
}

type xmlTokenKind int

const (
	xmlTagOpen xmlTokenKind = iota
	xmlTagClose
	xmlTagStandalone
	xmlText
)

type xmlToken struct {
	kind xmlTokenKind
	text string
}

// formatToolchainEntry re-indents a single toolchain entry with
// toolchainsIndent per level, one level deep inside the root element,
// leaving tag content untouched
func formatToolchainEntry(entry string) string {
	tokens := tokenizeXml(entry)
	lines := make([]string, 0, len(tokens))
	depth := 1
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		indent := strings.Repeat(toolchainsIndent, depth)
		switch token.kind {
		case xmlTagOpen:
			if i+2 < len(tokens) && tokens[i+1].kind == xmlText && tokens[i+2].kind == xmlTagClose {
				lines = append(lines, indent+token.text+tokens[i+1].text+tokens[i+2].text)
				i += 2
			} else if i+1 < len(tokens) && tokens[i+1].kind == xmlTagClose {
				lines = append(lines, indent+token.text+tokens[i+1].text)
				i++
			} else {
				lines = append(lines, indent+token.text)
				depth++
			}
		case xmlTagClose:
			depth--
			lines = append(lines, strings.Repeat(toolchainsIndent, depth)+token.text)
		default:
			lines = append(lines, indent+token.text)
		}
	}
	return strings.Join(lines, "\n")
}

func tokenizeXml(fragment string) []xmlToken {
	tokens := make([]xmlToken, 0)
	rest := fragment
	for len(rest) > 0 {
		tagStart := strings.Index(rest, "<")
		if tagStart == -1 {
			tokens = appendTextToken(tokens, rest)
			break
		}
		tokens = appendTextToken(tokens, rest[:tagStart])
		rest = rest[tagStart:]
		var tag string
		if strings.HasPrefix(rest, "<!--") {
			if end := strings.Index(rest, "-->"); end == -1 {
				tag = rest
			} else {
				tag = rest[:end+len("-->")]
			}
		} else {
			if end := strings.Index(rest, ">"); end == -1 {
				tag = rest
			} else {
				tag = rest[:end+1]
			}
		}
		rest = rest[len(tag):]
		tokens = append(tokens, xmlToken{kind: classifyTag(tag), text: tag})
	}
	return tokens
}

func appendTextToken(tokens []xmlToken, text string) []xmlToken {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return tokens
	}
	return append(tokens, xmlToken{kind: xmlText, text: trimmed})
}

func classifyTag(tag string) xmlTokenKind {
	if strings.HasPrefix(tag, "</") {
		return xmlTagClose
	}
	if strings.HasPrefix(tag, "<!") || strings.HasPrefix(tag, "<?") || strings.HasSuffix(tag, "/>") {
		return xmlTagStandalone
	}
	return xmlTagOpen
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
