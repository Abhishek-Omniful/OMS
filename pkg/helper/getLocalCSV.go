package helper

import (
	"log"
	"os"
)

func GetLocalCSV() ([]byte, error) {
	filePath := "C:/Users/Abhishek/Desktop/orders.csv"
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return fileBytes, nil
}
