package transaction

import (
	"bff/account"
	nats_client "bff/nats"
	"context"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/nats-io/nats.go"
)

type TransactionService struct {
	Repository        TransactionRepository
	AccountRepository account.AccountRepository
	NatsClient        *nats.Conn
	RedisClient       *redis.Client
	CacheDuration     int
}

func NewTransactionService(repository TransactionRepository, accountRepository account.AccountRepository, natsClient *nats.Conn, redisAddr string, cacheDuration int) *TransactionService {
	redis := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	return &TransactionService{Repository: repository, AccountRepository: accountRepository, NatsClient: natsClient, CacheDuration: 5, RedisClient: redis}
}

func (s *TransactionService) TransferMoney(transaction Transaction) (Transaction, error) {
	balance, err := s.verifyAccountBalance(transaction.FromAccount)
	if err != nil {
		return Transaction{}, err
	}

	if balance < transaction.Amount {
		return Transaction{}, errors.New("insufficient funds")
	}

	err = s.updateAccountBalance(transaction.FromAccount, -transaction.Amount)
	if err != nil {
		return Transaction{}, err
	}

	err = s.updateAccountBalance(transaction.ToAccount, transaction.Amount)
	if err != nil {
		return Transaction{}, err
	}

	transaction, err = s.Repository.TransferMoney(transaction)
	if err != nil {
		return Transaction{}, err
	}

	err = nats_client.PublishToStream(s.NatsClient, "transactions.new", []byte("Nueva transacción registrada"))
	if err != nil {
		log.Printf("Error publishing to NATS: %v", err)
	}

	return transaction, nil
}

func (s *TransactionService) GetTransactionHistory(accountID int) ([]Transaction, error) {
	return s.Repository.GetTransactionHistory(accountID)
}

func (s *TransactionService) verifyAccountBalance(accountID int) (int, error) {
	balance, err := s.AccountRepository.GetAccountBalance(accountID)
	if err != nil {
		log.Printf("Error verifying account balance: %v", err)
		return 0, err
	}
	return balance, nil
}

// Función para actualizar el saldo de una cuenta
func (s *TransactionService) updateAccountBalance(accountID int, newBalance int) error {
	balanceUpdated, err := s.AccountRepository.UpdateAccountBalance(accountID, newBalance)
	if err != nil {
		log.Printf("Error updating account balance: %v", err)
		return err
	}

	ctx := context.Background()
	s.RedisClient.Set(ctx, "account_balance_"+strconv.Itoa(accountID), balanceUpdated, time.Duration(s.CacheDuration)*time.Minute)

	return nil
}
