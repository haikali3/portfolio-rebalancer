package kafka

import (
	"context"
	"encoding/json"
	"log"
	"portfolio-rebalancer/internal/models"
	"portfolio-rebalancer/internal/services"
	"portfolio-rebalancer/internal/storage"

	"github.com/segmentio/kafka-go"
)

func StartRebalanceConsumer(ctx context.Context) {
	ConsumeMessage(ctx, func(msg kafka.Message) {

		// 1. unmarshal msg raw bytes from kafka
		var req models.UpdatedPortfolio
		if err := json.Unmarshal(msg.Value, &req); err != nil {
			log.Printf("Failed to unmarshal rebalance message: %v", err)
			return
		}
		// 2. get original allocation from elasticsearch
		original, err := storage.GetPortfolio(ctx, req.UserID)
		if err != nil {
			log.Printf("Failed to get portfolio for user %s: %v", req.UserID, err)
			return
		}

		// 3. calcu rebalance tx
		transactions := services.CalculateRebalance(req.UserID, req.NewAllocation, original.Allocation)

		// 4. save tx to elasticsearch
		if err := storage.SaveRebalanceTransactions(ctx, transactions); err != nil {
			log.Printf("Failed to save rebalance transactions for user %s: %v", req.UserID, err)
			return
		}

		log.Printf("Rebalanced user %s: %d transactions", req.UserID, len(transactions))
	})
}
