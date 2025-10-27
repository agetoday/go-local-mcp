package internal

import (
	"fmt"
	"net"
	"regexp"
	"sync"
)

var (
	domainRegex = regexp.MustCompile(`^([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$`)
	ipPortRegex = regexp.MustCompile(`^(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}):\d+$`)
	urlRegex    = regexp.MustCompile(`^(https?://)?([^/:]+)(:\d+)?(/.*)?$`)
)

// PreprocessInput 预处理输入数据
// 支持域名、IP:port和URL格式
func PreprocessInput(input string) (string, error) {
	// 尝试匹配URL格式
	if matches := urlRegex.FindStringSubmatch(input); len(matches) > 2 {
		host := matches[2] // 提取主机部分
		if domainRegex.MatchString(host) {
			ips, err := ResolveDomain(host)
			if err != nil {
				return "", err
			}
			if len(ips) == 0 {
				return "", fmt.Errorf("no IP found for domain %s", host)
			}
			return ips[0], nil
		}
		return host, nil // 直接返回IP部分
	}

	if domainRegex.MatchString(input) {
		ips, err := ResolveDomain(input)
		if err != nil {
			return "", err
		}
		if len(ips) == 0 {
			return "", fmt.Errorf("no IP found for domain %s", input)
		}
		return ips[0], nil
	}

	if matches := ipPortRegex.FindStringSubmatch(input); len(matches) > 1 {
		return matches[1], nil
	}

	return "", fmt.Errorf("invalid input format, expected domain, IP:port or URL")
}

// ResolveDomain 快速解析域名到IP地址列表
func ResolveDomain(domain string) ([]string, error) {
	return net.LookupHost(domain)
}

// DomainResult 表示域名解析结果
type DomainResult struct {
	Domain string
	IPs    []string
	Error  error
}

// ResolveDomainsConcurrently 并发解析多个域名
// 返回所有域名的解析结果，不会因为单个域名失败而中断
func ResolveDomainsConcurrently(domains []string, workers int) []DomainResult {
	var wg sync.WaitGroup
	sem := make(chan struct{}, workers) // 并发控制信号量
	results := make([]DomainResult, len(domains))
	var mu sync.Mutex

	for i, domain := range domains {
		wg.Add(1)
		sem <- struct{}{} // 获取信号量

		go func(idx int, d string) {
			defer func() {
				<-sem // 释放信号量
				wg.Done()
			}()

			var result DomainResult
			result.Domain = d

			// 预处理输入
			processed, e := PreprocessInput(d)
			if e != nil {
				result.Error = e
			} else {
				// 解析域名
				ips, e := net.LookupHost(processed)
				if e != nil {
					result.Error = e
				} else {
					result.IPs = ips
				}
			}

			mu.Lock()
			results[idx] = result
			mu.Unlock()
		}(i, domain)
	}

	wg.Wait()
	return results
}
