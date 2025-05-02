# ragkit: A Go package for document indexing and retrieval in RAG systems

![ragkit_logo](_asset/ragkit_logo_256.webp)

ragkit is a Go package designed to simplify the implementation of Retrieval-Augmented Generation (RAG) systems.
It includes the definition and implementation of a VectorStore interface that performs document indexing and retrieval,
providing tools for vectorization and semantic search capabilities.

- [Package documentation](https://pkg.go.dev/github.com/suapapa/go_ragkit).

## Installation

```sh
go get github.com/suapapa/go_ragkit
```

## Quick Start

```go
import (
    // ...
    ragkit "github.com/suapapa/go_ragkit"
	vstore_helper "github.com/suapapa/go_ragkit/vector_store/weaviate/helper"
)

func main() {
    // Initialize vector store (Weaviate + Ollama)
    vstore, err := vstore_helper.NewWeaviateOllamaVectorStore(
        "DoolyFamilyTree", // vector DB class name
        DefaultOllamaEmbedModel, 
    )
    if err != nil {
        panic(err)
    }

    // Create documents from text
    docs := ragkit.MakeDocsFromTexts(
        []string{
            "고길동의 집에는 둘리, 도우너, 또치, 희동이, 철수, 영희가 살고 있다.",
            "희동이는 고길동의 조카이다.",
        },
        nil,
    )

    // Index documents
    ctx := context.Background()
    for _, doc := range docs {
        _, err = vstore.Index(ctx, doc)
        if err != nil {
            panic(err)
        }
    }

    // Perform semantic search
    query := "고길동과 희동이의 관계?"
    results, err := vstore.RetrieveText(ctx, query, 1)
    if err != nil {
        panic(err)
    }
    fmt.Println(results[0].Text) // 희동이는 고길동의 조카이다.
}
```

## Examples

Pre-requirement - launch Weaviate for local vector DB:
```
docker run -it --rm -p 8080:8080 -p 50051:50051 cr.weaviate.io/semitechnologies/weaviate:1.30.2
```

- [Examples](examples/)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
