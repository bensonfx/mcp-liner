// Package tools 提供MCP工具实现
package tools

import (
	"encoding/json"
	"fmt"

	"github.com/bensonfx/mcp-liner/internal/config"
	"github.com/bensonfx/mcp-liner/internal/responses"
	"github.com/bensonfx/mcp-liner/internal/templates"
	"github.com/bensonfx/mcp-liner/internal/validation"
	"github.com/phuslu/log"
)

// GenerateLinerConfigParams generate_liner_config工具的参数
type GenerateLinerConfigParams struct {
	Template string                 `json:"template"` // http_forward, tunnel, dns, full
	Params   map[string]interface{} `json:"params"`   // 模板参数
}

// GenerateLinerConfig 生成完整的liner配置文件
func GenerateLinerConfig(arguments json.RawMessage) (string, error) {
	var params GenerateLinerConfigParams
	if err := json.Unmarshal(arguments, &params); err != nil {
		log.Error().Err(err).Msg("failed to parse parameters")
		return responses.ErrorResponse(
			fmt.Sprintf("Invalid parameters: %v", err),
			"Please provide valid JSON parameters with 'template' and 'params' fields",
		)
	}

	log.Info().
		Str("template", params.Template).
		Msg("generating liner config")

	var cfg config.Config

	// 根据模板类型生成配置
	switch params.Template {
	case "http_forward":
		cfg = generateHTTPForwardConfig(params.Params)
	case "tunnel_server":
		cfg = generateTunnelServerConfig(params.Params)
	case "tunnel_client":
		cfg = generateTunnelClientConfig(params.Params)
	case "dns":
		cfg = generateDNSConfig(params.Params)
	case "full":
		cfg = generateFullConfig(params.Params)
	default:
		return responses.ErrorResponse(
			fmt.Sprintf("Unknown template: %s", params.Template),
			"Supported templates: http_forward, tunnel_server, tunnel_client, dns, full",
		)
	}

	// 验证配置
	validationResult := validation.ValidateConfig(&cfg)
	if !validationResult.Valid {
		log.Warn().Int("errors", len(validationResult.Errors)).Msg("config validation failed")
		return responses.ValidationResponse(validationResult)
	}

	// 转换为YAML
	yamlContent, err := cfg.ToYAML()
	if err != nil {
		log.Error().Err(err).Msg("failed to convert config to YAML")
		return responses.ErrorResponse(
			fmt.Sprintf("Failed to generate YAML: %v", err),
			"Please check your configuration parameters",
		)
	}

	log.Info().Msg("config generated successfully")
	return responses.SuccessResponse(yamlContent, fmt.Sprintf("Generated %s configuration", params.Template))
}

// generateHTTPForwardConfig 生成HTTP转发配置
func generateHTTPForwardConfig(params map[string]interface{}) config.Config {
	listen := getStringSlice(params, "listen", []string{":443"})
	serverName := getStringSlice(params, "server_name", []string{"example.org"})
	dialerName := getString(params, "dialer", "local")
	dialerURL := getString(params, "dialer_url", "")

	return templates.SimpleHTTPForwardConfig(listen, serverName, dialerName, dialerURL)
}

// generateTunnelServerConfig 生成隧道服务端配置
func generateTunnelServerConfig(params map[string]interface{}) config.Config {
	listen := getStringSlice(params, "listen", []string{":443"})
	serverName := getStringSlice(params, "server_name", []string{"example.org"})
	authTable := getString(params, "auth_table", "auth_user.csv")

	return templates.SimpleTunnelConfig(
		"server",
		listen,
		serverName,
		authTable,
		nil, "", "", "",
	)
}

// generateTunnelClientConfig 生成隧道客户端配置
func generateTunnelClientConfig(params map[string]interface{}) config.Config {
	remoteListen := getStringSlice(params, "remote_listen", []string{"127.0.0.1:10022"})
	proxyPass := getString(params, "proxy_pass", "127.0.0.1:22")
	dialerName := getString(params, "dialer", "cloud")
	dialerURL := getString(params, "dialer_url", "")

	return templates.SimpleTunnelConfig(
		"client",
		nil, nil, "",
		remoteListen,
		proxyPass,
		dialerName,
		dialerURL,
	)
}

// generateDNSConfig 生成DNS配置
func generateDNSConfig(params map[string]interface{}) config.Config {
	listen := getStringSlice(params, "listen", []string{":53"})
	proxyPass := getString(params, "proxy_pass", "https://8.8.8.8/dns-query")

	return templates.SimpleDNSConfig(listen, proxyPass)
}

// generateFullConfig 生成完整配置
func generateFullConfig(params map[string]interface{}) config.Config {
	// 使用默认的完整配置模板
	global := config.NewDefaultGlobalConfig()

	// 允许自定义全局配置
	if logLevel, ok := params["log_level"].(string); ok {
		global.LogLevel = logLevel
	}
	if dnsServer, ok := params["dns_server"].(string); ok {
		global.DnsServer = dnsServer
	}
	if disableHttp3, ok := params["disable_http3"].(bool); ok {
		global.DisableHttp3 = disableHttp3
	}

	return config.Config{
		Global: global,
		Dialer: map[string]string{},
	}
}

// 辅助函数：从params中获取字符串
func getString(params map[string]interface{}, key string, defaultValue string) string {
	if val, ok := params[key]; ok {
		if strVal, ok := val.(string); ok {
			return strVal
		}
	}
	return defaultValue
}

// 辅助函数：从params中获取字符串数组
func getStringSlice(params map[string]interface{}, key string, defaultValue []string) []string {
	if val, ok := params[key]; ok {
		// 处理 JSON 解析后的 []interface{} 类型
		if arr, ok := val.([]interface{}); ok {
			result := make([]string, 0, len(arr))
			for _, item := range arr {
				if strVal, ok := item.(string); ok {
					result = append(result, strVal)
				}
			}
			if len(result) > 0 {
				return result
			}
		}
		// 处理直接的 []string 类型
		if strArr, ok := val.([]string); ok {
			return strArr
		}
	}
	return defaultValue
}
