package testutil

import (
	"github.com/pkk82/soft-ver-man/domain"
	"testing"
)

func AsVersion(version string, t *testing.T) domain.Version {
	v, err := domain.NewVersion(version)
	if err != nil {
		t.Errorf("Failed to create version: %s", err)
	}
	return v
}
