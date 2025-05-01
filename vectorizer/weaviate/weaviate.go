package weviate

import (
	"context"

	ragkit "github.com/suapapa/go_ragkit"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
)

var _ ragkit.Vectorizer = &Weaviate{}

type Weaviate struct {
	className string
	client    *weaviate.Client
	embedder  ragkit.Embedder
}

func NewWeaviate(client *weaviate.Client, className string, embedder ragkit.Embedder) *Weaviate {
	return &Weaviate{
		className: className,
		client:    client,
		embedder:  embedder,
	}
}

func (w *Weaviate) Index(ctx context.Context, docs []ragkit.Document) ([]string, error) {
	var ids []string
	for _, doc := range docs {
		var embedding []float32
		var err error

		if doc.ID == "" {
			doc.ID = ragkit.GenerateID(doc.Text)
		}

		// Use provided vector if available, otherwise generate one
		if doc.Vector != nil {
			embedding = doc.Vector
		} else {
			vectors, err := w.embedder.Embed(ctx, doc.Text)
			if err != nil {
				return nil, err
			}
			embedding = vectors[0]
		}

		// Create data object
		data := map[string]interface{}{
			"text":     doc.Text,
			"metadata": doc.Metadata,
		}

		// Create object with embedding
		_, err = w.client.Data().Creator().
			WithClassName(w.className).
			WithID(doc.ID).
			WithProperties(data).
			WithVector(embedding).
			Do(ctx)
		if err != nil {
			return nil, err
		}
		ids = append(ids, doc.ID)
	}
	return ids, nil
}

func (w *Weaviate) Delete(ctx context.Context, id string) error {
	return w.client.Data().Deleter().
		WithClassName(w.className).
		WithID(id).
		Do(ctx)
}

func (w *Weaviate) Exists(ctx context.Context, id string) (bool, error) {
	_, err := w.client.Data().ObjectsGetter().
		WithClassName(w.className).
		WithID(id).
		Do(ctx)
	if err != nil {
		if err.Error() == "object not found" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (w *Weaviate) Retrieve(ctx context.Context, query []float32, topK int) ([]ragkit.RetrievedDoc, error) {
	// Get near vector results
	response, err := w.client.GraphQL().Get().
		WithClassName(w.className).
		WithFields(graphql.Field{Name: "text"}, graphql.Field{Name: "metadata"}).
		WithNearVector(w.client.GraphQL().NearVectorArgBuilder().
			WithVector(query)).
		WithLimit(topK).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	// Parse results
	var results []ragkit.RetrievedDoc
	for _, obj := range response.Data["Get"].(map[string]interface{})["Document"].([]interface{}) {
		objMap := obj.(map[string]interface{})
		results = append(results, ragkit.RetrievedDoc{
			ID:       objMap["_id"].(string),
			Vector:   query, // Weaviate doesn't return the vector in the response
			Score:    objMap["_additional"].(map[string]interface{})["distance"].(float32),
			Metadata: objMap["metadata"].(map[string]interface{}),
		})
	}

	return results, nil
}

func (w *Weaviate) RetrieveText(ctx context.Context, text string, topK int) ([]ragkit.RetrievedDoc, error) {
	vectors, err := w.embedder.Embed(ctx, text)
	if err != nil {
		return nil, err
	}

	return w.Retrieve(ctx, vectors[0], topK)
} // Get near text results
