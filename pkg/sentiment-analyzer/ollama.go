package sentimentanalyzer

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ollama/ollama/api"
)

type OllamaSentimentAnalyzer struct {
	client       *api.Client
	model        string
	analysisChan chan analysisRequest
	ctx          context.Context
	cancel       context.CancelFunc
}

type analysisRequest struct {
	text     string
	source   string
	resultCh chan *SentimentResult
	errCh    chan error
}

func NewOllamaSentimentAnalyzer(ollamaHost string, model string) (*OllamaSentimentAnalyzer, error) {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		return nil, fmt.Errorf("failed to create Ollama client: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	analyzer := &OllamaSentimentAnalyzer{
		client:       client,
		model:        model,
		analysisChan: make(chan analysisRequest, 10),
		ctx:          ctx,
		cancel:       cancel,
	}

	go analyzer.worker()

	return analyzer, nil
}

// worker processes sentiment analysis requests asynchronously
func (a *OllamaSentimentAnalyzer) worker() {
	for {
		select {
		case <-a.ctx.Done():
			return
		case req := <-a.analysisChan:
			result, err := a.analyzeSync(req.text, req.source)
			if err != nil {
				req.errCh <- err
			} else {
				req.resultCh <- result
			}
		}
	}
}

// Analyze performs sentiment analysis on the given text
func (a *OllamaSentimentAnalyzer) Analyze(ctx context.Context, text string, source string) (*SentimentResult, error) {
	if strings.TrimSpace(text) == "" {
		return nil, fmt.Errorf("empty text provided")
	}

	// Create channels for result
	resultCh := make(chan *SentimentResult, 1)
	errCh := make(chan error, 1)

	// Send request to worker
	select {
	case a.analysisChan <- analysisRequest{
		text:     text,
		source:   source,
		resultCh: resultCh,
		errCh:    errCh,
	}:
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	// Wait for result
	select {
	case result := <-resultCh:
		return result, nil
	case err := <-errCh:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (a *OllamaSentimentAnalyzer) analyzeSync(text string, source string) (*SentimentResult, error) {
	prompt := fmt.Sprintf(`Analyze the sentiment of the following text and respond ONLY with a JSON object in this exact format:
{
  "sentiment": "positive" or "negative" or "neutral",
  "score": confidence score between 0.0 and 1.0,
  "emotions": {
    "joy": 0.0-1.0,
    "anger": 0.0-1.0,
    "sadness": 0.0-1.0,
    "fear": 0.0-1.0
  }
}

Text to analyze: "%s"

JSON response:`, text)

	req := &api.GenerateRequest{
		Model:  a.model,
		Prompt: prompt,
		Stream: new(bool),
		Options: map[string]any{
			"temperature": 0.1,
			"num_predict": 200,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var fullResponse strings.Builder
	err := a.client.Generate(ctx, req, func(resp api.GenerateResponse) error {
		fullResponse.WriteString(resp.Response)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("ollama generate error: %w", err)
	}

	responseText := strings.TrimSpace(fullResponse.String())
	jsonStart := strings.Index(responseText, "{")
	jsonEnd := strings.LastIndex(responseText, "}")
	if jsonStart == -1 || jsonEnd == -1 {
		return nil, fmt.Errorf("no valid JSON found in response: %s", responseText)
	}
	jsonText := responseText[jsonStart : jsonEnd+1]

	var parsed struct {
		Sentiment string             `json:"sentiment"`
		Score     float64            `json:"score"`
		Emotions  map[string]float64 `json:"emotions"`
	}

	if err := json.Unmarshal([]byte(jsonText), &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse sentiment response: %w, response: %s", err, jsonText)
	}

	sentiment := strings.ToLower(parsed.Sentiment)
	if sentiment != "positive" && sentiment != "negative" && sentiment != "neutral" {
		sentiment = "neutral"
	}

	return &SentimentResult{
		Text:      text,
		Sentiment: sentiment,
		Score:     parsed.Score,
		Emotions:  parsed.Emotions,
		Timestamp: time.Now(),
		Source:    source,
	}, nil
}

func (a *OllamaSentimentAnalyzer) Close() error {
	a.cancel()
	close(a.analysisChan)
	return nil
}
