package helper

import (
	"cmp"
	"log"
	"net/url"
	"time"

	oai "github.com/openai/openai-go"
	oai_option "github.com/openai/openai-go/option"
	oai_embedder "github.com/suapapa/go_ragkit/embedder/openai"
	weaviate_vstore "github.com/suapapa/go_ragkit/vector_store/weaviate"
	"github.com/weaviate/weaviate-go-client/v5/weaviate"

	ragkit "github.com/suapapa/go_ragkit"
)

func NewWeaviateOpenAIVectorStore(
	vectorDBClassName string,
	oaiEmbedModel string,
	oaiOptions ...oai_option.RequestOption,
) (ragkit.VectorStore, error) {
	// initialize openai
	opts := []oai_option.RequestOption{
		oai_option.WithAPIKey(oaiApiKey),
	}
	opts = append(opts, oaiOptions...)

	oaiClient := oai.NewClient(opts...)
	embedder := oai_embedder.New(&oaiClient, cmp.Or(oaiEmbedModel, DefaultOAIEmbedModel))

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
