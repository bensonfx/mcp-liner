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
	Listen         []string `json:"listen"`          // 监听地址，如 [":443"]
	ServerName     []string `json:"server_name"`     // 服务器名称，如 ["example.com"]
	ForwardPolicy  string   `json:"forward_policy"`  // 转发策略，如 "proxy_pass"
	PolicyTemplate string   `json:"policy_template"` // Policy模板（go template），如果提供则覆盖forward_policy
	Dialer         string   `json:"dialer"`          // 拨号器名称，如 "local"
	DialerURL      string   `json:"dialer_url"`      // 拨号器URL（如果需要配置dialer）
	EnableTunnel   bool     `json:"enable_tunnel"`   // 是否启用tunnel功能
	AuthTable      string   `json:"auth_table"`      // 认证表（tunnel模式使用）
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
	if len(params.ServerName) == 0 || (len(params.ServerName) == 1 && params.ServerName[0] == "*") {
		// Fix: don't use * as server name, use example.org if not provided or invalid
		if len(params.ServerName) == 0 {
			params.ServerName = []string{"example.org"}
		}
	}
	if params.ForwardPolicy == "" {
		params.ForwardPolicy = "proxy_pass"
	}
	if params.Dialer == "" {
		params.Dialer = "local"
	}

	// 如果提供了policy_template，使用它覆盖forward_policy
	policyValue := params.ForwardPolicy
	if params.PolicyTemplate != "" {
		policyValue = params.PolicyTemplate
	}

	var cfg config.Config

	if params.EnableTunnel {
		// 生成带tunnel的配置
		if params.AuthTable == "" {
			params.AuthTable = "auth_user.csv"
		}
		// 使用 TunnelServerTemplate 生成配置
		tunnelConfig := templates.TunnelServerTemplate(
			params.Listen,
			params.ServerName,
			params.AuthTable,
			[]string{"127.0.0.1", "240.0.0.0/8"},
		)

		cfg = config.Config{
			Global: config.NewDefaultGlobalConfig(),
			Dialer: map[string]string{},
			Https:  []config.HTTPConfig{tunnelConfig},
		}

	} else {
		// 生成普通HTTP转发配置
		httpConfig := templates.HTTPForwardTemplate(params.Listen, params.ServerName, params.Dialer, true)
		// 设置policy
		httpConfig.Forward.Policy = policyValue

		cfg = config.Config{
			Global: config.NewDefaultGlobalConfig(),
			Dialer: map[string]string{},
			Https:  []config.HTTPConfig{httpConfig},
		}

		// 如果提供了dialer URL，添加到配置
		if params.DialerURL != "" && params.Dialer != "local" {
			cfg.Dialer[params.Dialer] = params.DialerURL
		}
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
	if params.PolicyTemplate != "" {
		description += " with custom policy template"
	}

	log.Info().Msg("HTTP config generated successfully")
	return responses.SuccessResponse(yamlContent, description)
}
