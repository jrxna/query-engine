package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	SERVICE_NAME        string
	SERVICE_ENVIRONMENT string
	SERVICE_PORT        string
	DATABASE_URL        string
	JWT_SECRET_KEY      string
}

var config Config

func Init() {
	serviceEnvironment := os.Getenv("SERVICE_ENVIRONMENT")
	if len(serviceEnvironment) == 0 {
		serviceEnvironment = "development"
	}

	configFilePath := ".env"

	err := godotenv.Load(configFilePath)
	if err != nil {
		panic(err.Error())
	}
	config.SERVICE_NAME = os.Getenv("SERVICE_NAME")
	config.SERVICE_ENVIRONMENT = serviceEnvironment
	config.SERVICE_PORT = os.Getenv("SERVICE_PORT")
	config.DATABASE_URL = os.Getenv("DATABASE_URL")
	config.JWT_SECRET_KEY = os.Getenv("JWT_SECRET_KEY")

	fmt.Println(config)
}

func GetConfig() Config {
	return config
}
