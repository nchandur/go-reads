.PHONY: all

all: start-services restore pull-embedder embed

start-services:
	docker compose up -d

restore:
	docker exec -it mongodb mongorestore --drop --nsInclude "books.*" /app/backup_data/

pull-embedder:
	docker exec -it ollama ollama pull nomic-embed-text

embed:
	go run scripts/ingest/main.go

deep-clean:
	docker compose down -v

clean:
	docker compose down