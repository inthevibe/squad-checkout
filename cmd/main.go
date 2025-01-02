package main

import (
	"github.com/rivo/tview"
	"log"
	"squad-checkout/internal/handlers"
	"squad-checkout/internal/repositories"
)

func main() {
	db, err := repositories.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	app := tview.NewApplication()

	mainMenu := handlers.NewMainMenu(app, db)
	if err := app.SetRoot(mainMenu, true).Run(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}
