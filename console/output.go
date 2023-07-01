package console

import "os"

func Info(message string) {
	println(message)
}

func Fatal(error error) {
	println(error.Error())
	os.Exit(1)
}

func Error(error error) {
	println(error.Error())
}
