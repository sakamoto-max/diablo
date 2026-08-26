package config

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/sakamoto-max/diablo/internal/env"
)

type Config struct {
	Primary string `validate:"required"`
	Http    http
	Pg      pg
	Jwt     jwt
}

type http struct {
	Port string `validate:"required"`
}

type pg struct {
	Host         string `validate:"required"`
	Port         string `validate:"required"`
	User         string `validate:"required"`
	Password     string `validate:"required"`
	Database     string `validate:"required"`
	SSLMode      string `validate:"required"`
	MaxOpenConns int    `validate:"required"`
	MaxIdleConns int    `validate:"required"`
	MaxLifetime  int    `validate:"required"`
	MaxIdleTime  int    `validate:"required"`
}

type jwt struct {
	SecretKey string `validate:"required"`
}

func NewConfig() *Config {
	env.LoadEnv("../../app.env")

	primary := os.Getenv("PRIMARY")

	http := http{
		Port: os.Getenv("HTTP_PORT"),
	}

	pgMaxOpenConnsStr := os.Getenv("PG_MAX_OPEN_CONNS")
	pgMaxOpenConns, err := strconv.Atoi(pgMaxOpenConnsStr)
	if err != nil {
		log.Fatalf("failed to convert PG_MAX_OPEN_CONNS to int : %v", err)
	}

	pgMaxIdleConnsStr := os.Getenv("PG_MAX_IDLE_CONNS")
	pgMaxIdleConns, err := strconv.Atoi(pgMaxIdleConnsStr)
	if err != nil {
		log.Fatalf("failed to convert PG_MAX_IDLE_CONNS to int : %v", err)
	}

	pgMaxLifetimeStr := os.Getenv("PG_MAX_LIFETIME")
	pgMaxLifetime, err := strconv.Atoi(pgMaxLifetimeStr)
	if err != nil {
		log.Fatalf("failed to convert PG_MAX_LIFETIME to int : %v", err)
	}

	pgMaxIdleTimeStr := os.Getenv("PG_MAX_IDLE_TIME")
	pgMaxIdleTime, err := strconv.Atoi(pgMaxIdleTimeStr)
	if err != nil {
		log.Fatalf("failed to convert PG_MAX_IDLE_TIME to int : %v", err)
	}

	pg := pg{
		Host:         os.Getenv("PG_HOST"),
		Port:         os.Getenv("PG_PORT"),
		User:         os.Getenv("PG_USER"),
		Password:     os.Getenv("PG_PASSWORD"),
		Database:     os.Getenv("PG_DATABASE"),
		SSLMode:      os.Getenv("PG_SSL_MODE"),
		MaxOpenConns: pgMaxOpenConns,
		MaxIdleConns: pgMaxIdleConns,
		MaxLifetime:  pgMaxLifetime,
		MaxIdleTime:  pgMaxIdleTime,
	}

	jwt := jwt{
		SecretKey: os.Getenv("SECRET_KEY"),
	}

	config := Config{
		Primary: primary,
		Http:    http,
		Pg:      pg,
		Jwt:     jwt,
	}


	newValidator := validator.New(validator.WithRequiredStructEnabled())

	err = newValidator.Struct(config)
	if err != nil {
		var validatorErr validator.ValidationErrors
		if errors.As(err, &validatorErr) {
			log.Fatal(validatorErr.Error())
		}
	}

	return &config
}
