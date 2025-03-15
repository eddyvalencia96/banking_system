# Banking System API

This repository contains a RESTful API for a banking system that allows users to manage accounts, perform transactions, and more.

## Features

- User authentication and authorization
- Account management (create, read)
- Transaction operations (transfers, transaction history)

## Getting Started

### Prerequisites

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)

### Running the Application with Docker

1. Clone this repository:
   ```bash
   git clone https://github.com/eddyvalencia96/banking_system
   cd banking-system
   ```

2. Build and start the containers:
   ```bash
   docker-compose up --build
   ```

3. The API will be available at `http://localhost:8080`

## API Documentation

### Authentication

#### Login
```
POST /login
```
Request body:
```json
{
  "email": "edvalencia",
  "password": "960919edval"
}
```

### Accounts

#### Create account
```
POST /api/accounts
```

#### Get account by ID
```
GET /api/accounts/{id}
```

### Transactions

#### Transfer money
```
POST /api/transactions/transfer
```
Request body:
```json
{
  "fromAccountId": "123456",
  "toAccountId": "789012",
  "amount": 75.00,
  "description": "Rent payment"
}
```

#### Get transaction history
```
GET /api/transactions/history/{accountId}
```

## Environment Variables

The application uses the following environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_HOST` | Database host | `db` |
| `DB_PORT` | Database port | `5432` |
| `DB_NAME` | Database name | `banking_system` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | `postgres` |
| `JWT_SECRET` | Secret key for JWT tokens | `your_jwt_secret` |
| `PORT` | Application port | `8080` |

## Database Schema

The application uses a relational database with the following main tables:
- `users` - Stores user information
- `accounts` - Stores account details
- `transactions` - Records all transactions
