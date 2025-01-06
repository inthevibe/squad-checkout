package repositories

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./database/squad-checkout.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	query := `
    CREATE TABLE IF NOT EXISTS transactions (
        id TEXT PRIMARY KEY,
        description TEXT NOT NULL,
        transaction_date TEXT NOT NULL,
        amount_usd REAL NOT NULL
    );`
	_, err = db.Exec(query)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactions table: %v", err)
	}

	return db, nil
}
