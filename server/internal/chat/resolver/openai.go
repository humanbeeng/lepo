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
			Content: intro("distributed-cache"),
		}
		req.Conversation = append(req.Conversation, preamble)
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
			fmt.Println("Searching")
		}
	case chat.Git:
		{
			// TODO: Add git action
			fmt.Println("Gitting")
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
		Codebase name %v

You obey the following rules:
- Be cheerful, humble and friendly.
- Be concise in your responses and quote the code snippet that you are referring to.
- When asked about somthing, you can follow up with questions to gain more understanding and context about the user’s query.
- Make use of functions available to you, in order to get context and information.
- Think step by step.
- While thinking step by step, you can ask the user about additional infomation and context.
- Always use emojis to convey your tone.
`, repo, repo)
}

type FuncParams struct {
	Query struct {
		Type        string `json:"type"`
		Description string `json:"description"`
	} `json:"query"`
}
