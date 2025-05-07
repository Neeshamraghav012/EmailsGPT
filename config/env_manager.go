package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Panic("Please set the environment variable " + key)
	}
	return val
}

func init() {
	loadEnvs()
}

func loadEnvs() (err error) {
	err = godotenv.Load()
	return
}
