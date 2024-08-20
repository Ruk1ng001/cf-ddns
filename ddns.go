package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Config 配置结构体
type Config struct {
	APIToken      string `json:"api_token"`
	ZoneID        string `json:"zone_id"`
	RecordID      string `json:"record_id"`
	Domain        string `json:"domain"`
	CheckInterval int    `json:"check_interval"` // IP 检查时间间隔（分钟）
}

// CloudflareResponse Cloudflare API 的响应结构体
type CloudflareResponse struct {
	Success bool `json:"success"`
	Errors  []struct {
		Message string `json:"message"`
	} `json:"errors"`
	Result struct {
		ID      string `json:"id"`
		Content string `json:"content"`
	} `json:"result"`
}

func main() {
	log("程序启动")

	// 检查是否提供了配置文件路径
	if len(os.Args) < 2 {
		log("错误: 请提供配置文件路径")
		return
	}
	configFilePath := os.Args[1]

	// 读取配置
	config, err := loadConfig(configFilePath)
	if err != nil {
		log(fmt.Sprintf("读取配置文件出错: %v", err))
		return
	}

	var lastIP string

	// 程序启动时立即执行一次 IP 更新
	currentIP, err := getPublicIP()
	if err != nil {
		log(fmt.Sprintf("获取公网 IP 出错: %v", err))
	} else {
		// 检查 IP 是否变化并更新 DNS 记录
		if currentIP != lastIP {
			log(fmt.Sprintf("首次启动，更新 DNS 记录: %s", currentIP))
			err = updateDNSRecord(config, currentIP)
			if err != nil {
				log(fmt.Sprintf("更新 DNS 记录出错: %v", err))
			} else {
				log("DNS 记录更新成功")
				lastIP = currentIP
			}
		}
	}

	// 使用配置的检查时间间隔
	checkInterval := time.Duration(config.CheckInterval) * time.Minute
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	// 定时检测 IP 变化
	for {
		select {
		case <-ticker.C:
			// 获取当前公网 IP
			currentIP, err := getPublicIP()
			if err != nil {
				log(fmt.Sprintf("获取公网 IP 出错: %v", err))
				continue
			}

			// 检查 IP 是否发生变化
			if currentIP != lastIP {
				log(fmt.Sprintf("公网 IP 发生变化: %s -> %s", lastIP, currentIP))

				// 更新 Cloudflare DNS 记录
				err = updateDNSRecord(config, currentIP)
				if err != nil {
					log(fmt.Sprintf("更新 DNS 记录出错: %v", err))
				} else {
					log("DNS 记录更新成功")
					lastIP = currentIP
				}
			}
		}
	}
}

// 加载配置文件
func loadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &Config{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return nil, err
	}
	log("配置文件加载成功")
	return config, nil
}

// 获取公网 IP
func getPublicIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(ip), nil
}

// 更新 Cloudflare DNS 记录
func updateDNSRecord(config *Config, ip string) error {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", config.ZoneID, config.RecordID)

	// 构建请求体
	data := map[string]interface{}{
		"type":    "A",
		"name":    config.Domain,
		"content": ip,
		"ttl":     120,
		"proxied": true,
		"comment": "DDNS: " + time.Now().Format("2006-01-02 15:04:05"),
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+config.APIToken)
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 解析响应
	var cfResp CloudflareResponse
	if err := json.NewDecoder(resp.Body).Decode(&cfResp); err != nil {
		return err
	}

	if !cfResp.Success {
		return fmt.Errorf("cloudflare API 错误: %v", cfResp.Errors)
	}
	log(fmt.Sprintf("Cloudflare DNS 记录更新成功，新的 IP: %s", ip))
	return nil
}

// 简单的日志记录函数，添加时间戳
func log(message string) {
	fmt.Printf("[%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), message)
}
