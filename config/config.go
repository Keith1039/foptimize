package config

import (
	"github.com/joho/godotenv"
	"log"
)

var (
	ENV         string
	SERPAPI_KEY string
	DB_URL      string
)

func init() {
	env, err := godotenv.Read("../.env")
	if err != nil {
		log.Fatal(err)
	}
	ENV = env["ENV"]
	SERPAPI_KEY = env["SERPAPI_KEY"]
	DB_URL = env["DB_URL"]
}
