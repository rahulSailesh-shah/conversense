package sentimentanalyzer

import (
	"context"
	"fmt"
	"time"
)

type SentimentResult struct {
	Text      string             `json:"text"`
	Sentiment string             `json:"sentiment"`
	Score     float64            `json:"score"`
	Emotions  map[string]float64 `json:"emotions"`
	Timestamp time.Time          `json:"timestamp"`
	Source    string             `json:"source"`
}

type SentimentAnalyzer interface {
	Analyze(ctx context.Context, text string, source string) (*SentimentResult, error)
	Close() error
}

type SentimentAnalyzerType string

const (
	AnalyzerTypeOllama SentimentAnalyzerType = "ollama"
	AnalyzerTypeOpenAI SentimentAnalyzerType = "openai"
	AnalyzerTypeClaude SentimentAnalyzerType = "claude"
)

func NewSentimentAnalyzer(analyzerType SentimentAnalyzerType) (SentimentAnalyzer, error) {
	switch analyzerType {
	case AnalyzerTypeOllama:
		ollamaHost := "http://localhost:11434"
		model := "llama3.2:3b"
		return NewOllamaSentimentAnalyzer(ollamaHost, model)

	case AnalyzerTypeOpenAI:
		return nil, fmt.Errorf("openAI sentiment analyzer not yet implemented")

	case AnalyzerTypeClaude:
		return nil, fmt.Errorf("claude sentiment analyzer not yet implemented")

	default:
		return nil, fmt.Errorf("unknown analyzer type: %s", analyzerType)
	}
}
