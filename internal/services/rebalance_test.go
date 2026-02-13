package services

import (
	"portfolio-rebalancer/internal/models"
	"reflect"
	"sort"
	"testing"
)

func TestCalculateRebalance(t *testing.T) {
	tests := []struct {
		name               string
		userID             string
		originalAllocation map[string]float64
		updatedAllocation  map[string]float64
		expected           []models.RebalanceTransaction
	}{
		{
			name:               "rebalance stocks and bonds",
			userID:             "1",
			originalAllocation: map[string]float64{"stocks": 60, "bonds": 30, "gold": 10},
			updatedAllocation:  map[string]float64{"stocks": 70, "bonds": 20, "gold": 10},
			expected: []models.RebalanceTransaction{
				{UserID: "1", Asset: "stocks", Action: "SELL", RebalancePercent: 10},
				{UserID: "1", Asset: "bonds", Action: "BUY", RebalancePercent: 10},
			},
		},
		{
			name:               "rebalance all assets",
			userID:             "2",
			originalAllocation: map[string]float64{"stocks": 50, "bonds": 40, "gold": 10},
			updatedAllocation:  map[string]float64{"stocks": 30, "bonds": 50, "gold": 20},
			expected: []models.RebalanceTransaction{
				{UserID: "2", Asset: "stocks", Action: "BUY", RebalancePercent: 20},
				{UserID: "2", Asset: "bonds", Action: "SELL", RebalancePercent: 10},
				{UserID: "2", Asset: "gold", Action: "SELL", RebalancePercent: 10},
			},
		},
		{
			name:               "no rebalance needed",
			userID:             "3",
			originalAllocation: map[string]float64{"stocks": 70, "bonds": 20, "gold": 10},
			updatedAllocation:  map[string]float64{"stocks": 70, "bonds": 20, "gold": 10},
			expected:           nil,
		},
		{
			name:               "new asset in updated not in original",
			userID:             "4",
			originalAllocation: map[string]float64{"stocks": 70, "bonds": 20, "gold": 10},
			updatedAllocation:  map[string]float64{"stocks": 60, "bonds": 20, "gold": 10, "crypto": 10},
			expected: []models.RebalanceTransaction{
				{UserID: "4", Asset: "crypto", Action: "SELL", RebalancePercent: 10},
				{UserID: "4", Asset: "stocks", Action: "BUY", RebalancePercent: 10},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateRebalance(tt.userID, tt.updatedAllocation, tt.originalAllocation)
			sort.Slice(result, func(i, j int) bool {
				return result[i].Asset < result[j].Asset
			})
			sort.Slice(tt.expected, func(i, j int) bool {
				return tt.expected[i].Asset < tt.expected[j].Asset
			})

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}
