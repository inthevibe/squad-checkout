package handlers

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/rivo/tview"
	"squad-checkout/internal/models"
	"squad-checkout/internal/repositories"
	"squad-checkout/internal/services"
	"time"
)

func ShowStoreTransactionForm(app *tview.Application, db *sql.DB, returnHandler func()) {
	form := tview.NewForm()
	form.SetBorder(true).SetTitle("Store a Transaction").SetTitleAlign(tview.AlignLeft)

	var description, dateStr, amountStr string

	form.AddInputField("Description", "", 50, nil, func(text string) {
		description = text
	}).
		AddInputField("Date (YYYY-MM-DD)", "", 10, nil, func(text string) {
			dateStr = text
		}).
		AddInputField("Amount (USD)", "", 10, nil, func(text string) {
			amountStr = text
		}).
		AddButton("Save", func() {
			// Validate and save the transaction
			transaction, err := validateAndCreateTransaction(description, dateStr, amountStr)
			if err != nil {
				showError(app, form, err.Error())
				return
			}

			// Save the transaction to the database
			repo := repositories.NewTransactionRepository(db)
			service := services.NewTransactionService(repo)
			if err := service.StoreTransaction(*transaction); err != nil {
				showError(app, form, fmt.Sprintf("Failed to save transaction: %v", err))
				return
			}

			// Return to the main menu
			returnHandler()
		}).
		AddButton("Cancel", func() {
			returnHandler()
		})

	app.SetRoot(form, true).SetFocus(form)
}

func validateAndCreateTransaction(description, dateStr, amountStr string) (*models.Transaction, error) {
	// Validate description
	if len(description) == 0 || len(description) > 50 {
		return nil, fmt.Errorf("description must be between 1 and 50 characters")
	}

	// Validate date
	transactionDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid date format (expected YYYY-MM-DD)")
	}

	// Validate amount
	var amountUSD float64
	if _, err := fmt.Sscanf(amountStr, "%f", &amountUSD); err != nil || amountUSD <= 0 {
		return nil, fmt.Errorf("amount must be a positive number")
	}

	// Create the transaction
	return &models.Transaction{
		ID:              generateUUID(), // TODO: Implement UUID generation
		Description:     description,
		TransactionDate: transactionDate,
		AmountUSD:       amountUSD,
	}, nil
}

func showError(app *tview.Application, currentPage tview.Primitive, message string) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			app.SetRoot(currentPage, true).SetFocus(currentPage)
		})
	app.SetRoot(modal, false).SetFocus(modal)
}

func generateUUID() string {
	return uuid.New().String()
}
