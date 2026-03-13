package wordtranslationexercise

import (
	"encoding/json"
	"fmt"

	"gaudiot.com/fonli/base"
	aiservice "gaudiot.com/fonli/base/http_services/ai_service"
	"github.com/invopop/jsonschema"
)

func GenerateSchema[T any]() map[string]any {
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

type WordTranslationExerciseQuestion struct {
	Word        string `json:"word" jsonschema_description:"The word which the user should translate"`
	Translation string `json:"translation" jsonschema_description:"The translation of the word"`
}

type WordTranslationExercise struct {
	Questions []WordTranslationExerciseQuestion `json:"questions" jsonschema_description:"The questions for the exercise"`
}

type WordTranslation struct {
	aiService aiservice.AIService
}

func NewWordTranslation(aiService aiservice.AIService) *WordTranslation {
	return &WordTranslation{
		aiService: aiService,
	}
}

func (w *WordTranslation) NativeToForeignExercise(exercisesQuantity int, nativeLanguageCode string, foreignLanguageCode string) (*WordTranslationExercise, error) {
	nativeLanguage := base.LanguageFromCountryCode(nativeLanguageCode)
	foreignLanguage := base.LanguageFromCountryCode(foreignLanguageCode)
	exerciseSchema := GenerateSchema[WordTranslationExercise]()

	fmt.Println(nativeLanguage, foreignLanguage, nativeLanguageCode, foreignLanguageCode, "native to foreign")
	prompt := fmt.Sprintf(
		`Create %d exercises for the user to translate simple words from %s to %s, the words should be common and used in daily life.
		They can also be from less common scenarios like sports, olympics, vacation, party, etc.
		The response should be a JSON object, where the question must be in %s, and the translation mus be int %s
		`,
		exercisesQuantity,
		nativeLanguage,
		foreignLanguage,
		nativeLanguage,
		foreignLanguage,
	)

	fmt.Println(prompt)

	response, err := w.aiService.PromptWithStructuredResponse(prompt, exerciseSchema)
	if err != nil {
		return nil, err
	}

	var exercise WordTranslationExercise
	_ = json.Unmarshal([]byte(response), &exercise)

	return &exercise, nil
}

func (w *WordTranslation) ForeignToNativeExercise(exercisesQuantity int, foreignLanguageCode string, nativeLanguageCode string) (*WordTranslationExercise, error) {
	nativeLanguage := base.LanguageFromCountryCode(nativeLanguageCode)
	foreignLanguage := base.LanguageFromCountryCode(foreignLanguageCode)
	exerciseSchema := GenerateSchema[WordTranslationExercise]()

	fmt.Println(nativeLanguage, foreignLanguage, nativeLanguageCode, foreignLanguageCode, "foreign to native")
	prompt := fmt.Sprintf(
		`Create %d exercises for the user to translate simple words from %s to %s, the words should be common and used in daily life.
		They can also be from less common scenarios like sports, olympics, vacation, party, etc.
		The response should be a JSON object, where the question must be in %s, and the translation mus be int %s
		`,
		exercisesQuantity,
		foreignLanguage,
		nativeLanguage,
		foreignLanguage,
		nativeLanguage,
	)

	fmt.Println(prompt)

	response, err := w.aiService.PromptWithStructuredResponse(prompt, exerciseSchema)
	if err != nil {
		return nil, err
	}

	var exercise WordTranslationExercise
	_ = json.Unmarshal([]byte(response), &exercise)

	return &exercise, nil
}
