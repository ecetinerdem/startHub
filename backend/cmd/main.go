package main

import (
	"log"

	"github.com/ecetinerdem/starthub-backend/internal/app"
)

func main() {
	app := app.Init()

	log.Fatal(app.Listen(":3000"))
}
