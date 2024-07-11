package gemini

import (
	"context"
	"fmt"
	"log"
	
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"

	. "gitlab.com/orion-rep/gmmit/internal/pkg/common"
)

func RunPrompt(prompt string) (*genai.GenerateContentResponse) {

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(GetEnvArg("GMMIT_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	model := client.GenerativeModel(GetEnvArg("GMMIT_MODEL", "gemini-1.5-pro"))
	
	res := SendMessageToModel(ctx, model, prompt)

	return res
}

func SendMessageToModel(ctx context.Context, model *genai.GenerativeModel, msg string) *genai.GenerateContentResponse {
	cs := model.StartChat()
	res, err := cs.SendMessage(ctx, genai.Text(msg))
	if err != nil {
		log.Fatal(err)
	}
	return res
}

func ModelResponseToString(resp *genai.GenerateContentResponse)(string){
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

func PrintModelResponse(resp *genai.GenerateContentResponse) {
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