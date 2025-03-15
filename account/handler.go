package account

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type AccountHandler struct {
	Service AccountService
}

// NewAccountHandler creates a new instance of AccountHandler with the provided account service.
func NewAccountHandler(service AccountService) *AccountHandler {
	return &AccountHandler{Service: service}
}

// CreateAccount handles HTTP requests to create a new account.
// It parses the request body as JSON into an Account object, validates the account data,
// creates the account using the service, and returns the created account as JSON response.
func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var account Account
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate account fields
	if err := validateAccount(account); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newAccount, err := h.Service.CreateAccount(account)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newAccount)
}

// validateAccount checks if the account data is valid.
// It verifies that the name is not empty and the initial balance is not negative.
func validateAccount(account Account) error {
	if account.Name == "" {
		return errors.New("name is required")
	}

	if account.InitialBalance < 0 {
		return errors.New("initial balance cannot be negative")
	}

	return nil
}

// GetAccountBalance handles HTTP requests to retrieve an account's balance.
// It extracts the account ID from the URL parameters, retrieves the balance using the service,
// and returns the balance information as a JSON response.
func (h *AccountHandler) GetAccountBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	balance, err := h.Service.GetAccountBalance(idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(balance)
}
