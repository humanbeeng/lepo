package resolver

import (
	"fmt"

	"github.com/humanbeeng/lepo/server/internal/chat"
	"github.com/humanbeeng/lepo/server/internal/chat/function"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

type OpenAIChatResolver struct {
	openai      *openai.Client
	codeCtxFunc *function.CodeContextFunction
	logger      *zap.Logger
}

func NewOpenAIChatResolver(
	openai *openai.Client,
	codeCtxFunc *function.CodeContextFunction,
	logger *zap.Logger,
) *OpenAIChatResolver {
	return &OpenAIChatResolver{
		openai:      openai,
		codeCtxFunc: codeCtxFunc,
		logger:      logger,
	}
}

// TODO: Add context cancellation
func (o *OpenAIChatResolver) Resolve(req chat.ChatRequest) (chat.ChatResponse, error) {
	// if conversation is empty, then create a preamble
	if len(req.Conversation) == 0 {
		preamble := openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: intro("kyc-reports"),
		}
		req.Conversation = append(req.Conversation, preamble)
		req.Conversation = append(req.Conversation, openai.ChatCompletionMessage{
			Role: openai.ChatMessageRoleAssistant,
			Content: `Hey ! I'm Lepo. I'm here to answer all your queries regarding 'lepo-server' codebase.
			Please feel free to ask me about anything you want to know.`,
		})
		fmt.Println(preamble.Content)
	}

	// TODO: Add command validation

	switch req.ChatCommand {
	case chat.Ask:
		{
			// chatResp, err := o.handleAskQuery(req, req.Conversation)
		}

	case chat.Search:
		{
			// TODO: Add search
			o.logger.Info("Searching")
		}
	case chat.Git:
		{
			// TODO: Add git action
			o.logger.Info("Gitting")
		}
	case "":
		{
			chatResp, err := o.handleGeneralQuery(req, req.Conversation)
			if err != nil {
				return chat.ChatResponse{}, err
			}
			return chatResp, nil
		}

	}
	return chat.ChatResponse{
			Response:     "Invalid command. Please try again.",
			Conversation: req.Conversation,
			RepoID:       req.RepoID,
		}, fmt.Errorf(
			"error: Invalid command",
		)
}

func intro(repo string) string {
	return fmt.Sprintf(`ignore all previous instructions.
You are Lepo, an AI powered chatbot who can answer queries related to %v codebase and engage in a frienly conversation. 
Codebase name :%v

You obey the following rules:
- Be as brief and concise as possible without losing clarity. 
- Make extensive use of the functions available to you.
- Think step by step.
- Have a personality.
- Use emojis in your responses.
`, repo, repo)
}

type FuncParams struct {
	Query struct {
		Type        string `json:"type"`
		Description string `json:"description"`
	} `json:"query"`
}
