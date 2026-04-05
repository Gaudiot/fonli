package storytranslation

import (
	"errors"
	"testing"

	aiservice "gaudiot.com/fonli/base/http_services/ai_service"
	user_repository "gaudiot.com/fonli/base/repositories/user"
)

func testStoryTranslationUserRepo() *user_repository.UserRepositoryMock {
	const id = "test-user-id"
	return &user_repository.UserRepositoryMock{
		Users: map[string]*user_repository.User{
			id: {ID: id, LifestyleTopics: "nature, history"},
		},
	}
}

const testStoryTranslationUserID = "test-user-id"

func TestGenerateStory(t *testing.T) {
	mockAI := &aiservice.AIServiceMock{}
	st := NewStoryTranslation(mockAI, testStoryTranslationUserRepo())

	mockAI.PromptWithStructuredResponseFunc = func(prompt string, model map[string]any) (string, error) {
		return `{
			"story": "Era uma vez uma menina que morava em uma pequena vila perto da floresta. Todos os dias, ela levava pão fresco para a avó."
		}`, nil
	}

	got, err := st.GenerateStory("pt", "it", testStoryTranslationUserID)
	if err != nil {
		t.Errorf("GenerateStory() should not return an error, but got %v", err)
	}

	expectedStory := "Era uma vez uma menina que morava em uma pequena vila perto da floresta. Todos os dias, ela levava pão fresco para a avó."
	if got.Story != expectedStory {
		t.Errorf("GenerateStory() returned unexpected story. Got: %s, want: %s", got.Story, expectedStory)
	}
}

func TestGenerateStory_WithError(t *testing.T) {
	mockAI := &aiservice.AIServiceMock{}
	st := NewStoryTranslation(mockAI, testStoryTranslationUserRepo())

	mockAI.PromptWithStructuredResponseFunc = func(prompt string, model map[string]any) (string, error) {
		return "", errors.New("AI service failed")
	}

	got, err := st.GenerateStory("pt", "it", testStoryTranslationUserID)
	if err == nil {
		t.Errorf("GenerateStory() should return an error when AI service fails, but got nil")
	}
	if got != nil {
		t.Errorf("GenerateStory() should return nil when error occurs, but got: %#v", got)
	}
}

func TestEvaluateTranslation(t *testing.T) {
	mockAI := &aiservice.AIServiceMock{}
	st := NewStoryTranslation(mockAI, testStoryTranslationUserRepo())

	mockAI.PromptWithStructuredResponseFunc = func(prompt string, model map[string]any) (string, error) {
		return `{
			"score": 9,
			"errors": ["Verb agreement error in the second sentence.", "Incorrect usage of the definite article."],
			"correct_translation": "C'era una volta una ragazza che viveva in un piccolo villaggio vicino alla foresta. Ogni giorno portava pane fresco alla nonna."
		}`, nil
	}

	story := "Era uma vez uma menina que morava em uma pequena vila perto da floresta. Todos os dias, ela levava pão fresco para a avó."
	userTranslation := "C'era una volta una ragazza che viveva in un piccolo villaggio vicino alla foresta. Lei portava pane fresco alla sua nonna ogni giorno."
	got, err := st.EvaluateTranslation(story, userTranslation, "pt", "it")
	if err != nil {
		t.Errorf("EvaluateTranslation() should not return an error, but got %v", err)
	}

	if got.Score != 9 {
		t.Errorf("EvaluateTranslation() expected score 9, got %d", got.Score)
	}

	if len(got.Errors) != 2 {
		t.Errorf("EvaluateTranslation() expected 2 errors, got %d", len(got.Errors))
	}

	expectedCorrectTranslation := "C'era una volta una ragazza che viveva in un piccolo villaggio vicino alla foresta. Ogni giorno portava pane fresco alla nonna."
	if got.CorrectTranslation != expectedCorrectTranslation {
		t.Errorf("EvaluateTranslation() returned unexpected correct_translation. Got: %s, want: %s", got.CorrectTranslation, expectedCorrectTranslation)
	}
}

func TestEvaluateTranslation_WithManyErrors_EnforcesScoreLimit(t *testing.T) {
	mockAI := &aiservice.AIServiceMock{}
	st := NewStoryTranslation(mockAI, testStoryTranslationUserRepo())

	mockAI.PromptWithStructuredResponseFunc = func(prompt string, model map[string]any) (string, error) {
		return `{
			"score": 10,
			"errors": ["Error1", "Error2", "Error3", "Error4", "Error5", "Error6", "Error7", "Error8", "Error9", "Error10", "Error11"],
			"correct_translation": "correct translation"
		}`, nil
	}

	story := "dummy"
	userTranslation := "dummy"
	got, err := st.EvaluateTranslation(story, userTranslation, "pt", "it")
	if err != nil {
		t.Errorf("EvaluateTranslation() should not return an error, but got %v", err)
	}

	if got.Score > 10 {
		t.Errorf("EvaluateTranslation() should not return a score higher than 10 when there are many errors, got %d", got.Score)
	}
}

func TestEvaluateTranslation_WithError(t *testing.T) {
	mockAI := &aiservice.AIServiceMock{}
	st := NewStoryTranslation(mockAI, testStoryTranslationUserRepo())

	mockAI.PromptWithStructuredResponseFunc = func(prompt string, model map[string]any) (string, error) {
		return "", errors.New("AI service failed")
	}

	got, err := st.EvaluateTranslation("dummy", "dummy", "pt", "it")
	if err == nil {
		t.Errorf("EvaluateTranslation() should return an error when AI service fails, but got nil")
	}
	if got != nil {
		t.Errorf("EvaluateTranslation() should return nil when error occurs, but got: %#v", got)
	}
}
