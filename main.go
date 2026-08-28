package main

import (
	"encoding/json"
	"fmt"
	"github.com/Keith1039/foptimize/schema"
	"io"
	"log"
	"os"
)

func GetUsers() schema.Users {
	// Open our jsonFile
	jsonFile, err := os.Open("users.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Fatal()
	}

	fmt.Println("Successfully Opened users.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened jsonFile as a byte array.
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatal("Error reading file:", err)
	}

	// we initialize our Users array
	var users schema.Users

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	err = json.Unmarshal(byteValue, &users)
	if err != nil {
		log.Fatal("Error unmarshalling JSON:", err)
	}
	return users
}

func MarshalConfig(user schema.User) map[string]any {
	config, err := json.Marshal(user.Config)
	if err != nil {
		log.Fatal("Error marshalling JSON:", err)
	}
	var m map[string]any
	err = json.Unmarshal(config, &m)
	if err != nil {
		log.Fatal("Error unmarshalling JSON:", err)
	}
	return m
}

func main() {
	// Open our jsonFile
	jsonFile, err := os.Open("example.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Fatal()
	}

	fmt.Println("Successfully Opened users.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened jsonFile as a byte array.
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatal("Error reading file:", err)
	}

	var deals schema.Deals

	err = json.Unmarshal(byteValue, &deals)
	if err != nil {
		log.Fatal("Error unmarshalling JSON:", err)
	}
	for _, deal := range deals.Deals {
		indent, err := json.MarshalIndent(deal, "", "  ")
		if err != nil {
			return
		}
		fmt.Println(string(indent))
	}
}
