package main

import (
	"os"
	"fmt"
	"encoding/json"
	
	. "gitlab.com/orion-rep/gmmit/internal/pkg/common"
	. "gitlab.com/orion-rep/gmmit/internal/pkg/ai"
	"github.com/atotto/clipboard"
)

var prPrompt, gitPRDiff, gitDefaultBranch, gitPRBranch, repositoryName, repositoryProvider string = "","","","","",""

func RunPRGeneration() {
	Info("Getting context information")
	getPRContext()
    generatePRMessage()
}

func getPRContext() {
	
	repositoryName, repositoryProvider = GetRepositoryData()
	gitDefaultBranch = GetDefaultBranch()
	gitPRBranch = GetCurrentBranch()
	gitPRDiff = CalculateDiffBetweenBranches(gitDefaultBranch, gitPRBranch)
	
	if len(gitPRDiff) <= 0 {
		Warning("Git diff returned no files")
		Warning("Add some files to the staging area and run this command again")
		os.Exit(0)
	}
}

func generatePRMessage() {

	Info("Generating PR message")

	prPrompt = fmt.Sprintf("Create a Pull Request message with following sections: 'What changed?', 'Why/Context', 'How to test it?'. The title line should follow the 'Conventional Commits' standard. The Ticket ID MUST be present on the PR title line, look for it on the branch name: \"%s\". Respond with the pr message only. Title line can not be a generic line, must be a specific change. If there are many changes, list the rest at the end. Answer must be a valid json with no '`' characters, following this template: {title:'',description:''}.These are the changes to be merged:\n%s",
		gitPRBranch, gitPRDiff)
	
	Debug(prPrompt)
	res := RunPrompt(prPrompt)

	Debug("Model Response:\n%s",ModelResponseToString(res))

	var response map[string]string
	err := json.Unmarshal([]byte(ModelResponseToString(res)), &response)
	prTitle := response["title"]
	prDescription := response["description"]
	CheckIfError(err)

	Info("Text Generated")
	Info("PR Title:")
	fmt.Println(prTitle)
	Info("PR Description:")
	fmt.Println(prDescription)
	Info("---")

	if repositoryProvider == "Generic" {
		Debug("Repository provider not supported, PR creation dissabled")
		confirmCopyClipboard(prDescription)
		return
	}
	
	confirmPRCreation(prTitle, prDescription, repositoryProvider)
}

func confirmPRCreation(title, description, repoProvider string) {
	switch option := AskConfirmation("Do you want to create the PR(y)? or Regenerate the text(r)? [y/N/r]"); option {
		case 1:
			prURL := ""
			switch repoProvider {
				case GIT_PROVIDER_BITBUCKET:
					prURL = createPROnBitbucket(title, description, gitPRBranch, repositoryName)
				case GIT_PROVIDER_GITHUB:
					prURL = createPROnGithub(title, description, gitPRBranch, gitDefaultBranch, repositoryName)
				default:
					Error("Unexpected unknown repository provider: %s", repoProvider)
					os.Exit(1)
			}
			Info("PR created! You're good to go")
			OpenURL(prURL)
		case 2:
			generatePRMessage()
		default:
			confirmCopyClipboard(description)
	}
}

func confirmCopyClipboard(description string) {
	switch option := AskConfirmation("Copy this PR Description to your clipboard(y)? or Regenerate the text(r)? [y/N/r]"); option {
		case 1:
			clipboard.WriteAll(description)
			Info("PR description copied! You're good to go")
		case 2:
			generatePRMessage()
		default:
			os.Exit(0)
	}
}


func createPROnBitbucket(title string, message string, sourceBranch string, repo string) (string){
	
	url := "https://api.bitbucket.org/2.0/repositories/" + repo + "/pullrequests"
	
	payload := map[string]interface{}{
		"title":title, 
		"source": map[string]interface{}{
			"branch": map[string]string{
				"name": sourceBranch,
			},
		},
		"description": message,
	}
	
	resp, status, err := CallPost(url, payload ,GetEnvArg("GMMIT_BB_USER"), GetEnvArg("GMMIT_BB_PASS"))
	CheckIfError(err)

	response, err := ResponseJsonParser(resp)
	CheckIfError(err)
	Debug("Response: %s", response)

	if status != 201 {
		Error("PR creation failed with the following error message:")
		errorResp := response["error"].(map[string]interface{})
		Error(fmt.Sprint(errorResp["message"]))
		os.Exit(1)
	}

	newPRURL := fmt.Sprint(response["links"].(map[string]interface{})["html"].(map[string]interface{})["href"])
	Info("PR URL: %s", newPRURL)

	return newPRURL
}

func createPROnGithub(title string, message string, sourceBranch string, baseBranch string, repo string) (string){
	
	url := "https://api.github.com/repos/" + repo + "/pulls"
	
	payload := map[string]interface{}{
		"title":title, 
		"body": message,
		"head": sourceBranch,
		"base": baseBranch,
	}
	
	resp, status, err := CallPost(url, payload ,GetEnvArg("GMMIT_GH_USER"), GetEnvArg("GMMIT_GH_PASS"))
	CheckIfError(err)

	response, err := ResponseJsonParser(resp)
	CheckIfError(err)
	Debug("Response: %s", response)

	if status != 201 {
		Error("PR creation failed with the following error message:")
		Error(fmt.Sprint(response["message"]))
		if _, ok := response["errors"]; ok {
			Error(fmt.Sprint(response["errors"].([]interface{})[0].(map[string]interface{})["message"]))
		}
		os.Exit(1)
	}

	newPRURL := fmt.Sprint(response["html_url"])
	Info("PR URL: %s", newPRURL)

	return newPRURL
}
