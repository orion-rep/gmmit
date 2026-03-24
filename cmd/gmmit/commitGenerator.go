package main

import (
	"fmt"
	"strings"

	gemini "gitlab.com/orion-rep/gmmit/internal/pkg/ai"
	"gitlab.com/orion-rep/gmmit/internal/pkg/common"
)

var commitStandard, prompt, gitDiff, gitBranch string = "", "", "", ""

var runPromptFn = gemini.RunPrompt

const systemPrompt = `You're a Software developer with 15 years of experience working for Fortune 100 companies.
You're detail oriented, and put good effort on writing documentation and comments for your co-workers.`

func RunCommitGeneration() {
	common.Info("Getting context information")
	gitDiff, gitBranch = GetCommitContext()
	commitStandard = common.GetEnvArg("GMMIT_COMMIT_PATTERN", "<type>[optional scope]: <description> (#<ticket-id>)")
	GenerateCommitMessage()
}

func generatePrompt(commitStandard, gitBranch, gitDiff string) string {
	return fmt.Sprintf(`System: %s,

	User: Create a git commit message following the \"Conventional Commits\" standard: \"%s\".
	The Ticket ID MUST be present on the first line, look for it on the branch name: \"%s\".
	Respond with the commit message only. First line can not be a generic line, must be a specific change.
	If there are many changes, list the rest at the end.

	These are the file changes to be pushed:\n%s`, systemPrompt, commitStandard, gitBranch, gitDiff)
}

func GenerateCommitMessage() {

	common.Info("Generating commit message")

	prompt = generatePrompt(commitStandard, gitBranch, gitDiff)

	common.Debug(prompt)
	res := runPromptFn(prompt)

	common.Info("Text Generated")
	common.Info("Commit Message:")
	gemini.PrintModelResponse(res)
	common.Info("---")

	switch option := common.AskConfirmation("Do you want to use this commit message (y) or regenerate it (r)? [y/N/r]"); option {
	case 1:
		CreateCommit(gemini.ModelResponseToString(res))
		common.Info("Commit created.")
		if *runCommitPush {
			pushCommit()
		} else {
			common.Warning("Remember to run 'git push' !")
		}
	case 2:
		GenerateCommitMessage()
	default:
		common.PrintFinalLine()
	}
}

func GetCommitContext() (string, string) {

	diff := string(common.RunCommand("git", "diff", "--staged"))

	if len(diff) <= 0 {
		common.Warning("Git diff returned no files")
		common.Warning("Add some files to the staging area and run this command again")
		common.PrintFinalLine()
	}

	branch := strings.ReplaceAll(string(common.RunCommand("git", "rev-parse", "--abbrev-ref", "HEAD")), "\n", "")

	return diff, branch

}

func CreateCommit(msg string) {
	common.Info("Creating Commit")
	gitOptions := []string{"commit"}
	if *noVerifyFlag {
		common.Debug("Adding '--no-verify' option to git commit")
		gitOptions = append(gitOptions, "--no-verify")
	}
	gitOptions = append(gitOptions, "-m", msg)

	gitCommit := common.RunCommand("git", gitOptions...)

	common.Info("Git Command Log:")
	lines := strings.Split(string(gitCommit), "\n")
	for _, line := range lines {
		common.InfoH(line)
	}
}

func pushCommit() {
	common.Info("Pushing Commit")
	gitOptions := []string{"push"}
	common.RunCommand("git", gitOptions...)

	common.Info("Changes pushed to remote repo.")
}
