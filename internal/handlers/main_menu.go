package handlers

import (
	"database/sql"

	"github.com/rivo/tview"
)

func NewMainMenu(app *tview.Application, db *sql.DB) *tview.List {
	list := tview.NewList()
	list.AddItem("Store a Transaction", "", '1', func() {
		ShowStoreTransactionForm(app, db, func() {
			app.SetRoot(list, true).SetFocus(list)
		})
	})
	list.AddItem("Retrieve a Transaction", "", '2', func() {
		ShowRetrieveTransactionForm(app, db, func() {
			app.SetRoot(list, true).SetFocus(list)
		})
	})
	list.AddItem("Retrieve All Transactions", "", '3', func() {
		ShowAllTransactions(app, db, func() {
			app.SetRoot(list, true).SetFocus(list)
		})
	})
	list.AddItem("Exit", "", 'q', func() {
		app.Stop()
	})
	return list
}
