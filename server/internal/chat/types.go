package chat

import (
	"github.com/lepoai/lepo/server/internal/sync/extract"
	"github.com/sashabaranov/go-openai"
)

type ChatResolver interface {
	Resolve(req ChatRequest) (ChatResponse, error)
}

type ChatRequest struct {
	// TODO: Add other metadata once made available
	Query        string                         `json:"query" validate:"required"`
	Conversation []openai.ChatCompletionMessage `json:"conversation" validate:"required"`
	ChatCommand  ChatCommand                    `json:"chat_command"`
	RepoID       string                         `json:"repo_id" validate:"required"`
}

type ChatResponse struct {
	Response     string                         `json:"response"`
	Conversation []openai.ChatCompletionMessage `json:"conversation"`
	RepoID       string                         `json:"repo_id"`
	// RepoID string
	// TODO: Add Source[] once made available
}

type CodeContext struct {
	Code     string            `json:"code"`
	Language extract.Language  `json:"language"`
	Module   string            `json:"module"`
	File     string            `json:"file"`
	CodeType extract.ChunkType `json:"codeType"`
}

type FetchCodeContextRequest struct {
	Query string `json:"query"`
}

type CodeContextGraphQLResponse struct {
	Data struct {
		Get struct {
			CodeSnippets []struct {
				Code     string            `json:"code"`
				CodeType extract.ChunkType `json:"codeType"`
				File     string            `json:"file"`
				Module   string            `json:"module"`
				Language extract.Language  `json:"language"`
			} `json:"CodeSnippets"`
		} `json:"Get"`
	} `json:"Data"`
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
)
