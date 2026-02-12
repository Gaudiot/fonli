package src

import (
	"encoding/json"
	"fmt"

	"gaudiot.com/fonli/base"
	httpservices "gaudiot.com/fonli/base/http_services"
	"github.com/invopop/jsonschema"
)

type NativeToForeignExerciseQuestion struct {
	Word        string `json:"word" jsonschema_description:"The word which the user should translate"`
	Translation string `json:"translation" jsonschema_description:"The translation of the word"`
}

type NativeToForeignExercise struct {
	Questions []NativeToForeignExerciseQuestion `json:"questions" jsonschema_description:"The questions for the exercise"`
}

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

func CreateNativeToForeignExercise(exercisesQuantity int) (*NativeToForeignExercise, error) {
	nativeLanguage := base.Languages.Portuguese
	foreignLanguage := base.Languages.Italian
	exerciseSchema := GenerateSchema[NativeToForeignExercise]()

	prompt := fmt.Sprintf(
		"Create %d exercises for the user to translate simple words from %s to %s, the words should be common and used in daily life. They can also be from less common scenarios like sports, olympics, vacation, party, etc.",
		exercisesQuantity,
		nativeLanguage,
		foreignLanguage,
	)

	response, err := httpservices.GetOpenAIStructuredResponse(prompt, exerciseSchema)
	if err != nil {
		return nil, err
	}

	var exercise NativeToForeignExercise
	_ = json.Unmarshal([]byte(response), &exercise)

	return &exercise, nil
}
