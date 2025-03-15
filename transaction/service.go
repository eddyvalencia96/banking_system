package transaction

import (
	"bff/account"
	nats_client "bff/nats"
	"errors"
	"log"

	"github.com/nats-io/nats.go"
)

type TransactionService struct {
	Repository        TransactionRepository
	AccountRepository account.AccountRepository
	NatsClient        *nats.Conn
}

func NewTransactionService(repository TransactionRepository, accountRepository account.AccountRepository, natsClient *nats.Conn) *TransactionService {
	return &TransactionService{Repository: repository, AccountRepository: accountRepository, NatsClient: natsClient}
}

// Función para transferir dinero entre cuentas
func (s *TransactionService) TransferMoney(transaction Transaction) (Transaction, error) {
	// Verificar saldo antes de la transferencia
	balance, err := s.verifyAccountBalance(transaction.FromAccount)
	if err != nil {
		return Transaction{}, err
	}
	if balance < transaction.Amount {
		return Transaction{}, errors.New("insufficient funds")
	}

	// Realizar la transferencia
	transaction, err = s.Repository.TransferMoney(transaction)
	if err != nil {
		return Transaction{}, err
	}

	// Publicar evento en JetStream
	err = nats_client.PublishToStream(s.NatsClient, "transactions.new", []byte("Nueva transacción registrada"))
	if err != nil {
		log.Printf("Error publishing to NATS: %v", err)
	}

	return transaction, nil
}

// Función para obtener el historial de transacciones de una cuenta
func (s *TransactionService) GetTransactionHistory(accountID int) ([]Transaction, error) {
	return s.Repository.GetTransactionHistory(accountID)
}

// Función para verificar el saldo de una cuenta
func (s *TransactionService) verifyAccountBalance(accountID int) (int, error) {
	// Obtener el saldo de la cuenta
	balance, err := s.AccountRepository.GetAccountBalance(accountID)
	if err != nil {
		log.Printf("Error verifying account balance: %v", err)
		return 0, err
	}
	return balance, nil
}
