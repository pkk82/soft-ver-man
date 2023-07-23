package shell

import (
	"os/user"
)

type dirFinder interface {
	HomeDir() (string, error)

	SoftDir() (string, error)
}

type prodDirFinder struct {
	softDir string
}

func (receiver prodDirFinder) HomeDir() (string, error) {
	current, err := user.Current()
	if err != nil {
		return "", err
	}
	return current.HomeDir, nil
}

func (receiver prodDirFinder) SoftDir() (string, error) {
	return receiver.softDir, nil
}
