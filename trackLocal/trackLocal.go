package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// NetworkStats 存储网络接口统计信息
type NetworkStats struct {
	InterfaceName string
	BytesSent     uint64
	BytesRecv     uint64
	Timestamp     time.Time
}

func main() {
	fmt.Println("Starting local traffic monitor...")

	// 创建日志目录
	if err := os.MkdirAll("logs", 0755); err != nil {
		fmt.Printf("Error creating logs directory: %v\n", err)
		return
	}

	// 初始化日志文件
	logFile, err := os.OpenFile("logs/traffic.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening log file: %v\n", err)
		return
	}
	defer logFile.Close()

	// 初始化前一次统计信息map
	prevStatsMap := make(map[string]NetworkStats)

	// 获取网络接口列表
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("Error getting network interfaces: %v\n", err)
		return
	}

	// 监控循环
	for {
		for _, iface := range interfaces {
			// 跳过回环接口和非活动接口
			if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
				continue
			}

			// 获取接口统计信息
			stats, err := getInterfaceStats(iface.Name)
			if err != nil {
				fmt.Printf("Error getting stats for interface %s: %v\n", iface.Name, err)
				continue
			}

			// 获取所有连接信息
			connections, err := getConnectionsForInterface(iface.Name)
			if err != nil {
				fmt.Printf("Error getting connections for interface %s: %v\n", iface.Name, err)
				continue
			}
			// 初始化CSV日志文件
			dateTime := time.Now().Format("2006010215")
			if _, err := os.Stat("recCsv"); os.IsNotExist(err) {
				os.Mkdir("recCsv", 0755)
			}
			csvFile, err := os.OpenFile("recCsv/"+dateTime+".csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Printf("Error opening CSV log file: %v\n", err)
				continue
			}

			// 定义常见端口白名单
			commonPorts := map[string]bool{
				"80": true, "443": true, "22": true, "21": true, "25": true,
				"53": true, "110": true, "143": true, "3306": true, "3389": true,
			}

			// 记录完整的连接信息到CSV和日志
			for _, conn := range connections {
				// 写入CSV: 时间,本地IP,远端IP,端口,进程名,PID
				csvEntry := fmt.Sprintf("%s,%s,%s,%s,%s,%s\n",
					time.Now().Format("2006-01-02 15:04:05"),
					conn.LocalIP,
					conn.RemoteIP,
					conn.Port,
					conn.Process,
					conn.PID)

				if _, err := csvFile.WriteString(csvEntry); err != nil {
					fmt.Printf("Error writing to CSV file: %v\n", err)
				}

				// 写入日志文件
				logEntry := fmt.Sprintf("[CONNECTION] %s -> %s:%s, Process: %s (PID: %s)\n",
					conn.LocalIP,
					conn.RemoteIP,
					conn.Port,
					conn.Process,
					conn.PID)

				if _, err := logFile.WriteString(logEntry); err != nil {
					fmt.Printf("Error writing to log file: %v\n", err)
				}

				// 检测非标准端口
				if !commonPorts[conn.Port] {
					warningMsg := fmt.Sprintf("[SECURITY WARNING] Uncommon port %s used by %s (PID: %s)\n",
						conn.Port, conn.Process, conn.PID)
					fmt.Print(warningMsg)
					if _, err := logFile.WriteString(warningMsg); err != nil {
						fmt.Printf("Error writing warning to log: %v\n", err)
					}
				}
			}

			// 在所有写入完成后关闭文件
			if err := csvFile.Close(); err != nil {
				fmt.Printf("Error closing CSV file: %v\n", err)
			}

			// 检测异常流量

			// 更新前一次统计信息
			prevStatsMap[iface.Name] = stats

			// 记录统计信息到日志文件
			logEntry := fmt.Sprintf("[%s] %s: Sent %d bytes, Received %d bytes\n",
				stats.Timestamp.Format("2006-01-02 15:04:05"),
				stats.InterfaceName,
				stats.BytesSent,
				stats.BytesRecv)

			if _, err := logFile.WriteString(logEntry); err != nil {
				fmt.Printf("Error writing to log file: %v\n", err)
			}

			// 打印统计信息到控制台
			// fmt.Print("====================================================", logEntry)
		}

		// 每5秒采样一次
		time.Sleep(5 * time.Second)
	}
}

// parseBytesFromNetsh 从netsh命令输出中解析字节数
func parseBytesFromNetsh(output []byte, metric string) uint64 {
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, metric) {
			// 更精确的解析逻辑
			colonIndex := strings.Index(line, ":")
			if colonIndex == -1 {
				continue
			}

			valueStr := strings.TrimSpace(line[colonIndex+1:])
			parts := strings.Fields(valueStr)
			if len(parts) == 0 {
				continue
			}

			// 提取数字部分，去除可能的分隔符
			numStr := strings.ReplaceAll(parts[0], ",", "")
			if bytes, err := strconv.ParseUint(numStr, 10, 64); err == nil {
				return bytes
			}
		}
	}
	return 0
}

// getInterfaceStats 获取指定网络接口的统计信息
func getInterfaceStats(ifaceName string) (NetworkStats, error) {
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return NetworkStats{}, err
	}

	stats := NetworkStats{
		InterfaceName: iface.Name,
		Timestamp:     time.Now(),
	}

	// 使用更可靠的命令获取网络统计
	cmd := exec.Command("netsh", "interface", "ipv4", "show", "interface", fmt.Sprintf("name=\"%s\"", ifaceName))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return stats, fmt.Errorf("failed to get interface stats: %v", err)
	}

	// 解析发送和接收的字节数
	stats.BytesSent = parseBytesFromNetsh(output, "Bytes Sent")
	stats.BytesRecv = parseBytesFromNetsh(output, "Bytes Received")

	// 如果仍然获取不到数据，尝试备用方法
	if stats.BytesSent == 0 && stats.BytesRecv == 0 {
		cmd = exec.Command("netstat", "-e")
		output, err = cmd.CombinedOutput()
		if err == nil {
			stats.BytesSent = parseBytesFromNetsh(output, "Bytes Sent")
			stats.BytesRecv = parseBytesFromNetsh(output, "Bytes Received")
		}
	}

	return stats, nil
}

// detectAnomalies 检测流量异常

// ConnectionInfo 存储连接信息
type ConnectionInfo struct {
	LocalIP     string
	RemoteIP    string
	Port        string
	PID         string
	Process     string
	CreateTime  string
	DNSResolved bool
	Hostname    string
}

// getConnectionsForInterface 获取指定网络接口的所有连接信息
func getConnectionsForInterface(ifaceName string) ([]ConnectionInfo, error) {
	// 使用netstat命令获取详细网络连接信息
	// 设置控制台代码页为UTF-8
	chcpCmd := exec.Command("cmd", "/c", "chcp", "65001")
	if err := chcpCmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to set UTF-8 codepage: %v", err)
	}

	// 执行netstat命令并捕获UTF-8格式输出
	cmd := exec.Command("cmd", "/c", "netstat", "-ano")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to get network connections: %v", err)
	}

	// 将输出转换为UTF-8字符串
	outputStr := string(output)

	var connections []ConnectionInfo

	lines := strings.Split(outputStr, "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if strings.Contains(line, "SYN_SENT") || strings.Contains(line, "SYN_RECV") || strings.Contains(line, "ESTABLISHED") {
			if len(fields) >= 5 {
				// 解析本地和远程地址
				addrParts := strings.Split(fields[1], ":")
				remoteParts := strings.Split(fields[2], ":")

				if len(addrParts) >= 2 && len(remoteParts) >= 2 {
					info := ConnectionInfo{
						LocalIP:    addrParts[0],
						RemoteIP:   remoteParts[0],
						Port:       remoteParts[1],
						PID:        fields[len(fields)-1],
						CreateTime: time.Now().Format("2006-01-02 15:04:05.000"),
					}
					// 跨平台获取进程名
					var procCmd *exec.Cmd
					switch runtime.GOOS {
					case "windows":
						// Windows平台使用修正后的tasklist命令
						procCmd = exec.Command("tasklist", "/FI", fmt.Sprintf("PID eq %s", info.PID), "/FO", "CSV", "/NH")
					case "linux", "darwin":
						// Linux/macOS平台使用ps命令
						procCmd = exec.Command("ps", "-p", info.PID, "-o", "comm=")
					default:
						info.Process = "unknown"
						continue
					}

					procOutput, err := procCmd.CombinedOutput()
					if err == nil {
						if runtime.GOOS == "windows" {
							// Windows平台输出格式为CSV
							procInfo := strings.Split(string(procOutput), ",")
							if len(procInfo) > 0 {
								info.Process = strings.Trim(procInfo[0], "\"")
							}
						} else {
							// Linux/macOS平台直接获取命令名
							info.Process = strings.TrimSpace(string(procOutput))
						}
					} else {
						info.Process = "unknown"
					}

					// 尝试DNS解析
					hostnames, err := net.LookupAddr(info.RemoteIP)
					if err == nil && len(hostnames) > 0 {
						info.DNSResolved = true
						info.Hostname = hostnames[0]
					} else {
						info.DNSResolved = false
						info.Hostname = "unresolved"
					}

					connections = append(connections, info)
				}
			}
		}
	}

	return connections, nil
}
