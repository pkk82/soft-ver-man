package shell

import (
	"os/user"
)

type homeDirFinder interface {
	HomeDir() (string, error)
}

type osHomeDirFinder struct {
}

func (receiver osHomeDirFinder) HomeDir() (string, error) {
	current, err := user.Current()
	if err != nil {
		return "", err
	}
	return current.HomeDir, nil
}
