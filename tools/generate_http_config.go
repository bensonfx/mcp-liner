package tools

import (
	"encoding/json"
	"fmt"

	"github.com/bensonfx/mcp-liner/internal/config"
	"github.com/bensonfx/mcp-liner/internal/responses"
	"github.com/bensonfx/mcp-liner/internal/templates"
	"github.com/phuslu/log"
)

// GenerateHTTPConfigParams generate_http_config工具的参数
type GenerateHTTPConfigParams struct {
	Listen        []string `json:"listen"`         // 监听地址，如 [":443"]
	ServerName    []string `json:"server_name"`    // 服务器名称，如 ["example.com"]
	ForwardPolicy string   `json:"forward_policy"` // 转发策略，如 "proxy_pass"
	Dialer        string   `json:"dialer"`         // 拨号器名称，如 "local"
	DialerURL     string   `json:"dialer_url"`     // 拨号器URL（如果需要配置dialer）
	EnableTunnel  bool     `json:"enable_tunnel"`  // 是否启用tunnel功能
	AuthTable     string   `json:"auth_table"`     // 认证表（tunnel模式使用）
}

// GenerateHTTPConfig 生成HTTP/HTTPS配置
func GenerateHTTPConfig(arguments json.RawMessage) (string, error) {
	var params GenerateHTTPConfigParams
	if err := json.Unmarshal(arguments, &params); err != nil {
		log.Error().Err(err).Msg("failed to parse parameters")
		return responses.ErrorResponse(
			fmt.Sprintf("Invalid parameters: %v", err),
			"Please provide valid parameters for HTTP configuration",
		)
	}

	log.Info().
		Strs("listen", params.Listen).
		Strs("server_name", params.ServerName).
		Bool("enable_tunnel", params.EnableTunnel).
		Msg("generating HTTP config")

	// 设置默认值
	if len(params.Listen) == 0 {
		params.Listen = []string{":443"}
	}
	if len(params.ServerName) == 0 {
		params.ServerName = []string{"example.org"}
	}
	if params.ForwardPolicy == "" {
		params.ForwardPolicy = "proxy_pass"
	}
	if params.Dialer == "" {
		params.Dialer = "local"
	}

	var cfg config.Config

	if params.EnableTunnel {
		// 生成带tunnel的配置
		if params.AuthTable == "" {
			params.AuthTable = "auth_user.csv"
		}
		cfg = templates.SimpleTunnelConfig(
			"server",
			params.Listen,
			params.ServerName,
			params.AuthTable,
			nil, "", "", "",
		)
	} else {
		// 生成普通HTTP转发配置
		cfg = templates.SimpleHTTPForwardConfig(
			params.Listen,
			params.ServerName,
			params.Dialer,
			params.DialerURL,
		)
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

	description := "Generated HTTP/HTTPS configuration"
	if params.EnableTunnel {
		description = "Generated HTTP/HTTPS configuration with tunnel support"
	}

	log.Info().Msg("HTTP config generated successfully")
	return responses.SuccessResponse(yamlContent, description)
}
