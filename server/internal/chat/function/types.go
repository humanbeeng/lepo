package function

import (
	"github.com/humanbeeng/lepo/server/internal/sync/extract"
)

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
