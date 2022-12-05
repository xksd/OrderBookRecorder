package storage

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/xksd/OrderBookRecorder/data"
)

// If a CSV file not exists for the current day, it will be created
// if file exists, new lines will be appended to the end of file.
func SaveToCSV(ob *data.SymbolObSnapshot, folderName string) error {
	// datetime YYMMDD for csv-file name
	year, month, day := time.Now().Date()
	dateName := fmt.Sprintf("%d_%d_%d__", year, month, day)

	// Compose a file name
	fileName := dateName + ob.Symbol + ".csv"

	// Open CSV file or create a new one
	csvFile, err := OpenCSVFile(folderName, fileName)
	if err != nil {
		log.Fatal("OpenCSVFile() FAILED : ", err)
	}
	defer csvFile.Close()

	// Create writer
	w := csv.NewWriter(csvFile)
	defer w.Flush()

	// Convert (ob SymbolObSnapshot) to CSV format:
	// (timestamp, side, price, quantity)
	csvLines, err := ob.PrepareForCsv()
	if err != nil {
		log.Fatal("SaveToCSV() Failed")
		log.Fatal(err)
	}

	if err := w.WriteAll(csvLines); err != nil {
		log.Fatal("Error writing file!")
		return errors.New("SaveToCSV() => Error writing file!")
	}
	return nil

}

// Open OR Create a CSV file and return pointer to File
func OpenCSVFile(folderName string, fileName string) (*os.File, error) {

	// Check if folder exists
	if _, err := os.Stat(folderName); err != nil {
		os.Mkdir(folderName, 0755)
	}
	// Create Writer to file
	var csvFile *os.File
	var err error

	if _, err := os.Stat(folderName + "/" + fileName); err != nil {
		// Create CSV file
		// Check if file exists
		fmt.Println("Created new CSV File: " + folderName + "/" + fileName)
		csvFile, err = os.Create(folderName + "/" + fileName)
		if err != nil {
			log.Fatalf("Failed creating file: %s", err)
			return nil, errors.New("Failed creating file: %s")
		}
	} else {
		csvFile, err = os.OpenFile(folderName+"/"+fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Fatalf("Failed creating file: %s", err)
			return nil, errors.New("Failed creating file: %s")
		}
	}
	if err != nil {
		log.Fatal("Uncaught error while opening the CSV file!")
		log.Fatal(err)
		return nil, errors.New("Uncaught error while opening the CSV file!")
	}

	return csvFile, nil
}
