package main

import (
	"log"

	"github.com/nchandur/go-reads/internal/api"
	"github.com/nchandur/go-reads/internal/db"
	"github.com/nchandur/go-reads/internal/vectordb"
)

func main() {

	err := db.Connect()

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := db.Disconnect(); err != nil {
			log.Fatal(err)
		}
	}()

	err = vectordb.Connect()

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := vectordb.Disconnect(); err != nil {
			log.Fatal(err)
		}
	}()

	r := api.SetUpRouter()

	log.Println("Server running at port :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

}
