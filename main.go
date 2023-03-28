package main

import (
	"fmt"
	"log"
	"os"

	ProcessExcelFile "github.com/soudengwu/SortProvider"
) 

func main() {
	if len(os.Args); != 3 {
		log.Fatalf("Usage: %s input.xlsx output.xlsx", os.Arga[0])
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	if err := ExcelReader.ProcessExcelFile(inputFile, outputFile); err != nil{
		log.Fatalf("Error processing Excel file: %v", err)
	}

	fmt.Println("Email Providers have been added to the output file")
}