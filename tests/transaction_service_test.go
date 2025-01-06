package tests

import (
	"database/sql"
	"squad-checkout/internal/models"
	"squad-checkout/internal/repositories"
	"squad-checkout/internal/services"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	_ "modernc.org/sqlite"
)

func TestTransactionService_StoreAndRetrieve(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE transactions (
			id TEXT PRIMARY KEY,
			description TEXT NOT NULL,
			transaction_date TEXT NOT NULL,
			amount_usd REAL NOT NULL
		);
	`)
	if err != nil {
		t.Fatalf("Failed to create transactions table: %v", err)
	}

	repo := repositories.NewTransactionRepository(db)
	service := services.NewTransactionService(repo)

	transaction := models.Transaction{
		ID:              "123e4567-e89b-12d3-a456-426614174000",
		Description:     "Groceries",
		TransactionDate: time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
		AmountUSD:       50.00,
	}

	err = service.StoreTransaction(transaction)
	assert.NoError(t, err)

	retrievedTransaction, err := service.RetrieveTransaction(transaction.ID)
	assert.NoError(t, err)
	assert.Equal(t, transaction.ID, retrievedTransaction.ID)
	assert.Equal(t, transaction.Description, retrievedTransaction.Description)
	assert.Equal(t, transaction.TransactionDate, retrievedTransaction.TransactionDate)
	assert.Equal(t, transaction.AmountUSD, retrievedTransaction.AmountUSD)
}
