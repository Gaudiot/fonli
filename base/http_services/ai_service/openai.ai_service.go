package aiservice

import (
	"context"
	"errors"

	"gaudiot.com/fonli/core"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

type OpenAIAIService struct{}

func NewOpenAIAIService() *OpenAIAIService {
	return &OpenAIAIService{}
}

func (s *OpenAIAIService) Prompt(prompt string) (string, error) {
	ctx := context.Background()
	envConfig := core.GetEnvConfig()
	if envConfig == nil {
		return "", errors.New("environment configuration is not loaded")
	}

	client := openai.NewClient(option.WithAPIKey(envConfig.OpenAIKey))
	resp, err := client.Responses.New(ctx, responses.ResponseNewParams{
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String(prompt),
		},
		Model: openai.ChatModelGPT4_1,
	})
	if err != nil {
		return "", err
	}

	return resp.OutputText(), nil
}

func (s *OpenAIAIService) PromptWithStructuredResponse(prompt string, model map[string]any) (string, error) {
	ctx := context.Background()
	envConfig := core.GetEnvConfig()
	if envConfig == nil {
		return "", errors.New("environment configuration is not loaded")
	}

	client := openai.NewClient(option.WithAPIKey(envConfig.OpenAIKey))
	resp, err2 := client.Responses.New(ctx, responses.ResponseNewParams{
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String(prompt),
		},
		Text: responses.ResponseTextConfigParam{
			Format: responses.ResponseFormatTextConfigParamOfJSONSchema(
				"model",
				model,
			),
		},
		Model: openai.ChatModelGPT4_1,
	})

	if err2 != nil {
		return "", err2
	}

	return resp.OutputText(), nil
}
