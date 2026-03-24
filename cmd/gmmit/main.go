package main

import (
	"flag"
	"fmt"
	"os"

	"gitlab.com/orion-rep/gmmit/internal/pkg/common"
)

var (
	Version string = "[unknown]"

	noVerifyFlag    = flag.Bool("no-verify", false, "Runs the 'git commit' command with '--no-verify'.")
	generatePR      = flag.Bool("pr", false, "Generates a PR Message for changes in branch to be merged into default branch.")
	runCommitPush   = flag.Bool("pu", false, "Runs 'git push' after creating the commit.")
	autoConfirm     = flag.Bool("y", false, "Auto-confirm all prompts, run in non-interactive mode.")
	commitType      = flag.String("type", "", "Force the commit type (feat, fix, docs, style, refactor, perf, test, build, ci, chore, revert).")
	hintFlag        = flag.String("hint", "", "Free-text hint for the AI model about what changes to focus on.")
	addAllFlag      = flag.Bool("add-all", false, "Stage all modified files before committing (like 'git add .').")
	addAllShortFlag = flag.Bool("a", false, "Stage all modified files before committing (like 'git add .').")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: gmmit [options]\n\nOptions:\n")
		flag.VisitAll(func(f *flag.Flag) {
			if f.Name == "a" {
				return // shown together with --add-all
			}
			var prefix string
			if f.Name == "add-all" {
				prefix = "-a, --add-all"
			} else if len(f.Name) == 1 {
				prefix = "-" + f.Name
			} else {
				prefix = "--" + f.Name
			}
			fmt.Fprintf(os.Stderr, "  %-20s%s\n", prefix, f.Usage)
		})
	}
}

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
