package node

import (
	"reflect"
	"testing"
)

func TestSupportedFiles(t *testing.T) {
	filesPerVersions := []PackagesPerVersion{
		{Version: "v0.1.14", Date: "2011-08-26", Files: []string{"src"}},
		{Version: "v0.6.12", Date: "2012-03-02", Files: []string{"osx-x64-pkg", "src", "win-x86-exe"}},
		{Version: "v0.8.6", Date: "2012-08-06", Files: []string{"linux-x64",
			"linux-x86",
			"osx-x64-pkg",
			"osx-x64-tar",
			"osx-x86-tar",
			"src",
			"sunos-x64",
			"sunos-x86",
			"win-x64-exe",
			"win-x86-exe",
			"win-x86-msi"}},
		{Version: "v20.2.0", Date: "2023-05-16", Files: []string{"aix-ppc64",
			"headers",
			"linux-arm64",
			"linux-armv7l",
			"linux-ppc64le",
			"linux-s390x",
			"linux-x64",
			"osx-arm64-tar",
			"osx-x64-pkg",
			"osx-x64-tar",
			"src",
			"win-arm64-7z",
			"win-arm64-zip",
			"win-x64-7z",
			"win-x64-exe",
			"win-x64-msi",
			"win-x64-zip",
			"win-x86-7z",
			"win-x86-exe",
			"win-x86-msi",
			"win-x86-zip"}},
	}

	/*
	    "darwin": Indicates the macOS operating system.
	    "linux": Indicates a Linux-based operating system.
	    "windows": Indicates the Windows operating system.
	    "freebsd": Indicates the FreeBSD operating system.
	    "netbsd": Indicates the NetBSD operating system.
	    "openbsd": Indicates the OpenBSD operating system.
	    "dragonfly": Indicates the DragonFly BSD operating system.
	    "plan9": Indicates the Plan 9 operating system.
	    "solaris": Indicates the Oracle Solaris operating system.

	   "amd64": Indicates the x86-64 architecture, also known as Intel 64 or AMD64.
	   "386": Indicates the 32-bit x86 architecture.
	   "arm": Indicates the ARM architecture.
	   "arm64": Indicates the 64-bit ARM architecture.
	   "ppc64": Indicates the 64-bit PowerPC architecture.
	   "ppc64le": Indicates the little-endian 64-bit PowerPC architecture.
	   "mips": Indicates the MIPS architecture.
	   "mipsle": Indicates the little-endian MIPS architecture.
	   "mips64": Indicates the 64-bit MIPS architecture.
	   "mips64le": Indicates the little-endian 64-bit MIPS architecture.
	*/

	actual := supportedPackages(&filesPerVersions, "linux", "amd64")
	expected := []Package{
		{Version: "v0.8.6", FileName: "node-v0.8.6-linux-x64.tar.gz"},
		{Version: "v20.2.0", FileName: "node-v20.2.0-linux-x64.tar.gz"},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected: %v, got: %v", expected, actual)
	}

	actual = supportedPackages(&filesPerVersions, "windows", "amd64")
	expected = []Package{
		{Version: "v20.2.0", FileName: "node-v20.2.0-win-x64.zip"},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected: %v, got: %v", expected, actual)
	}

	actual = supportedPackages(&filesPerVersions, "darwin", "amd64")
	expected = []Package{
		{Version: "v0.8.6", FileName: "node-v0.8.6-darwin-x64.tar.gz"},
		{Version: "v20.2.0", FileName: "node-v20.2.0-darwin-x64.tar.gz"},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected: %v, got: %v", expected, actual)
	}

	actual = supportedPackages(&filesPerVersions, "darwin", "arm64")
	expected = []Package{
		{Version: "v20.2.0", FileName: "node-v20.2.0-darwin-arm64.tar.gz"},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected: %v, got: %v", expected, actual)
	}
}

func TestDownloadLink(t *testing.T) {
	versions := []Package{
		{Version: "v0.8.6", FileName: "node-v0.8.6-darwin-x64.tar.gz"},
		{Version: "v20.2.0", FileName: "node-v20.2.0-darwin-arm64.tar.gz"},
		{Version: "v20.2.0", FileName: "node-v20.2.0-win-x64.zip"},
	}
	expected := []string{
		"https://nodejs.org/dist/v0.8.6/node-v0.8.6-darwin-x64.tar.gz",
		"https://nodejs.org/dist/v20.2.0/node-v20.2.0-darwin-arm64.tar.gz",
		"https://nodejs.org/dist/v20.2.0/node-v20.2.0-win-x64.zip",
	}

	for i, version := range versions {
		actual := version.DownloadLink()
		if actual != expected[i] {
			t.Errorf("Expected: %v, got: %v", expected[i], actual)
		}
	}
}
