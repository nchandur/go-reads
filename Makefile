include .env
export $(shell sed 's/=.*//' .env)

up: start-services restore preprocess

start-services:
	docker compose up --build -d
	docker exec -it ollama ollama pull ${OLLAMA_EMBEDDING_MODEL}

restore:
	docker exec -it mongodb mongorestore --authenticationDatabase admin -u ${MONGO_INITDB_ROOT_USERNAME} -p ${MONGO_INITDB_ROOT_PASSWORD} --drop --db books /app/local_data/books/

preprocess:
	docker exec -it go_webapp go run /app/scripts/preprocess/main.go
	docker exec -it go_webapp go run /app/scripts/ingest/main.go

down:
	docker compose down

clean:
	docker compose down -v
	docker volume prune -a -f