package main

import (
	"log"

	"github.com/ecetinerdem/starthub-backend/internal/app"
	"github.com/ecetinerdem/starthub-backend/internal/database"
)

func main() {
	db := database.ConnectDB()
	database.RunMigrations(db)
	app := app.Init(db)

	log.Fatal(app.Listen(":3000"))
}
