package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	gemini "gitlab.com/orion-rep/gmmit/internal/pkg/ai"
	"gitlab.com/orion-rep/gmmit/internal/pkg/common"
)

var prPrompt, gitPRDiff, gitDefaultBranch, gitPRBranch, repositoryName, repositoryProvider string = "", "", "", "", "", ""

var callPostFn = common.CallPost

func RunPRGeneration() {
	common.Info("Getting context information")
	getPRContext()
	generatePRMessage()
}

func getPRContext() {

	repositoryName, repositoryProvider = common.GetRepositoryData()
	gitDefaultBranch = common.GetDefaultBranch()
	gitPRBranch = common.GetCurrentBranch()
	gitPRDiff = common.CalculateDiffBetweenBranches(gitDefaultBranch, gitPRBranch)

	if len(gitPRDiff) <= 0 {
		common.Warning("Git diff returned no files")
		common.Warning("Push your commits to the remote branch and run this command again.")
		common.PrintFinalLine()
	}
}

func generatePRPrompt(branch, diff, hint string) string {
	hintConstraint := ""
	if hint != "" {
		hintConstraint = fmt.Sprintf(" Focus on the following when writing the PR message: \"%s\".", hint)
	}
	return fmt.Sprintf("Create a Pull Request message with following sections: 'What changed?', 'Why/Context', 'How to test it?'. The title line should follow the 'Conventional Commits' standard. The Ticket ID MUST be present on the PR title line, look for it on the branch name: \"%s\". Respond with the pr message only. Title line can not be a generic line, must be a specific change. If there are many changes, list the rest at the end.%s Answer must be a valid json with no '`' characters, following this template: {title:'',description:''}.These are the changes to be merged:\n%s",
		branch, hintConstraint, diff)
}

func generatePRMessage() {

	common.Info("Generating PR message")

	prPrompt = generatePRPrompt(gitPRBranch, gitPRDiff, *hintFlag)

	common.Debug(prPrompt)
	res := runPromptFn(prPrompt)

	stringRes := gemini.ModelResponseToString(res)
	common.Debug("Model Response:\n%s", stringRes)

	var response map[string]string
	err := json.Unmarshal([]byte(stringRes), &response)
	common.CheckIfError(err)
	prTitle := response["title"]
	prDescription := response["description"]

	common.Info("Text Generated")
	common.Info("PR Title:")
	common.InfoH(prTitle)
	common.Info("PR Description:")
	prDescriptionLines := strings.Split(string(prDescription), "\n")
	for _, line := range prDescriptionLines {
		common.InfoH(line)
	}
	common.Info("---")

	if repositoryProvider == "Generic" {
		common.Debug("Repository provider not supported, PR creation dissabled")
		confirmCopyClipboard(prDescription)
		return
	}

	confirmPRCreation(prTitle, prDescription, repositoryProvider)
}

func confirmPRCreation(title, description, repoProvider string) {
	switch option := common.AskConfirmation("Do you want to create the PR with this description(y) or regenerate it (r)? [y/N/r]"); option {
	case 1:
		prURL := ""
		switch repoProvider {
		case common.GIT_PROVIDER_BITBUCKET:
			prURL = createPROnBitbucket(title, description, gitPRBranch, repositoryName)
		case common.GIT_PROVIDER_GITHUB:
			prURL = createPROnGithub(title, description, gitPRBranch, gitDefaultBranch, repositoryName)
		default:
			common.Error("Unexpected unknown repository provider: %s", repoProvider)
			common.PrintFailLine()
		}
		common.Info("PR created! You're good to go")
		err := common.OpenURL(prURL)
		common.CheckIfError(err)
	case 2:
		generatePRMessage()
	default:
		confirmCopyClipboard(description)
	}
}

func confirmCopyClipboard(description string) {
	switch option := common.AskConfirmation("Do you want to copy this PR description to your clipboard(y) or regenerate the text(r)? [y/N/r]"); option {
	case 1:
		err := clipboard.WriteAll(description)
		common.CheckIfError(err)
		common.Info("PR description copied! You're good to go")
	case 2:
		generatePRMessage()
	default:
		common.PrintFinalLine()
	}
}

func createPROnBitbucket(title string, message string, sourceBranch string, repo string) string {

	url := "https://api.bitbucket.org/2.0/repositories/" + repo + "/pullrequests"

	payload := map[string]interface{}{
		"title": title,
		"source": map[string]interface{}{
			"branch": map[string]string{
				"name": sourceBranch,
			},
		},
		"description": message,
	}

	resp, status, err := callPostFn(url, payload, common.GetEnvArg("GMMIT_BB_USER"), common.GetEnvArg("GMMIT_BB_PASS"))
	common.CheckIfError(err)

	response, err := common.ResponseJsonParser(resp)
	common.CheckIfError(err)
	common.Debug("Response: %s", response)

	if status != 201 {
		common.Error("PR creation failed with the following error message:")
		errorResp := response["error"].(map[string]interface{})
		common.Error(fmt.Sprint(errorResp["message"]))
		common.PrintFailLine()
	}

	newPRURL := fmt.Sprint(response["links"].(map[string]interface{})["html"].(map[string]interface{})["href"])
	common.Info("PR URL: %s", newPRURL)

	return newPRURL
}

func createPROnGithub(title string, message string, sourceBranch string, baseBranch string, repo string) string {

	url := "https://api.github.com/repos/" + repo + "/pulls"

	payload := map[string]interface{}{
		"title": title,
		"body":  message,
		"head":  sourceBranch,
		"base":  baseBranch,
	}

	resp, status, err := callPostFn(url, payload, common.GetEnvArg("GMMIT_GH_USER"), common.GetEnvArg("GMMIT_GH_PASS"))
	common.CheckIfError(err)

	response, err := common.ResponseJsonParser(resp)
	common.CheckIfError(err)
	common.Debug("Response: %s", response)

	if status != 201 {
		common.Error("PR creation failed with the following error message:")
		common.Error(fmt.Sprint(response["message"]))
		if _, ok := response["errors"]; ok {
			common.Error(fmt.Sprint(response["errors"].([]interface{})[0].(map[string]interface{})["message"]))
		}
		common.PrintFailLine()
	}

	newPRURL := fmt.Sprint(response["html_url"])
	common.Info("PR URL: %s", newPRURL)

	return newPRURL
}
