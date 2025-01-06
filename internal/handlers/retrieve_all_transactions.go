package handlers

import (
	"database/sql"
	"fmt"
	"squad-checkout/internal/repositories"
	"squad-checkout/internal/services"
	"squad-checkout/internal/utils"

	"github.com/rivo/tview"
)

func ShowAllTransactions(app *tview.Application, db *sql.DB, returnHandler func()) {
	repo := repositories.NewTransactionRepository(db)
	service := services.NewTransactionService(repo)
	transactions, err := service.RetrieveAllTransactions()
	if err != nil {
		showError(app, nil, fmt.Sprintf("Failed to retrieve transactions: %v", err))
		return
	}

	currencies, err := utils.GetSupportedCurrencies()
	if err != nil {
		showError(app, nil, fmt.Sprintf("Failed to fetch supported currencies: %v", err))
		return
	}

	table := tview.NewTable()
	table.SetBorders(true)
	table.SetTitle(" All Transactions ")
	table.SetTitleAlign(tview.AlignLeft)

	headers := []string{"ID", "Description", "Date", "Amount (USD)", "Currency", "Converted Amount"}
	for col, header := range headers {
		table.SetCell(0, col, tview.NewTableCell(header).
			SetSelectable(false).
			SetAlign(tview.AlignCenter).
			SetTextColor(tview.Styles.SecondaryTextColor).
			SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor))
	}

	for i, transaction := range transactions {
		row := i + 1
		table.SetCell(row, 0, tview.NewTableCell(transaction.ID).SetAlign(tview.AlignLeft))
		table.SetCell(row, 1, tview.NewTableCell(transaction.Description).SetAlign(tview.AlignLeft))
		table.SetCell(row, 2, tview.NewTableCell(transaction.TransactionDate.Format("02-01-2006")).SetAlign(tview.AlignCenter))
		table.SetCell(row, 3, tview.NewTableCell(fmt.Sprintf("%.2f", transaction.AmountUSD)).SetAlign(tview.AlignRight))

		selectedCurrency := currencies[28]
		exchangeRate, err := utils.GetExchangeRate(selectedCurrency, transaction.TransactionDate)
		if err != nil {
			table.SetCell(row, 4, tview.NewTableCell("N/A").SetAlign(tview.AlignCenter))
			table.SetCell(row, 5, tview.NewTableCell("N/A").SetAlign(tview.AlignRight))
		} else {
			convertedAmount := transaction.AmountUSD * exchangeRate
			table.SetCell(row, 4, tview.NewTableCell(selectedCurrency).SetAlign(tview.AlignCenter))
			table.SetCell(row, 5, tview.NewTableCell(fmt.Sprintf("%.2f", convertedAmount)).SetAlign(tview.AlignRight))
		}
	}

	nav := tview.NewList().
		AddItem("Back to Main Menu", "", 'b', func() {
			returnHandler()
		})

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(table, 0, 9, false).
		AddItem(nav, 3, 1, true)

	app.SetRoot(flex, true).SetFocus(nav)
}
