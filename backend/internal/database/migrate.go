package database

import (
	"context"
	"fmt"
	"os"
)

func RunMigrations() {
	fmt.Println("📄 Running database migrations...")

	content, err := os.ReadFile("sql/schema.sql")

	if err != nil {
		panic("Failed to read migration file from sql: " + err.Error())
	}

	_, err = DB.Exec(context.Background(), string(content))

	if err != nil {
		panic("Failed to read migration file from sql: " + err.Error())
	}

	fmt.Println("✅ Migrations completed successfully!")

}
