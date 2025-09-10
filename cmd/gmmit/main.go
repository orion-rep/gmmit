package main

import (
	"flag"

	. "gitlab.com/orion-rep/gmmit/internal/pkg/common"
)

var (
	noVerifyFlag = flag.Bool("no-verify", false, "Runs the 'git commit' command with '--no-verify'.")
	generatePR   = flag.Bool("pr", false, "Generates a PR Message for changes in branch to be merged into default branch.")
)

func main() {
	PrintStartupLines()
	flag.Parse()
	LoadEnvironment()

	if *generatePR {
		RunPRGeneration()
	} else {
		RunCommitGeneration()
	}
	PrintFinalLine()
}
