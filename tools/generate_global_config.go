package tools

import (
	"encoding/json"
	"fmt"

	"github.com/bensonfx/mcp-liner/internal/config"
	"github.com/bensonfx/mcp-liner/internal/responses"
	"github.com/phuslu/log"
)

// GenerateGlobalConfigParams generate_global_config工具的参数
type GenerateGlobalConfigParams struct {
	LogLevel     string `json:"log_level"`     // info, debug, warn, error
	DnsServer    string `json:"dns_server"`    // DNS服务器地址
	DisableHttp3 bool   `json:"disable_http3"` // 是否禁用HTTP3
	DialTimeout  int    `json:"dial_timeout"`  // 拨号超时（秒）
}

// GenerateGlobalConfig 生成liner全局配置
func GenerateGlobalConfig(arguments json.RawMessage) (string, error) {
	var params GenerateGlobalConfigParams
	if err := json.Unmarshal(arguments, &params); err != nil {
		log.Error().Err(err).Msg("failed to parse parameters")
		return responses.ErrorResponse(
			fmt.Sprintf("Invalid parameters: %v", err),
			"Please provide valid parameters for global configuration",
		)
	}

	log.Info().
		Str("log_level", params.LogLevel).
		Str("dns_server", params.DnsServer).
		Msg("generating global config")

	// 创建全局配置
	globalCfg := config.NewDefaultGlobalConfig()

	// 应用自定义参数
	if params.LogLevel != "" {
		globalCfg.LogLevel = params.LogLevel
	}
	if params.DnsServer != "" {
		globalCfg.DnsServer = params.DnsServer
	}
	globalCfg.DisableHttp3 = params.DisableHttp3
	if params.DialTimeout > 0 {
		globalCfg.DialTimeout = params.DialTimeout
	}

	// 创建完整配置（仅包含global部分）
	cfg := config.Config{
		Global: globalCfg,
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

	log.Info().Msg("global config generated successfully")
	return responses.SuccessResponse(yamlContent, "Generated global configuration")
}
