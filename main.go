package main

import (
	"log"
	"os"

	"github.com/jibbolo/cwdash/internal/app"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	log.Println("build", Build)
	log.Fatal(app.New(Build, port).Run())
}
