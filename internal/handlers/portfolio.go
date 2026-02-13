package handlers

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"portfolio-rebalancer/internal/kafka"
	"portfolio-rebalancer/internal/models"
	"portfolio-rebalancer/internal/storage"
)

// HandlePortfolio handles new portfolio creation requests (feel free to update the request parameter/model)
// Sample Request (POST /portfolio):
//
//	{
//	    "user_id": "1",
//	    "allocation": {"stocks": 60, "bonds": 30, "gold": 10}
//	}
func HandlePortfolio(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return NewAPIError(http.StatusMethodNotAllowed, "method not allowed")
	}

	var p models.Portfolio
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		return InvalidJSON()
	}

	// TODO: Add Logic here
	if p.UserID == "" {
		return NewAPIError(http.StatusBadRequest, "user_id is required")
	}

	if len(p.Allocation) == 0 {
		return NewAPIError(http.StatusBadRequest, "allocation is required")
	}

	if err := validateAllocationSum(p.Allocation); err != nil {
		return err
	}

	if err := storage.SavePortfolio(r.Context(), p); err != nil {
		log.Printf("Failed to save portfolio: %v", err)
		return NewAPIError(http.StatusInternalServerError, "failed to save portfolio")
	}

	return writeJSON(w, http.StatusCreated, map[string]any{
		"status_code": http.StatusCreated,
		"data":        p,
	})
}

// HandleRebalance handles portfolio rebalance requests from 3rd party provider (feel free to update the request parameter/model)
// Sample Request (POST /rebalance):
//
//	{
//	    "user_id": "1",
//	    "new_allocation": {"stocks": 70, "bonds": 20, "gold": 10}
//	}
func HandleRebalance(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return NewAPIError(http.StatusMethodNotAllowed, "method not allowed")
	}

	var req models.UpdatedPortfolio
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return InvalidJSON()
	}

	log.Println("HandleRebalance==", req)

	// TODO: Add Logic here
	if req.UserID == "" {
		return NewAPIError(http.StatusBadRequest, "user_id is required")
	}

	if len(req.NewAllocation) == 0 {
		return NewAPIError(http.StatusBadRequest, "new_allocation is required")
	}

	if err := validateAllocationSum(req.NewAllocation); err != nil {
		return err
	}

	// publish rebalance request to kafka for async processing
	payload, err := json.Marshal(req)
	if err != nil {
		return NewAPIError(http.StatusInternalServerError, "failed to marshal rebalance request")
	}

	if err := kafka.PublishMessage(r.Context(), payload); err != nil {
		log.Printf("Failed to publish rebalance message to Kafka: %v", err)
		return NewAPIError(http.StatusInternalServerError, "failed to publish rebalance message")
	}

	return writeJSON(w, http.StatusOK, map[string]any{
		"status_code": http.StatusOK,
		"msg":         "rebalance request received and being processed",
	})
}

func validateAllocationSum(allocation map[string]float64) error {
	var sum float64
	for _, pct := range allocation {
		if pct < 0 {
			return NewAPIError(http.StatusBadRequest, "allocation percentages must be non-negative")
		}
		sum += pct
	}
	if math.Abs(sum-100) > 0.01 {
		return NewAPIError(http.StatusBadRequest, "allocation percentages must sum to 100")
	}
	return nil
}
