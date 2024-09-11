package main

import "github.com/eugenedhz/auth_service_test/internal/app"

func main() {
	app := app.NewApp()
	err := app.Run()
	if err != nil {
		panic(err)
	}
}
