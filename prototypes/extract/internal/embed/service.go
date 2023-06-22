package embed

import "github.com/humanbeeng/lepo/prototypes/extract/internal/extract"

type Embedding struct{}

func (e *Embedding) Embed(chunk extract.Chunk) error {
	return nil
}
