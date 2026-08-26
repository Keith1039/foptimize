package main

import (
	"encoding/json"
	"fmt"
	"github.com/Keith1039/foptimize/schema"
	"github.com/Keith1039/foptimize/serp"
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
	client := serp.GetClient()
	users := GetUsers()
	fmt.Println("Users:")
	fmt.Printf("%+v", users)
	var allDeals []any

	for _, user := range users.Users {
		fmt.Printf("email: %s\n", user.Email)
		config := MarshalConfig(user)
		for _, airport := range user.RelevantAirports {
			params := map[string]string{
				"departure_id": airport,
				"currency":     "CAD",
				"gl":           "ca",
				"hl":           "en",
			}
			for key, value := range config {
				if value.(string) != "" {
					params[key] = value.(string)
				}
			}
			fmt.Println(params)
			result, err := client.Search(params)
			if err != nil {
				log.Fatal(err)
			}
			deals, ok := result["deals"].([]any)
			if !ok {
				log.Fatal("type conversion failed")
			}
			allDeals = append(allDeals, deals...)
		}
	}
	jsonData, err := json.MarshalIndent(map[string][]any{"alldeals": allDeals}, "", " ")
	if err != nil {
		log.Fatal("Error marshalling JSON:", err)
	}
	err = os.WriteFile("alldeals.json", jsonData, os.ModePerm)
}
