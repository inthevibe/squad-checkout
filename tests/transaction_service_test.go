package tests

import (
	"database/sql"
	"squad-checkout/internal/models"
	"squad-checkout/internal/repositories"
	"squad-checkout/internal/services"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	_ "modernc.org/sqlite" // Import the pure Go SQLite driver
)

func TestTransactionService_StoreAndRetrieve(t *testing.T) {
	// Set up an in-memory SQLite database for testing
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Create the transactions table
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

	// Initialize the repository and service
	repo := repositories.NewTransactionRepository(db)
	service := services.NewTransactionService(repo)

	// Test data
	transaction := models.Transaction{
		ID:              "123e4567-e89b-12d3-a456-426614174000",
		Description:     "Groceries",
		TransactionDate: time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
		AmountUSD:       50.00,
	}

	// Test storing a transaction
	err = service.StoreTransaction(transaction)
	assert.NoError(t, err)

	// Test retrieving the transaction
	retrievedTransaction, err := service.RetrieveTransaction(transaction.ID)
	assert.NoError(t, err)
	assert.Equal(t, transaction.ID, retrievedTransaction.ID)
	assert.Equal(t, transaction.Description, retrievedTransaction.Description)
	assert.Equal(t, transaction.TransactionDate, retrievedTransaction.TransactionDate)
	assert.Equal(t, transaction.AmountUSD, retrievedTransaction.AmountUSD)
}
