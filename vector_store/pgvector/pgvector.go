package pgvector

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/pgvector/pgvector-go"
	ragkit "github.com/suapapa/go_ragkit"
)

var _ ragkit.VectorStore = &PGVector{}

type PGVector struct {
	className string
	conn      *pgx.Conn
	embedder  ragkit.Embedder
	dimension int

	mu sync.Mutex
}

func New(connStr string, dimension int, className string, embedder ragkit.Embedder) *PGVector {
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	// defer conn.Close(context.Background())

	ret := &PGVector{
		className: className,
		conn:      conn,
		embedder:  embedder,
		dimension: dimension,
	}

	err = ret.ensureTable(context.Background())
	if err != nil {
		log.Fatalf("Failed to ensure table: %v", err)
	}

	return ret
}

func (p *PGVector) Close() error {
	return p.conn.Close(context.Background())
}

// ensureTable creates the table if it doesn't exist
func (p *PGVector) ensureTable(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Create the pgvector extension if it doesn't exist
	_, err := p.conn.Exec(ctx, `CREATE EXTENSION IF NOT EXISTS vector SCHEMA public;`)
	if err != nil {
		return fmt.Errorf("failed to create pgvector extension: %w", err)
	}

	// Create the table if it doesn't exist
	_, err = p.conn.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id TEXT PRIMARY KEY,
			text TEXT NOT NULL,
			metadata JSONB,
			embedding vector(%d)
			-- embedding 컬럼의 타입이 어느 스키마의 vector 타입인지 명확히 하려면, 
			-- "public.vector(1536)"처럼 스키마를 명시할 수 있습니다.
			-- 예: embedding public.vector(1536)
			-- vector 확장이 어느 스키마에 설치되어 있는지 확인하려면 아래 쿼리를 사용할 수 있습니다:
			-- SELECT n.nspname FROM pg_extension e JOIN pg_namespace n ON e.extnamespace = n.oid WHERE e.extname = 'vector';
		)
	`, p.className, p.dimension))
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	// Create index for vector similarity search
	_, err = p.conn.Exec(ctx, fmt.Sprintf(`
		CREATE INDEX IF NOT EXISTS %s_embedding_idx ON %s 
		USING ivfflat (embedding vector_cosine_ops)
		WITH (lists = 100)
	`, p.className, p.className))
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}

	return nil
}

func (p *PGVector) Index(ctx context.Context, docs ...ragkit.Document) ([]string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	var ids []string
	for _, doc := range docs {
		var embedding []float32
		var err error

		if doc.ID == "" {
			doc.ID = ragkit.GenerateID(doc.Text, doc.Metadata)
		}

		// Use provided vector if available, otherwise generate one
		if doc.Vector != nil {
			embedding = doc.Vector
		} else {
			vectors, err := p.embedder.EmbedTexts(ctx, doc.Text)
			if err != nil {
				return nil, err
			}
			embedding = vectors[0]
		}

		// Insert document into database
		_, err = p.conn.Exec(ctx, fmt.Sprintf(`
			INSERT INTO %s (id, text, metadata, embedding) 
			VALUES ($1, $2, $3, $4)
		`, p.className), doc.ID, doc.Text, doc.Metadata, pgvector.NewVector(embedding))
		if err != nil {
			return nil, err
		}
		ids = append(ids, doc.ID)
	}
	return ids, nil
}

func (p *PGVector) Delete(ctx context.Context, id string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	_, err := p.conn.Exec(ctx, fmt.Sprintf(`
		DELETE FROM %s WHERE id = $1
	`, p.className), id)
	return err
}

func (p *PGVector) Exists(ctx context.Context, id string) (bool, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	var exists bool
	err := p.conn.QueryRow(ctx, fmt.Sprintf(`
		SELECT EXISTS(SELECT 1 FROM %s WHERE id = $1)
	`, p.className), id).Scan(&exists)
	return exists, err
}

func (p *PGVector) Retrieve(ctx context.Context, query []float32, topK int, metadataFieldNames ...string) ([]ragkit.RetrievedDoc, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Perform vector similarity search
	rows, err := p.conn.Query(ctx, fmt.Sprintf(`
		SELECT text, metadata, embedding 
		FROM %s 
		ORDER BY embedding <=> $1 
		LIMIT $2
	`, p.className), pgvector.NewVector(query), topK)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []ragkit.RetrievedDoc
	for rows.Next() {
		var doc ragkit.RetrievedDoc
		var embedding pgvector.Vector
		err := rows.Scan(&doc.Text, &doc.Metadata, &embedding)
		if err != nil {
			return nil, err
		}
		doc.Vector = embedding.Slice()
		results = append(results, doc)
	}
	return results, rows.Err()
}

func (p *PGVector) RetrieveText(ctx context.Context, text string, topK int, metadataFieldNames ...string) ([]ragkit.RetrievedDoc, error) {
	vectors, err := p.embedder.EmbedTexts(ctx, text)
	if err != nil {
		return nil, err
	}
	return p.Retrieve(ctx, vectors[0], topK, metadataFieldNames...)
}

func (p *PGVector) String() string {
	return fmt.Sprintf("PGVector(table: %s, embedder: %s)", p.className, p.embedder)
}
