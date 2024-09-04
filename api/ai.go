package api

import (
	"context"
	"fmt"
	"os"

	"github.com/emmaahmads/summafy/util"
	"github.com/sashabaranov/go-openai"
)

func (server *Server) SummarizeTextFile(filePath string) (string, error) {
	client := openai.NewClient(server.apiKey)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	resp, err := client.CreateCompletion(
		context.Background(),
		openai.CompletionRequest{
			Model:     openai.GPT3Babbage002,
			MaxTokens: 120,
			Prompt:    fmt.Sprintf("Summarize the following text:\n\n%s", string(content)),
		},
	)
	if err != nil {
		return "", err
	}
	util.MyGinLogger("resp0", resp.Choices[0].Text)
	return resp.Choices[0].Text, nil
}
