package wordconjugationexercise

import (
	"testing"

	"gaudiot.com/fonli/base"
	aiservice "gaudiot.com/fonli/base/http_services/ai_service"
)

var mockAiService = &aiservice.AIServiceMock{}
var wordConjugation = NewWordConjugation(mockAiService)

func TestWordConjugationExercise(t *testing.T) {
	mockAiService.PromptWithStructuredResponseFunc = func(prompt string, model map[string]any) (string, error) {
		return `{
			"word": "parlare",
			"tense": "presente",
			"conjugations": [
				{
					"person": "1st",
					"number": "singular",
					"conjugation": "parlo"
				},
				{
					"person": "2nd",
					"number": "singular",
					"conjugation": "parli"
				},
				{
					"person": "3rd",
					"number": "singular",
					"conjugation": "parla"
				},
				{
					"person": "1st",
					"number": "plural",
					"conjugation": "parliamo"
				},
				{
					"person": "2nd",
					"number": "plural",
					"conjugation": "parlate"
				},
				{
					"person": "3rd",
					"number": "plural",
					"conjugation": "parlano"
				}
			]
		}`, nil
	}

	got, err := wordConjugation.GenerateExercise(base.PresentSimple, "it")

	if err != nil {
		t.Errorf("GenerateExercise() should not return an error, but got %v", err)
	}

	if got == nil {
		t.Errorf("GenerateExercise() should return an exercise, but got nil")
	}

	if got.Word != "parlare" {
		t.Errorf("GenerateExercise() should return a word, but got %s", got.Word)
	}

	if got.Tense != "presente" {
		t.Errorf("GenerateExercise() should return a tense, but got %s", got.Tense)
	}

	if len(got.Conjugations) != 6 {
		t.Errorf("GenerateExercise() should return 6 conjugations, but got %d", len(got.Conjugations))
	}

	if got.Conjugations[0].Person != "1st" {
		t.Errorf("GenerateExercise() should return a person, but got %s", got.Conjugations[0].Person)
	}

	if got.Conjugations[0].Number != "singular" {
		t.Errorf("GenerateExercise() should return a number, but got %s", got.Conjugations[0].Number)
	}

	if got.Conjugations[0].Conjugation != "parlo" {
		t.Errorf("GenerateExercise() should return a conjugation, but got %s", got.Conjugations[0].Conjugation)
	}
}
