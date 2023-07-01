package console

import "os"

func Info(message string) {
	println(message)
}

func Fatal(error error) {
	println(error)
	os.Exit(1)
}
