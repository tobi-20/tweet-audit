package model

type ModelResponse struct {
	Decision string `json:"decision"`
}

type AIClient interface {
	Analyze(prompt string) (string, error)
}
