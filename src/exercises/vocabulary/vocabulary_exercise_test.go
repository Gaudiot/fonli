package vocabularyexercise

import (
	"errors"
	"strings"
	"testing"

	aiservice "gaudiot.com/fonli/base/http_services/ai_service"
	user_repository "gaudiot.com/fonli/base/repositories/user"
)

func testVocabularyUserRepo() *user_repository.UserRepositoryMock {
	const id = "test-user-id"
	return &user_repository.UserRepositoryMock{
		Users: map[string]*user_repository.User{
			id: {ID: id, LifestyleTopics: "music, travel"},
		},
	}
}

const testVocabularyUserID = "test-user-id"

func TestGenerate(t *testing.T) {
	mockAI := &aiservice.AIServiceMock{}
	ve := NewVocabularyExercise(mockAI, testVocabularyUserRepo())

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
	got, err := ve.Generate(exercisesQuantity, "pt", "it", testVocabularyUserID)

	if err != nil {
		t.Errorf("Generate(%d) should not return an error, but got %v", exercisesQuantity, err)
	}

	if got == nil {
		t.Fatalf("Generate(%d) should return an exercise, but got nil", exercisesQuantity)
	}

	if len(got.Questions) != exercisesQuantity {
		t.Errorf("Generate(%d) should have %d questions, but has %d", exercisesQuantity, exercisesQuantity, len(got.Questions))
	}

	if got.Questions[1].Word != "carro" {
		t.Errorf("Generate: expected second word to be 'carro', got '%s'", got.Questions[1].Word)
	}

	if got.Questions[1].Translation != "macchina" {
		t.Errorf("Generate: expected second translation to be 'macchina', got '%s'", got.Questions[1].Translation)
	}
}

func TestGenerate_PromptUsesUserData(t *testing.T) {
	mockAI := &aiservice.AIServiceMock{}
	ve := NewVocabularyExercise(mockAI, testVocabularyUserRepo())

	var gotPrompt string
	mockAI.PromptWithStructuredResponseFunc = func(prompt string, model map[string]any) (string, error) {
		gotPrompt = prompt
		return `{"questions": []}`, nil
	}

	if _, err := ve.Generate(5, "pt", "it", testVocabularyUserID); err != nil {
		t.Fatalf("Generate() should not return an error, but got %v", err)
	}

	for _, want := range []string{"pt", "it", "5 exercises", "music, travel"} {
		if !strings.Contains(gotPrompt, want) {
			t.Errorf("Generate: expected prompt to contain '%s', but it does not", want)
		}
	}
}

func TestGenerate_WithAIError(t *testing.T) {
	mockAI := &aiservice.AIServiceMock{}
	ve := NewVocabularyExercise(mockAI, testVocabularyUserRepo())

	mockAI.PromptWithStructuredResponseFunc = func(prompt string, model map[string]any) (string, error) {
		return "", errors.New("AI service failed")
	}

	exercisesQuantity := 3
	got, err := ve.Generate(exercisesQuantity, "pt", "it", testVocabularyUserID)

	if err == nil {
		t.Errorf("Generate(%d) should return an error when AI service fails, but got nil", exercisesQuantity)
	}

	if got != nil {
		t.Errorf("Generate(%d) should return nil exercise when error occurs, but got %v", exercisesQuantity, got)
	}
}

func TestGenerate_WithInvalidAIResponse(t *testing.T) {
	mockAI := &aiservice.AIServiceMock{}
	ve := NewVocabularyExercise(mockAI, testVocabularyUserRepo())

	mockAI.PromptWithStructuredResponseFunc = func(prompt string, model map[string]any) (string, error) {
		return "not a json", nil
	}

	got, err := ve.Generate(3, "pt", "it", testVocabularyUserID)

	if err != nil {
		t.Errorf("Generate() should not return an error for an invalid AI response, but got %v", err)
	}

	if got == nil {
		t.Fatalf("Generate() should return an empty exercise for an invalid AI response, but got nil")
	}

	if len(got.Questions) != 0 {
		t.Errorf("Generate() should return no questions for an invalid AI response, but got %d", len(got.Questions))
	}
}
