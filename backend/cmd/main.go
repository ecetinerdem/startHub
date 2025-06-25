package main

import (
	"log"

	"github.com/ecetinerdem/starthub-backend/internal/app"
	"github.com/ecetinerdem/starthub-backend/internal/database"
)

func main() {
	database.ConnectDB()
	database.RunMigrations()
	app := app.Init()

	log.Fatal(app.Listen(":3000"))
}
