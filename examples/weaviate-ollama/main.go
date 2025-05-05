package main

import (
	"cmp"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	ollama_api "github.com/ollama/ollama/api"
	ollama_embedder "github.com/suapapa/go_ragkit/embedder/ollama"
	"github.com/suapapa/go_ragkit/examples/common"
	weaviate_vstore "github.com/suapapa/go_ragkit/vector_store/weaviate"
	"github.com/weaviate/weaviate-go-client/v5/weaviate"
)

var (
	ollamaAddr       = cmp.Or(os.Getenv("OLLAMA_ADDR"), "http://localhost:11434")
	ollamaEmbedModel = cmp.Or(os.Getenv("OLLAMA_EMBED_MODEL"), "bge-m3:latest")
	weaviateAddr     = cmp.Or(os.Getenv("WEAVIATE_ADDR"), "http://localhost:8080")
)

func main() {
	// define embedder
	log.Println("defining embedder")
	ollamaURL, err := url.Parse(ollamaAddr)
	if err != nil {
		panic(err)
	}
	ollamaClient := ollama_api.NewClient(ollamaURL, http.DefaultClient)
	embedder := ollama_embedder.New(ollamaClient, ollamaEmbedModel)

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
