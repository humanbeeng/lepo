package function

import (
	"context"
	"encoding/json"
	"log"

	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
)

type CodeContextFunction struct {
	weaviate *weaviate.Client
}

func NewCodeContextFunction(weaviate *weaviate.Client) *CodeContextFunction {
	return &CodeContextFunction{weaviate: weaviate}
}

func (f *CodeContextFunction) FetchCodeContext(req FetchCodeContextRequest) ([]CodeContext, error) {
	fields := []graphql.Field{
		{Name: "code"},
		{Name: "language"},
		{Name: "module"},
		{Name: "file"},
		{Name: "codeType"},
	}

	// TODO: Remove this query and move this into a test maybe ?
	nearText := f.weaviate.GraphQL().
		NearTextArgBuilder().
		WithConcepts([]string{req.Query})

	result, err := f.weaviate.GraphQL().Get().
		WithClassName("CodeSnippets").
		WithLimit(10).
		WithNearText(nearText).
		WithFields(fields...).
		Do(context.Background())
	if err != nil {
		return []CodeContext{}, err
	}

	s, err := result.MarshalBinary()
	if err != nil {
		return nil, err
	}

	var resp CodeContextGraphQLResponse
	err = json.Unmarshal(s, &resp)
	if err != nil {
		return nil, err
	}

	ctx := make([]CodeContext, 0)
	for _, v := range resp.Data.Get.CodeSnippets {
		ctx = append(ctx, CodeContext{
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

func (f *CodeContextFunction) NewCodeContextFunctionDefinition() openai.FunctionDefinition {
	// create openai function call params
	params := jsonschema.Definition{
		Type: jsonschema.Object,
		Properties: map[string]jsonschema.Definition{
			"query": {
				Type:        jsonschema.String,
				Description: "context free natural language query string",
			},
		},
		Required: []string{"query"},
	}

	askFunc := openai.FunctionDefinition{
		Name:        "get_code_context",
		Description: "Fetches code snippets, context for a given query",
		Parameters:  params,
	}
	return askFunc
}
