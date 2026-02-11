package services

import (
	"portfolio-rebalancer/internal/models"
)

func CalculateRebalance(UserID string, updatedAllocation, currentAlloctaion map[string]float64) []models.RebalanceTransaction {
	var result []models.RebalanceTransaction

	// TODO: create rebalance transactions and update portfolio
	for asset, originalPercentage := range currentAlloctaion {
		diff := originalPercentage - updatedAllocation[asset]
		if diff > 0 {
			// user has too little, buy more
			result = append(result, models.RebalanceTransaction{
				UserID:           UserID,
				Asset:            asset,
				Action:           "BUY",
				RebalancePercent: diff,
			})
		} else if diff < 0 {
			// user has too much, sell some
			result = append(result, models.RebalanceTransaction{
				UserID:           UserID,
				Asset:            asset,
				Action:           "SELL",
				RebalancePercent: -diff,
			})
		}
	}

	return result
}
