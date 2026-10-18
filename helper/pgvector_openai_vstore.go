package helper

import (
	"cmp"
	"fmt"

	oai "github.com/openai/openai-go"
	oai_option "github.com/openai/openai-go/option"
	oai_embedder "github.com/suapapa/go_ragkit/embedder/openai"
	pgvector_vstore "github.com/suapapa/go_ragkit/vector_store/pgvector"

	ragkit "github.com/suapapa/go_ragkit"
)

func NewPGVectorOpenAIVectorStore(
	vectorDBClassName string,
	oaiEmbedModel string,
	oaiOptions ...oai_option.RequestOption,
) (ragkit.VectorStore, error) {
	connStr := pgvectorConn()
	if connStr == "" {
		return nil, fmt.Errorf("PGVECTOR_CONN environment variable is not set")
	}

	// initialize openai (OPENAI_ENDPOINT / OPENAI_APIKEY / OPENAI_EMBED_MODEL)
	opts := []oai_option.RequestOption{
		oai_option.WithAPIKey(oaiApiKey()),
	}
	if ep := oaiEndpoint(); ep != "" {
		opts = append(opts, oai_option.WithBaseURL(ep))
	}
	opts = append(opts, oaiOptions...)

	oaiClient := oai.NewClient(opts...)
	embedder := oai_embedder.New(&oaiClient, cmp.Or(oaiEmbedModel, DefaultOAIEmbedModel()))

	return pgvector_vstore.New(connStr, DefaultOAIEmbedDim(), vectorDBClassName, embedder), nil
}
