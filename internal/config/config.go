package config

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Configs struct {
	DbUser         string
	DbHost         string
	DbPassword     string
	DbPort         int
	DbName         string
	SigningKey     string
	AccessTokenTTL time.Duration
	ServerPort     string
	LogLevel       string
}

func LoadConfigs() (*Configs, error) {
	if err := godotenv.Load(); err != nil {
		return &Configs{}, errors.New("Failed to load env variables")
	}

	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return &Configs{}, errors.New("Invalid DB_PORT variable")
	}

	tokenTTL, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_TTL"))
	if err != nil {
		return &Configs{}, errors.New("Invalid ACCESS_TOKEN_TTL variable")
	}

	return &Configs{
		DbUser:         os.Getenv("DB_USER"),
		DbHost:         os.Getenv("DB_HOST"),
		DbPassword:     os.Getenv("DB_PASSWORD"),
		DbPort:         port,
		DbName:         os.Getenv("DB_NAME"),
		SigningKey:     os.Getenv("SIGNING_KEY"),
		AccessTokenTTL: time.Second * time.Duration(tokenTTL),
		ServerPort:     os.Getenv("SERVER_PORT"),
		LogLevel:       os.Getenv("LOG_LEVEL"),
	}, nil
}
