package common

import (
	"fmt"
	"os"
	"time"
)

var dateFormat string = "2006-01-02 - 15:04:05"

var OsExit = os.Exit

func logLine(logText string) {
	fmt.Printf("%s ║ %s\n", time.Now().Format(dateFormat), logText)
}

// Debug should be used to display debugging information
func Debug(format string, args ...interface{}) {
	debug := GetEnvArg("GMMIT_DEBUG", "false")
	if debug == "true" {
		logLine(fmt.Sprintf("\x1b[35;1m[D] %s\x1b[0m", fmt.Sprintf(format, args...)))
	}
}

// Info should be used to describe the example commands that are about to run.
func Info(format string, args ...interface{}) {
	logLine(fmt.Sprintf("\x1b[34;1m[I] %s\x1b[0m", fmt.Sprintf(format, args...)))
}

func InfoH(format string, args ...interface{}) {
	logLine(fmt.Sprintf("\x1b[34;1m[I]\x1b[0m %s", fmt.Sprintf(format, args...)))
}

// Question should be used to display a message that the user will have to interact with.
func Question(format string, args ...interface{}) {
	logLine(fmt.Sprintf("\x1b[36;1m[I] %s\x1b[0m ", fmt.Sprintf(format, args...)))
	fmt.Printf("                      ╚══ Answer ═> ")
}

// Warning should be used to display a message about an unexpected situation
func Warning(format string, args ...interface{}) {
	logLine(fmt.Sprintf("\x1b[33;1m[W] %s\x1b[0m", fmt.Sprintf(format, args...)))
}

// Error should be used to display an error that prevents the tool from continuing
func Error(format string, args ...interface{}) {
	logLine(fmt.Sprintf("\x1b[31;1m[E] %s\x1b[0m", fmt.Sprintf(format, args...)))
}

func DeleteLastLine() {
	fmt.Fprint(os.Stdout, "\033[1A\033[2K")
}

func PrintStartLine() {
	fmt.Printf("%s ╔══[ \x1b[34;1mProcess\x1b[0m ]\n", time.Now().Format(dateFormat))
}

func PrintFinalLine() {
	fmt.Printf("%s ╚══[ \x1b[34;1mDone\x1b[0m ]\n", time.Now().Format(dateFormat))
	OsExit(0)
}

func PrintFailLine() {
	fmt.Printf("%s ╚══[ \x1b[31;1mFail\x1b[0m ]\n", time.Now().Format(dateFormat))
	OsExit(1)
}
