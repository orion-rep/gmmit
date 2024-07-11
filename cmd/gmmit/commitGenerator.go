package main

import (
	"os"
	"fmt"
	"strings"
	
	. "gitlab.com/orion-rep/gmmit/internal/pkg/common"
	. "gitlab.com/orion-rep/gmmit/internal/pkg/ai"
)

var commitStandard, prompt string = "",""

func GenerateCommitMessage() {

	Info("Getting context information")
	gitDiff, gitBranch := GetCommitContext()
    commitStandard = GetEnvArg("GMMIT_COMMIT_PATTERN", "<type>[optional scope]: <description> (#<ticket-id>)")
	
	Info("Generating commit message")

	prompt = fmt.Sprintf("Create a git commit message following the \"Conventional Commits\" standard: \"%s\". The Ticket ID MUST be present on the first line, look for it on the branch name: \"%s\". Respond with the commit message only. First line can not be a generic line, must be a specific change. If there are many changes, list the rest at the end. These are the file changes to be pushed:\n%s",
				commitStandard, gitBranch, gitDiff)
	
	Debug(prompt)

	res := RunPrompt(prompt)

	PrintModelResponse(res)

	if (AskConfirmation("Create a commit with this message? [y/N]") == 1) {
		CreateCommit(ModelResponseToString(res))
	} 
}

func GetCommitContext()(string, string) {
	
	diff := string(RunCommand("git","diff","--staged"))

	if len(diff) <= 0 {
		Warning("Git diff returned no files.")
		Warning("Add some files to the staging area and run this command again.")
		os.Exit(0)
	}
	
	branch := strings.ReplaceAll(string(RunCommand("git", "rev-parse", "--abbrev-ref", "HEAD")), "\n", "")

	return diff, branch
    
}

func CreateCommit(msg string) {
	Info("Creating Commit...")
	gitCommit := RunCommand("git","commit","-m",msg)

	Info("Git Command Log:")
	fmt.Println(string(gitCommit))
}