package repositories

import (
	"database/sql"
	"fmt"
	"squad-checkout/internal/models"
	"time"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Save(transaction models.Transaction) error {
	query := `
    INSERT INTO transactions (id, description, transaction_date, amount_usd)
    VALUES (?, ?, ?, ?);`
	_, err := r.db.Exec(query, transaction.ID, transaction.Description, transaction.TransactionDate.Format(time.RFC3339), transaction.AmountUSD)
	if err != nil {
		return fmt.Errorf("failed to save transaction: %v", err)
	}
	return nil
}

func (r *TransactionRepository) FindByID(id string) (*models.Transaction, error) {
	query := `
    SELECT id, description, transaction_date, amount_usd
    FROM transactions
    WHERE id = ?;`
	row := r.db.QueryRow(query, id)

	var transaction models.Transaction
	var dateStr string
	err := row.Scan(&transaction.ID, &transaction.Description, &dateStr, &transaction.AmountUSD)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("transaction not found")
		}
		return nil, fmt.Errorf("failed to retrieve transaction: %v", err)
	}

	// Parse the transaction date
	transaction.TransactionDate, err = time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse transaction date: %v", err)
	}

	return &transaction, nil
}

func (r *TransactionRepository) FindAll() ([]models.Transaction, error) {
	query := `
    SELECT id, description, transaction_date, amount_usd
    FROM transactions;`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve transactions: %v", err)
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		var dateStr string
		err := rows.Scan(&transaction.ID, &transaction.Description, &dateStr, &transaction.AmountUSD)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %v", err)
		}

		// Parse the transaction date
		transaction.TransactionDate, err = time.Parse(time.RFC3339, dateStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse transaction date: %v", err)
		}

		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over transactions: %v", err)
	}

	return transactions, nil
}
