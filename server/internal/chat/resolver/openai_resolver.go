package resolver

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/lepoai/lepo/server/internal/chat"
	"github.com/sashabaranov/go-openai"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"go.uber.org/zap"
)

type OpenAIChatResolver struct {
	openai   *openai.Client
	weaviate *weaviate.Client
	logger   *zap.Logger
}

func NewOpenAIChatResolver(openai *openai.Client, weaviate *weaviate.Client, logger *zap.Logger) *OpenAIChatResolver {
	return &OpenAIChatResolver{
		openai:   openai,
		weaviate: weaviate,
		logger:   logger,
	}
}

// TODO: Add context cancellation
func (o *OpenAIChatResolver) Resolve(req chat.ChatRequest) (chat.ChatResponse, error) {
	// if conversation is empty, then create a preamble
	if len(req.Conversation) == 0 {
		preamble := openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: intro(),
		}
		req.Conversation = append(req.Conversation, preamble)
	}
	// validate command
	cmd := req.ChatCommand

	switch cmd {
	case chat.Ask:
		{
			fmt.Println("Asking")
			chatResp, err := o.handleAskQuery(req, req.Conversation)
			if err != nil {
				return chat.ChatResponse{}, err
			}
			return chatResp, nil
		}

	case chat.Search:
		{
			fmt.Println("Searching")
		}
	case "":
		{
			fmt.Println("General query")
		}

	default:
		{
			return chat.ChatResponse{}, fmt.Errorf("error: Invalid command")
		}
	}
	// validate query string

	// based on command, populate functions

	// construct completion request
	// receive response
	// if function_call, then perform function call
	// optional(update user about performing function call)
	// feed function call response to llm and check if end_result == none or function call
	// loop
	// append to conversation,
	// construct response
	return chat.ChatResponse{}, nil
}

func intro() string {
	return `ignore all previous instructions.
You are Lepo, an AI powered assistant who can answer queries regarding 'go-lb' codebase. Assume all queries from user is for go-lb codebase. 

You obey the following rules:

- Be cheerful and humble.
- Quote the code snippet that you are referring to, in your response.
- Be concise in your responses and quote your sources.
- When asked about somthing, you always follow up with questions to gain more understanding and context about the user’s query.
- Make use of functions available to you, in order to get context and information.
- Before you write any code or provide any examples, always make sure you have answers to all your follow up questions.
- Think step by step.
- While thinking step by step, you can ask the user about additional infomation and context.
- Always use emojis to convey your tone.
`
}

func (o *OpenAIChatResolver) handleAskQuery(req chat.ChatRequest, conversation []openai.ChatCompletionMessage) (chat.ChatResponse, error) {
	askFunc := &openai.FunctionDefine{
		Name:        "fetch_context",
		Description: "Fetches context about all queries.",
		Parameters: &openai.FunctionParams{
			Type: openai.JSONSchemaTypeObject,
			Properties: map[string]*openai.JSONSchemaDefine{
				"query": {
					Type:        openai.JSONSchemaTypeString,
					Description: "Context free natural language string to query vector database.",
				},
			},
			Required: []string{"query"},
		},
	}
	functions := make([]*openai.FunctionDefine, 0)
	functions = append(functions, askFunc)

	message := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Name:    "humanbeeng",
		Content: req.Query,
	}

	conversation = append(conversation, message)

	res, err := o.openai.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       openai.GPT3Dot5Turbo16K0613,
			Messages:    conversation,
			Functions:   functions,
			Temperature: 0.0,
			// Revisit: function call force ?
			FunctionCall: "auto",
		},
	)
	if err != nil {
		log.Printf("ChatCompletion error: %v\n", err)
	}
	var choice openai.ChatCompletionChoice
	fmt.Printf("First choice %+v\n", choice)
	choice = res.Choices[0]

	if choice.FinishReason != openai.FinishReasonStop {
		if choice.FinishReason == openai.FinishReasonFunctionCall {
			log.Printf("Function call invoked")
			f := choice.Message.FunctionCall
			fmt.Printf("Calling function %v with args %v", f.Name, f.Arguments)

			var queryArgs chat.FetchCodeContextRequest
			err := json.Unmarshal([]byte(f.Arguments), &queryArgs)
			if err != nil {
				log.Println("Unable to parse query args")
			}

			fmt.Printf("%+v\n", queryArgs)

			codeCtx, err := o.fetchCodeContext(queryArgs.Query)
			if err != nil {
				log.Printf("error: Unable to resolve code context\n")
			}

			funcMessage := openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleFunction,
				Content: fmt.Sprintf("Context:\n%+v", codeCtx),
				Name:    f.Name,
			}

			conversation = append(conversation, funcMessage)

			// fxns := make([]*openai.FunctionDefine, 0)
			// fxns = append(fxns, askFunc)
			completionReq := openai.ChatCompletionRequest{
				Model:       openai.GPT3Dot5Turbo16K0613,
				Messages:    conversation,
				Temperature: 0.0,
			}

			fmt.Printf("Function completion request %+v", completionReq)

			funcRes, err := o.openai.CreateChatCompletion(context.Background(), completionReq)
			if err != nil {
				log.Println(err)
			}

			funcChoice := funcRes.Choices[0]

			choice = funcChoice
		} else {
			log.Println("finish reason", choice.FinishReason)
		}
	}
	conversation = append(conversation, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: choice.Message.Content,
		Name:    "lepo",
	})

	return chat.ChatResponse{
		Response:     choice.Message.Content,
		Conversation: conversation,
		RepoID:       "1",
	}, nil
}

func (o *OpenAIChatResolver) fetchCodeContext(query string) ([]chat.CodeContext, error) {
	fields := []graphql.Field{
		{Name: "code"},
		{Name: "language"},
		{Name: "module"},
		{Name: "file"},
		{Name: "codeType"},
	}

	// TODO: Remove this query and move this into a test maybe ?
	nearText := o.weaviate.GraphQL().
		NearTextArgBuilder().
		WithConcepts([]string{query})

	result, err := o.weaviate.GraphQL().Get().
		WithClassName("CodeSnippets").
		WithLimit(10).
		WithNearText(nearText).
		WithFields(fields...).
		Do(context.Background())
	if err != nil {
		return []chat.CodeContext{}, err
	}

	s, err := result.MarshalBinary()
	if err != nil {
		return nil, err
	}

	var resp chat.CodeContextGraphQLResponse
	err = json.Unmarshal(s, &resp)
	if err != nil {
		return nil, err
	}

	ctx := make([]chat.CodeContext, 0)
	for _, v := range resp.Data.Get.CodeSnippets {
		ctx = append(ctx, chat.CodeContext{
			Code:     v.Code,
			CodeType: v.CodeType,
			Module:   v.Module,
			File:     v.File,
			Language: v.Language,
		})
	}
	log.Println("Context fetched")

	return ctx, nil
}
