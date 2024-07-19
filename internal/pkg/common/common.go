package common

import (
	"bufio"
	"os"
	"os/exec"
	"strings"
	"bytes"
	"runtime"
)

// RunCommand executes a command-line command with the provided parameters
// and returns the output of the command as a string.
func RunCommand(command string, params ...string) (string) {
    // Log the command and parameters being executed
    Debug("Running command: \"%s\", with params: %s", command, params)

    // Create a new command object with the provided command and parameters
    cmd := exec.Command(command, params...)

    // Create buffers to store the standard output and standard error
    var out bytes.Buffer
    var stderr bytes.Buffer

    // Assign the buffers to the command's output and error streams
    cmd.Stdout = &out
    cmd.Stderr = &stderr

    // Execute the command
    err := cmd.Run()

    // Log the command output
    Debug("Command output: %s", out.String())

    // Check if an error occurred during the command execution
    // and handle it accordingly
    CheckIfError(err, stderr.String())

    // Return the command output as a string
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

    scanner := bufio.NewScanner(os.Stdin)
    scanner.Scan()
    input := strings.ToLower(strings.TrimSpace(scanner.Text()))

    switch input {
    case "y", "yes":
        return 1 //Confirmed
    case "r":
        return 2 //Retry
    default:
        return 0 //Canceled
    }
}

func CommandExists(cmd string) {
	Debug("Checking %s installation", cmd)
	_, err := exec.LookPath(cmd)
	if err != nil {
		Warning("%s could not be found. Install it and run this command again.", cmd)
	} 
	CheckIfError(err)
}

// https://stackoverflow.com/questions/39320371/how-start-web-server-to-open-page-in-browser-in-golang
// open opens the specified URL in the default browser of the user.
func OpenURL(url string) error {
    var cmd string
    var args []string

    switch runtime.GOOS {
    case "windows":
        cmd = "cmd"
        args = []string{"/c", "start"}
    case "darwin":
        cmd = "open"
    default: // "linux", "freebsd", "openbsd", "netbsd"
        cmd = "xdg-open"
    }
    args = append(args, url)
    return exec.Command(cmd, args...).Start()
}
