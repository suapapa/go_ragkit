package weviate

import (
	"context"

	ragkit "github.com/suapapa/go_ragkit"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
)

var _ ragkit.IndexerRetriever = &Weaviate{}

type Weaviate struct {
	client   *weaviate.Client
	embedder ragkit.Embeder
}

func NewWeaviate(client *weaviate.Client, embedder ragkit.Embeder) *Weaviate {
	return &Weaviate{
		client:   client,
		embedder: embedder,
	}
}

func (w *Weaviate) Index(ctx context.Context, docs []ragkit.Document) error {
	for _, doc := range docs {
		var embedding []float32
		var err error

		// Use provided vector if available, otherwise generate one
		if doc.Vector != nil {
			embedding = doc.Vector
		} else {
			vectors, err := w.embedder.Embed(ctx, doc.Text)
			if err != nil {
				return err
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
			WithClassName("Document").
			WithID(doc.ID).
			WithProperties(data).
			WithVector(embedding).
			Do(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *Weaviate) Delete(ctx context.Context, id string) error {
	return w.client.Data().Deleter().
		WithClassName("Document").
		WithID(id).
		Do(ctx)
}

func (w *Weaviate) Exists(ctx context.Context, id string) (bool, error) {
	_, err := w.client.Data().ObjectsGetter().
		WithClassName("Document").
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
		WithClassName("Document").
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
