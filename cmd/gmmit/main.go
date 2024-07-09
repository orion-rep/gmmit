package main

import (
	"bufio"
	"context"
	"os"
	"fmt"
	"strings"
	"log"

	"github.com/joho/godotenv"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"

	. "gitlab.com/orion-rep/gmmit/internal/pkg/common"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	Info("Getting context information.")
	
	GetEnvArg("GMMIT_API_KEY")

	//gitDiff, err := exec.Command("git","diff","--staged").Output()
	gitDiff := RunCommand("git","diff","--staged")

	if len(gitDiff) <= 0 {
		Warning("Git diff returned no files.")
		Warning("Add some files to the staging area and run this command again.")
		os.Exit(0)
	}
	
	gitBranch := RunCommand("git", "rev-parse", "--abbrev-ref", "HEAD")
    
	Info("Generating commit message.")

    commitStandard := "Conventional Commits"
	prompt := fmt.Sprintf("Create a git commit message following the \"%s\" standard for branch \"%s\". Respond with the commit message only. First line can not be a generic line, must be a specific change. If there are many changes, list the rest at the end. These are the file changes to be pushed:%s",
				commitStandard, gitBranch, gitDiff)

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GMMIT_API_KEY")))
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
    	CheckIfError(err)
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