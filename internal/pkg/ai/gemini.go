package gemini

import (
	"context"
	"fmt"
	"log"
	"time"
	"strings"
	"strconv"
	
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"

	. "gitlab.com/orion-rep/gmmit/internal/pkg/common"
)

func RunPrompt(prompt string) (*genai.GenerateContentResponse) {
	ctx := context.Background()
	Debug("Getting GenAI Client")
	client, err := genai.NewClient(ctx, option.WithAPIKey(GetEnvArg("GMMIT_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	Debug("Getting GenAI Model")
	model := client.GenerativeModel(GetEnvArg("GMMIT_MODEL", "gemini-1.5-pro"))
	model.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockOnlyHigh,
		},
	}
	
	res := SendMessageToModel(ctx, model, prompt)
	
	return res
}

// SendMessageToModel sends a message to the specified GenerativeModel and returns the response.
// It retries the SendMessage call up to maxRetries times if the error is a 500 Internal Server Error,
// with a delay of retryDelay between each retry attempt.
//
// ctx is the context.Context to use for the operation.
// model is a pointer to the genai.GenerativeModel to use for sending the message.
// msg is the message to send to the model.
//
// Returns a pointer to the genai.GenerateContentResponse containing the model's response.
// If an error occurs and cannot be retried, it logs the error using log.Fatal and exits the program.
func SendMessageToModel(ctx context.Context, model *genai.GenerativeModel, msg string) *genai.GenerateContentResponse {
    Debug("Starting GenAI Chat")
    cs := model.StartChat()
    Debug("Sending GenAI Message")

    var res *genai.GenerateContentResponse
    var err error
    maxRetries, err := strconv.Atoi(GetEnvArg("GMMIT_MAX_RETRIES", "5"))
	CheckIfError(err)
    retryDelay, err := strconv.Atoi(GetEnvArg("GMMIT_RETRY_DELAY", "5"))
	CheckIfError(err)

	retryDelayDuration := time.Duration(retryDelay) * time.Second

    for i := 0; i < maxRetries; i++ {
        res, err = cs.SendMessage(ctx, genai.Text(msg))
        if err == nil {
            break
        }

        // Check if the error is a 500 Internal Server Error
        if strings.Contains(err.Error(), "500") {
            Debug(fmt.Sprintf("Received 500 error, retrying in %v (attempt %d/%d)", retryDelayDuration, i+1, maxRetries))
            time.Sleep(retryDelayDuration)
            continue
        }

        // For other errors, break the loop
        log.Fatal(err)
    }

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
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				fmt.Println(part)
			}
		}
	}
}