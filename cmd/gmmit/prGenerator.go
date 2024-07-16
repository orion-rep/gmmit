package main

import (
	"os"
	"fmt"
	"strings"
	
	. "gitlab.com/orion-rep/gmmit/internal/pkg/common"
	//. "gitlab.com/orion-rep/gmmit/internal/pkg/ai"
)

var prPrompt, gitPRDiff, gitPRBranch string = "","",""

func RunPRGeneration() {
	Info("Getting context information")
	gitDiff, gitBranch = GetPRContext()
    GeneratePRMessage()
}


func GetPRContext()(string, string) {
	
	defaultBranch := strings.ReplaceAll(string(RunCommand("git", "rev-parse", "--abbrev-ref", "origin/HEAD")), "\n", "")

	Debug("Checking changes against branch '%s'", defaultBranch)
	diff := string(RunCommand("git","diff",fmt.Sprintf("%s...", defaultBranch)))

	if len(diff) <= 0 {
		Warning("Git diff returned no files")
		Warning("Add some files to the staging area and run this command again")
		os.Exit(0)
	}
	
	branch := strings.ReplaceAll(string(RunCommand("git", "rev-parse", "--abbrev-ref", "HEAD")), "\n", "")

	return diff, branch
    
}


func GeneratePRMessage() {

	Info("Generating PR message")

	prPrompt = fmt.Sprintf("Create a Pull Request message with following sections: 'What changed?', 'Why/Context', 'How to test it?'. The Ticket ID MUST be present on the PR title line, look for it on the branch name: \"%s\". Respond with the pr message only. Title line can not be a generic line, must be a specific change. If there are many changes, list the rest at the end. These are the changes to be merged:\n%s",
		gitPRBranch, gitPRDiff)
	
	Debug(prPrompt)
	//res := RunPRPrompt(prPrompt)

	//PrintModelResponse(res)

	switch option := AskConfirmation("Copy this PR Message to your clipboard? [y/N/r]"); option {
		
		default:
			os.Exit(0)
	}
}
