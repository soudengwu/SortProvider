package main

import (
	"fmt"
	"log"
	"os"

	ExcelReader "github.com/soudengwu/SortProvider/Reader"
)

func main() {
	if len(os.Args) != 4 {
		log.Fatalf("Usage: %s input.xlsx output.xlsx noMxOutput.xlsx", os.Args[0])
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]
	noMxOutputFile := os.Args[3]

	if err := ExcelReader.ProcessExcelFile(inputFile, outputFile, noMxOutputFile); err != nil {
		log.Fatalf("Error processing Excel file: %v", err)
	}

	fmt.Println("Email Providers have been added to the output file")
}
