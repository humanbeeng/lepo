package resolver

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/humanbeeng/lepo/server/internal/chat"
	"github.com/humanbeeng/lepo/server/internal/chat/function"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

func (o *OpenAIChatResolver) handleGeneralQuery(
	req chat.ChatRequest,
	conversation []openai.ChatCompletionMessage,
) (chat.ChatResponse, error) {
	ctxFunc := o.codeCtxFunc.NewCodeContextFunctionDefinition()

	functions := make([]openai.FunctionDefinition, 0)
	functions = append(functions, ctxFunc)

	message := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Name:    "humanbeeng",
		Content: req.Query,
	}

	conversation = append(conversation, message)

	initReq := openai.ChatCompletionRequest{
		Model:       openai.GPT3Dot5Turbo16K0613,
		Messages:    conversation,
		Functions:   functions,
		Temperature: 0.1,
		// Revisit: function call force ?
		FunctionCall: "auto",
	}

	initRes, err := o.openai.CreateChatCompletion(
		context.Background(),
		initReq,
	)
	if err != nil {
		o.logger.Error("ChatCompletion error", zap.Error(err))
		return createErrorResponse(req.Conversation, req.RepoID), err
	}

	if len(initRes.Choices) == 0 {
		return createErrorResponse(req.Conversation, req.RepoID), err
	}

	choice := initRes.Choices[0]

	if choice.FinishReason == openai.FinishReasonFunctionCall {
		f := choice.Message.FunctionCall

		var args function.FetchCodeContextRequest

		err := json.Unmarshal([]byte(f.Arguments), &args)
		if err != nil {
			o.logger.Error("Unable to parse query args", zap.Error(err))

			return createErrorResponse(conversation, req.RepoID), err
		}

		o.logger.Info("Fetching context with args", zap.Any("args", args))

		codeCtx, err := o.codeCtxFunc.FetchCodeContext(args)
		if err != nil {
			o.logger.Error("Unable to resolve code context", zap.Error(err))
		}
		o.logger.Info("Number of snippets fetched", zap.Int("count", len(codeCtx)))

		funcCtxMsg := openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleFunction,
			Content: fmt.Sprintf("Context:\n%+v", codeCtx),
			Name:    ctxFunc.Name,
		}

		conversation = append(conversation, funcCtxMsg)

		funcCtxReq := openai.ChatCompletionRequest{
			Model:       openai.GPT3Dot5Turbo16K0613,
			Messages:    conversation,
			Functions:   functions,
			Temperature: 0.2,
		}

		funcCtxRes, err := o.openai.CreateChatCompletion(context.Background(), funcCtxReq)
		if err != nil || len(funcCtxRes.Choices) == 0 {
			return createErrorResponse(conversation, req.RepoID), err
		}

		funcCtxResChoice := funcCtxRes.Choices[0]

		conversation = append(conversation, funcCtxResChoice.Message)

		m := make(map[string]bool, 0)

		for _, s := range codeCtx {
			_, exists := m[s.File]
			if !exists {
				m[s.File] = true
			}
		}
		sources := make([]string, 0)

		for k := range m {
			sources = append(sources, k)
		}

		return chat.ChatResponse{
			Response:     funcCtxResChoice.Message.Content,
			Conversation: conversation,
			RepoID:       req.RepoID,
			Sources:      sources,
		}, nil
	}

	conversation = append(conversation, choice.Message)

	return chat.ChatResponse{
		Response:     choice.Message.Content,
		Conversation: conversation,
		RepoID:       req.RepoID,
	}, nil
}

func createErrorResponse(
	conversation []openai.ChatCompletionMessage,
	repoID string,
) chat.ChatResponse {
	return chat.ChatResponse{
		Response:     "Something went wrong 😢. Please try again 🙇.",
		Conversation: conversation,
		RepoID:       repoID,
	}
}
