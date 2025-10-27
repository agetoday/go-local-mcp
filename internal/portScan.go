package internal

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

// 扫描单个端口
func ScanPort(ip string, port int, timeout time.Duration, results chan<- int) {
	target := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", target, timeout)
	if err == nil {
		conn.Close()
		results <- port
	}
}

// 扫描单个IP的所有端口
func ScanIP(ip string, start, end, workers int, timeout time.Duration, resultWriter *os.File) {
	var wg sync.WaitGroup
	results := make(chan int)
	totalPorts := end - start + 1
	progress := make(chan int, totalPorts)
	openPorts := make([]int, 0)

	// 启动结果收集器
	go func() {
		for port := range results {
			openPorts = append(openPorts, port)
			msg := fmt.Sprintf("%s:%d is open\n", ip, port)
			fmt.Print(msg)
			if resultWriter != nil {
				resultWriter.WriteString(msg)
			}
		}
	}()

	// 启动进度显示器
	go func() {
		scanned := 0
		for range progress {
			scanned++
			fmt.Printf("\rScanning progress: %d/%d (%.1f%%)", scanned, totalPorts, float64(scanned)/float64(totalPorts)*100)
		}
		fmt.Println() // 换行
	}()

	// 创建worker池
	ports := make(chan int, workers*2)
	for i := 0; i < workers; i++ {
		go func() {
			for port := range ports {
				ScanPort(ip, port, timeout, results)
				progress <- port
				wg.Done()
			}
		}()
	}

	// 分发端口扫描任务
	for port := start; port <= end; port++ {
		wg.Add(1)
		ports <- port
	}

	wg.Wait()
	close(results)
	close(ports)
	close(progress)

	// 对发现的开放端口进行nmap服务识别
	if len(openPorts) > 0 && resultWriter != nil {
		fmt.Println("\nStarting nmap service detection for open ports...")
		for i, port := range openPorts {
			wg.Add(1)
			go func(i int, port int) {
				defer wg.Done()
				fmt.Printf("\rProcessing port %d/%d", i+1, len(openPorts))
				err := NmapScan(ip, port, port, resultWriter)
				if err != nil {
					fmt.Printf("\nNmap scan failed for port %d: %v\n", port, err)
					// 继续扫描其他端口
				}
			}(i, port)
		}
		wg.Wait()
		fmt.Println("\nNmap service detection completed")
	}
}

