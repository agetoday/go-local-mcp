package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"gocode/config"
	"gocode/internal"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	// 使用绝对路径读取config.yaml文件
	startTime := time.Now()
	conf, err := config.ParseYAMLConfig("config.yaml")
	if err != nil {
		log.Printf("Error reading config: %v\n", err)
		return
	}

	// 读取输入文件
	domains, err := readDomainsFromFile(conf.InputFile)
	if err != nil {
		log.Printf("Error reading input file: %v\n", err)
		return
	}

	// 解析域名
	results := internal.ResolveDomainsConcurrently(domains, conf.Workers)

	// 写入输出文件
	if err := writeResultsToFile(conf.OutputFile, results); err != nil {
		fmt.Printf("Error writing output file: %v\n", err)
		return
	}

	successCount := 0
	for _, r := range results {
		if r.Error == nil {
			successCount++
		}
	}
	fmt.Printf("Processed %d domains (%d success, %d failed) to %s\n",
		len(results), successCount, len(results)-successCount, conf.OutputFile)
	fmt.Printf("Total time: %v\n", time.Since(startTime))
}

func readDomainsFromFile(filePath string) ([]string, error) {
	fmt.Printf("Reading domains from %s\n", filePath)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var domains []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		domain := strings.TrimSpace(scanner.Text())
		if domain != "" {
			domains = append(domains, domain)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return domains, nil
}

func writeResultsToFile(filePath string, results []internal.DomainResult) error {
	file, err := os.Create(filePath + ".csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	if err := writer.Write([]string{"name", "ips", "error"}); err != nil {
		return err
	}

	for _, result := range results {
		var record []string
		record = append(record, result.Domain) // name field
		if result.Error != nil {
			record = append(record, "", result.Error.Error())
		} else {
			record = append(record, strings.Join(result.IPs, ", "), "")
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}
