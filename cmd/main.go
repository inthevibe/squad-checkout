package main

import (
	"github.com/rivo/tview"
	"log"
	"squad-checkout/internal/handlers"
)

func main() {
	app := tview.NewApplication()

	mainMenu() := handlers.NewMainMenu(app)
	if err := app.SetRoot(mainMenu, true).Run(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}
