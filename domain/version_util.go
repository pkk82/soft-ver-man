package domain

import "testing"

func Ver(version string, t *testing.T) Version {
	v, err := NewVersion(version)
	if err != nil {
		t.Errorf("Failed to create version: %s", err)
	}
	return v
}
