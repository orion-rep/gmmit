package common

import (
	"fmt"
	"os"
	"strings"
)

func logLine(logText string) {
	lineTopLength := 101
	if len(logText) < lineTopLength {
		fmt.Printf("%-s%*s║\n", logText, 103-len(logText), "")
	} else {
		lineBreakIndex := strings.LastIndex(logText, " ")
		fmt.Printf("%-s%*s║\n", logText[:lineBreakIndex], 99-len(logText[:lineBreakIndex]), "")
		fmt.Printf(" ║     %-s%*s║\n", logText[lineBreakIndex:], 87-len(logText[lineBreakIndex:]), "")
	}
}

// Debug should be used to display debugging information
func Debug(format string, args ...interface{}) {
	debug := GetEnvArg("GMMIT_DEBUG", "false")
	if debug == "true" {
		logLine(fmt.Sprintf(" ║ \x1b[35;1m[D] %s\x1b[0m", fmt.Sprintf(format, args...)))
	}
}

// Info should be used to describe the example commands that are about to run.
func Info(format string, args ...interface{}) {
	logLine(fmt.Sprintf(" ║ \x1b[34;1m[I] %s\x1b[0m", fmt.Sprintf(format, args...)))
}

func InfoH(format string, args ...interface{}) {
	logLine(fmt.Sprintf(" ║ \x1b[34;1m[I]\x1b[0m %s", fmt.Sprintf(format, args...)))
}

// Question should be used to display a message that the user will have to interact with.
func Question(format string, args ...interface{}) {
	logLine(fmt.Sprintf(" ║ \x1b[36;1m[I] %s\x1b[0m ", fmt.Sprintf(format, args...)))
	fmt.Printf(" ╚══ Answer ═> ")
}

// Warning should be used to display a message about an unexpected situation
func Warning(format string, args ...interface{}) {
	logLine(fmt.Sprintf(" ║ \x1b[33;1m[W] %s\x1b[0m", fmt.Sprintf(format, args...)))
}

// Error should be used to display an error that prevents the tool from continuing
func Error(format string, args ...interface{}) {
	logLine(fmt.Sprintf(" ║ \x1b[31;1m[E] %s\x1b[0m", fmt.Sprintf(format, args...)))
}

func DeleteLastLine() {
	fmt.Fprint(os.Stdout, "\033[1A\033[2K")
}

func PrintStartLine() {
	fmt.Println("    ╔═══════════╗")
	fmt.Println(" ╔══╝  \x1b[34;1mProcess\x1b[0m  ╚═════════════════════════════════════════════════════════════════════════╗")
}

func PrintFinalLine() {
	fmt.Println(" ╚═════════════════════════════════════════════════════════════════════════════╗  \x1b[34;1mDone\x1b[0m  ╔═╝")
	fmt.Println("                                                                               ╚════════╝")
	os.Exit(0)
}
func PrintFailLine() {
	fmt.Println(" ╚═══════════════════════════════════════════════════════════════════════════════╗  \x1b[31;1mFail\x1b[0m  ║")
	fmt.Println("                                                                                 ╚════════╝")
	os.Exit(1)
}
