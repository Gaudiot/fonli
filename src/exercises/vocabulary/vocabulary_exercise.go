package vocabularyexercise

import (
	"encoding/json"
	"fmt"

	aiservice "gaudiot.com/fonli/base/http_services/ai_service"
	user_repository "gaudiot.com/fonli/base/repositories/user"
	"gaudiot.com/fonli/core/utils"
)

type vocabularyQuestionSchema struct {
	Word        string `json:"word" jsonschema_description:"The word which the user should translate"`
	Translation string `json:"translation" jsonschema_description:"The translation of the word"`
}

type VocabularyExerciseSchema struct {
	Questions []vocabularyQuestionSchema `json:"questions" jsonschema_description:"The questions for the exercise"`
}

type VocabularyExercise struct {
	aiService      aiservice.AIService
	userRepository user_repository.UserRepository
}

func NewVocabularyExercise(aiService aiservice.AIService, userRepository user_repository.UserRepository) *VocabularyExercise {
	return &VocabularyExercise{
		aiService:      aiService,
		userRepository: userRepository,
	}
}

func (w *VocabularyExercise) Generate(exercisesQuantity int, baseLanguage, targetLanguage, userID string) (*VocabularyExerciseSchema, error) {
	user, err := w.userRepository.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	userLifestyleTopics := user.LifestyleTopics

	exerciseSchema := utils.GenerateSchema[VocabularyExerciseSchema]()

	prompt := fmt.Sprintf(
		`I am a person who speaks %s and I want to learn to speak %s.
		At this moment I want to improve my vocabulary in the language %s.
		Therefore, I want to create exercises that require me to translate words from %s to %s. 
		I also want that half of the exercises needs me to translate the other way around.
		There must be a total of %d exercises.
		For each exercise the word and its translation must be single word.
		It should be a mix of common words and lifestyle words. Here are some lifestyle topics about me: %s.
		`,
		baseLanguage,
		targetLanguage,
		targetLanguage,
		baseLanguage,
		targetLanguage,
		exercisesQuantity,
		userLifestyleTopics,
	)

	response, err := w.aiService.PromptWithStructuredResponse(prompt, exerciseSchema)
	if err != nil {
		return nil, err
	}

	var exercise VocabularyExerciseSchema
	_ = json.Unmarshal([]byte(response), &exercise)

	return &exercise, nil
}
