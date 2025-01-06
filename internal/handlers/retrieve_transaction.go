package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"squad-checkout/internal/models"
	"squad-checkout/internal/repositories"
	"squad-checkout/internal/services"
	"squad-checkout/internal/utils"

	"github.com/rivo/tview"
)

func ShowRetrieveTransactionForm(app *tview.Application, db *sql.DB, returnHandler func()) {
	form := tview.NewForm()
	form.SetBorder(true).SetTitle("Retrieve a Transaction").SetTitleAlign(tview.AlignLeft)

	var id string
	var selectedCurrency string
	var isProcessing bool

	currencies, err := utils.GetSupportedCurrencies()
	if err != nil {
		showError(app, form, fmt.Sprintf("Failed to fetch supported currencies: %v", err))
		return
	}

	form.AddInputField("Transaction ID", "", 36, nil, func(text string) {
		id = text
	}).
		AddDropDown("Currency", currencies, 28, func(option string, index int) {
			selectedCurrency = option
		}).
		AddButton("Retrieve", func() {
			if isProcessing {
				return
			}
			isProcessing = true

			loadingModal := tview.NewModal().
				SetText("Retrieving transaction... Please wait.").
				AddButtons([]string{}).
				SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor)
			app.SetRoot(loadingModal, false).SetFocus(loadingModal)

			go func() {
				repo := repositories.NewTransactionRepository(db)
				service := services.NewTransactionService(repo)
				transaction, err := service.RetrieveTransaction(id)
				if err != nil {
					app.QueueUpdateDraw(func() {
						isProcessing = false
						showError(app, form, fmt.Sprintf("Failed to retrieve transaction: %v", err))
					})
					return
				}

				exchangeRate, err := utils.GetExchangeRate(selectedCurrency, transaction.TransactionDate)
				if err != nil {
					app.QueueUpdateDraw(func() {
						isProcessing = false
						if errors.Is(err, errors.New("API unresponsive: request timed out")) {
							showError(app, form, "The API is unresponsive. Please try again later.")
						} else {
							showError(app, form, fmt.Sprintf("Failed to fetch exchange rate: %v", err))
						}
					})
					return
				}

				convertedAmount := transaction.AmountUSD * exchangeRate

				app.QueueUpdateDraw(func() {
					isProcessing = false
					showTransactionDetails(app, form, transaction, selectedCurrency, exchangeRate, convertedAmount)
				})
			}()
		}).
		AddButton("Cancel", func() {
			returnHandler()
		})

	app.SetRoot(form, true).SetFocus(form)
}

func showTransactionDetails(app *tview.Application, currentPage tview.Primitive, transaction *models.Transaction, currency string, exchangeRate, convertedAmount float64) {
	details := fmt.Sprintf(
		"ID: %s\nDescription: %s\nDate: %s\nAmount (USD): %.2f\nExchange Rate: %.4f\nAmount (%s): %.2f",
		transaction.ID,
		transaction.Description,
		transaction.TransactionDate.Format("02-01-2006"),
		transaction.AmountUSD,
		exchangeRate,
		currency,
		convertedAmount,
	)

	modal := tview.NewModal().
		SetText(details).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			app.SetRoot(currentPage, true).SetFocus(currentPage)
		})
	app.SetRoot(modal, false).SetFocus(modal)
}
