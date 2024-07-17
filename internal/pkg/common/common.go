package common

import (
	"bufio"
	"os"
	"os/exec"
	"strings"
	"bytes"
	"log"

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

func AskConfirmation(message string) (int) {

	Question(message)

	reader := bufio.NewReader(os.Stdin)
	confirmation, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
	confirmation = strings.ToLower(strings.TrimSpace(confirmation))
	if confirmation == "y" || confirmation == "yes" {
		return 1
	} else if confirmation == "r" {
		return 2
	}
	return 0
}

func CommandExists(cmd string) {
	Debug("Checking %s installation", cmd)
	_, err := exec.LookPath(cmd)
	if err != nil {
		Warning("%s could not be found. Install it and run this command again.", cmd)
	} 
	CheckIfError(err)
}