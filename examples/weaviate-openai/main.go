package main

import (
	"cmp"
	"fmt"
	"log"
	"net/url"
	"os"

	oai "github.com/openai/openai-go"
	oai_option "github.com/openai/openai-go/option"
	oai_embedder "github.com/suapapa/go_ragkit/embedder/openai"
	"github.com/suapapa/go_ragkit/examples/common"
	weaviate_vstore "github.com/suapapa/go_ragkit/vector_store/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

var (
	oaiSecretKey  = cmp.Or(os.Getenv("OPENAI_SECRET_KEY"), "")
	oaiEmbedModel = cmp.Or(os.Getenv("OPEAI_EMBED_MODEL"), "text-embedding-3-small")
	weaviateAddr  = cmp.Or(os.Getenv("WEAVIATE_ADDR"), "http://localhost:8080")
)

func main() {
	// define embedder
	log.Println("defining embedder")
	log.Println("oaiSecretKey", oaiSecretKey)

	oaiClient := oai.NewClient(
		// oai_option.WithEnvironmentProduction(),
		oai_option.WithAPIKey(oaiSecretKey),
	)
	embedder := oai_embedder.New(&oaiClient, oaiEmbedModel)

	// define vstore
	log.Println("defining vstore")
	weaviateURL, err := url.Parse(weaviateAddr)
	if err != nil {
		panic(err)
	}
	weaviateClient, err := weaviate.NewClient(weaviate.Config{
		Host:   weaviateURL.Host,
		Scheme: weaviateURL.Scheme,
	})
	if err != nil {
		panic(err)
	}
	vstore := weaviate_vstore.New(weaviateClient, "FamilyTree", embedder)
	fmt.Printf("vstore: %s\n", vstore)

	common.IndexAndRetriveExample(vstore)
}
