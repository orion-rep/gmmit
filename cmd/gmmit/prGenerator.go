package main

import (
	"os"
	"fmt"
	"encoding/json"
	
	. "gitlab.com/orion-rep/gmmit/internal/pkg/common"
	. "gitlab.com/orion-rep/gmmit/internal/pkg/ai"
	"github.com/atotto/clipboard"
)

var prPrompt, gitPRDiff, gitDefaultBranch, gitPRBranch, repositoryName string = "","","","",""

func RunPRGeneration() {
	Info("Getting context information")
	getPRContext()
    generatePRMessage()
}

func getPRContext() {
	
	repositoryName = GetRepositoryName()
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

	switch option := AskConfirmation("Do you want to create the PR? [y/N/r]"); option {
		case 1:
			prURL := createPRBitbucket(prTitle,gitPRBranch, prDescription, repositoryName)
			Info("PR created! You're good to go")
			OpenURL(prURL)
		case 2:
			generatePRMessage()
		default:
			switch option := AskConfirmation("Copy this PR Message to your clipboard? [y/N/r]"); option {
				case 1:
					clipboard.WriteAll(ModelResponseToString(res))
					Info("PR Message copied! You're good to go")
				case 2:
					generatePRMessage()
				default:
					os.Exit(0)
			}
	}
}

func createPRBitbucket(title string, sourceBranch string, message string, repo string) (string){
	
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
