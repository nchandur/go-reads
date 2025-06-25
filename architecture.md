graph TD
```
    A[User] -->|1. Sends Recommendation Request| B(Go API Service)

    B -->|2. Sends Query to Ollama for Embedding| C[Ollama Container]
    C -->|3. Returns Query Vector Embedding| B

    B -->|4. Sends Query Vector to Qdrant for Similarity Search| D[Qdrant Container]
    D -->|5. Returns Top K Book IDs (and maybe minimal payload)| B

    B -->|6. Fetches Full Book Details by ID| E[MongoDB Container]
    E -->|7. Returns Full Book Details| B

    B -->|8. Returns Top K Recommended Books| A

    subgraph Initial Data Ingestion (One-Time Script)
        F[Book Data Source] -->|a. Read Book Summaries| G(Go Ingestion Script)
        G -->|b. Send Summary to Ollama for Embedding| C
        C -->|c. Returns Book Vector Embedding| G
        G -->|d. Insert Vector + Book ID| D
        G -->|e. Insert Full Book Details + Vector (optional)| E
    end
```
```
book-recommendation-system/
├── cmd/
│   └── book-recommender/      # Main application executable
│       └── main.go            # Entry point for the API service
├── internal/
│   ├── api/                   # HTTP handlers and routing
│   │   └── handlers.go
│   │   └── routes.go
│   ├── config/                # Application configuration loading
│   │   └── config.go
│   ├── model/                 # Data structures (structs for books, queries, API requests/responses)
│   │   └── book.go
│   │   └── query.go
│   ├── ollama/                # Client for interacting with the Ollama embedding service
│   │   └── client.go
│   ├── qdrant/                # Client for interacting with the Qdrant vector database
│   │   └── client.go
│   ├── mongodb/               # Client for interacting with the MongoDB database
│   │   └── client.go
│   ├── service/               # Core business logic (orchestrates calls to clients, handles recommendations)
│   │   └── recommender.go
│   └── app/                   # Application bootstrapping and dependency injection
│       └── app.go
├── scripts/
│   └── ingest/                # One-time data ingestion script
│       └── main.go            # Entry point for the ingestion process
│   └── data/                  # Placeholder for raw book data files (e.g., .csv, .json)
│       └── books.json
├── pkg/                       # Reusable packages not specific to this application (if any)
│   └── utils/
│       └── string_helpers.go
├── Dockerfile                 # For building your Go API service container
├── docker-compose.yml         # For orchestrating your services (Go API, Ollama, Qdrant, MongoDB)
├── go.mod                     # Go module definitions
├── go.sum                     # Go module checksums
├── .env.example               # Example environment variables
└── README.md                  # Project documentation
```