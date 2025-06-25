.PHONY: all

all: start-services restore pull-embedder

start-services:
	docker compose up -d

restore:
	docker exec -it mongodb mongorestore --drop --nsInclude "books.*" /app/backup_data/

pull-embedder:
	docker exec -it ollama ollama pull nomic-embed-text

deep-clean:
	docker compose down -v

clean:
	docker compose down