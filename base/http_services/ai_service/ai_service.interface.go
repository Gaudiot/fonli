package aiservice

type AIService interface {
	Prompt(prompt string) (string, error)
	PromptWithStructuredResponse(prompt string, model map[string]any) (string, error)
}
