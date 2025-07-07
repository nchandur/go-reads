# GoReads!

GoReads! is a robust and scalable book recommendation system built entirely in Go. It leverages a combination of a traditional database for book and author data, a vector database for semantic search, and a Go web application to provide a powerful and interactive API for recommendations.

# Features

- **Content-Based Recommendations**: Utilizes vector embeddings to provide recommendations based on book content.

- **Scalable Architecture**: Built with Docker Compose, making it easy to deploy and scale.

- **Go Web API**: A high-performance Gin web framework provides a clean API for user interaction.

- **Persistent Data Storage**: Uses MongoDB for reliable storage of book and author information.

- **Efficient Vector Search**: Qdrant VectorDB enables fast and accurate similarity searches for recommendations.

# Technologies Used

- **Go**: The primary language for the web application.

- **Gin Web Framework**: For building the RESTful API endpoints.

- **MongoDB**: A NoSQL database for storing structured book and author data.

- **Qdrant VectorDB**: A vector similarity search engine for storing and querying book embeddings.

- **Docker & Docker Compose**: For containerization and orchestration of all services.

**Make**: For simplifying common development tasks like starting and stopping services.

# Prerequisites

Before getting started, ensure the following are installed:

- **Git**

- **Docker & Docker Compose**

# Getting Started

Follow these steps to get GoReads! up and running locally.

1. **Clone the Repository**
```{bash}
git clone https://github.com/nchandur/go-reads
cd go-reads
```

2. **Prepare the Data**

You'll need the data/ directory containing the most updated database. Please request this from me [Contact Me](https://nchandur.github.io/portfolio/). Once you have it, place this directory in the root of your project folder.

Your project structure should look something like this:

```
.
├── data
│   └── books
│       ├── prelude.json
│       ├── works.bson
│       └── works.metadata.json
├── docker-compose.yaml
├── Makefile
├── README.md
└── webapp
        ...
```

3. **Configure Environment Variables**

Create a `.env` file in the root of your project directory by copying the example:

```
QDRANT_API_URL=http://qdrant:6333
MONGO_DB_URI=mongodb://mongodb:27017
```

4. **Start the Services**

Once your `data/` directory is in place and your `.env` file is configured, you can start all the services using Docker Compose:

```{bash}
make up
```

This command will build the necessary Docker images (if not already built), create the containers for MongoDB, Qdrant, and your Go Web App, and start them. It might take a few minutes for all services to be fully operational, especially on the first run as images are downloaded and built.

# Usage (API Endpoints)

Once all services are running, your Go Web App will expose several API endpoints that you can use to interact with the recommendation system. You can access the Go Web App at http://localhost:8080

## Landing Page

### 1. `GET /`

**Description**: A landing page endpoint that indicates if all services are up and running, and all databases have been loaded correctly.

**Expected Response**: Status 200 (Success) if everything is operational.

## Health Checks

### 2. `GET /health/mongodb`

**Description**: Pings the MongoDB instance to check its connectivity.

**Expected Response**: Status 200 (Success) if MongoDB is reachable.

## 3. `GET /health/vectordb`

**Description**: Pings the Qdrant VectorDB instance to check its connectivity.

**Expected Response**: Status 200 (Success) if Qdrant is reachable.


## Retrieve Data

### 4. `GET /books/{id}`

**Description**: Retrieves a single book document by its unique ID.

**Parameters**:

- `:id` (path parameter, integer): The unique identifier of the book.

**Expected Response**:

- **Status 200**: Returns the book document.

- **Status 400**: If `id` is malformed or invalid.

- **Status 500**: If there's an internal server error.

### 5. `GET /books`

**Description**: Searches for books by title, supporting fuzzy matching.

**Query Parameters**:

- `title` (string): The title of the book to search for.

**Expected Response**:

- **Status 200**: Returns a document (or slice of documents) matching the title.

- **Status 400**: If `title` parameter is missing or invalid.

- **Status 500**: If there's an internal server error.

### 6. `GET /authors/{:id}`

**Description**: Retrieves a single author document by their unique ID.

**Parameters**:

- `:id` (path parameter, integer): The unique identifier of the author.

**Expected Response**:

- **Status 200**: Returns the author document.

- **Status 400**: If `id` is malformed or invalid.

- **Status 500**: If there's an internal server error.

### 7. `GET /authors`

**Description**: Searches for authors by name, supporting fuzzy matching.

**Query Parameters**:

- `name` (string): The name of the author to search for.

**Expected Response**:

- **Status 200**: Returns a document (or slice of documents) matching the author's name.

- **Status 400**: If `name` parameter is missing or invalid.

- **Status 500**: If there's an internal server error.

## Recommendations

### 8. `GET /books/recommendations`

**Description**: Provides book recommendations based on a given book title.

**Query Parameters**:

- `title` (string): The title of the book to get recommendations for.

- `n` (integer): The number of recommendations to return.

**Expected Response**:

- **Status 200**: Returns the matched book title and a slice of recommended book documents.

- **Status 400**: If `title` or `n` parameters are missing or invalid.

- **Status 500**: If there's an internal server error.

### 9. `GET /authors/recommendations`

**Description**: Provides author recommendations based on a given author's name.

**Query Parameters**:

- `name` (string): The name of the author to get recommendations for.

- `n` (integer): The number of recommendations to return.

**Expected Response**:

- **Status 200**: Returns the matched author and a slice of recommended author documents.

- **Status 400**: If `name` or `n` parameters are missing or invalid.

- **Status 500**: If there's an internal server error.

### 10. `GET /genres/books`

**Description**: Provides book recommendations from a specific genre.

**Query Parameters**:

- `genre` (string): The genre to get book recommendations from.

- `n` (integer): The number of recommendations to return.

**Expected Response**:

- **Status 200**: Returns a slice of recommended book documents from the specified genre.

- **Status 400**: If `genre` or `n` parameters are missing or invalid.

- **Status 500**: If there's an internal server error.

### 11. `GET /genres/authors`

**Description**: Provides author recommendations from a specific genre.

**Query Parameters**:

- `genre` (string): The genre to get author recommendations from.

- `n` (integer): The number of recommendations to return.

**Expected Response**:

- **Status 200**: Returns a slice of recommended author documents from the specified genre.

- **Status 400**: If genre or n parameters are missing or invalid.

- **Status 500**: If there's an internal server error.


# Stopping and Cleaning

You can manage the Docker services using the `make` commands:

**Stop Services**
To stop all running Docker containers gracefully:
```{bash}
make down
```

**Clean Up (Wipe Volumes and Containers)**
To stop all services and completely remove all containers, networks, and volumes (which will delete your MongoDB and Qdrant data), run:
```{bash}
make clean
```
Use `make clean` with caution, as it will permanently delete all data stored in your MongoDB and Qdrant volumes.
