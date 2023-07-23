package shell

import (
	"fmt"
	"github.com/pkk82/soft-ver-man/config"
)

func bashToLoad(fileName string) string {
	return fmt.Sprintf("[[ -s \"$HOME/%v/%v\" ]] && source \"$HOME/%v/%v\"",
		config.HomeConfigDir, fileName, config.HomeConfigDir, fileName)
}

func makeRcName(name string) string {
	return fmt.Sprintf(".%vrc", name)
}
