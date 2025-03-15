package account

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

// TestCreateAccount prueba la creación de una cuenta
func TestCreateAccount(t *testing.T) {
	// Configurar la base de datos en memoria
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// Crear la tabla de cuentas
	createTableQuery := `
	CREATE TABLE accounts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		initial_balance INTEGER NOT NULL,
		balance INTEGER NOT NULL
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		t.Fatalf("Error creating table: %v", err)
	}

	// Crear el repositorio
	repo := &AccountRepository{DB: db}

	// Crear una cuenta
	account := Account{Name: "Test Account", InitialBalance: 1000}
	createdAccount, err := repo.CreateAccount(account)
	if err != nil {
		t.Fatalf("Error creating account: %v", err)
	}

	// Verificar la cuenta creada
	assert.Equal(t, 1, createdAccount.ID)
	assert.Equal(t, "Test Account", createdAccount.Name)
	assert.Equal(t, 1000, createdAccount.Balance)
}

// TestGetAccountBalance prueba la obtención del saldo de una cuenta
func TestGetAccountBalance(t *testing.T) {
	// Configurar la base de datos en memoria
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// Crear la tabla de cuentas
	createTableQuery := `
	CREATE TABLE accounts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		initial_balance INTEGER NOT NULL,
		balance INTEGER NOT NULL
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		t.Fatalf("Error creating table: %v", err)
	}

	// Crear el repositorio
	repo := &AccountRepository{DB: db}

	// Insertar una cuenta de prueba
	insertQuery := `INSERT INTO accounts (name, initial_balance, balance) VALUES (?, ?, ?)`
	_, err = db.Exec(insertQuery, "Test Account", 1000, 1000)
	if err != nil {
		t.Fatalf("Error inserting test account: %v", err)
	}

	// Obtener el saldo de la cuenta
	balance, err := repo.GetAccountBalance(1)
	if err != nil {
		t.Fatalf("Error getting account balance: %v", err)
	}

	// Verificar el saldo
	assert.Equal(t, 1000, balance)
}
