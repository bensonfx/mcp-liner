package tools

import (
	"encoding/json"
	"fmt"

	"github.com/bensonfx/mcp-liner/internal/responses"
	"github.com/bensonfx/mcp-liner/internal/templates"
	"github.com/phuslu/log"
)

// GenerateDNSConfigParams generate_dns_config工具的参数
type GenerateDNSConfigParams struct {
	Listen    []string `json:"listen"`     // 监听地址，如 [":53"]
	Policy    string   `json:"policy"`     // DNS策略，如 "forward"
	ProxyPass string   `json:"proxy_pass"` // 上游DNS服务器
	CacheSize int      `json:"cache_size"` // 缓存大小
}

// GenerateDNSConfig 生成DNS配置
func GenerateDNSConfig(arguments json.RawMessage) (string, error) {
	var params GenerateDNSConfigParams
	if err := json.Unmarshal(arguments, &params); err != nil {
		log.Error().Err(err).Msg("failed to parse parameters")
		return responses.ErrorResponse(
			fmt.Sprintf("Invalid parameters: %v", err),
			"Please provide valid parameters for DNS configuration",
		)
	}

	log.Info().
		Strs("listen", params.Listen).
		Str("policy", params.Policy).
		Str("proxy_pass", params.ProxyPass).
		Msg("generating DNS config")

	// 设置默认值
	if len(params.Listen) == 0 {
		params.Listen = []string{":53"}
	}
	if params.ProxyPass == "" {
		params.ProxyPass = "https://8.8.8.8/dns-query"
	}

	// 生成DNS配置
	cfg := templates.SimpleDNSConfig(params.Listen, params.ProxyPass)

	// 转换为YAML
	yamlContent, err := cfg.ToYAML()
	if err != nil {
		log.Error().Err(err).Msg("failed to convert config to YAML")
		return responses.ErrorResponse(
			fmt.Sprintf("Failed to generate YAML: %v", err),
			"",
		)
	}

	log.Info().Msg("DNS config generated successfully")
	return responses.SuccessResponse(yamlContent, "Generated DNS configuration")
}
