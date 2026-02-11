package models

type Portfolio struct {
	UserID     string             `json:"user_id"`
	Allocation map[string]float64 `json:"allocation"` // Current user allocation in percentage terms
}

type UpdatedPortfolio struct {
	UserID        string             `json:"user_id"`
	NewAllocation map[string]float64 `json:"new_allocation"` // Updated user allocation from provider in percentage terms
}

type RebalanceTransaction struct {
	// TODO: Add model
	UserID               string  `json:"user_id"`
	Action               string  `json:"action"` // "BUY" or "SELL"
	Asset                string  `json:"asset"`  // e.g., "stocks", "bonds", etc.
	RebalanceTransaction float64 `json:"amount"` // percentage to buy/sell
}
