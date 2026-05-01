package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func main() {
	fmt.Println("Database utils")
	file, err := os.Open("../data/uscities.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Records type: %T\n", records)
}
