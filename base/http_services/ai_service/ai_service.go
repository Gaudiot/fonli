package aiservice

type AIService interface {
	Prompt(prompt string) (string, error)
	PromptWithStructuredResponse(prompt string, model map[string]any) (string, error)
}

type AIServiceMock struct {
	PromptFunc                       func(prompt string) (string, error)
	PromptWithStructuredResponseFunc func(prompt string, model map[string]any) (string, error)
}

func (s *AIServiceMock) Prompt(prompt string) (string, error) {
	return s.PromptFunc(prompt)
}

func (s *AIServiceMock) PromptWithStructuredResponse(prompt string, model map[string]any) (string, error) {
	return s.PromptWithStructuredResponseFunc(prompt, model)
}
