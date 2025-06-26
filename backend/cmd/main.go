package main

import (
	"log"
	"os"

	"github.com/ecetinerdem/starthub-backend/internal/app"
	"github.com/ecetinerdem/starthub-backend/internal/database"
)

func main() {
	db := database.ConnectDB()
	database.RunMigrations(db)
	app := app.Init(db)

	PORT := os.Getenv("PORT")

	log.Fatal(app.Listen(":" + PORT))
}
