package mcpTools

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/tealeg/xlsx/v3"
)

// ReadExcelFiles 读取指定路径下的Excel文件
func ProcessFiles(dirPath string) ([][]string, error) {
	results := make([][]string, 0) // 修改这里，将results定义为[][]string类型的切片
	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Println("read dir err !")
		return nil, err
	}
	for _, file := range files {
		fmt.Println("dir have file: ", file.Name())
		if filepath.Ext(file.Name()) == ".xlsx" {
			data, err := readExcelFile(filepath.Join(dirPath, file.Name()))
			if err != nil {
				log.Println("read xlsx err !", filepath.Join(dirPath, file.Name()), err)
				return nil, err
			}
			results = append(results, data) // 这里不需要修改
		}
	}
	return results, nil
}

func readExcelFile(filePath string) ([]string, error) {
	// 验证文件格式
	if !strings.HasSuffix(filePath, ".xlsx") {
		return nil, fmt.Errorf("invalid file format, expected .xlsx")
	}
	// 打开Excel文件
	xlFile, err := xlsx.FileToSlice(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open Excel file: %v ,%v", filePath, err)
	}
	var results []string
	// 读取第一个工作表
	sheet := xlFile[0] // 直接使用返回的切片
	for _, row := range sheet {
		var rowData []string
		for _, cell := range row {
			rowData = append(rowData, cell)
		}
		results = append(results, strings.Join(rowData, ","))
	}
	return results, nil
}

// ConvertXmlToJson 将XML文件转换为JSON并保存
func ConvertXmlToJson(xmlPath, jsonDir string) error {
	xmlData, err := ioutil.ReadFile(xmlPath)
	if err != nil {
		return err
	}

	var data interface{}
	if err := xml.Unmarshal(xmlData, &data); err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(jsonDir, os.ModePerm); err != nil {
		return err
	}

	jsonPath := filepath.Join(jsonDir, filepath.Base(xmlPath)+".json")
	return ioutil.WriteFile(jsonPath, jsonData, 0644)
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
