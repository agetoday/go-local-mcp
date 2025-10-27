package internal

import (
	"bufio"
	"fmt"
	"os"
)

// 从文件读取IP列表
func ReadIPsFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var ips []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ip := scanner.Text()
		if ip != "" {
			ips = append(ips, ip)
		}
	}
	return ips, scanner.Err()
}

// 获取当前工作目录
func GetCurrentDirectory() string {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Sprintf("(error getting current directory: %v)", err)
	}
	return dir
}
