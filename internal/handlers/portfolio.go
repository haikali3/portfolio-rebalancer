package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"portfolio-rebalancer/internal/models"
	"portfolio-rebalancer/internal/services"
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

	if err := storage.SavePortfolio(r.Context(), p); err != nil {
		log.Printf("Failed to save portfolio: %v", err)
		return NewAPIError(http.StatusInternalServerError, "failed to save portfolio")
	}

	return writeJSON(w, http.StatusCreated, p)
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

	// 1. get user's allocation
	original, err := storage.GetPortfolio(r.Context(), req.UserID)
	if err != nil {
		log.Printf("Failed to get portfolio: %v", err)
		return NewAPIError(http.StatusInternalServerError, "failed to get portfolio")
	}

	// 2. calc rebalance transaction (buy/sell) to move from new allowcation back to original allocation
	transactions := services.CalculateRebalance(req.NewAllocation, original.Allocation)

	// 3. save rebalance transaction to db
	storage.SaveRebalanceTransactions(r.Context(), transactions)

	return writeJSON(w, http.StatusOK, transactions)
}
