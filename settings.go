package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func envVariable(key string) string {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

// Env var
var Env = map[string]string{
	"GO_ENV":     envVariable("GO_ENV"),
	"INSTANCE":   envVariable("INSTANCE"),
	"URL":        envVariable("URL"),
	"PORT":       envVariable("PORT"),
	"JWT_ISSUER": envVariable("JWT_ISSUER"),
	"JWT_SECRET": envVariable("JWT_SECRET"),
}
