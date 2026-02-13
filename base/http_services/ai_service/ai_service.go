package aiservice

type AIService interface {
	PromptWithStructuredResponse(prompt string, model map[string]any) (string, error)
}

type AIServiceMock struct {
	PromptWithStructuredResponseFunc func(prompt string, model map[string]any) (string, error)
}

func (s *AIServiceMock) PromptWithStructuredResponse(prompt string, model map[string]any) (string, error) {
	return s.PromptWithStructuredResponseFunc(prompt, model)
}
