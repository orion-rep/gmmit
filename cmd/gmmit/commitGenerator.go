package main

import (
	"fmt"
	"os"
	"strings"

	. "gitlab.com/orion-rep/gmmit/internal/pkg/ai"
	. "gitlab.com/orion-rep/gmmit/internal/pkg/common"
)

var commitStandard, prompt, gitDiff, gitBranch string = "", "", "", ""

func RunCommitGeneration() {
	Info("Getting context information")
	gitDiff, gitBranch = GetCommitContext()
	commitStandard = GetEnvArg("GMMIT_COMMIT_PATTERN", "<type>[optional scope]: <description> (#<ticket-id>)")
	GenerateCommitMessage()
}

func generatePrompt(commitStandard, gitBranch, gitDiff string) string {
	return fmt.Sprintf("Create a git commit message following the \"Conventional Commits\" standard: \"%s\". The Ticket ID MUST be present on the first line, look for it on the branch name: \"%s\". Respond with the commit message only. First line can not be a generic line, must be a specific change. If there are many changes, list the rest at the end. These are the file changes to be pushed:\n%s", commitStandard, gitBranch, gitDiff)
}

func GenerateCommitMessage() {

	Info("Generating commit message")

	prompt = generatePrompt(commitStandard, gitBranch, gitDiff)

	Debug(prompt)
	res := RunPrompt(prompt)

	Info("Text Generated")
	Info("Commit Message:")
	PrintModelResponse(res)
	Info("---")

	switch option := AskConfirmation("Create a commit with this message? or Regenerate the text(r)? [y/N/r]"); option {
	case 1:
		CreateCommit(ModelResponseToString(res))
		Info("Commit created, remember to run 'git push'.")
	case 2:
		GenerateCommitMessage()
	default:
		os.Exit(0)
	}
}

func GetCommitContext() (string, string) {

	diff := string(RunCommand("git", "diff", "--staged"))

	if len(diff) <= 0 {
		Warning("Git diff returned no files")
		Warning("Add some files to the staging area and run this command again")
		os.Exit(0)
	}

	branch := strings.ReplaceAll(string(RunCommand("git", "rev-parse", "--abbrev-ref", "HEAD")), "\n", "")

	return diff, branch

}

func CreateCommit(msg string) {
	Info("Creating Commit")
	gitOptions := []string{"commit"}
	if *noVerifyFlag {
		Debug("Adding '--no-verify' option to git commit")
		gitOptions = append(gitOptions, "--no-verify")
	}
	gitOptions = append(gitOptions, "-m", msg)

	gitCommit := RunCommand("git", gitOptions...)

	Info("Git Command Log:")
	fmt.Println(string(gitCommit))
}
