package chat

import (
	"errors"

	"github.com/sashabaranov/go-openai"
)

type ChatResolver interface {
	Resolve(req ChatRequest) (ChatResponse, error)
}

type ChatRequest struct {
	// Note: Add other metadata once made available
	Query        string                         `json:"query"        validate:"required"`
	Conversation []openai.ChatCompletionMessage `json:"conversation" validate:"required"`
	ChatCommand  ChatCommand                    `json:"chat_command"`
	RepoID       string                         `json:"repo_id"      validate:"required"`
}

type ChatResponse struct {
	Response     string                         `json:"response"`
	Conversation []openai.ChatCompletionMessage `json:"conversation"`
	RepoID       string                         `json:"repo_id"`
	// TODO: Add Source[] once made available
}

type FunctionCall struct {
	Name string `json:"name"`
}

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

type (
	ChatCommand string
	Role        string
)

const (
	User Role = "user"
	Lepo Role = "lepo"
)

const (
	Ask    ChatCommand = "/ask"
	Search ChatCommand = "/search"
	Git    ChatCommand = "/git"
)

var (
	InvalidCommand error = errors.New("Invalid command")
	OpenAIError    error = errors.New("Something went wrong while calling OpenAI")
)
