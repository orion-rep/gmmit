package main

import (
	"fmt"

	. "gitlab.com/orion-rep/gmmit/internal/pkg/common"
	. "gitlab.com/orion-rep/gmmit/internal/pkg/ai"
)

var commitStandard, prompt string = "",""

func main() {

	LoadEnvironment()
	gitDiff, gitBranch := GetCommitContext()
	
	Info("Generating commit message.")

    commitStandard = GetEnvArg("GMMIT_COMMIT_PATTERN", "<type>[optional scope]: <description> (#<ticket-id>)")
	prompt = fmt.Sprintf("Create a git commit message following the \"Conventional Commits\" standard: \"%s\". The Ticket ID MUST be present on the first line, look for it on the branch name: \"%s\". Respond with the commit message only. First line can not be a generic line, must be a specific change. If there are many changes, list the rest at the end. These are the file changes to be pushed:\n%s",
				commitStandard, gitBranch, gitDiff)
	Debug(prompt)


	res := RunPrompt(prompt)
	PrintModelResponse(res)

	if AskConfirmation("Create a commit with this message? [y/N]") {
		Info("Creating Commit...")
		gitCommit := RunCommand("git","commit","-m",ModelResponseToString(res))
    	Info("Git Command Log:")
		fmt.Println(string(gitCommit))
	} 
}
