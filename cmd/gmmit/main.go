package main

import (
	"bufio"
	"context"
	"os"
	"fmt"
	"strings"
	"log"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"

	. "gitlab.com/orion-rep/gmmit/internal/pkg/common"
)

var commitStandard, prompt string = "",""

func main() {

	LoadEnvironment()
	gitDiff, gitBranch := GetCommitContext()
	
	Info("Generating commit message.")

    //commitStandard := "Conventional Commits"
    commitStandard = GetEnvArg("GMMIT_COMMIT_PATTERN", "<type>[optional scope]: <description> (#<ticket-id>)")
	prompt = fmt.Sprintf("Create a git commit message following the \"Conventional Commits\" standard: \"%s\". The Ticket ID MUST be present on the first line, look for it on the branch name: \"%s\". Respond with the commit message only. First line can not be a generic line, must be a specific change. If there are many changes, list the rest at the end. These are the file changes to be pushed:\n%s",
				commitStandard, gitBranch, gitDiff)

	Debug(prompt)

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(GetEnvArg("GMMIT_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	model := client.GenerativeModel("gemini-1.5-pro")
	cs := model.StartChat()

	send := func(msg string) *genai.GenerateContentResponse {
		res, err := cs.SendMessage(ctx, genai.Text(msg))
		if err != nil {
			log.Fatal(err)
		}
		return res
	}

	res := send(prompt)
	stringRes := responseToString(res)
	printResponse(res)
	Info("Create a commit with this message? [y/N]")
	reader := bufio.NewReader(os.Stdin)
	confirmation, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
	confirmation = strings.ToLower(strings.TrimSpace(confirmation))
	if confirmation == "y" || confirmation == "yes" {
		Info("Creating Commit...")
		gitCommit := RunCommand("git","commit","-m",stringRes)
    	Info("Git Command Log:")
		fmt.Println(string(gitCommit))
	} 
}

func responseToString(resp *genai.GenerateContentResponse)(string){
	stringResponse := ""
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				stringResponse += fmt.Sprintln(part)
			}
		}
	}
	return stringResponse
}
func printResponse(resp *genai.GenerateContentResponse) {
	Info("Generated Text:")
	Info("---")
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				fmt.Println(part)
			}
		}
	}
	Info("---")
}