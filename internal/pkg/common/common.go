package common

import (
	"os"
	"os/exec"
	"strings"
	"bytes"

)

func RunCommand(command string, params ...string) (string) {
	Debug("Running command: \"%s\", with params: %s",command, params)
	cmd := exec.Command(command, params...)
	
	var out bytes.Buffer
	var stderr bytes.Buffer
	
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	Debug("Command output: %s", out.String())
    CheckIfError(err, stderr.String())
	
	return out.String()
}

// CheckArgs should be used to ensure the right command line arguments are
// passed before executing an example.
func CheckArgs(arg ...string) {
	if len(os.Args) < len(arg)+1 {
		Error("Usage: %s %s", os.Args[0], strings.Join(arg, " "))
		os.Exit(1)
	}
}

// CheckIfError should be used to naively panics if an error is not nil.
func CheckIfError(err error, context ...string) {
	if err == nil {
		return
	}
	Error(err.Error())
	if len(context) > 0 {
		for _, c := range context { 
			Error(c)
		} 
	}
	os.Exit(1)
}
