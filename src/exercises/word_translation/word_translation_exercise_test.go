package wordtranslationexercise

import (
	"errors"
	"testing"

	aiservice "gaudiot.com/fonli/base/http_services/ai_service"
	user_repository "gaudiot.com/fonli/base/repositories/user"
)

func testWordTranslationUserRepo() *user_repository.UserRepositoryMock {
	const testUserID = "test-user-id"
	return &user_repository.UserRepositoryMock{
		Users: map[string]*user_repository.User{
			testUserID: {
				ID:              testUserID,
				LifestyleTopics: "music, travel",
			},
		},
	}
}

const testWordTranslationUserID = "test-user-id"

func TestNativeToForeignExercise(t *testing.T) {
	mockAI := &aiservice.AIServiceMock{}
	wt := NewWordTranslation(mockAI, testWordTranslationUserRepo())

	mockAI.PromptWithStructuredResponseFunc = func(prompt string, model map[string]any) (string, error) {
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
	got, err := wt.NativeToForeignExercise(exercisesQuantity, "pt", "it", testWordTranslationUserID)

	if err != nil {
		t.Errorf("NativeToForeignExercise(%d) should not return an error, but got %v", exercisesQuantity, err)
	}

	if len(got.Questions) != exercisesQuantity {
		t.Errorf("NativeToForeignExercise(%d) should have %d questions, but has %d", exercisesQuantity, exercisesQuantity, len(got.Questions))
	}

	if got.Questions[1].Word != "carro" {
		t.Errorf("NativeToForeignExercise: expected second word to be 'carro', got '%s'", got.Questions[1].Word)
	}

	if got.Questions[1].Translation != "macchina" {
		t.Errorf("NativeToForeignExercise: expected second translation to be 'macchina', got '%s'", got.Questions[1].Translation)
	}
}

func TestForeignToNativeExercise(t *testing.T) {
	mockAI := &aiservice.AIServiceMock{}
	wt := NewWordTranslation(mockAI, testWordTranslationUserRepo())

	mockAI.PromptWithStructuredResponseFunc = func(prompt string, model map[string]any) (string, error) {
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
	got, err := wt.ForeignToNativeExercise(exercisesQuantity, "it", "pt", testWordTranslationUserID)

	if err != nil {
		t.Errorf("ForeignToNativeExercise(%d) should not return an error, but got %v", exercisesQuantity, err)
	}

	if len(got.Questions) != exercisesQuantity {
		t.Errorf("ForeignToNativeExercise(%d) should have %d questions, but has %d", exercisesQuantity, exercisesQuantity, len(got.Questions))
	}

	if got.Questions[0].Word != "casa" {
		t.Errorf("ForeignToNativeExercise: expected first word to be 'casa', got '%s'", got.Questions[0].Word)
	}

	if got.Questions[0].Translation != "casa" {
		t.Errorf("ForeignToNativeExercise: expected first translation to be 'casa', got '%s'", got.Questions[0].Translation)
	}
}

func TestNativeToForeignExercise_WithError(t *testing.T) {
	mockAI := &aiservice.AIServiceMock{}
	wt := NewWordTranslation(mockAI, testWordTranslationUserRepo())

	mockAI.PromptWithStructuredResponseFunc = func(prompt string, model map[string]any) (string, error) {
		return "", errors.New("AI service failed")
	}

	exercisesQuantity := 3
	got, err := wt.NativeToForeignExercise(exercisesQuantity, "pt", "it", testWordTranslationUserID)

	if err == nil {
		t.Errorf("NativeToForeignExercise(%d) should return an error when AI service fails, but got nil", exercisesQuantity)
	}

	if got != nil {
		t.Errorf("NativeToForeignExercise(%d) should return nil exercise when error occurs, but got %v", exercisesQuantity, got)
	}
}

func TestForeignToNativeExercise_WithError(t *testing.T) {
	mockAI := &aiservice.AIServiceMock{}
	wt := NewWordTranslation(mockAI, testWordTranslationUserRepo())

	mockAI.PromptWithStructuredResponseFunc = func(prompt string, model map[string]any) (string, error) {
		return "", errors.New("AI service failed")
	}

	exercisesQuantity := 3
	got, err := wt.ForeignToNativeExercise(exercisesQuantity, "it", "pt", testWordTranslationUserID)

	if err == nil {
		t.Errorf("ForeignToNativeExercise(%d) should return an error when AI service fails, but got nil", exercisesQuantity)
	}

	if got != nil {
		t.Errorf("ForeignToNativeExercise(%d) should return nil exercise when error occurs, but got %v", exercisesQuantity, got)
	}
}
