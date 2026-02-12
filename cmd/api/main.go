package main

import (
	"context"
	"log"
	"net/http"
	"portfolio-rebalancer/internal/handlers"
	"portfolio-rebalancer/internal/kafka"
	"portfolio-rebalancer/internal/storage"
)

func main() {

	// Initializing elasticsearch if needed
	if err := storage.InitElastic(); err != nil {
		log.Fatalf("Failed to initialize Elasticsearch: %v", err)
	}

	if err := kafka.InitKafka(); err != nil {
		log.Fatalf("Failed to initialize Kafka producer: %v", err)
	}
	kafka.StartRebalanceConsumer(context.Background())

	http.HandleFunc("/portfolio", handlers.Make(handlers.HandlePortfolio))
	http.HandleFunc("/rebalance", handlers.Make(handlers.HandleRebalance))

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
