package utils

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/xksd/OrderBookRecorder/storage"
)

type Settings struct {
	Timeframe       int
	Prod            bool
	CSVFolder       string
	UserSymbolsList []string
	ApiKey          string
	ApiSecret       string
}

type settingsJson struct {
	Timeframe int      `json:"timeframeInSeconds"`
	CSVFolder string   `json:"csvFolder"`
	Symbols   []string `json:"symbols"`
}

func ReadSettings() Settings {
	// declare empty struct
	var settings Settings

	// Paths to files
	envFilePath := "./.env"
	settingsFilePath := "./settings.json"

	// Check if settings files exist
	// if not - try to find alternative path up to 3 levels above
	storage.SearchAlternativeFilePath(&envFilePath, 3)
	storage.SearchAlternativeFilePath(&settingsFilePath, 3)

	// Init godotenv library
	err := godotenv.Load(envFilePath)
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

	content, err := os.ReadFile(settingsFilePath)
	if err != nil {
		log.Fatal("Error loading JSON Settings file!")
		log.Fatal(err)
	}

	// Unmarshall the data from settings.json
	var payload settingsJson
	if err = json.Unmarshal(content, &payload); err != nil {
		log.Fatal("Error unmarshalling settings.json")
		log.Fatal(err)
	}

	settings.UserSymbolsList = payload.Symbols
	settings.Timeframe = payload.Timeframe
	settings.CSVFolder = payload.CSVFolder

	return settings
}
