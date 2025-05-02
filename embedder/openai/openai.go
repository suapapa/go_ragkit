package openai

import (
	"context"
	"fmt"

	oai "github.com/openai/openai-go"
	ragkit "github.com/suapapa/go_ragkit"
)

var _ ragkit.Embedder = &OpenAI{}

type OpenAI struct {
	client oai.Client
	model  string
}

func New(client oai.Client, model string) *OpenAI {
	return &OpenAI{
		client: client,
		model:  model,
	}
}

func (o *OpenAI) String() string {
	return fmt.Sprintf("OpenAI(%s)", o.model)
}

func (o *OpenAI) EmbedText(ctx context.Context, text string) ([]float32, error) {
	embeddings, err := o.EmbedTexts(ctx, text)
	if err != nil {
		return nil, err
	}
	return embeddings[0], nil
}

func (o *OpenAI) EmbedTexts(ctx context.Context, texts ...string) ([][]float32, error) {
	embeddings := make([][]float32, len(texts))
	resp, err := o.client.Embeddings.New(ctx, oai.EmbeddingNewParams{
		Model: o.model,
		Input: oai.EmbeddingNewParamsInputUnion{
			OfArrayOfStrings: texts,
		},
	})
	if err != nil {
		return nil, err
	}

	for i, embedding := range resp.Data {
		// Convert []float64 to []float32
		embeddings[i] = make([]float32, len(embedding.Embedding))
		for j, v := range embedding.Embedding {
			embeddings[i][j] = float32(v)
		}
	}

	return embeddings, nil
}
