package config

import (
	"errors"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/sakamoto-max/diablo/internal/env"
)

type Config struct {
	Primary string `validate:"required"`
	Http    http
	Db      db
}

type http struct {
	Port string `validate:"required"`
}

type db struct {
	Up   string `validate:"required"`
	Down string `validate:"required"`
}

func NewConfig() *Config {
	env.LoadEnv("../../app.env")

	primary := os.Getenv("PRIMARY")

	http := http{
		Port: os.Getenv("HTTP_PORT"),
	}

	db := db{
		Up:   os.Getenv("DB_UP"),
		Down: os.Getenv("DB_DOWN"),
	}

	config := Config{
		Primary: primary,
		Http:    http,
		Db:      db,
	}

	newValidator := validator.New(validator.WithRequiredStructEnabled())

	err := newValidator.Struct(config)
	if err != nil {
		var validatorErr validator.ValidationErrors
		if errors.As(err, &validatorErr) {
			log.Fatal(validatorErr.Error())
		}
	}

	return &config
}
