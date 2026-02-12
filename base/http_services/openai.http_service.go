package httpservices

import (
	"context"

	"gaudiot.com/fonli/core"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

func GetOpenAIStructuredResponse(prompt string, model map[string]any) (string, error) {
	ctx := context.Background()
	envConfig, err := core.LoadEnvConfig()
	if err != nil {
		return "", err
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
