package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"portfolio-rebalancer/internal/models"
	"portfolio-rebalancer/internal/services"
	"portfolio-rebalancer/internal/storage"

	"github.com/segmentio/kafka-go"
)

const maxRetries = 3

func StartRebalanceConsumer(ctx context.Context) {
	ConsumeMessage(ctx, func(msg kafka.Message) error {
		var lastErr error
		for attempt := 1; attempt <= maxRetries; attempt++ {
			lastErr = processRebalanceMessage(ctx, msg)
			if lastErr == nil {
				return nil
			}
			log.Printf("Retry %d/%d failed for message: %v", attempt, maxRetries, lastErr)
			if attempt < maxRetries {
				time.Sleep(time.Duration(attempt) * time.Second)
			}
		}

		log.Printf("All retries exhausted, sending message to DLQ: %v", lastErr)
		if dlqErr := PublishToDLQ(ctx, msg.Value); dlqErr != nil {
			log.Printf("Failed to publish to DLQ: %v", dlqErr)
		}
		return lastErr
	})
}

func processRebalanceMessage(ctx context.Context, msg kafka.Message) error {
	var req models.UpdatedPortfolio
	if err := json.Unmarshal(msg.Value, &req); err != nil {
		return err
	}

	original, err := storage.GetPortfolio(ctx, req.UserID)
	if err != nil {
		return err
	}

	transactions := services.CalculateRebalance(req.UserID, req.NewAllocation, original.Allocation)

	if len(transactions) == 0 {
		log.Printf("No rebalancing needed for user %s", req.UserID)
		return nil
	}

	if err := storage.SaveRebalanceTransactions(ctx, transactions); err != nil {
		return err
	}

	log.Printf("Rebalanced user %s: %d transactions", req.UserID, len(transactions))
	return nil
}
