package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"bff/account"
	"bff/middleware"
	nats_client "bff/nats"
	"bff/transaction"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	natsClient, err := nats_client.NewClient(os.Getenv("NATS_URL"))
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	defer natsClient.Close()

	err = nats_client.CreateStream(natsClient, "TRANSACTIONS", []string{"transactions.*"})
	if err != nil {
		log.Fatalf("Error creating stream: %v", err)
	}

	accountRepo := account.NewAccountRepository()
	transactionRepo := transaction.NewTransactionRepository()

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		log.Fatal("REDIS_ADDR environment variable is not set")
	}

	cacheDuration := 5
	redisDuration := os.Getenv("REDIS_DURATION")
	if redisDuration != "" {
		cacheDuration, err = strconv.Atoi(redisDuration)
		if err != nil {
			log.Fatalf("Error parsing REDIS_DURATION: %v", err)
		}
	}

	accountService := account.NewAccountService(*accountRepo, redisAddr, cacheDuration)
	transactionService := transaction.NewTransactionService(*transactionRepo, *accountRepo, natsClient, redisAddr, cacheDuration)

	accountHandler := account.NewAccountHandler(*accountService)
	transactionHandler := transaction.NewTransactionHandler(*transactionService)

	router := mux.NewRouter()

	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		var user account.User
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		login, err := accountRepo.ValidateUserPassword(user.Username, user.Password)
		if err != nil || !login {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		token, _ := GenerateJWT(user.Username)
		w.Write([]byte(token))
	}).Methods("POST")

	protected := router.PathPrefix("/api").Subrouter()
	protected.Use(middleware.JWTAuthMiddleware)
	protected.HandleFunc("/accounts", accountHandler.CreateAccount).Methods("POST")
	protected.HandleFunc("/accounts/{id}", accountHandler.GetAccountBalance).Methods("GET")
	protected.HandleFunc("/transactions", transactionHandler.TransferMoney).Methods("POST")
	protected.HandleFunc("/transactions/{account_id}", transactionHandler.GetTransactionHistory).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func GenerateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &middleware.Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(middleware.JWTKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
