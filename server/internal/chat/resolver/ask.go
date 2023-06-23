package resolver

import (
	"github.com/humanbeeng/lepo/server/internal/chat"
	"github.com/sashabaranov/go-openai"
)

// TODO: Implement focused ask (Inline selection/Code Lens)
func (o *OpenAIChatResolver) handleAskQuery(
	req chat.ChatRequest,
	conversation []openai.ChatCompletionMessage,
) (chat.ChatResponse, error) {
	return chat.ChatResponse{}, nil
}
