package main

import (
	"log"
	"os"

	"github.com/jibbolo/cwdash/pkg/app"
	"github.com/jibbolo/cwdash/pkg/server"
)

func main() {
	log.Println("build", Build)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	a, err := app.New(Build)
	if err != nil {
		log.Fatal(err)
	}

	tsAuthKey := os.Getenv("TS_AUTHKEY")
	if tsAuthKey != "" {
		log.Fatal(server.TailscaleServer(tsAuthKey, "cwdash", a.Handler()))
	} else {
		log.Fatal(server.HTTPStandardServer(":"+port, a.Handler()))
	}

}
