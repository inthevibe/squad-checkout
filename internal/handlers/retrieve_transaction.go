package handlers

import (
	"database/sql"
	"fmt"
	"github.com/rivo/tview"
	"squad-checkout/internal/models"
	"squad-checkout/internal/repositories"
	"squad-checkout/internal/services"
)

func ShowRetrieveTransactionForm(app *tview.Application, db *sql.DB, returnHandler func()) {
	form := tview.NewForm()
	form.SetBorder(true).SetTitle("Retrieve a Transaction").SetTitleAlign(tview.AlignLeft)

	var id string

	form.AddInputField("Transaction ID", "", 36, nil, func(text string) {
		id = text
	}).
		AddButton("Retrieve", func() {
			// Retrieve the transaction from the database
			repo := repositories.NewTransactionRepository(db)
			service := services.NewTransactionService(repo)
			transaction, err := service.RetrieveTransaction(id)
			if err != nil {
				showError(app, form, fmt.Sprintf("Failed to retrieve transaction: %v", err))
				return
			}

			// Display the transaction details
			showTransactionDetails(app, form, transaction)
		}).
		AddButton("Cancel", func() {
			returnHandler()
		})

	app.SetRoot(form, true).SetFocus(form)
}

func showTransactionDetails(app *tview.Application, currentPage tview.Primitive, transaction *models.Transaction) {
	details := fmt.Sprintf(
		"ID: %s\nDescription: %s\nDate: %s\nAmount (USD): %.2f",
		transaction.ID,
		transaction.Description,
		transaction.TransactionDate.Format("2006-01-02"),
		transaction.AmountUSD,
	)

	modal := tview.NewModal().
		SetText(details).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			app.SetRoot(currentPage, true).SetFocus(currentPage)
		})
	app.SetRoot(modal, false).SetFocus(modal)
}
