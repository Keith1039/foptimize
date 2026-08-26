package config

import (
	"github.com/joho/godotenv"
	"log"
)

type ProjectConfig struct {
	ENV         string
	SERPAPI_KEY string
}

var Config *ProjectConfig

func init() {
	env, err := godotenv.Read()
	if err != nil {
		log.Fatal(err)
	}
	Config = &ProjectConfig{
		ENV:         env["config"],
		SERPAPI_KEY: env["SERPAPI_KEY"],
	}
}
