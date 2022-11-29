package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Settings struct {
	Prod            bool
	UserSymbolsList []string
	ApiKey          string
	ApiSecret       string
}

type userSymbolsList struct {
	Symbols []string
}

func GetSettings() Settings {
	// declare empty struct
	var settings Settings

	// Init godotenv library
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// extract PROD env variable to bool
	prodString := strings.ToLower(os.Getenv("PROD"))
	if prodString == "true" {
		settings.Prod = true
	} else if prodString == "false" {
		settings.Prod = false
	} else {
		log.Fatal("Somethings wrong with .ENV / PROD variable!")
		log.Fatal("Please make sure to set PROD to true or false")
	}

	// extract API_KEY and API_SECRET
	settings.ApiKey = os.Getenv("API_KEY")
	settings.ApiSecret = os.Getenv("API_SECRET")

	// extract userSymbolsList from Json file
	content, err := ioutil.ReadFile("./settings.json")
	if err != nil {
		log.Fatal("Error loading JSON Settings file!")
		log.Fatal(err)
	}

	// Unmarshall the data from settings.json
	var payload userSymbolsList
	if err = json.Unmarshal(content, &payload); err != nil {
		log.Fatal("Error unmarshalling settings.json")
		log.Fatal(err)
	}

	settings.UserSymbolsList = payload.Symbols

	return settings
}
