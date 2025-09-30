package main

import (
	"log"

	"github.com/pavel97go/npa-service/internal/app"
)

func main() {
	srv := app.BuildServer()
	if err := srv.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
