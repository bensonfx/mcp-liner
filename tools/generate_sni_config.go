package tools

import (
	"encoding/json"
	"fmt"

	"github.com/bensonfx/mcp-liner/internal/config"
	"github.com/bensonfx/mcp-liner/internal/responses"
	"github.com/phuslu/log"
)

// GenerateSniConfigParams generate_sni_config工具的参数
type GenerateSniConfigParams struct {
	Enabled     bool   `json:"enabled"`      // 是否启用SNI转发
	Policy      string `json:"policy"`       // 转发策略（可以是go template）
	Dialer      string `json:"dialer"`       // 拨号器名称
	DialerURL   string `json:"dialer_url"`   // 拨号器URL（可选）
	DisableIpv6 bool   `json:"disable_ipv6"` // 禁用IPv6
	PreferIpv6  bool   `json:"prefer_ipv6"`  // 优先IPv6
	Log         bool   `json:"log"`          // 是否启用日志
}

// GenerateSniConfig 生成SNI配置
func GenerateSniConfig(arguments json.RawMessage) (string, error) {
	var params GenerateSniConfigParams
	if err := json.Unmarshal(arguments, &params); err != nil {
		log.Error().Err(err).Msg("failed to parse parameters")
		return responses.ErrorResponse(
			fmt.Sprintf("Invalid parameters: %v", err),
			"Please provide valid parameters for SNI configuration",
		)
	}

	log.Info().
		Bool("enabled", params.Enabled).
		Str("policy", params.Policy).
		Str("dialer", params.Dialer).
		Bool("log", params.Log).
		Msg("generating SNI config")

	// 设置默认值
	if params.Dialer == "" {
		params.Dialer = "local"
	}
	if params.Policy == "" {
		params.Policy = "proxy_pass"
	}

	// 构建配置
	cfg := config.Config{
		Global: config.NewDefaultGlobalConfig(),
		Dialer: map[string]string{
			"local": "local",
		},
		Sni: config.SniConfig{
			Enabled: params.Enabled,
			Forward: config.SniForwardConfig{
				Policy:      params.Policy,
				Dialer:      params.Dialer,
				DisableIpv6: params.DisableIpv6,
				PreferIpv6:  params.PreferIpv6,
				Log:         params.Log,
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

	description := "Generated SNI-based routing configuration\n\n"
	description += "SNI (Server Name Indication) routing intercepts TLS ClientHello and routes based on SNI.\n\n"
	description += "Features:\n"
	description += "- Routes traffic before TLS handshake completes\n"
	description += "- Works with encrypted HTTPS traffic\n"
	description += "- Lower overhead than HTTP proxy\n"
	description += "- Supports policy templates for dynamic routing\n\n"

	if params.Policy != "" && params.Policy != "proxy_pass" {
		description += "Policy template example:\n"
		description += "  {{ if hasSuffixes \"google.com|youtube.com\" .ServerName }}google_dialer{{ else }}direct{{ end }}\n\n"
		description += "Available context:\n"
		description += "  .ServerName: SNI server name\n"
		description += "  .ClientHello: TLS ClientHello info\n"
		description += "  .RemoteAddr: Client address\n"
	}

	log.Info().Msg("SNI config generated successfully")
	return responses.SuccessResponse(yamlContent, description)
}
