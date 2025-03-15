package transaction

type Transaction struct {
	TransactionID string `json:"transaction_id"`
	FromAccount   int    `json:"from_account"`
	ToAccount     int    `json:"to_account"`
	Amount        int    `json:"amount"`
	Timestamp     string `json:"timestamp"`
}
