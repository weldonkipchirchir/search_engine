# Search Engine

A distributed search engine built with Go and Rust, featuring web crawling, indexing, and full-text search capabilities.

## Architecture

- **Crawler (Go)**: Web crawler that collects pages and stores them in PostgreSQL
- **Indexer (Rust)**: Builds inverted index for fast searching
- **Search API (Go)**: REST API for searching indexed content

## Technologies

- Go 1.21+
- Rust 1.75+
- PostgreSQL 15+
- Docker & Docker Compose

## Quick Start

### Using Docker (Recommended)

```bash
# Start all services
docker-compose up -d

# Setup database
./scripts/setup_db.sh

# Add seed URLs
./scripts/seed_data.sh

# View logs
docker-compose logs -f
```

Access the search engine at: http://localhost:8080

### Manual Setup

1. **Setup PostgreSQL**

```bash
createdb search_engine
psql search_engine < migrations/001_create_documents.up.sql
psql search_engine < migrations/002_create_index.up.sql
```

2. **Run Crawler**

```bash
cd crawler
go mod download
go run cmd/crawler/main.go
```

3. **Run Indexer**

```bash
cd indexer
cargo build --release
cargo run --release
```

4. **Run Search API**

```bash
cd search-api
go mod download
go run cmd/api/main.go
```

## API Endpoints

- `GET /` - Search homepage
- `GET /api/v1/search?q=query` - Search for documents
- `GET /api/v1/health` - Health check
- `GET /api/v1/stats` - System statistics

## Configuration

Copy `.env.example` to `.env` and adjust settings.

## Development

```bash
# Run tests
make test

# View logs
./scripts/logs.sh [service]

# Rebuild specific service
docker-compose build [service]
```

## License

MIT
