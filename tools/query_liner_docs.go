package tools

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bensonfx/mcp-liner/internal/responses"
	"github.com/phuslu/log"
)

// QueryLinerDocsParams query_liner_docs工具的参数
type QueryLinerDocsParams struct {
	Topic string `json:"topic"` // global, http, tunnel, dns, dialer, policy
}

// QueryLinerDocs 查询liner文档和使用说明
func QueryLinerDocs(arguments json.RawMessage) (string, error) {
	var params QueryLinerDocsParams
	if err := json.Unmarshal(arguments, &params); err != nil {
		log.Error().Err(err).Msg("failed to parse parameters")
		return responses.ErrorResponse(
			fmt.Sprintf("Invalid parameters: %v", err),
			"Please provide 'topic' parameter",
		)
	}

	log.Info().Str("topic", params.Topic).Msg("querying liner docs")

	// 根据topic返回相应文档
	var content string
	switch strings.ToLower(params.Topic) {
	case "global":
		content = getGlobalDocs()
	case "http":
		content = getHTTPDocs()
	case "tunnel":
		content = getTunnelDocs()
	case "dns":
		content = getDNSDocs()
	case "dialer":
		content = getDialerDocs()
	case "policy":
		content = getPolicyDocs()
	default:
		return responses.ErrorResponse(
			fmt.Sprintf("Unknown topic: %s", params.Topic),
			"Supported topics: global, http, tunnel, dns, dialer, policy",
		)
	}

	log.Info().Str("topic", params.Topic).Msg("docs retrieved successfully")
	return responses.DocumentationResponse(params.Topic, content)
}

// getGlobalDocs 获取全局配置文档
func getGlobalDocs() string {
	doc := `全局配置用于设置liner的基础运行参数

常用配置项:
- log_level: 日志级别 (info/debug/warn/error)
- dns_server: DNS服务器地址，支持DoH (如: https://8.8.8.8/dns-query)
- disable_http3: 是否禁用HTTP3，建议设为false
- dial_timeout: 拨号超时时间（秒）

示例:
global:
  log_level: info
  dns_server: https://8.8.8.8/dns-query
  disable_http3: false
`
	return doc
}

// getHTTPDocs 获取HTTP配置文档
func getHTTPDocs() string {
	doc := `HTTP/HTTPS配置用于设置HTTP服务器，支持转发、隧道等功能

基本配置:
- listen: 监听地址列表，如 [":443"]
- server_name: 服务器名称（HTTPS必需），如 ["example.org"]

Forward转发示例:
https:
  - listen: [":443"]
    server_name: ["example.org"]
    forward:
      policy: proxy_pass
      dialer: local
      log: true

Tunnel服务端配置示例:
https:
  - listen: [":443"]
    server_name: ["tunnel.example.org"]
    tunnel:
      enabled: true
      auth_table: auth_user.csv
      allow_listens: ["127.0.0.1"]
      log: true
`
	return doc
}

// getTunnelDocs 获取隧道配置文档
func getTunnelDocs() string {
	doc := fmt.Sprintf(`Tunnel功能用于实现内网穿透

服务端配置（公网服务器）:
https:
  - listen: [":443"]
    server_name: ["tunnel.example.org"]
    tunnel:
      enabled: true
      auth_table: auth_user.csv
      allow_listens: ["127.0.0.1"]

客户端配置（内网机器）:
dialer:
  cloud: https://user:pass@tunnel.example.org

tunnel:
  - remote_listen: ['127.0.0.1:10022']
    proxy_pass: '127.0.0.1:22'
    dialer: cloud
    dial_timeout: 5
    log: true

使用: ssh -p 10022 user@公网服务器IP
`)
	return doc
}

// getDNSDocs 获取DNS配置文档
func getDNSDocs() string {
	doc := `DNS服务配置支持DNS/DoT/DoH

基本配置:
dns:
  - listen: [":53"]
    policy: forward
    proxy_pass: https://8.8.8.8/dns-query
    cache_size: 4096
    log: true

上游服务器支持:
- DoH: https://8.8.8.8/dns-query
- DoT: tls://8.8.8.8:853
- UDP: 8.8.8.8:53
`
	return doc
}

// getDialerDocs 获取拨号器配置文档
func getDialerDocs() string {
	doc := `Dialer定义如何连接到上游服务器

支持的类型:
1. local - 本地直连
2. socks5 - SOCKS5代理
3. http2 - HTTP/2代理
4. http3 - HTTP/3代理
5. ssh - SSH隧道
6. wss - WebSocket Secure

配置示例:
dialer:
  local: local
  proxy1: socks5://proxy.example.com:1080
  proxy2: http2://user:pass@proxy.example.com:443
  cloud: https://username:password@tunnel.example.org

使用示例:
https:
  - listen: [":443"]
    server_name: ["example.org"]
    forward:
      policy: proxy_pass
      dialer: cloud
      log: true
`
	return doc
}

// getPolicyDocs 获取策略配置文档
func getPolicyDocs() string {
	doc := fmt.Sprintf(`Policy是liner的核心特性，使用Go template实现灵活转发

简单用法:
forward:
  policy: proxy_pass
  dialer: local

固定目标:
forward:
  policy: http://backend:8080
  dialer: local

Go Template条件转发:
forward:
  policy: |
    {{if .Host | contains "api"}}
      http://api-server:8080
    {{else}}
      http://web-server:8080
    {{end}}
  dialer: local

可用变量:
- .Host - 请求的Host头
- .Method - HTTP方法
- .URL - 请求URL
- .RemoteAddr - 客户端地址
`)
	return doc
}
