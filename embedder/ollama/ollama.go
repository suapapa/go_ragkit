package ollama

import (
	"context"

	ollama_api "github.com/ollama/ollama/api"
	ragkit "github.com/suapapa/go_ragkit"
)

var _ ragkit.Embeder = &Ollama{}

type Ollama struct {
	client *ollama_api.Client
	model  string
}

func NewOllama(client *ollama_api.Client, model string) *Ollama {
	return &Ollama{
		client: client,
		model:  model,
	}
}

func (o *Ollama) Embed(ctx context.Context, texts ...string) ([][]float32, error) {
	embeddings := make([][]float32, len(texts))

	for i, text := range texts {
		req := &ollama_api.EmbedRequest{
			Model: o.model,
			Input: text,
		}

		resp, err := o.client.Embed(ctx, req)
		if err != nil {
			return nil, err
		}

		// Ollama can return multiple embeddings, but we only need one embedding per text
		if len(resp.Embeddings) > 0 {
			embeddings[i] = resp.Embeddings[0]
		}
	}

	return embeddings, nil
}

// Dimension returns the dimension of the embedding vector
// Ollama embedding models typically use 4096 dimensions
// func (o *Ollama) Dimension() int {
// 	return 4096
// }
