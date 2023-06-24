package node

import (
	"reflect"
	"testing"
)

func TestSupportedFiles(t *testing.T) {
	versions := []Version{
		{Id: "v0.1.14", Date: "2011-08-26", Files: []string{"src"}},
		{Id: "v0.6.12", Date: "2012-03-02", Files: []string{"osx-x64-pkg", "src", "win-x86-exe"}},
		{Id: "v0.8.6", Date: "2012-08-06", Files: []string{"linux-x64",
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
		{Id: "v20.2.0", Date: "2023-05-16", Files: []string{"aix-ppc64",
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

	supported := supportedVersions(&versions, "linux", "amd64")
	if !reflect.DeepEqual(supported, []Version{versions[2], versions[3]}) {
		t.Errorf("Expected %v, got %v", []Version{versions[2], versions[3]}, supported)
	}

	supported = supportedVersions(&versions, "windows", "amd64")
	if !reflect.DeepEqual(supported, []Version{versions[3]}) {
		t.Errorf("Expected %v, got %v", []Version{versions[3]}, supported)
	}

	supported = supportedVersions(&versions, "darwin", "amd64")
	if !reflect.DeepEqual(supported, []Version{versions[2], versions[3]}) {
		t.Errorf("Expected %v, got %v", []Version{versions[2], versions[3]}, supported)
	}

	supported = supportedVersions(&versions, "darwin", "arm64")
	if !reflect.DeepEqual(supported, []Version{versions[3]}) {
		t.Errorf("Expected %v, got %v", []Version{versions[3]}, supported)
	}
}
