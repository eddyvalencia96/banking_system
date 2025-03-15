package transaction

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type TransactionRepository struct {
	DB *sql.DB
}

func NewTransactionRepository() *TransactionRepository {
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
	return &TransactionRepository{DB: db}
}

func (r *TransactionRepository) TransferMoney(transaction Transaction) (Transaction, error) {
	query := `INSERT INTO transactions (from_account, to_account, amount, timestamp) VALUES ($1, $2, $3, $4) RETURNING transaction_id`
	err := r.DB.QueryRow(query, transaction.FromAccount, transaction.ToAccount, transaction.Amount, transaction.Timestamp).Scan(&transaction.TransactionID)
	if err != nil {
		return Transaction{}, err
	}
	return transaction, nil
}

func (r *TransactionRepository) GetTransactionHistory(accountID int) ([]Transaction, error) {
	rows, err := r.DB.Query("SELECT * FROM transactions WHERE from_account=$1 OR to_account=$1", accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var t Transaction
		err := rows.Scan(&t.TransactionID, &t.FromAccount, &t.ToAccount, &t.Amount, &t.Timestamp)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}
	return transactions, nil
}
