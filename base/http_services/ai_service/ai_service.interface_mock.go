package aiservice

import "errors"

type AIServiceMock struct {
	PromptFunc      func(prompt string) (string, error)
	PromptCallCount int

	PromptWithStructuredResponseFunc      func(prompt string, model map[string]any) (string, error)
	PromptWithStructuredResponseCallCount int
}

func (m *AIServiceMock) Prompt(prompt string) (string, error) {
	m.PromptCallCount++
	if m.PromptFunc != nil {
		return m.PromptFunc(prompt)
	}
	return "", errors.New("[Mock] not implemented")
}

func (m *AIServiceMock) PromptWithStructuredResponse(prompt string, model map[string]any) (string, error) {
	m.PromptWithStructuredResponseCallCount++
	if m.PromptWithStructuredResponseFunc != nil {
		return m.PromptWithStructuredResponseFunc(prompt, model)
	}
	return "", errors.New("[Mock] not implemented")
}
