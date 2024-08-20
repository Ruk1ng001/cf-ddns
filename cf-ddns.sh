#!/bin/bash

# 配置参数
API_TOKEN="your_api_token"
ZONE_ID="your_zone_id"
RECORD_ID="your_record_id"
DOMAIN="your.domain.com"

# 存储上次的IP的文件
LAST_IP_FILE="/tmp/last_ip.txt"

# 获取当前公网IP
get_public_ip() {
    curl -s https://api.ipify.org
}

# 更新Cloudflare DNS记录
update_dns_record() {
    local ip=$1
    local update_url="https://api.cloudflare.com/client/v4/zones/${ZONE_ID}/dns_records/${RECORD_ID}"

    local data=$(cat <<EOF
{
    "type": "A",
    "name": "$DOMAIN",
    "content": "$ip",
    "ttl": 120,
    "proxied": true
}
EOF
)

    curl -s -X PUT "$update_url" \
        -H "Authorization: Bearer $API_TOKEN" \
        -H "Content-Type: application/json" \
        --data "$data"
}

# 初始化或读取上次的公网IP
if [ -f "$LAST_IP_FILE" ]; then
    last_ip=$(cat "$LAST_IP_FILE")
else
    last_ip=""
fi

# 获取当前IP
current_ip=$(get_public_ip)

# 检查IP是否变化
if [ "$current_ip" != "$last_ip" ]; then
    echo "[$(date)] 公网IP发生变化: $last_ip -> $current_ip"

    response=$(update_dns_record "$current_ip")

    if echo "$response" | grep -q '"success":true'; then
        echo "[$(date)] DNS记录更新成功"
        echo "$current_ip" > "$LAST_IP_FILE"
    else
        echo "[$(date)] 更新DNS记录失败: $response"
    fi
else
    echo "[$(date)] 公网IP未变化: $current_ip"
fi
