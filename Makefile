.PHONY:

all: start-services restore preprocess embed

start-services:
	docker compose up -d
	docker exec -it ollama ollama pull nomic-embed-text

restore:
	docker exec -it mongodb mongorestore --drop --nsInclude "books.*" /app/backup_data/

embed:
	go run scripts/ingest/main.go

scrape:
	go run scripts/scrape/main.go

deep-clean:
	docker compose down -v
	docker volume prune -a -f

clean:
	docker compose down

preprocess:
	go run scripts/preprocess/main.go