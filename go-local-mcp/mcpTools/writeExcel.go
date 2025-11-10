package mcpTools

import (
	"fmt"
	"os"
	"time"

	"github.com/tealeg/xlsx/v3"
)

func WriteExcelFile(content [][]string) error {

	// Open file for writing
	dirPath := "output"
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		os.Mkdir(dirPath, 0755)
	}
	timeTamp := time.Now().Format("20060102_150405")
	filePath := dirPath + fmt.Sprintf("/%s.xlsx", timeTamp)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a new Excel file
	xlFile := xlsx.NewFile()

	// Create a new sheet in the Excel file
	sheet, err := xlFile.AddSheet("Sheet1")
	if err != nil {
		return err
	}

	// Write data to the sheet
	for _, rowData := range content {
		row := sheet.AddRow()
		for _, cellData := range rowData {
			cell := row.AddCell()
			cell.Value = cellData
		}
	}

	return xlFile.Save(file.Name())
}
