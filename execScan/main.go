package main

import (
	"fmt"
	_ "gocode/config"
	"gocode/internal"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	// 从.env加载配置

	inputFile := os.Getenv("INPUT_FILE")
	outputFile := os.Getenv("OUTPUT_FILE")

	start, err := strconv.Atoi(os.Getenv("START_PORT"))
	if err != nil {
		start = 1
	}

	end, err := strconv.Atoi(os.Getenv("END_PORT"))
	if err != nil {
		end = 1024
	}

	workers, err := strconv.Atoi(os.Getenv("WORKERS"))
	if err != nil {
		workers = 100
	}

	timeout, err := time.ParseDuration(os.Getenv("TIMEOUT"))
	if err != nil {
		timeout = 1 * time.Second
	}

	// 读取IP列表
	var ips []string
	if inputFile != "" {
		var err error
		ips, err = internal.ReadIPsFromFile(inputFile)
		if err != nil {
			log.Printf("Error reading input file: %v\n", err)
			return
		}
	} else {
		ips = []string{inputFile}
	}

	// 扫描每个IP
	for _, ip := range ips {
		fmt.Printf("Scanning %s ports %d-%d with %d workers...\n",
			ip, start, end, workers)
		timeStart := time.Now()

		// 准备结果输出
		var resultWriter *os.File
		if outputFile != "" {
			dateStamp := time.Now().Format("20060102")
			timestamp := time.Now().Unix()
			filename := fmt.Sprintf("%d_%s.txt", timestamp, ip)
			fmt.Printf("Creating output file: %s in directory: %s\n",
				filename, internal.GetCurrentDirectory())
			var err error
			if _, err = os.Stat(outputFile); os.IsNotExist(err) {
				os.Mkdir(outputFile, 0755)
			}
			if _, err = os.Stat(outputFile + "/" + dateStamp); os.IsNotExist(err) {
				os.Mkdir(outputFile+"/"+dateStamp, 0755)
			}
			resultWriter, err = os.Create(outputFile + "/" + dateStamp + "/" + filename)
			if err != nil {
				fmt.Printf("Error creating output file '%s': %v\n", filename, err)
				continue
			}
			defer resultWriter.Close()
		}

		internal.ScanIP(ip, start, end, workers, timeout, resultWriter)
		timeEnd := time.Now()
		fmt.Printf("Scanning %s ports %d-%d with %d workers took %s\n",
			ip, start, end, workers, timeEnd.Sub(timeStart))
	}
}
