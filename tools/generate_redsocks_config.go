package tools

import (
	"encoding/json"
	"fmt"

	"github.com/bensonfx/mcp-liner/internal/config"
	"github.com/bensonfx/mcp-liner/internal/responses"
	"github.com/phuslu/log"
)

// GenerateRedsocksConfigParams generate_redsocks_config工具的参数
type GenerateRedsocksConfigParams struct {
	Listen    []string `json:"listen"`     // 监听地址，如 [":12345"]
	Dialer    string   `json:"dialer"`     // 拨号器名称，如 "proxy"
	DialerURL string   `json:"dialer_url"` // 拨号器URL（可选）
	Log       bool     `json:"log"`        // 是否启用日志
}

// GenerateRedsocksConfig 生成Redsocks透明代理配置
func GenerateRedsocksConfig(arguments json.RawMessage) (string, error) {
	var params GenerateRedsocksConfigParams
	if err := json.Unmarshal(arguments, &params); err != nil {
		log.Error().Err(err).Msg("failed to parse parameters")
		return responses.ErrorResponse(
			fmt.Sprintf("Invalid parameters: %v", err),
			"Please provide valid parameters for Redsocks configuration",
		)
	}

	log.Info().
		Strs("listen", params.Listen).
		Str("dialer", params.Dialer).
		Bool("log", params.Log).
		Msg("generating redsocks config")

	// 设置默认值
	if len(params.Listen) == 0 {
		params.Listen = []string{":12345"}
	}
	if params.Dialer == "" {
		params.Dialer = "proxy"
	}

	// 构建配置
	cfg := config.Config{
		Global: config.NewDefaultGlobalConfig(),
		Dialer: map[string]string{},
		Redsocks: []config.RedsocksConfig{
			{
				Listen: params.Listen,
				Forward: config.RedsocksForwardConfig{
					Dialer: params.Dialer,
					Log:    params.Log,
				},
			},
		},
	}

	// 如果提供了dialer URL，添加到配置
	if params.DialerURL != "" && params.Dialer != "local" {
		cfg.Dialer[params.Dialer] = params.DialerURL
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

	description := "Generated Redsocks transparent proxy configuration\n\n"
	description += "⚠️  IMPORTANT: This configuration requires:\n"
	description += "1. Linux operating system (redsocks uses SO_ORIGINAL_DST)\n"
	description += "2. iptables rules to redirect traffic to redsocks port\n"
	description += "3. Use generate_redsocks_iptables tool to generate firewall rules\n\n"
	description += "Example usage:\n"
	description += "1. Save this config to liner.yaml\n"
	description += "2. Generate iptables rules with generate_redsocks_iptables\n"
	description += "3. Apply iptables rules: sudo iptables-restore < rules.v4\n"
	description += "4. Start liner: sudo ./liner -c liner.yaml\n"

	log.Info().Msg("redsocks config generated successfully")
	return responses.SuccessResponse(yamlContent, description)
}
