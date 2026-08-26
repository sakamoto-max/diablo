package main

import (
	"github.com/sakamoto-max/diablo/internal/app"
	"github.com/sakamoto-max/diablo/internal/config"
)

func main() {
	config := config.NewConfig()

	app := app.NewApp(config)
	app.Run()
}