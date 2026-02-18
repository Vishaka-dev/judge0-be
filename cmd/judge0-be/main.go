package main

import (
	"log"

	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/config"
	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/database"
	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/server"
)

func main() {
	config.Load()
	database.Init()

	app := server.NewRouter()
	log.Println("Running local server on http://localhost:3000")
	app.Run(":3000")
}
