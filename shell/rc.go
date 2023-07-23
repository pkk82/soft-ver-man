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

func exportVariable(name, value string) string {
	return fmt.Sprintf("export %v=\"%v\"", name, value)
}

func exportRefPathVariable(name, refVar, path string) string {
	return fmt.Sprintf("export %v=\"$%v/%v\"", name, refVar, path)
}

func exportHomeVariable(name, refVar, path string) string {
	return fmt.Sprintf("export %v_HOME=\"$%v/%v\"", name, refVar, path)
}

func exportHomeMajorVersionVariable(name string, version int, refVar, path string) string {
	return fmt.Sprintf("export %v_%v_HOME=\"$%v/%v\"", name, version, refVar, path)
}
