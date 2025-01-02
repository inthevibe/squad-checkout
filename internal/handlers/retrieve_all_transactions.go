package handlers

import (
	"database/sql"
	"fmt"
	"github.com/rivo/tview"
	"squad-checkout/internal/repositories"
	"squad-checkout/internal/services"
)

func ShowAllTransactions(app *tview.Application, db *sql.DB, returnHandler func()) {
	// Retrieve all transactions
	repo := repositories.NewTransactionRepository(db)
	service := services.NewTransactionService(repo)
	transactions, err := service.RetrieveAllTransactions()
	if err != nil {
		showError(app, nil, fmt.Sprintf("Failed to retrieve transactions: %v", err))
		return
	}

	// Create table and set its properties
	table := tview.NewTable().
		SetBorders(true)

	// Add headers
	headers := []string{"ID", "Description", "Date", "Amount (USD)"}
	for col, header := range headers {
		table.SetCell(0, col, tview.NewTableCell(header).
			SetSelectable(false).
			SetAlign(tview.AlignCenter).
			SetTextColor(tview.Styles.SecondaryTextColor))
	}

	// Add transaction data
	for i, transaction := range transactions {
		row := i + 1
		table.SetCell(row, 0, tview.NewTableCell(transaction.ID))
		table.SetCell(row, 1, tview.NewTableCell(transaction.Description))
		table.SetCell(row, 2, tview.NewTableCell(transaction.TransactionDate.Format("2006-01-02")))
		table.SetCell(row, 3, tview.NewTableCell(fmt.Sprintf("%.2f", transaction.AmountUSD)))
	}

	// Create a list for navigation
	nav := tview.NewList().
		AddItem("Back to Main Menu", "", 'b', func() {
			returnHandler() // Invoke the return handler to go back to the main menu
		})

	// Create the main layout
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(table, 0, 9, false). // Set the table not to grab focus
		AddItem(nav, 3, 1, true)     // Set the nav list to grab focus

	// Set the focus explicitly to the navigation list
	app.SetRoot(flex, true).SetFocus(nav)
}
