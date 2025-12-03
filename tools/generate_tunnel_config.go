package tools

import (
	"encoding/json"
	"fmt"

	"github.com/bensonfx/mcp-liner/internal/config"
	"github.com/bensonfx/mcp-liner/internal/responses"
	"github.com/bensonfx/mcp-liner/internal/templates"
	"github.com/phuslu/log"
)

// GenerateTunnelConfigParams generate_tunnel_config工具的参数
type GenerateTunnelConfigParams struct {
	Role         string   `json:"role"`          // server 或 client
	Listen       []string `json:"listen"`        // 监听地址（server模式）
	ServerName   []string `json:"server_name"`   // 服务器名称（server模式）
	AuthTable    string   `json:"auth_table"`    // 认证表（server模式）
	RemoteListen []string `json:"remote_listen"` // 远程监听（client模式）
	ProxyPass    string   `json:"proxy_pass"`    // 代理目标（client模式）
	Dialer       string   `json:"dialer"`        // 拨号器名称（client模式）
	DialerURL    string   `json:"dialer_url"`    // 拨号器URL（client模式）
}

// GenerateTunnelConfig 生成内网穿透配置
func GenerateTunnelConfig(arguments json.RawMessage) (string, error) {
	var params GenerateTunnelConfigParams
	if err := json.Unmarshal(arguments, &params); err != nil {
		log.Error().Err(err).Msg("failed to parse parameters")
		return responses.ErrorResponse(
			fmt.Sprintf("Invalid parameters: %v", err),
			"Please provide valid parameters for tunnel configuration",
		)
	}

	log.Info().
		Str("role", params.Role).
		Msg("generating tunnel config")

	// 验证role参数
	if params.Role != "server" && params.Role != "client" {
		return responses.ErrorResponse(
			fmt.Sprintf("Invalid role: %s", params.Role),
			"Role must be either 'server' or 'client'",
		)
	}

	var cfg config.Config

	if params.Role == "server" {
		// 服务端配置
		if len(params.Listen) == 0 {
			params.Listen = []string{":443"}
		}
		if len(params.ServerName) == 0 {
			params.ServerName = []string{"tunnel.example.org"}
		}
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
		// 客户端配置
		if len(params.RemoteListen) == 0 {
			params.RemoteListen = []string{"127.0.0.1:10022"}
		}
		if params.ProxyPass == "" {
			params.ProxyPass = "127.0.0.1:22"
		}
		if params.Dialer == "" {
			params.Dialer = "cloud"
		}

		cfg = templates.SimpleTunnelConfig(
			"client",
			nil, nil, "",
			params.RemoteListen,
			params.ProxyPass,
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

	description := fmt.Sprintf("Generated tunnel %s configuration", params.Role)
	log.Info().Str("role", params.Role).Msg("tunnel config generated successfully")

	return responses.SuccessResponse(yamlContent, description)
}
