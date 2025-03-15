package account

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type AccountRepository struct {
	DB *sql.DB
}

// NewAccountRepository creates a new repository by establishing a connection to the PostgreSQL database
// using environment variables for connection details. It returns a pointer to the AccountRepository.
// It will panic if the database connection cannot be established.
func NewAccountRepository() *AccountRepository {
	connStr := fmt.Sprintf(
		"user=%s dbname=%s sslmode=disable password=%s host=%s port=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	return &AccountRepository{DB: db}
}

// CreateAccount inserts a new account into the database with the given name and initial balance.
// The balance is set to the initial balance. It returns the created account with the ID assigned by the database
// or an error if the operation fails.
func (r *AccountRepository) CreateAccount(account Account) (Account, error) {
	account.Balance = account.InitialBalance
	query := `INSERT INTO accounts (name, initial_balance, balance) VALUES ($1, $2, $3) RETURNING id`
	err := r.DB.QueryRow(query, account.Name, account.InitialBalance, account.Balance).Scan(&account.ID)
	if err != nil {
		return Account{}, err
	}
	return account, nil
}

// GetAccountBalance retrieves the current balance of an account identified by the provided ID.
// It returns the balance and nil error if successful, or 0 and an error if the account is not found
// or if any other database error occurs.
func (r *AccountRepository) GetAccountBalance(id int) (int, error) {
	var balance int
	query := `SELECT balance FROM accounts WHERE id=$1`
	err := r.DB.QueryRow(query, id).Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}

// UpdateAccountBalance modifies the balance of an account identified by the provided ID.
// It adds the amount to the current balance (use negative amount for subtraction).
// It returns the new balance and nil error if successful, or 0 and an error if the account is not found
// or if any other database error occurs.
func (r *AccountRepository) UpdateAccountBalance(id int, amount int) (int, error) {
	var newBalance int
	query := `UPDATE accounts SET balance = balance + $2 WHERE id=$1 RETURNING balance`
	err := r.DB.QueryRow(query, id, amount).Scan(&newBalance)
	if err != nil {
		return 0, err
	}
	return newBalance, nil
}

func (r *AccountRepository) ValidateUserPassword(username, password string) (bool, error) {
	var storedPassword string
	query := `SELECT password FROM users WHERE "username"=$1`
	err := r.DB.QueryRow(query, username).Scan(&storedPassword)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return storedPassword == password, nil
}
