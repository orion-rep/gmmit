package common

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// RunCommand executes a command-line command with the provided parameters
// and returns the output of the command as a string.
func RunCommand(command string, params ...string) string {
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

// CheckArgs ensures that the correct command-line arguments are provided
// before executing an example. It checks if the number of arguments provided
// is less than the required number of arguments. If so, it prints the usage
// message and exits the program with a non-zero status code.
func CheckArgs(arg ...string) {
	if len(os.Args) < len(arg)+1 {
		Error("Usage: %s %s", os.Args[0], strings.Join(arg, " "))
		PrintFailLine()
	}
}

// CheckIfError checks if an error occurred and handles it accordingly.
// If the error is not nil, it logs the error message and any additional context
// provided, and then exits the program with a non-zero status code.
func CheckIfError(err error, context ...string) {
	if err == nil {
		return
	}
	Error("Somenthing went wrong")
	Error("Error: %s", err.Error())
	if e, ok := err.(*json.SyntaxError); ok {
		Error("at byte offset %d", e.Offset)
	}
	if len(context) > 0 {
		for _, c := range context {
			lines := strings.Split(string(c), "\n")
			for _, line := range lines { // Using blank identifier for index if not needed
				Error(line)
			}
		}
	}
	PrintFailLine()
}

// AskConfirmation prompts the user with a message and waits for user input.
// It returns an integer value based on the user's response:
// 1 for "y" or "yes" (confirmed)
// 2 for "r" (retry)
// 0 for any other input (canceled)
func AskConfirmation(message string) int {
	Question(message)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := strings.ToLower(strings.TrimSpace(scanner.Text()))
	DeleteLastLine()

	switch input {
	case "y", "yes":
		return 1 //Confirmed
	case "r":
		return 2 //Retry
	default:
		return 0 //Canceled
	}
}

// CommandExists checks if a given command is installed on the system.
// It logs a message indicating that it's checking for the command's installation,
// attempts to find the command's path using exec.LookPath, and if the command
// is not found, it logs a warning message and prompts the user to install the command.
// If an error occurs during the process, it handles the error using CheckIfError.
func CommandExists(cmd string) {
	Debug("Checking %s installation", cmd)
	_, err := exec.LookPath(cmd)
	if err != nil {
		Warning("%s could not be found. Install it and run this command again.", cmd)
	}
	CheckIfError(err)
}

// OpenURL opens the specified URL in the default browser of the user's operating system.
// It determines the appropriate command and arguments based on the user's operating system
// and executes the command to open the URL in the default browser.
// This function is based on the solution from the following Stack Overflow question:
// https://stackoverflow.com/questions/39320371/how-start-web-server-to-open-page-in-browser-in-golang
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
