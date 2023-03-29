package ExcelReader

import (
	"fmt"
	"log"
	"sync"

	CheckProvider "github.com/soudengwu/SortProvider/SortProvider"
	"github.com/xuri/excelize/v2"
)

func ProcessExcelFile(inputFile, outputFile, noMXOutputFile string) error {
	f, err := excelize.OpenFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open input file: %v", err)
	}

	srcSheet := f.GetSheetName(0)

	noMxFile := excelize.NewFile()
	noMxSheet := "NoMxRecords"
	noMxFile.NewSheet(noMxSheet)
	for _, colName := range []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N"} {
		cellValue, err := f.GetCellValue(srcSheet, fmt.Sprintf("%s1", colName))
		if err != nil {
			fmt.Errorf("Failed to read Cell Value %v", err)
			return err
		}
		noMxFile.SetCellValue(noMxSheet, fmt.Sprintf("%s1", colName), cellValue)
	}

	f.SetCellValue(srcSheet, "01", "Email Provider")

	rows, err := f.GetRows(srcSheet)
	if err != nil {
		return fmt.Errorf("failed to get rows: %v", err)
	}

	var wg sync.WaitGroup
	ch := make(chan bool, 5)

	noMxRowCounter := 2

	for i, row := range rows {
		if i == 0 {
			continue
		}

		wg.Add(1)
		ch <- true

		go func(i int, email string) {
			defer wg.Done()

			provider, err := CheckProvider.GetEmailProvider(email)
			if err != nil {
				log.Printf("Error getting email provider for %s: %v", email, err)
				copyRowToNoMxFile(noMxFile, f, "NoMxRecords", "Sheet1", i+1, noMxRowCounter)
				noMxRowCounter++
			} else {
				f.SetCellValue("Sheet1", fmt.Sprintf("O%d", i+1), provider)
			}
			<-ch
		}(i, row[5])
	}

	wg.Wait()

	if err := f.SaveAs(outputFile); err != nil {
		return fmt.Errorf("Failed to save output file: %v", err)
	}

	if err := noMxFile.SaveAs(noMXOutputFile); err != nil {
		return fmt.Errorf("failed to save no MX output file: %v", err)
	}

	return nil

}

func copyRowToNoMxFile(dst, src *excelize.File, dstSheet, srcSheet string, srcRow, dstRow int) {
	for _, colName := range []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N"} {
		cellValue, err := src.GetCellValue(srcSheet, fmt.Sprintf("%s%d", colName, srcRow))
		if err != nil {
			fmt.Errorf("Failed to read Cell Value %v", err)
		}
		dst.SetCellValue(dstSheet, fmt.Sprintf("%s%d", colName, dstRow), cellValue)
	}
}
