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
	"github.com/pkk82/soft-ver-man/domain"
	"github.com/pkk82/soft-ver-man/util/file"
	"github.com/pkk82/soft-ver-man/util/test"
	"os"
	"path/filepath"
	"testing"
)

const jdk17Entry = `  <toolchain>
    <type>jdk</type>
    <provides>
      <version>17</version>
    </provides>
    <configuration>
      <jdkHome>/opt/jdk-17.0.2</jdkHome>
    </configuration>
  </toolchain>`

const netbeansEntry = `  <toolchain>
    <type>netbeans</type>
    <provides>
      <version>12</version>
    </provides>
    <configuration>
      <installDir>/opt/netbeans</installDir>
    </configuration>
  </toolchain>`

func Test_upsertToolchainEntry_createsFile(t *testing.T) {
	toolchainsPath := filepath.Join(test.CreateTestDir(t), ".m2", "toolchains.xml")

	err := upsertToolchainEntry(toolchainsPath, 21, "/opt/jdk-21.0.4")
	if err != nil {
		t.Fatalf("Failed to upsert toolchain entry: %s", err)
	}

	expected := xmlDeclaration + "\n" +
		defaultToolchainsOpeningTag + `
  <toolchain>
    <type>jdk</type>
    <provides>
      <version>21</version>
      <vendor>zulu</vendor>
    </provides>
    <configuration>
      <jdkHome>/opt/jdk-21.0.4</jdkHome>
    </configuration>
  </toolchain>
</toolchains>
`
	assertFileContent(t, toolchainsPath, expected)
}

func Test_upsertToolchainEntry_addsEntryToExistingFile(t *testing.T) {
	toolchainsPath := filepath.Join(test.CreateTestDir(t), "toolchains.xml")
	existing := xmlDeclaration + "\n" +
		defaultToolchainsOpeningTag + "\n" +
		netbeansEntry + "\n" +
		jdk17Entry + "\n" +
		"</toolchains>\n"
	writeFile(t, toolchainsPath, existing)

	err := upsertToolchainEntry(toolchainsPath, 21, "/opt/jdk-21.0.4")
	if err != nil {
		t.Fatalf("Failed to upsert toolchain entry: %s", err)
	}

	expected := xmlDeclaration + "\n" +
		defaultToolchainsOpeningTag + "\n" +
		netbeansEntry + "\n" +
		jdk17Entry + `
  <toolchain>
    <type>jdk</type>
    <provides>
      <version>21</version>
      <vendor>zulu</vendor>
    </provides>
    <configuration>
      <jdkHome>/opt/jdk-21.0.4</jdkHome>
    </configuration>
  </toolchain>
</toolchains>
`
	assertFileContent(t, toolchainsPath, expected)
}

func Test_upsertToolchainEntry_overridesEntryWithSameMajorVersion(t *testing.T) {
	toolchainsPath := filepath.Join(test.CreateTestDir(t), "toolchains.xml")
	existing := xmlDeclaration + "\n" +
		defaultToolchainsOpeningTag + `
  <toolchain>
    <type>jdk</type>
    <provides>
      <version>21</version>
      <vendor>zulu</vendor>
    </provides>
    <configuration>
      <jdkHome>/opt/jdk-21.0.1</jdkHome>
    </configuration>
  </toolchain>
` + jdk17Entry + "\n" +
		"</toolchains>\n"
	writeFile(t, toolchainsPath, existing)

	err := upsertToolchainEntry(toolchainsPath, 21, "/opt/jdk-21.0.4")
	if err != nil {
		t.Fatalf("Failed to upsert toolchain entry: %s", err)
	}

	expected := xmlDeclaration + "\n" +
		defaultToolchainsOpeningTag + `
  <toolchain>
    <type>jdk</type>
    <provides>
      <version>21</version>
      <vendor>zulu</vendor>
    </provides>
    <configuration>
      <jdkHome>/opt/jdk-21.0.4</jdkHome>
    </configuration>
  </toolchain>
` + jdk17Entry + "\n" +
		"</toolchains>\n"
	assertFileContent(t, toolchainsPath, expected)
}

func Test_removeToolchainEntry_removesEntry(t *testing.T) {
	toolchainsPath := filepath.Join(test.CreateTestDir(t), "toolchains.xml")
	existing := xmlDeclaration + "\n" +
		defaultToolchainsOpeningTag + "\n" +
		netbeansEntry + "\n" +
		jdk17Entry + `
  <toolchain>
    <type>jdk</type>
    <provides>
      <version>21</version>
      <vendor>zulu</vendor>
    </provides>
    <configuration>
      <jdkHome>/opt/jdk-21.0.4</jdkHome>
    </configuration>
  </toolchain>
</toolchains>
`
	writeFile(t, toolchainsPath, existing)

	err := removeToolchainEntry(toolchainsPath, 21)
	if err != nil {
		t.Fatalf("Failed to remove toolchain entry: %s", err)
	}

	expected := xmlDeclaration + "\n" +
		defaultToolchainsOpeningTag + "\n" +
		netbeansEntry + "\n" +
		jdk17Entry + "\n" +
		"</toolchains>\n"
	assertFileContent(t, toolchainsPath, expected)
}

func Test_removeToolchainEntry_missingFile(t *testing.T) {
	toolchainsPath := filepath.Join(test.CreateTestDir(t), "toolchains.xml")

	err := removeToolchainEntry(toolchainsPath, 21)
	if err != nil {
		t.Fatalf("Failed to remove toolchain entry: %s", err)
	}

	exists, err := file.FileExists(toolchainsPath)
	if err != nil {
		t.Fatalf("Failed to check file existence: %s", err)
	}
	if exists {
		t.Errorf("Expected no file to be created")
	}
}

func Test_upsertToolchainEntry_preservesCustomOpeningTag(t *testing.T) {
	toolchainsPath := filepath.Join(test.CreateTestDir(t), "toolchains.xml")
	existing := "<toolchains>\n" +
		jdk17Entry + "\n" +
		"</toolchains>\n"
	writeFile(t, toolchainsPath, existing)

	err := upsertToolchainEntry(toolchainsPath, 17, "/opt/jdk-17.0.9")
	if err != nil {
		t.Fatalf("Failed to upsert toolchain entry: %s", err)
	}

	expected := xmlDeclaration + `
<toolchains>
  <toolchain>
    <type>jdk</type>
    <provides>
      <version>17</version>
      <vendor>zulu</vendor>
    </provides>
    <configuration>
      <jdkHome>/opt/jdk-17.0.9</jdkHome>
    </configuration>
  </toolchain>
</toolchains>
`
	assertFileContent(t, toolchainsPath, expected)
}

func Test_upsertToolchainEntry_reformatsExistingEntries(t *testing.T) {
	toolchainsPath := filepath.Join(test.CreateTestDir(t), "toolchains.xml")
	existing := xmlDeclaration + "\n" +
		defaultToolchainsOpeningTag + `
    <toolchain>
	<type>netbeans</type>
            <provides><version>12</version></provides>
        <configuration>
	        <!-- custom install dir -->
            <installDir>/opt/netbeans</installDir></configuration>
      </toolchain>
</toolchains>
`
	writeFile(t, toolchainsPath, existing)

	err := upsertToolchainEntry(toolchainsPath, 21, "/opt/jdk-21.0.4")
	if err != nil {
		t.Fatalf("Failed to upsert toolchain entry: %s", err)
	}

	expected := xmlDeclaration + "\n" +
		defaultToolchainsOpeningTag + `
  <toolchain>
    <type>netbeans</type>
    <provides>
      <version>12</version>
    </provides>
    <configuration>
      <!-- custom install dir -->
      <installDir>/opt/netbeans</installDir>
    </configuration>
  </toolchain>
  <toolchain>
    <type>jdk</type>
    <provides>
      <version>21</version>
      <vendor>zulu</vendor>
    </provides>
    <configuration>
      <jdkHome>/opt/jdk-21.0.4</jdkHome>
    </configuration>
  </toolchain>
</toolchains>
`
	assertFileContent(t, toolchainsPath, expected)
}

func Test_findLatestWithMajor(t *testing.T) {
	items := []domain.InstalledPackage{
		{Version: domain.Ver("21.0.1", t), Path: "/opt/jdk-21.0.1"},
		{Version: domain.Ver("17.0.2", t), Path: "/opt/jdk-17.0.2"},
		{Version: domain.Ver("21.0.4", t), Path: "/opt/jdk-21.0.4"},
	}

	latest := findLatestWithMajor(items, 21)
	if latest == nil || latest.Path != "/opt/jdk-21.0.4" {
		t.Errorf("Expected /opt/jdk-21.0.4, got %v", latest)
	}

	latest = findLatestWithMajor(items, 11)
	if latest != nil {
		t.Errorf("Expected nil, got %v", latest)
	}
}

func writeFile(t *testing.T, path, content string) {
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create file: %s", err)
	}
}

func assertFileContent(t *testing.T, path, expected string) {
	actual, err := file.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read file: %s", err)
	}
	if actual != expected {
		t.Errorf("Expected content:\n%s\nActual content:\n%s", expected, actual)
	}
}
