package tools

import (
	"encoding/json"
	"fmt"

	"github.com/bensonfx/mcp-liner/internal/config"
	"github.com/bensonfx/mcp-liner/internal/responses"
	"github.com/bensonfx/mcp-liner/internal/templates"
	"github.com/phuslu/log"
)

// GenerateDNSConfigParams generate_dns_config工具的参数
type GenerateDNSConfigParams struct {
	Listen         []string `json:"listen"`          // 监听地址，如 [":53"]
	ProxyPass      string   `json:"proxy_pass"`      // 上游DNS服务器，如 "https://8.8.8.8/dns-query"
	PolicyTemplate string   `json:"policy_template"` // Policy模板（go template），用于自定义DNS路由
	CacheSize      int      `json:"cache_size"`      // DNS缓存大小
	Log            bool     `json:"log"`             // 是否启用日志
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
		Str("proxy_pass", params.ProxyPass).
		Bool("log", params.Log).
		Msg("generating DNS config")

	// 设置默认值
	if len(params.Listen) == 0 {
		params.Listen = []string{":53"}
	}
	if params.ProxyPass == "" {
		params.ProxyPass = "https://8.8.8.8/dns-query"
	}
	if params.CacheSize == 0 {
		params.CacheSize = 4096
	}

	// 构建DNS配置
	dnsConfig := templates.DNSForwardTemplate(params.Listen, params.ProxyPass)
	dnsConfig.CacheSize = params.CacheSize
	dnsConfig.Log = params.Log

	// 如果提供了policy template，使用它
	if params.PolicyTemplate != "" {
		dnsConfig.Policy = params.PolicyTemplate
	}

	cfg := config.Config{
		Global: config.NewDefaultGlobalConfig(),
		Dns:    []config.DnsConfig{dnsConfig},
	}

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
