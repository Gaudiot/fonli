package wordtranslationexercise

import (
	"errors"
	"testing"

	aiservice "gaudiot.com/fonli/base/http_services/ai_service"
)

var mockAiService = &aiservice.AIServiceMock{}
var wordTranslation = NewWordTranslation(mockAiService)

func TestNativeToForeignExercise(t *testing.T) {
	mockAiService.PromptWithStructuredResponseFunc = func(prompt string, model map[string]any) (string, error) {
		return `{
			"questions": [
				{
					"word": "casa",
					"translation": "casa"
				},
				{
					"word": "carro",
					"translation": "macchina"
				},
				{
					"word": "moto",
					"translation": "moto"
				}
			]
		}`, nil
	}

	exercisesQuantity := 3
	got, err := wordTranslation.NativeToForeignExercise(exercisesQuantity)

	if err != nil {
		t.Errorf("NativeToForeignExercise(%d) should not return an error, but got %v", exercisesQuantity, err)
	}

	if len(got.Questions) != exercisesQuantity {
		t.Errorf("NativeToForeignExercise(%d) should have %d questions, but has %d", exercisesQuantity, exercisesQuantity, len(got.Questions))
	}

	// Verify first question
	if got.Questions[1].Word != "carro" {
		t.Errorf("NativeToForeignExercise: expected first word to be 'casa', got '%s'", got.Questions[0].Word)
	}

	if got.Questions[1].Translation != "macchina" {
		t.Errorf("NativeToForeignExercise: expected first translation to be 'casa', got '%s'", got.Questions[0].Translation)
	}
}

func TestForeignToNativeExercise(t *testing.T) {
	mockAiService.PromptWithStructuredResponseFunc = func(prompt string, model map[string]any) (string, error) {
		return `{
			"questions": [
				{
					"word": "casa",
					"translation": "casa"
				},
				{
					"word": "macchina",
					"translation": "carro"
				},
				{
					"word": "moto",
					"translation": "moto"
				}
			]
		}`, nil
	}

	exercisesQuantity := 3
	got, err := wordTranslation.ForeignToNativeExercise(exercisesQuantity)

	if err != nil {
		t.Errorf("ForeignToNativeExercise(%d) should not return an error, but got %v", exercisesQuantity, err)
	}

	if len(got.Questions) != exercisesQuantity {
		t.Errorf("ForeignToNativeExercise(%d) should have %d questions, but has %d", exercisesQuantity, exercisesQuantity, len(got.Questions))
	}

	// Verify first question
	if got.Questions[0].Word != "casa" {
		t.Errorf("ForeignToNativeExercise: expected first word to be 'casa', got '%s'", got.Questions[0].Word)
	}

	if got.Questions[0].Translation != "casa" {
		t.Errorf("ForeignToNativeExercise: expected first translation to be 'casa', got '%s'", got.Questions[0].Translation)
	}
}

func TestNativeToForeignExercise_WithError(t *testing.T) {
	mockAiService.PromptWithStructuredResponseFunc = func(prompt string, model map[string]any) (string, error) {
		return "", errors.New("AI service failed")
	}

	exercisesQuantity := 3
	got, err := wordTranslation.NativeToForeignExercise(exercisesQuantity)

	if err == nil {
		t.Errorf("NativeToForeignExercise(%d) should return an error when AI service fails, but got nil", exercisesQuantity)
	}

	if got != nil {
		t.Errorf("NativeToForeignExercise(%d) should return nil exercise when error occurs, but got %v", exercisesQuantity, got)
	}
}

func TestForeignToNativeExercise_WithError(t *testing.T) {
	mockAiService.PromptWithStructuredResponseFunc = func(prompt string, model map[string]any) (string, error) {
		return "", errors.New("AI service failed")
	}

	exercisesQuantity := 3
	got, err := wordTranslation.ForeignToNativeExercise(exercisesQuantity)

	if err == nil {
		t.Errorf("ForeignToNativeExercise(%d) should return an error when AI service fails, but got nil", exercisesQuantity)
	}

	if got != nil {
		t.Errorf("ForeignToNativeExercise(%d) should return nil exercise when error occurs, but got %v", exercisesQuantity, got)
	}
}
