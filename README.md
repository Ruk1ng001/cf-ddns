# 动态 DNS 更新工具

> 一个使用 Go 编写的自动化动态 DNS 更新工具，支持多平台编译和发布。

## 项目简介

这个项目是一个动态 DNS 更新工具，旨在通过定期检查公网 IP 的变化，并自动更新 DNS 记录。该工具特别适合使用动态 IP 的用户，帮助他们保持域名与当前 IP 的同步。

## 功能特性

- 自动检测公网 IP 变化
- 自动更新 DNS 记录（支持 Cloudflare API）
- 支持多平台编译（Linux, macOS, Windows）

## 安装与使用

### 1. 获取可执行文件

你可以从 [GitHub Releases](https://github.com/Ruk1ng001/cf-ddns/releases) 页面下载适用于你平台的可执行文件。

### 2. 配置文件

在使用前，你需要创建一个配置文件 `config.json`，示例如下：

'''json
{
"api_token": "your_cloudflare_api_token",
"zone_id": "your_cloudflare_zone_id",
"record_id": "your_dns_record_id",
"domain": "yourdomain.com",
"check_interval": 10
}
'''

### 3. 运行工具

将配置文件放置在与可执行文件相同的目录下，然后运行：

'''bash
./cf-ddns -config config.json
'''

## 配合 Cron 运行

你可以使用 Cron 来定期运行该工具，以实现自动化的动态 DNS 更新。以下是一个每 10 分钟运行一次的示例：

1. 编辑 Crontab 文件：

'''bash
crontab -e
'''

2. 添加如下内容，将路径替换为你可执行文件的实际路径：

'''bash
*/10 * * * * /path/to/cf-ddns/cf-ddns.sh >> /var/log/cf-ddns.log 2>&1
'''

这行命令表示每 10 分钟运行一次工具，并将输出日志保存到 `/var/log/cf-ddns.log` 文件中。

3. 修改`cf-ddns.sh`脚本：

'''bash
API_TOKEN="your_api_token"
ZONE_ID="your_zone_id"
RECORD_ID="your_record_id"
DOMAIN="your.domain.com"
'''

## 构建与发布

### 1. 手动构建

如果你希望自行构建，确保你已安装 Go 环境，然后运行以下命令以构建不同平台的二进制文件：

'''bash
GOOS=linux GOARCH=amd64 go build -o dist/cf-ddns-linux-amd64
GOOS=darwin GOARCH=amd64 go build -o dist/cf-ddns-darwin-amd64
GOOS=windows GOARCH=amd64 go build -o dist/cf-ddns-windows-amd64.exe
'''

## changeLog

`v0.0.1` 初始版本

## 如何贡献

欢迎贡献代码和提出意见！请先 fork 本项目，创建一个新的分支进行开发，然后提交 Pull Request。
