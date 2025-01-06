package handlers

import (
	"database/sql"
	"fmt"
	"squad-checkout/internal/models"
	"squad-checkout/internal/repositories"
	"squad-checkout/internal/services"
	"time"

	"github.com/google/uuid"
	"github.com/rivo/tview"
)

func ShowStoreTransactionForm(app *tview.Application, db *sql.DB, returnHandler func()) {
	form := tview.NewForm()
	form.SetBorder(true).SetTitle("Store a Transaction").SetTitleAlign(tview.AlignLeft)

	var description, dateStr, amountStr string

	form.AddInputField("Description", "", 50, nil, func(text string) {
		description = text
	}).
		AddInputField("Date (DD-MM-YYYY)", "", 10, nil, func(text string) {
			dateStr = text
		}).
		AddInputField("Amount (USD)", "", 10, nil, func(text string) {
			amountStr = text
		}).
		AddButton("Save", func() {
			transaction, err := validateAndCreateTransaction(description, dateStr, amountStr)
			if err != nil {
				showError(app, form, err.Error())
				return
			}

			repo := repositories.NewTransactionRepository(db)
			service := services.NewTransactionService(repo)
			if err := service.StoreTransaction(*transaction); err != nil {
				showError(app, form, fmt.Sprintf("Failed to save transaction: %v", err))
				return
			}

			returnHandler()
		}).
		AddButton("Cancel", func() {
			returnHandler()
		})

	app.SetRoot(form, true).SetFocus(form)
}

func validateAndCreateTransaction(description, dateStr, amountStr string) (*models.Transaction, error) {
	if len(description) == 0 || len(description) > 50 {
		return nil, fmt.Errorf("description must be between 1 and 50 characters")
	}

	transactionDate, err := time.Parse("02-01-2006", dateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid date format (expected DD-MM-YYYY)")
	}

	var amountUSD float64
	if _, err := fmt.Sscanf(amountStr, "%f", &amountUSD); err != nil || amountUSD <= 0 {
		return nil, fmt.Errorf("amount must be a positive number")
	}

	// Create the transaction
	return &models.Transaction{
		ID:              generateUUID(),
		Description:     description,
		TransactionDate: transactionDate,
		AmountUSD:       amountUSD,
	}, nil
}

func generateUUID() string {
	return uuid.New().String()
}
