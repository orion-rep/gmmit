package common

import (
	"fmt"
)

// Debug should be used to display debugging information
func Debug(format string, args ...interface{}) {
	debug := GetEnvArg("GMMIT_DEBUG", "false")
	if debug == "true" {
		fmt.Printf("\x1b[35;1m[DEBG] %s\x1b[0m\n", fmt.Sprintf(format, args...))
	}
}
// Info should be used to describe the example commands that are about to run.
func Info(format string, args ...interface{}) {
	fmt.Printf("\x1b[34;1m[INFO] %s\x1b[0m\n", fmt.Sprintf(format, args...))
}
// Question should be used to display a message that the user will have to interact with.
func Question(format string, args ...interface{}) {
	fmt.Printf("\x1b[36;1m[INFO] %s\x1b[0m ", fmt.Sprintf(format, args...))
}
// Warning should be used to display a message about an unexpected situation
func Warning(format string, args ...interface{}) {
	fmt.Printf("\x1b[33;1m[WARN] %s\x1b[0m\n", fmt.Sprintf(format, args...))
}

// Error should be used to display an error that prevents the tool from continuing
func Error(format string, args ...interface{}) {
	fmt.Printf("\x1b[31;1m[ERRO] %s\x1b[0m\n", fmt.Sprintf(format, args...))
}