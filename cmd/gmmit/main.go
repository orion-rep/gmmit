package main

import (
	"flag"
	"fmt"

	. "gitlab.com/orion-rep/gmmit/internal/pkg/common"
)

var (
	Version string = "[unknown]"

	noVerifyFlag = flag.Bool("no-verify", false, "Runs the 'git commit' command with '--no-verify'.")
	generatePR   = flag.Bool("pr", false, "Generates a PR Message for changes in branch to be merged into default branch.")
)

func PrintHeader() {
	fmt.Println(" ╔════════════════════════════════════════════════════════════════════════════════════════╗")
	fmt.Println(" ║                                 ▗▄▄▖▄▄▄▄  ▄▄▄▄  ▄    ■                                 ║")
	fmt.Println(" ║                                ▐▌   █ █ █ █ █ █ ▄ ▗▄▟▙▄▖                               ║")
	fmt.Println(" ║                                ▐▌▝▜▌█   █ █   █ █   ▐▌                                 ║")
	fmt.Println(" ║                                ▝▚▄▞▘            █   ▐▌                                 ║")
	fmt.Println(" ║                                                     ▐▌                                 ║")
	fmt.Println(" ║                                   Version: " + fmt.Sprintf("%-s%*s", Version, 44-len(Version), "") + "║")
	fmt.Println(" ╚════════════════════════════════════════════════════════════════════════════════════════╝")
}

func main() {
	PrintHeader()
	flag.Parse()

	PrintStartLine()
	LoadEnvironment()

	if *generatePR {
		RunPRGeneration()
	} else {
		RunCommitGeneration()
	}
	PrintFinalLine()
}
