package weviate

import (
	"context"
	"fmt"
	"strings"

	ragkit "github.com/suapapa/go_ragkit"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
)

var _ ragkit.VectorStore = &Weaviate{}

type Weaviate struct {
	className string
	client    *weaviate.Client
	embedder  ragkit.Embedder
}

func New(client *weaviate.Client, className string, embedder ragkit.Embedder) *Weaviate {
	return &Weaviate{
		className: ragkit.ToCamelCase(className),
		client:    client,
		embedder:  embedder,
	}
}

func (w *Weaviate) Index(ctx context.Context, docs ...ragkit.Document) ([]string, error) {
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
			vectors, err := w.embedder.EmbedTexts(ctx, doc.Text)
			if err != nil {
				return nil, err
			}
			embedding = vectors[0]
		}

		// Create data object
		data := map[string]any{
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
		if err.Error() == "object not found" || strings.Contains(err.Error(), "404") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (w *Weaviate) Retrieve(ctx context.Context, query []float32, topK int, metadataFieldNames ...string) ([]ragkit.RetrievedDoc, error) {
	var metadataFields []graphql.Field
	for _, name := range metadataFieldNames {
		metadataFields = append(metadataFields, graphql.Field{Name: name})
	}

	// Get near vector results
	response, err := w.client.GraphQL().Get().
		WithClassName(w.className).
		WithFields(
			graphql.Field{Name: "text"},
			graphql.Field{
				Name:   "metadata",
				Fields: metadataFields,
			},
		).
		WithNearVector(w.client.GraphQL().NearVectorArgBuilder().
			WithVector(query)).
		WithLimit(topK).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	// Parse results
	var results []ragkit.RetrievedDoc
	if resultData, ok := response.Data["Get"].(map[string]any); ok {
		if classResultData, ok := resultData[w.className].([]any); ok {
			for _, obj := range classResultData {
				objMap := obj.(map[string]any)

				var metadata map[string]any
				if m, ok := objMap["metadata"].(map[string]any); ok {
					metadata = m
				}

				results = append(results, ragkit.RetrievedDoc{
					// ID:       objMap["_id"].(string),
					// Score:    objMap["_additional"].(map[string]any)["distance"].(float32),
					Vector:   query, // Weaviate doesn't return the vector in the response
					Text:     objMap["text"].(string),
					Metadata: metadata,
				})
			}
		} else {
			return nil, fmt.Errorf("no results found in class: %s", w.className)
		}
	} else {
		return nil, fmt.Errorf("no results found")
	}

	return results, nil
}

func (w *Weaviate) RetrieveText(ctx context.Context, text string, topK int, metadataFieldNames ...string) ([]ragkit.RetrievedDoc, error) {
	vectors, err := w.embedder.EmbedTexts(ctx, text)
	if err != nil {
		return nil, err
	}

	return w.Retrieve(ctx, vectors[0], topK, metadataFieldNames...)
} // Get near text results

func (w *Weaviate) String() string {
	return fmt.Sprintf("Weaviate(class: %s, embedder: %s)", w.className, w.embedder)
}
