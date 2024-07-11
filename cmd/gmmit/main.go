package main

import (
	"flag"

	. "gitlab.com/orion-rep/gmmit/internal/pkg/common"
)

var(
    noVerifyFlag = flag.Bool("no-verify", false, "Runs the 'git commit' command with '--no-verify'.")
)

func main() {
	flag.Parse()
	LoadEnvironment()

	RunCommitGeneration()
}
