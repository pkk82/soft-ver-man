package intellij_test

import (
	"fmt"
	"github.com/pkk82/soft-ver-man/domain"
	"github.com/pkk82/soft-ver-man/software/intellij"
	"testing"
)

func Test_calculateDownloadUrl(t *testing.T) {
	plugin := domain.GetPlugin(intellij.Name)

	tests := []struct {
		version      domain.Version
		os           string
		arch         string
		expectedPath string
		expectedType domain.Type
	}{
		{
			version:      version("2024.1", t),
			os:           "linux",
			arch:         "amd64",
			expectedPath: "https://download.jetbrains.com/idea/ideaIU-2024.1.tar.gz",
			expectedType: domain.TAR_GZ,
		}, {
			version:      version("2023.3.6", t),
			os:           "linux",
			arch:         "arm64",
			expectedPath: "https://download.jetbrains.com/idea/ideaIU-2023.3.6-aarch64.tar.gz",
			expectedType: domain.TAR_GZ,
		},
		{
			version:      version("2023.1.6", t),
			os:           "windows",
			arch:         "amd64",
			expectedPath: "https://download.jetbrains.com/idea/ideaIU-2023.1.6.win.zip",
			expectedType: domain.ZIP,
		},
		{
			version:      version("2022.2.5", t),
			os:           "darwin",
			arch:         "amd64",
			expectedPath: "https://download.jetbrains.com/idea/ideaIU-2022.2.5.dmg",
			expectedType: domain.DMG,
		},
		{
			version:      version("2022.2.5", t),
			os:           "darwin",
			arch:         "arm64",
			expectedPath: "https://download.jetbrains.com/idea/ideaIU-2022.2.5-aarch64.dmg",
			expectedType: domain.DMG,
		},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("Checking download path for version %s, os %s and arch %s", tt.version.Value, tt.os, tt.arch)
		t.Run(name, func(t *testing.T) {
			actualPath, actualType := plugin.CalculateDownloadUrl(tt.version, tt.os, tt.arch)
			if actualPath != tt.expectedPath {
				t.Errorf("calculateDownloadPath().path = %v, want %v", actualPath, tt.expectedPath)
			}
			if actualType != tt.expectedType {
				t.Errorf("calculateDownloadPath().type = %v, want %v", actualType, tt.expectedType)
			}
		})
	}
}

func version(version string, t *testing.T) domain.Version {
	newVersion, err := domain.NewVersion(version)
	if err != nil {
		t.Errorf("Failed to create version: %s", err)
	}
	return newVersion
}
