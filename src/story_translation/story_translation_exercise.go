package storytranslation

import (
	"encoding/json"
	"fmt"

	"gaudiot.com/fonli/base"
	aiservice "gaudiot.com/fonli/base/http_services/ai_service"
	"github.com/invopop/jsonschema"
)

func generateSchema[T any]() map[string]any {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)

	data, _ := json.Marshal(schema)
	var result map[string]any
	json.Unmarshal(data, &result)
	return result
}

type GenerateStoryResponse struct {
	Story string `json:"story" jsonschema_description:"A medium-sized, engaging story in Portuguese for translation to Italian"`
}

type EvaluateTranslationResponse struct {
	Score              int      `json:"score" jsonschema_description:"Score from 0 to 10 representing translation quality"`
	Errors             []string `json:"errors" jsonschema_description:"A list of the main translation mistakes"`
	CorrectTranslation string   `json:"correct_translation" jsonschema_description:"The correct translation of the original story in Italian"`
}

type StoryTranslation struct {
	aiService aiservice.AIService
}

func NewStoryTranslation(aiService aiservice.AIService) *StoryTranslation {
	return &StoryTranslation{
		aiService: aiService,
	}
}

// GenerateStory generates a medium-sized story in Portuguese for the user to translate
func (h *StoryTranslation) GenerateStory() (*GenerateStoryResponse, error) {
	nativeLanguage := base.Languages.Portuguese
	foreignLanguage := base.Languages.Italian

	schema := generateSchema[GenerateStoryResponse]()
	prompt := fmt.Sprintf(
		`Create a story in %s, of medium length, to be translated into %s by a student learning it.
		The story should be interesting, appropriate for students (not too easy, not too hard), and contain between 3 to 4 paragraphs.
		The result should be only a JSON object that contains the field "story".`,
		nativeLanguage,
		foreignLanguage,
	)

	response, err := h.aiService.PromptWithStructuredResponse(prompt, schema)
	if err != nil {
		return nil, err
	}

	var res GenerateStoryResponse
	if err := json.Unmarshal([]byte(response), &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// EvaluateTranslation sends the original story and the user translation to AI to receive evaluation and feedback
func (h *StoryTranslation) EvaluateTranslation(originalStory, userTranslation string) (*EvaluateTranslationResponse, error) {
	nativeLanguage := base.Languages.Portuguese
	foreignLanguage := base.Languages.Italian

	schema := generateSchema[EvaluateTranslationResponse]()
	prompt := fmt.Sprintf(
		`You will evaluate the translation made by a student from %s to %s.
		You will receive the original text, followed by the student's translation/response.
		Evaluate the response and assign a score from 0 to 10 (field "score").
		List the main errors found (field "errors", which should be a list of strings).
		The corrections should be made in the user's native language.
		If there are many errors or the translation is poor, limit the maximum score to 10.
		In the field "correct_translation", provide the correct translation of the original text into %s.
		The result should be a JSON object with the fields "score", "errors" (list), and "correct_translation".
		---
		Original text (%s):
		%s

		User's translation:
		%s
		`,
		nativeLanguage,
		foreignLanguage,
		foreignLanguage,
		nativeLanguage,
		originalStory,
		userTranslation,
	)

	response, err := h.aiService.PromptWithStructuredResponse(prompt, schema)
	if err != nil {
		return nil, err
	}

	var res EvaluateTranslationResponse
	if err := json.Unmarshal([]byte(response), &res); err != nil {
		return nil, err
	}

	// Garantia de score máximo a 10, embora dependa do LLM:
	if res.Score > 10 {
		res.Score = 10
	}
	if res.Score < 0 {
		res.Score = 0
	}

	return &res, nil
}
