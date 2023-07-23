package shell

import (
	"os/user"
)

type DirFinder interface {
	HomeDir() (string, error)

	SoftDir() (string, error)
}

type ProdDirFinder struct {
	SoftwareDir string
}

func (receiver ProdDirFinder) HomeDir() (string, error) {
	current, err := user.Current()
	if err != nil {
		return "", err
	}
	return current.HomeDir, nil
}

func (receiver ProdDirFinder) SoftDir() (string, error) {
	return receiver.SoftwareDir, nil
}
