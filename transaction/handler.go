package transaction

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type TransactionHandler struct {
	Service TransactionService
}

func NewTransactionHandler(service TransactionService) *TransactionHandler {
	return &TransactionHandler{Service: service}
}

func (h *TransactionHandler) TransferMoney(w http.ResponseWriter, r *http.Request) {
	var transaction Transaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validateTransaction(transaction); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	transaction.Timestamp = time.Now().Format(time.RFC3339)

	result, err := h.Service.TransferMoney(transaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// validateTransaction validates the transaction structure
func validateTransaction(transaction Transaction) error {
	if transaction.FromAccount <= 0 {
		return errors.New("invalid source account ID")
	}

	if transaction.ToAccount <= 0 {
		return errors.New("invalid destination account ID")
	}

	if transaction.Amount <= 0 {
		return errors.New("transaction amount must be positive")
	}

	return nil
}

func (h *TransactionHandler) GetTransactionHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountIDStr := vars["account_id"]

	accountID, err := strconv.Atoi(accountIDStr)
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	history, err := h.Service.GetTransactionHistory(accountID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}
