package common

import (
	"fmt"
)

// Info should be used to describe the example commands that are about to run.
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

// Warning should be used to display a warning
func Warning(format string, args ...interface{}) {
	fmt.Printf("\x1b[33;1m[WARN] %s\x1b[0m\n", fmt.Sprintf(format, args...))
}