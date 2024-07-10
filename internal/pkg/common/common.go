package common

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

)

func RunCommand(command string, params ...string) ([]byte) {
	Debug("Running command: \"%s\", with params: %s",command, params)
	output, err := exec.Command(command, params...).Output()
	Debug("Command output: %s", output)
    CheckIfError(err)

	return output
}

// CheckArgs should be used to ensure the right command line arguments are
// passed before executing an example.
func CheckArgs(arg ...string) {
	if len(os.Args) < len(arg)+1 {
		Warning("Usage: %s %s", os.Args[0], strings.Join(arg, " "))
		os.Exit(1)
	}
}

// CheckIfError should be used to naively panics if an error is not nil.
func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}
