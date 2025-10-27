package internal

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 解析IP范围 (如192.168.1.1-254)
func parseIPRange(ipRange string) ([]string, error) {
	parts := strings.Split(ipRange, ".")
	if len(parts) != 4 {
		return nil, fmt.Errorf("invalid IP range format")
	}

	lastPart := parts[3]
	if !strings.Contains(lastPart, "-") {
		return []string{ipRange}, nil
	}

	rangeParts := strings.Split(lastPart, "-")
	if len(rangeParts) != 2 {
		return nil, fmt.Errorf("invalid IP range format")
	}

	start, err := strconv.Atoi(rangeParts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid start IP")
	}

	end, err := strconv.Atoi(rangeParts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid end IP")
	}

	if start > end {
		return nil, fmt.Errorf("start IP must be <= end IP")
	}

	var ips []string
	base := strings.Join(parts[:3], ".")
	for i := start; i <= end; i++ {
		ips = append(ips, fmt.Sprintf("%s.%d", base, i))
	}
	return ips, nil
}

// 实现内网存活主机扫描
func ScanIPisWorking(ipRange string, timeout time.Duration, resultWriter *os.File) {
	// 解析IP范围
	ips, err := parseIPRange(ipRange)
	if err != nil {
		fmt.Printf("Error parsing IP range: %v\n", err)
		return
	}

	var wg sync.WaitGroup
	aliveHosts := make(chan string)

	// 结果收集器
	go func() {
		for host := range aliveHosts {
			msg := fmt.Sprintf("%s is alive\n", host)
			fmt.Print(msg)
			if resultWriter != nil {
				resultWriter.WriteString(msg)
			}
		}
	}()

	// 进度显示
	total := len(ips)
	progress := make(chan int, total)
	go func() {
		scanned := 0
		for range progress {
			scanned++
			fmt.Printf("\rScanning progress: %d/%d (%.1f%%)", scanned, total, float64(scanned)/float64(total)*100)
		}
		fmt.Println()
	}()

	// 并发ping检测
	for _, ip := range ips {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			if ping(ip, timeout) {
				aliveHosts <- ip
			}
			progress <- 1
		}(ip)
	}

	wg.Wait()
	close(aliveHosts)
	close(progress)
}

// 执行ping命令检测主机存活
func ping(ip string, timeout time.Duration) bool {
	cmd := exec.Command("ping", "-n", "1",  strconv.Itoa(int(timeout/time.Millisecond)), ip)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}

	// 检查ping输出是否包含"TTL="表示主机存活
	return strings.Contains(string(output), "TTL=")
}
