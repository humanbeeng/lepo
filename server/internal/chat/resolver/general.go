package resolver

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/humanbeeng/lepo/server/internal/chat"
	"github.com/humanbeeng/lepo/server/internal/chat/function"
	"github.com/sashabaranov/go-openai"
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

	s, _ := json.Marshal(&chat.FunctionCall{Name: ctxFunc.Name})

	fmt.Printf("Unmarshalled func call %s", s)

	funcReq := openai.ChatCompletionRequest{
		Model:       openai.GPT3Dot5Turbo16K0613,
		Messages:    conversation,
		Functions:   functions,
		Temperature: 0.1,
		// Revisit: function call force ?
		FunctionCall: "auto",
	}

	fmt.Printf("%+v\n", funcReq)

	res, err := o.openai.CreateChatCompletion(
		context.Background(),
		funcReq,
	)

	fmt.Printf("Func call first response %+v\n", res)

	fmt.Printf("%v\n", res.Choices[0].Message.FunctionCall)

	if err != nil {
		log.Printf("ChatCompletion error: %v\n", err)
		return chat.ChatResponse{Response: "error", Conversation: conversation, RepoID: "1"}, err
	}
	if len(res.Choices) == 0 {
		return chat.ChatResponse{
			Response:     "Something went wrong. Please try again",
			Conversation: req.Conversation,
			RepoID:       req.RepoID,
		}, nil
	}
	choice := res.Choices[0]

	if choice.FinishReason == openai.FinishReasonFunctionCall {
		f := choice.Message.FunctionCall

		var args function.FetchCodeContextRequest
		err := json.Unmarshal([]byte(f.Arguments), &args)
		if err != nil {
			log.Println("Unable to parse query args")

			return chat.ChatResponse{
				Response:     "error",
				Conversation: conversation,
				RepoID:       req.RepoID,
			}, err
		}

		codeCtx, err := o.codeCtxFunc.FetchCodeContext(args)
		if err != nil {
			log.Printf("error: Unable to resolve code context\n")
		}

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
			return chat.ChatResponse{
				Response:     "Something went wrong. Please try again",
				Conversation: conversation,
				RepoID:       req.RepoID,
			}, nil
		}

		funcCtxResChoice := funcCtxRes.Choices[0]

		conversation = append(conversation, funcCtxResChoice.Message)

		return chat.ChatResponse{
			Response:     funcCtxResChoice.Message.Content,
			Conversation: conversation,
			RepoID:       req.RepoID,
		}, nil
	}

	conversation = append(conversation, choice.Message)

	return chat.ChatResponse{
		Response:     choice.Message.Content,
		Conversation: conversation,
		RepoID:       req.RepoID,
	}, nil
}
