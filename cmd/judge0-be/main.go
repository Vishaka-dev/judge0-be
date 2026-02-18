package main

import (
	"log"

	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/internal/config"
	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/internal/supabase"
	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/server"
)

func main() {
	config.Load()
	supabase.Init()

	app := server.NewRouter()
	log.Println("Running local server on http://localhost:3000")
	app.Run(":3000")
}
