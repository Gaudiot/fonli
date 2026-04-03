package wordconjugationexercise

import (
	"encoding/json"
	"fmt"

	"gaudiot.com/fonli/base"
	aiservice "gaudiot.com/fonli/base/http_services/ai_service"
	wordtranslationexercise "gaudiot.com/fonli/src/exercises/word_translation"
)

type Conjugation struct {
	Person      string `json:"person" jsonschema_description:"The person of the conjugation"`
	Number      string `json:"number" jsonschema_description:"The number of the conjugation"`
	Conjugation string `json:"conjugation" jsonschema_description:"The conjugation of the word"`
}

type WordConjugationExercise struct {
	Word         string        `json:"word" jsonschema_description:"The word which the user should conjugate"`
	Tense        string        `json:"tense" jsonschema_description:"The tense of the conjugation"`
	Conjugations []Conjugation `json:"conjugations" jsonschema_description:"The conjugations of the word"`
}

type WordConjugation struct {
	aiService aiservice.AIService
}

func NewWordConjugation(aiService aiservice.AIService) *WordConjugation {
	return &WordConjugation{
		aiService: aiService,
	}
}

func (w *WordConjugation) GenerateExercise(tense base.Tense, foreignLanguageCode string) (*WordConjugationExercise, error) {
	foreignLanguage := base.LanguageFromCountryCode(foreignLanguageCode)
	exerciseSchema := wordtranslationexercise.GenerateSchema[WordConjugationExercise]()

	prompt := fmt.Sprintf(
		`Create an exercise for the user to conjugate a verb used in daily life in %s.
		The tense must be %s. And the tense should be returned in %s.
		If the tense is not valid for the verb choose the most similar tense.
		Other acceptable verbs are the ones that are used in common scenarios such as: 
		- Going on vacation
		- Summer Olympics
		- Partying with friends
		- Practice sports
		- Ask for something at the restaurant
		`,
		foreignLanguage,
		tense,
		foreignLanguage,
	)

	response, err := w.aiService.PromptWithStructuredResponse(prompt, exerciseSchema)
	if err != nil {
		return nil, err
	}

	var exercise WordConjugationExercise
	_ = json.Unmarshal([]byte(response), &exercise)

	return &exercise, nil
}
