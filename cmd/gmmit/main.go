package main

import (
	"flag"
	"fmt"

	"gitlab.com/orion-rep/gmmit/internal/pkg/common"
)

var (
	Version string = "[unknown]"

	noVerifyFlag  = flag.Bool("no-verify", false, "Runs the 'git commit' command with '--no-verify'.")
	generatePR    = flag.Bool("pr", false, "Generates a PR Message for changes in branch to be merged into default branch.")
	runCommitPush = flag.Bool("pu", false, "Runs 'git push' after creating the commit.")
	autoConfirm   = flag.Bool("y", false, "Auto-confirm all prompts, run in non-interactive mode.")
	commitType    = flag.String("type", "", "Force the commit type (feat, fix, docs, style, refactor, perf, test, build, ci, chore, revert).")
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
	common.AutoConfirm = *autoConfirm

	common.PrintStartLine()
	common.LoadEnvironment()

	if *generatePR {
		RunPRGeneration()
	} else {
		RunCommitGeneration()
	}
	common.PrintFinalLine()
}
