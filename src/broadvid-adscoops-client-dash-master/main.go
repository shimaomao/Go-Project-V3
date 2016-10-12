package main

import (
	"os"

	"app"
)

func main() {
	env := os.Getenv("MARTINI_ENV")

	appPort := os.Getenv("PORT")

	if appPort == "" {
		appPort = "3000"
	}

	isTesting := true

	if env == "production" {
		isTesting = false
	}
	if m, err := app.App(isTesting); err == nil { // Pass 'false' to app if you're not testing

		m.Run()
	}
}
