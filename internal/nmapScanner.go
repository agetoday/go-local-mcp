package internal

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// NmapScan 使用nmap扫描指定IP和端口范围
func NmapScan(ip string, startPort, endPort int, resultWriter *os.File) error {
	// 构造nmap命令
	args := []string{
		"-sV", // 服务版本检测
		"-T4", // 时序模板，平衡速度和准确性
		"-p", fmt.Sprintf("%d-%d", startPort, endPort),
		ip,
	}

	// 执行nmap命令
	cmd := exec.Command("nmap", args...)

	// 捕获标准输出和错误输出
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("nmap scan failed: %v\n%s", err, stderr.String())
	}

	// 解析并输出结果
	results := parseNmapOutput(stdout.String())
	for _, result := range results {
		msg := fmt.Sprintf("%s\n", result)
		fmt.Print(msg)
		if resultWriter != nil {
			if _, err := resultWriter.WriteString(msg); err != nil {
				return fmt.Errorf("failed to write results: %v", err)
			}
		}
	}

	return nil
}

// parseNmapOutput 解析nmap输出结果
func parseNmapOutput(output string) []string {
	var results []string
	lines := strings.Split(output, "\n")
	var currentPort string

	for _, line := range lines {
		// 匹配端口行，如: "22/tcp open  ssh"
		if strings.Contains(line, "/tcp") || strings.Contains(line, "/udp") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				portInfo := strings.Split(parts[0], "/")
				if len(portInfo) == 2 {
					currentPort = fmt.Sprintf("%s:%s %s", parts[len(parts)-1], portInfo[0], portInfo[1])
					results = append(results, currentPort)
				}
			}
		} else if strings.HasPrefix(line, "|") && currentPort != "" {
			// 匹配服务详细信息行，如: "| ssh-hostkey: 2048 SHA256:..."
			serviceInfo := strings.TrimSpace(strings.TrimPrefix(line, "|"))
			results = append(results, fmt.Sprintf("  %s", serviceInfo))
		}
	}

	return results
}
