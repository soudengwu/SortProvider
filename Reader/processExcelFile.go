package ExcelReader

import (
	"fmt"
	"log"
	"sync"

	CheckProvider "github.com/soudengwu/SortProvider"
	"github.com/xuri/excelize/v2"
)

func ProcessExcelFile(inputFile, outputFile string) error {
	f, err := excelize.OpenFile(inputFile)
	if err != nil {
		return fmt.Errorf("Failed to open file %v", err)
	}

	f.SetCellValue("Sheet1", "01", "Email Provider")

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return fmt.Errorf("failed to get rows: %v", err)
	}

	var wg sync.WaitGroup
	ch := make(chan bool, 5)

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

	return nil

}
