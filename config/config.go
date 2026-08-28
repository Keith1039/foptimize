package config

import (
	"github.com/joho/godotenv"
	"log"
)

var (
	ENV         string
	SERPAPI_KEY string
)

func init() {
	env, err := godotenv.Read()
	if err != nil {
		log.Fatal(err)
	}
	ENV = env["ENV"]
	SERPAPI_KEY = env["SERPAPI_KEY"]
}
