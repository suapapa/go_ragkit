# AGENTS.md — go_ragkit

RAG용 VectorStore / Embedder Go 라이브러리. AI 에이전트·기여자용 작업 규칙입니다.

## 프로젝트 한 줄

텍스트를 임베딩해 Weaviate·pgvector 등에 인덱싱하고 의미 검색하는 Go 패키지.

## 읽기 순서

1. [`README.md`](README.md)
2. [`interface.go`](interface.go) — `Embedder`, `VectorStore` 계약
3. [`helper/`](helper/) — Weaviate/OpenAI 조합 헬퍼
4. [`examples/`](examples/) — 로컬 실행 예제

## 레이아웃

| 경로 | 책임 |
|------|------|
| `embedder/` | OpenAI·Ollama 임베더 구현 |
| `vector_store/` | Weaviate·pgvector 스토어 |
| `helper/` | 자주 쓰는 조합 팩토리 (`NewWeaviateOpenAIVectorStore` 등) |
| `examples/` | 독립 실행 예제 |

## 환경변수 (OpenAI 경로)

헬퍼는 **호출 시점**에 env를 읽습니다 (패키지 init이 아님).

| 변수 | 용도 |
|------|------|
| `OPENAI_ENDPOINT` | OpenAI 호환 base URL (비우면 SDK 기본) |
| `OPENAI_APIKEY` | API 키 |
| `OPENAI_EMBED_MODEL` | 기본 `text-embedding-3-small` |
| `WEAVIATE_ADDR` | 기본 `http://localhost:8080` |

Ollama 전용 env (`OLLAMA_ADDR` 등)는 사용하지 않습니다. Ollama 헬퍼는 하드코딩 기본값만 둡니다.

## 코딩 규칙

- Go **1.27** 기준. `go build ./...` 로 검증.
- 공개 API 변경 시 README Quick Start / examples를 함께 갱신.
- 새 임베더·스토어는 `interface.go` 계약을 만족시키고 examples에 최소 하나 추가.
- 시크릿·API 키를 소스나 커밋에 넣지 말 것.
- 의존성 업그레이드 시 weaviate `v1.40.x` grpc 모듈 깨짐을 주의 (`v1.39.x` 계열 유지 권장).

## 하지 말 것

- `OPENAI_SECRET_KEY` / `OPEAI_EMBED_MODEL` / `OLLAMA_*` env 이름을 다시 도입
- helper에서 env를 패키지 레벨 `var`로 캐시 (`.env` 로드 순서와 충돌)
- 사용자 요청 없이 커밋·푸시

## 커밋

사용자 요청 시에만. 메시지 스타일: `feat|fix|refactor|chore: …`
