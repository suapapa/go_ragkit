package helper

import (
	"cmp"
	"log"
	"net/http"
	"net/url"
	"time"

	ollama_api "github.com/ollama/ollama/api"
	ollama_embedder "github.com/suapapa/go_ragkit/embedder/ollama"
	weaviate_vstore "github.com/suapapa/go_ragkit/vector_store/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"

	ragkit "github.com/suapapa/go_ragkit"
)

func NewWeaviateOllamaVectorStore(
	vectorDBClassName string,
	ollamaEmbedModel string,
) (ragkit.VectorStore, error) {
	// initialize ollama
	ollamaURL, err := url.Parse(ollamaAddr)
	if err != nil {
		return nil, err
	}
	ollamaClient := ollama_api.NewClient(ollamaURL, http.DefaultClient)
	embedder := ollama_embedder.New(ollamaClient, cmp.Or(ollamaEmbedModel, defaultOllamaEmbedModel))

	// initialize weaviate
	weaviateURL, err := url.Parse(weaviateAddr)
	if err != nil {
		return nil, err
	}
	weaviateClient, err := weaviate.NewClient(weaviate.Config{
		Host:           weaviateURL.Host,
		Scheme:         weaviateURL.Scheme,
		StartupTimeout: 3 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	// define vector_store
	log.Println("defining vector_store")
	return weaviate_vstore.New(weaviateClient, vectorDBClassName, embedder), nil
}
