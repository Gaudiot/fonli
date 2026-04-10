package wordconjugationexercise

import (
	"encoding/json"
	"fmt"

	"gaudiot.com/fonli/base"
	aiservice "gaudiot.com/fonli/base/http_services/ai_service"
	user_repository "gaudiot.com/fonli/base/repositories/user"
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
	aiService      aiservice.AIService
	userRepository user_repository.UserRepository
}

func NewWordConjugation(aiService aiservice.AIService, userRepository user_repository.UserRepository) *WordConjugation {
	return &WordConjugation{
		aiService:      aiService,
		userRepository: userRepository,
	}
}

func (w *WordConjugation) GenerateExercise(tense base.Tense, foreignLanguageCode, userID string) (*WordConjugationExercise, error) {
	user, err := w.userRepository.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	userLifestyleTopics := user.LifestyleTopics

	foreignLanguage := base.LanguageFromCountryCode(foreignLanguageCode)
	exerciseSchema := wordtranslationexercise.GenerateSchema[WordConjugationExercise]()

	prompt := fmt.Sprintf(
		`Create an exercise for the user to conjugate a verb used in daily life in %s.
		The tense must be %s. And the tense should be returned in %s.
		If the tense is not valid for the verb choose the most similar tense.
		There are some lifestyle topics that the user likes to use in his daily life, these are: %s.
		Prefer verbs and contexts that relate to these topics when relevant.
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
		userLifestyleTopics,
	)

	response, err := w.aiService.PromptWithStructuredResponse(prompt, exerciseSchema)
	if err != nil {
		return nil, err
	}

	var exercise WordConjugationExercise
	_ = json.Unmarshal([]byte(response), &exercise)

	return &exercise, nil
}
