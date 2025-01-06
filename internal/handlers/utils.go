package handlers

import (
	"github.com/rivo/tview"
)

func showError(app *tview.Application, currentPage tview.Primitive, message string) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			app.SetRoot(currentPage, true).SetFocus(currentPage)
		})
	app.SetRoot(modal, false).SetFocus(modal)
}
