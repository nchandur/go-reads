.PHONY:

all: start-services pull-embedder embed backup

start-services:
	docker compose up -d

restore:
	docker exec -it mongodb mongorestore --drop --nsInclude "books.*" /app/backup_data/
	docker cp data/backup/qdrant_backup/. vectordb:/qdrant/data/


pull-embedder:
	docker exec -it ollama ollama pull nomic-embed-text

embed:
	go run scripts/ingest/main.go

scrape:
	go run scripts/scrape/main.go

deep-clean:
	docker compose down -v
	docker volume prune -a -f

clean:
	docker compose down

backup:
	mongodump --host localhost --port 9001 --db books --out data/backup
	sudo cp -r /var/lib/docker/volumes/goreads_qdrant_data/_data data/backup/qdrant_backup

preprocess:
	go run scripts/preprocess/main.go