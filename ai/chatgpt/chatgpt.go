package chatgpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/luizhenriquees/go-http-rpa/entity"
)

const (
	openAIEndpoint = "https://api.openai.com/v1/chat/completions"
	openAIModel    = "gpt-4.1-mini"
)

// Request represents the request structure for OpenAI API
type Request struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Temperature float64       `json:"temperature"`
}

// ChatMessage represents a message in the ChatGPT conversation
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Response represents the response structure from OpenAI API
type Response struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type Agent struct {
	apiKey string
}

func NewAgent() *Agent {
	return &Agent{
		apiKey: os.Getenv("CHATGPT_API_KEY"),
	}
}

// GetAnswerIndex sends the question to ChatGPT and gets the answer index
func (a *Agent) GetAnswerIndex(question entity.Question) (int, error) {
	if a.apiKey == "" {
		return 0, fmt.Errorf("CHATGPT_API_KEY environment variable not set")
	}

	prompt := a.buildPrompt(question)
	fmt.Println("Sending question to AI model...")

	reqBody := Request{
		Model: openAIModel,
		Messages: []ChatMessage{
			{
				Role:    "system",
				Content: "You are an assistant that helps answer multiple-choice questions. Reply ONLY with the number (index) of the correct option, nothing else.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: 0.3,
	}

	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return 0, fmt.Errorf("error marshalling request: %w", err)
	}

	req, err := http.NewRequest("POST", openAIEndpoint, bytes.NewBuffer(reqJSON))
	if err != nil {
		return 0, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+a.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error sending request to OpenAI: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("OpenAI API returned error: %s", string(body))
	}

	var chatResponse Response
	if err := json.Unmarshal(body, &chatResponse); err != nil {
		return 0, fmt.Errorf("error parsing response: %w", err)
	}

	if len(chatResponse.Choices) == 0 {
		return 0, fmt.Errorf("no answer choices returned from OpenAI")
	}

	// Extract the answer index from response
	answerContent := chatResponse.Choices[0].Message.Content
	fmt.Printf("AI response: %s\n", answerContent)

	answerIndex, err := strconv.Atoi(answerContent)
	if err != nil {
		for _, char := range answerContent {
			if char >= '0' && char <= '9' {
				answerIndex, _ = strconv.Atoi(string(char))
				break
			}
		}
	}

	if answerIndex < 0 || answerIndex >= len(question.Options) {
		fmt.Printf("Warning: AI returned invalid index %d, defaulting to 0\n", answerIndex)
		answerIndex = 0
	}
	return answerIndex, nil
}

// buildPrompt creates a prompt for the AI with the question and options
func (a *Agent) buildPrompt(question entity.Question) string {
	prompt := fmt.Sprintf("Question: %s\n\nOptions:\n", question.Question)

	for i, option := range question.Options {
		prompt += fmt.Sprintf("%d: %s\n", i, option)
	}

	prompt += "\nRespond only with the number of the correct option."
	return prompt
}
